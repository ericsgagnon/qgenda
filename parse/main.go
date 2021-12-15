package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	// "golang.org/x/net/html"

	"gopkg.in/yaml.v3"
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
		switch v.Name {
		case "Response Objects":
			itemArray := map[string]string{}
			for _, u := range v.Item {
				name := scrubName(u.Name)
				itemArray[name] = u.Description
			}
			items[v.Name] = itemArray
			// fmt.Println(items)
		case "API Calls":
			itemArray := map[string]string{}
			for _, u := range v.Item {
				name := scrubName(u.Name)
				fmt.Sprintf("%s:\n%s\n", name, u)

				// itemArray[name] = u.Description
			}
			items[v.Name] = itemArray
		case "Introduction":
			fmt.Println(v.Name)
		case "Appendix":
			fmt.Println(v.Name)
		}
	}

	responseObjectSlice := map[string][]map[string]string{}
	for k, v := range items["Response Objects"] {
		if !regexp.MustCompile(`(?i)(Encounter(Phi)?FieldSetting)|(TaskLocation)`).MatchString(k) {
			// fmt.Printf("%+v\n", k)
			responseObjectSlice[k] = parseResponseObjectToSlice(v)
			// fmt.Printf("%+v\n", responseObjects[k])
		}
	}
	// fmt.Println(items["API Calls"])
	// jsonTypes := map[string]string{}
	// for ro, flds := range responseObjectSlice {
	// 	fmt.Sprintf("------------------------------------------\n%s:\n", ro)
	// 	for fld, attributes := range flds {
	// 		fmt.Sprintf("\t%3d\t%v\n", fld, attributes)
	// 		jsonTypes[attributes["jsonType"]] = attributes["goType"]
	// 	}
	// }
	// for jsonType, goType := range jsonTypes {
	// 	fmt.Sprintf("%30s\t%s\n", jsonType, goType)
	// }

	tpl, err := os.ReadFile(`qgenda.tmpl`)
	if err != nil {
		log.Fatalln(err)
	}

	// the following are either unimplement or inaccessible by our login
	delete(responseObjectSlice, "Profile")
	delete(responseObjectSlice, "RequestApproved")
	delete(responseObjectSlice, "RequestLimit")
	delete(responseObjectSlice, "DailyConfiguration")
	delete(responseObjectSlice, "DailyCase")
	delete(responseObjectSlice, "Room")
	delete(responseObjectSlice, "PayRate")
	delete(responseObjectSlice, "TimeEvent")
	delete(responseObjectSlice, "StaffLocation")
	// delete(responseObjectSlice, "Location")
	delete(responseObjectSlice, "StaffTarget")
	delete(responseObjectSlice, "User")
	delete(responseObjectSlice, "NotificationList")
	delete(responseObjectSlice, "StaffMemberDetail")
	delete(responseObjectSlice, "TagDetailsByCompany")

	var buf bytes.Buffer
	t := template.Must(template.New("letter").Parse(string(tpl)))
	if err := t.Execute(&buf, responseObjectSlice); err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(buf.String())

	os.WriteFile("generated/qgenda.go", buf.Bytes(), os.ModePerm)
	// goCode, err := format.Source(buf.Bytes())
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Print(goCode)
	// for roName, ro := range responseObjects {
	// 	fmt.Printf(`%s:\t%s\n`, roName, ro)
	// }
	// fmt.Println(strings.Title(`we known how me you`))

	// apiCalls := map[string][]map[string]string{}
	// for k, v := range items["API Calls"] {
	// 	fmt.Printf("%s:\t%v\n", k, v)
	// }
	// fmt.Println(items["API Calls"])

	// m := make(map[interface{}]interface{})
	// m := make(map[string]interface{})

	data, err := os.ReadFile("../samples/qgenda_restapi.postman_collection.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	err = yaml.Unmarshal(data, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// fmt.Printf("--- m:\n%+v\n\n", m.Alias)
	fmt.Println(m.Content[0].Content[0])

}

type Config struct {
	Info struct {
		PostmanID   string `yaml:"_postman_id"`
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		Schema      string `yaml:"schema"`
	} `yaml:"info"`
	APICalls []APICall
	Item     []struct {
		Name        string        `yaml:""`
		Item        []interface{} `yaml:""`
		Description string        `yaml:""`
		Event       []interface{} `yaml:""`
	} `yaml:"item"`
	Auth struct {
		Type   string `yaml:"type"`
		Bearer []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
			Type  string `yaml:"type"`
		} `yaml:"bearer"`
	}
	Event []struct {
		Listen string   `yam:"listen"`
		Script string   `yam:"script"`
		Type   string   `yam:"type"`
		Exec   []string `yam:"exec"`
	} `yaml:"event"`
}

type APICall struct {
	Name string
	Item []APICallItem
}

type APICallItem struct {
	Name     string
	Request  APICAllRequest
	Response APICallResponse
}

type APICAllRequest struct {
	Method string
	Header []KeyTypeValue
	URL    APIURL
}

type APIURL struct {
	Raw         string
	Protocol    string
	Host        []string
	Path        []string
	Query       []KeyTypeValue
	Description string
}

type KeyTypeValue struct {
	Key   string
	Type  string
	Value string
}

type APICallResponse struct {
}

// PostmanItem is another small step for dealing with Postman
type PostmanItem struct {
	Name        string
	Item        []PostmanItem
	Description string
	Event       []PostmanEvent
	Request     interface{}
	Response    interface{}
}

type ConfigEvent struct {
	Listen string   `yam:"listen"`
	Script string   `yam:"script"`
	Type   string   `yam:"type"`
	Exec   []string `yam:"exec"`
}

type PostmanEventScript struct {
	Type string
	Exec []string
}

type stringWrapper string

func (w stringWrapper) regexReplaceAll(match string, replacement string) stringWrapper {
	v := regexp.MustCompile(match).ReplaceAllString(string(w), replacement)
	return stringWrapper(v)

}

func (w stringWrapper) toString() string {
	return string(w)
}

func scrubName(s string) string {
	return stringWrapper(strings.Title(s)).
		regexReplaceAll(`\s`, "").
		regexReplaceAll(`(?i)\(\s*for\s*User\s*\)`, "").
		regexReplaceAll(`(?i)\[\]`, "").
		regexReplaceAll(`(?i)\(array\)`, "").
		regexReplaceAll(`^\s+`, "").
		regexReplaceAll(`\s+$`, "").
		regexReplaceAll(`:`, "").
		regexReplaceAll(`[*]`, "").
		regexReplaceAll(`(?i)\(for User\)`, "").
		toString()
}

func parseResponseObjectToSlice(rawDef string) []map[string]string {

	// cleanup and 'standardize' the format - qgenda's docs
	// mix and match html and md tables for their definitions in docs
	// wipDef := stringWrapper(rawDef)
	wipDef := stringWrapper(rawDef).
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
		regexReplaceAll(`(?im)^([^|]*\|[^|]*)$`, "$1|").
		toString()

	rows := strings.Split(wipDef, "\n")
	def := []map[string]string{}
	// fldNum := 0
	for _, v := range rows {
		rowDef := map[string]string{}
		vs := strings.Split(v, "|")
		if len(vs) < 3 {
			continue
		}
		// name := stringWrapper(vs[0]).
		// 	regexReplaceAll(`^\s+`, "").
		// 	regexReplaceAll(`\s+$`, "").
		// 	regexReplaceAll(`:`, "").
		// 	regexReplaceAll(`[*]`, "").
		// 	regexReplaceAll(`(?i)\(for User\)`, "").
		// 	toString()
		name := scrubName(vs[0])
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
		regexReplaceAll(`\[\](.*)Tags$`, "[]Tag").
		regexReplaceAll(`Credentials$`, "Credential").
		regexReplaceAll(`Companies$`, "Company").
		regexReplaceAll(`Profiles$`, "Profile").
		regexReplaceAll(`Shifts$`, "Shift").
		regexReplaceAll(`Organizations$`, "Organization").
		regexReplaceAll(`EncounterRoles$`, "EncounterRole").
		regexReplaceAll(`ContactEmails$`, "ContactEmail").
		regexReplaceAll(`Assignments$`, "Assignment").
		regexReplaceAll(`Tasks$`, "Task").
		regexReplaceAll(`StaffMembers$`, "StaffMember").
		regexReplaceAll(`StaffMember$`, "StaffMemberDetail").
		regexReplaceAll(`Locations$`, "Location").
		regexReplaceAll(`room`, "Room").
		toString()

	typeMap := map[string]string{
		"int":       "int",
		"string":    "string",
		"boolean":   "bool",
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
	Request     interface{}
	Response    interface{}
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
