package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {

	f, err := os.Open("../samples/qgenda_restapi.postman_collection.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fc, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal(fc, &data); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)

	// for k, _ := range data {
	// 	fmt.Println(k)
	// }
	di := data["item"]
	dib, err := json.Marshal(di)
	if err != nil {
		log.Fatal(err)
	}
	dim := []interface{}{}
	if err := json.Unmarshal(dib, &dim); err != nil {
		log.Fatal(err)
	}
	// for k, _ := range dim {
	// 	fmt.Println(k)
	// 	// fmt.Println(v)
	// }
	pc := PostmanCollection{}
	if err := json.Unmarshal(fc, &pc); err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Postman Collection: %+v", pc.Item)

	// fmt.Printf("%T\n", pc.Item[2])
	for k, v := range pc.Item[2].Item {
		desc := v.Description
		reg := regexp.MustCompile("-|\\*\\*|`")
		desc = reg.ReplaceAllString(desc, "")
		// desc = strings.ReplaceAll(desc, "-", "")
		// desc = strings.ReplaceAll(desc, "**`", "")
		// desc = strings.ReplaceAll(desc, "`**", "")
		desc = strings.ReplaceAll(desc, "  ", " ")
		fmt.Printf("%+v\n#############################################\n", desc)
		// desc = strings.ReplaceAll(desc, "-", "")
		// fmt.Println(k)
		if k == 30 {
			for ri, rv := range strings.Split(desc, "\n") {
				if ri != 0 {
					for ei, ev := range strings.Split(rv, "|") {
						if ev != "" {
							fmt.Printf("%+v:\t%+v\n", ei, ev)

						}

					}

				}
			}

			// fmt.Printf("%v:\n%v\n", v.Name, v.Description)
		}
	}

	// fmt.Println(data["item"])
	// dit := data["item"].(map[string]interface{})
	// ditr := reflect.TypeOf(dit)
	// fmt.Println(ditr)

	// fmt.Println(reflect.TypeOf(di).Kind())
	// dis := reflect.ValueOf(data["item"])

	// for i := 0; i < dis.Len(); i++ {
	// 	fmt.Println(reflect.TypeOf(dis.Index(i).Kind()))
	// }
	// for i, v := range data["item"] {
	// 	fmt.Println(i)
	// }
	// fmt.Println(dit)

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
