package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	fi, err := os.Stat("./.scratch")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v\n", fi)
	fmt.Println(FileExists("./.scratch"))
	fmt.Println(FileExists("./heehaw"))

	fs, err := os.Stat(".scratch")
	fmt.Println(!os.IsNotExist(err))
	fmt.Println(fs.ModTime())
	fmt.Println(fs.Mode())

	tt := Time{
		Created: ptr(time.Now().UTC()),
	}
	fmt.Println(tt)
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

func toTimePointer(a *any) *time.Time {
	return time.Time(a)
}

type Time struct {
	Created *time.Time
}
