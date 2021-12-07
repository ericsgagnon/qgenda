package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	// "golang.org/x/net/html"
)

func main() {

	// load json file
	f, err := os.Open("../samples/qgenda_restapi.postman_collection.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fc, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	pc := PostmanCollection{}
	if err := json.Unmarshal(fc, &pc); err != nil {
		log.Fatal(err)
	}

	// recursively copy raw definitions to maps
	items := map[string]map[string]string{}
	for _, v := range pc.Item {
		itemArray := map[string]string{}
		for _, u := range v.Item {
			name := scrubName(u.Name)

			itemArray[name] = u.Description
		}
		items[v.Name] = itemArray
	}

	// let's work on the 'response objects'
	// responsObjects["object name"]["field name"]["jsonType"/"description"/etc]"value"
	// responseObjects := map[string]map[string]map[string]string{}
	// for k, v := range items["Response Objects"] {
	// 	if !regexp.MustCompile(`(?i)(Encounter(Phi)?FieldSetting)|(TaskLocation)`).MatchString(k) {
	// 		// fmt.Printf("%+v\n", k)
	// 		responseObjects[k] = parseResponseObject(v)
	// 		// fmt.Printf("%+v\n", responseObjects[k])
	// 	}
	// }

	responseObjectSlice := map[string][]map[string]string{}
	for k, v := range items["Response Objects"] {
		if !regexp.MustCompile(`(?i)(Encounter(Phi)?FieldSetting)|(TaskLocation)`).MatchString(k) {
			// fmt.Printf("%+v\n", k)
			responseObjectSlice[k] = parseResponseObjectToSlice(v)
			// fmt.Printf("%+v\n", responseObjects[k])
		}
	}

	jsonTypes := map[string]string{}
	for ro, flds := range responseObjectSlice {
		fmt.Sprintf("------------------------------------------\n%s:\n", ro)
		for fld, attributes := range flds {
			fmt.Sprintf("\t%3d\t%v\n", fld, attributes)
			jsonTypes[attributes["jsonType"]] = attributes["goType"]
		}
	}
	for jsonType, goType := range jsonTypes {
		fmt.Sprintf("%30s\t%s\n", jsonType, goType)
	}

	tpl := `
	{{- printf "package qgenda\n\n" -}}
	{{- range $roName, $roValue := . -}}
	{{- printf "\n\ntype %s struct {\n" $roName -}}
		{{- range $i, $field := $roValue -}}
		{{- printf "\t%s\t%s\n" (index $field "name") (index $field "goType") }}		
		{{- end -}}
	{{- printf "}\n" -}}
	{{- end -}}
	`
	var buf bytes.Buffer
	t := template.Must(template.New("letter").Parse(tpl))
	if err := t.Execute(&buf, responseObjectSlice); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(buf.String())

	os.WriteFile("generated/qgenda.go", buf.Bytes(), os.ModePerm)
	// goCode, err := format.Source(buf.Bytes())
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Print(goCode)
	// for roName, ro := range responseObjects {
	// 	fmt.Printf(`%s:\t%s\n`, roName, ro)
	// }
}

type stringWrapper string

func (w stringWrapper) regexReplaceAll(match string, replacement string) stringWrapper {
	v := regexp.MustCompile(match).ReplaceAllString(string(w), replacement)
	return stringWrapper(v)

}

func (w stringWrapper) toString() string {
	return string(w)
}

func parseResponseObject(rawDef string) map[string]map[string]string {

	// cleanup and 'standardize' the format - qgenda's docs
	// mix and match html and md tables for their definitions in docs
	wipDef := stringWrapper(rawDef)
	wipDef = wipDef.
		regexReplaceAll(`(?i)</?table>|<td>|<tr>|</?strong>|</?code>|\n|\t|</?p>|</?br ?/?>|</?tbody>`, "").
		// regexReplaceAll(`(?i)</?table>|<td>|<tr>|</?strong>|</?code>|\n|\t|<p>|<br>|</?tbody>`, "").
		regexReplaceAll(`\s\s+`, " ").
		regexReplaceAll("`?\\*\\*`?", "").
		regexReplaceAll("`", "").
		regexReplaceAll(`^\s*\|\s*`, "").
		regexReplaceAll(`\s*\|\s*$`, "").
		regexReplaceAll(`\s*\|\s*`, "|").
		regexReplaceAll("</td>", "|").
		regexReplaceAll("</tr>", "\n").
		regexReplaceAll("&nbsp;", "").
		regexReplaceAll("&emsp;", "").
		regexReplaceAll("&gt;", "").
		regexReplaceAll("&lt;", "").
		regexReplaceAll(`\|\|`, "\n").
		regexReplaceAll(`(?im)^(\s*)|(\|*)$`, "").
		regexReplaceAll(`(?i)(^\|?\s*name\s*\|\s*type\s*\|\s*description.*)|(^\|?\s*-+\s*\|\s*-+\s*(\|\s*-+\s*\|?.*)?)`, "").
		regexReplaceAll(`:?\s*-+\s*:?\|:?\s*-*\s*:?\|:?\s*-*\s*:?\|?\s*`, "").
		regexReplaceAll(`(?im)^\s*\|\s*`, "").
		regexReplaceAll(`(?im)^([^|]*\|[^|]*)$`, "$1|")

	rows := strings.Split(string(wipDef), "\n")
	def := map[string]map[string]string{}
	fldNum := 0
	for _, v := range rows {
		rowDef := map[string]string{}
		vs := strings.Split(v, "|")
		if len(vs) < 3 {
			continue
		}
		name := stringWrapper(vs[0]).
			regexReplaceAll(`^\s+`, "").
			regexReplaceAll(`\s+$`, "").
			regexReplaceAll(`:`, "").
			regexReplaceAll(`[*]`, "").
			toString()
		jsonType := stringWrapper(strings.ToLower(vs[1])).
			regexReplaceAll(`^\s+`, "").
			regexReplaceAll(`\s+$`, "").
			regexReplaceAll(`\*|:`, "").
			regexReplaceAll(`[<*>]`, "").
			regexReplaceAll(`integer`, "int").
			regexReplaceAll(`(guid)/?(u?uid)?`, "string").
			regexReplaceAll(`string/?string`, "string").
			regexReplaceAll(`time/timespan`, "time").
			regexReplaceAll(`date/time`, "date").
			regexReplaceAll(``, "").
			toString()

		// jsonType = string(jsonType)
		if strings.Contains(jsonType, "rray") {
			rowDef["jsonType"] = fmt.Sprintf("%s[%s]", jsonType, name)
		} else {
			rowDef["jsonType"] = jsonType
		}
		rowDef["description"] = strings.TrimSpace(vs[2])
		rowDef["structTag"] = fmt.Sprintf("`json:\"%s\"`", name)
		rowDef["goType"] = mapTypeJSONToGo(jsonType)
		rowDef["fieldNumber"] = fmt.Sprint(fldNum)
		fldNum++
		def[name] = rowDef
		// if jsonType == "" {
		// 	fmt.Printf("%s:\t%s\n", name, rowDef)
		// }
		// fmt.Printf("%25s:\t%25s\t%s\n", name, rowDef["jsonType"], rowDef["description"])
		// fmt.Printf("%+v:\t%+v\n", i, strings.Split(v, "|"))
	}

	return def
}

func scrubName(s string) string {
	return stringWrapper(s).
		regexReplaceAll(`\s`, "").
		regexReplaceAll(`(?i)\(\s*for\s*User\s*\)`, "").
		toString()
}

func parseResponseObjectToSlice(rawDef string) []map[string]string {

	// cleanup and 'standardize' the format - qgenda's docs
	// mix and match html and md tables for their definitions in docs
	wipDef := stringWrapper(rawDef)
	wipDef = wipDef.
		regexReplaceAll(`(?i)</?table>|<td>|<tr>|</?strong>|</?code>|\n|\t|</?p>|</?br ?/?>|</?tbody>`, "").
		// regexReplaceAll(`(?i)</?table>|<td>|<tr>|</?strong>|</?code>|\n|\t|<p>|<br>|</?tbody>`, "").
		regexReplaceAll(`\s\s+`, " ").
		regexReplaceAll("`?\\*\\*`?", "").
		regexReplaceAll("`", "").
		regexReplaceAll(`^\s*\|\s*`, "").
		regexReplaceAll(`\s*\|\s*$`, "").
		regexReplaceAll(`\s*\|\s*`, "|").
		regexReplaceAll("</td>", "|").
		regexReplaceAll("</tr>", "\n").
		regexReplaceAll("&nbsp;", "").
		regexReplaceAll("&emsp;", "").
		regexReplaceAll("&gt;", "").
		regexReplaceAll("&lt;", "").
		regexReplaceAll(`\|\|`, "\n").
		regexReplaceAll(`(?im)^(\s*)|(\|*)$`, "").
		regexReplaceAll(`(?i)(^\|?\s*(property)?\s*name\s*\|\s*value\s*\|\s*description.*)|(^\|?\s*-+\s*\|\s*-+\s*(\|\s*-+\s*\|?.*)?)`, "").
		regexReplaceAll(`(?i)(^\|?\s*name\s*\|\s*type\s*\|\s*description.*)|(^\|?\s*-+\s*\|\s*-+\s*(\|\s*-+\s*\|?.*)?)`, "").
		regexReplaceAll(`:?\s*-+\s*:?\|:?\s*-*\s*:?\|:?\s*-*\s*:?\|?\s*`, "").
		regexReplaceAll(`(?im)^\s*\|\s*`, "").
		regexReplaceAll(`(?im)^([^|]*\|[^|]*)$`, "$1|")

	rows := strings.Split(string(wipDef), "\n")
	def := []map[string]string{}
	// fldNum := 0
	for _, v := range rows {
		rowDef := map[string]string{}
		vs := strings.Split(v, "|")
		if len(vs) < 3 {
			continue
		}
		name := stringWrapper(vs[0]).
			regexReplaceAll(`^\s+`, "").
			regexReplaceAll(`\s+$`, "").
			regexReplaceAll(`:`, "").
			regexReplaceAll(`[*]`, "").
			regexReplaceAll(`(?i)\(for User\)`, "").
			toString()

		jsonType := stringWrapper(strings.ToLower(vs[1])).
			regexReplaceAll(`^\s+`, "").
			regexReplaceAll(`\s+$`, "").
			regexReplaceAll(`\*|:`, "").
			regexReplaceAll(`[<*>]`, "").
			regexReplaceAll(`integer`, "int").
			regexReplaceAll(`(guid)/?(u?uid)?`, "string").
			regexReplaceAll(`string/?string`, "string").
			regexReplaceAll(`time/timespan`, "time").
			regexReplaceAll(`date/time`, "date").
			regexReplaceAll(`(?i).*string.*`, "[]string").
			regexReplaceAll(`(?i)\[\]`, "").
			regexReplaceAll(`(?i)\(array\)`, "").
			regexReplaceAll(`(?im)^[^[]+\[(?P<test>[^]]+)\].*$`, "[]$test").
			toString()

		if strings.Contains(jsonType, "rray") {
			jsonType = fmt.Sprintf("%s[%s]", jsonType, name)
		}

		// jsonType = string(jsonType)
		// if strings.Contains(jsonType, "rray") {
		// 	rowDef["jsonType"] = fmt.Sprintf("%s[%s]", jsonType, name)
		// } else {
		// 	rowDef["jsonType"] = jsonType
		// }
		jsonType = stringWrapper(jsonType).
			regexReplaceAll(`^\s+`, "").
			regexReplaceAll(`\s+$`, "").
			regexReplaceAll(`\*|:`, "").
			regexReplaceAll(`[<*>]`, "").
			regexReplaceAll(`integer`, "int").
			regexReplaceAll(`(guid)/?(u?uid)?`, "string").
			regexReplaceAll(`string/?string`, "string").
			regexReplaceAll(`time/timespan`, "time").
			regexReplaceAll(`date/time`, "date").
			regexReplaceAll(`(?i)array.*string.*`, "[]string").
			regexReplaceAll(`(?i)\[\]`, "").
			regexReplaceAll(`(?i)\(array\)`, "").
			regexReplaceAll(`(?im)^[^[]+\[(?P<test>[^]]+)\].*$`, "[]$test").
			// regexReplaceAll(`(?im).*`, "hahahaha").
			toString()

		rowDef["jsonType"] = jsonType
		rowDef["description"] = strings.TrimSpace(vs[2])
		rowDef["structTag"] = fmt.Sprintf("`json:\"%s\"`", name)
		rowDef["goType"] = mapTypeJSONToGo(jsonType)
		rowDef["name"] = name
		def = append(def, rowDef)
		// if jsonType == "" {
		// 	fmt.Printf("%s:\t%s\n", name, rowDef)
		// }
		// fmt.Printf("%25s:\t%25s\t%s\n", name, rowDef["jsonType"], rowDef["description"])
		// fmt.Printf("%+v:\t%+v\n", i, strings.Split(v, "|"))
	}

	return def
}

func mapTypeJSONToGo(s string) string {
	s = stringWrapper(s).
		toString()
	typeMap := map[string]string{
		"int":       "int",
		"string":    "string",
		"boolean":   "boolean",
		"decimal":   "float64",
		"number":    "float64",
		"time":      "time.Time",
		"date":      "time.Time",
		"datetime":  "time.Time",
		"timestamp": "time.Time",
	}
	if _, ok := typeMap[s]; ok {
		return typeMap[s]
	}
	return s
}

// PostmanCollection represent the json structure of postman's collection format
// https://schema.postman.com/collection/json/v2.1.0/draft-07/docs/index.html
// https://schema.postman.com/collection/json/v2.1.0/draft-07/collection.json
type PostmanCollection struct {
	Info     map[string]string
	Item     []PostmanItem
	Event    []PostmanEvent
	Variable interface{}
	Auth     interface{}
}

// PostmanItem is another small step for dealing with Postman
type PostmanItem struct {
	Name        string
	Item        []PostmanItem
	Description string
	Event       []PostmanEvent
}

type PostmanEvent struct {
	Listen string
	Script PostmanEventScript
}

type PostmanEventScript struct {
	Type string
	Exec []string
}

// ResponseFieldDefinition is the definition of a field/element
// of a qgenda response object as defined by its api docs
type ResponseFieldDefinition struct {
	Name        string
	Type        string
	Description string
}

// ResponseObjectDefinition is the definition of a qgenda response
// object as defined by its api docs
type ResponseObjectDefinition []ResponseFieldDefinition

func ParseResponseObjects(items ...PostmanItem) (map[string]ResponseObjectDefinition, error) {

	ros := map[string]ResponseObjectDefinition{}

	return ros, nil
}
