package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	fi, err := os.Stat("./.scratch")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("ModTime:\t%s\n", fi.ModTime())

	fmt.Printf("%#v\n", fi)

	fmt.Println(FileExists("./.scratch"))
	fmt.Println(FileExists("./heehaw"))

	fs, err := os.Stat(".scratch")
	fmt.Println(!os.IsNotExist(err))
	fmt.Println(fs.ModTime())
	fmt.Println(fs.Mode())
	t := time.Now().UTC()
	tt := Time{
		Created: &t,
	}
	fmt.Println(tt)

	// var b []byte
	b, err := os.ReadFile("/home/liveware/.cache/qgenda/authtoken.json")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
	tkn := &AuthToken{
		Token: &http.Header{},
	}
	if err := json.Unmarshal(b, tkn); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(tkn.Token)
	fmt.Println(tkn.Expires)
	// fi, err := os.Stat("types.tmpl")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
}

func FileExists(filepath string) bool {

	fileinfo, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}
	// Return false if the fileinfo says the file path is a directory.
	return !fileinfo.IsDir()
}

func ptr(a any) *any {
	return &a
}

type Time struct {
	Created *time.Time
}

type AuthToken struct {
	Token         *http.Header  `json:"token"`
	Timestamp     time.Time     `json:"timeStamp"` // defaults to response's date header
	ValidDuration time.Duration `json:"validDuration"`
	Expires       time.Time     `json:"expires"`
	Cache         *any          `json:"-"`
}
