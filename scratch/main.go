package main

import (
	"fmt"
	"log"
	"os"
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
}

func FileExists(filepath string) bool {

	fileinfo, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}
	// Return false if the fileinfo says the file path is a directory.
	return !fileinfo.IsDir()
}
