package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

//https://restapi.qgenda.com/?version=latest

var err error

func main() {

	ctx := context.Background()

	// use environment variables to provide credentials
	q, err := NewQgendaClient(
		QgendaClientConfig{
			BaseURL:       "https://api.qgenda.com/v2",
			ClientTimeout: time.Second * 10,
			// grab credentials from environment variables
			Email:      os.Getenv("QGENDA_EMAIL"),
			CompanyKey: os.Getenv("QGENDA_COMPANY_KEY"),
			Password:   os.Getenv("QGENDA_PASSWORD"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx = q.Auth(ctx)

	fmt.Println("---------------------------------------------------------")
	// for _, v := range q.Authorization.Cookies(q.BaseURL) {
	// 	fmt.Printf("%v: %v\n%v\n", v.Name, v.Expires, v.Value)
	// }

	fmt.Println("---------------------------------------------------------")
	fmt.Println(ctx)
	fmt.Println("---------------------------------------------------------")
	// res, err := q.Client.Get("https://api.qgenda.com/v2/company")
	// if err != nil {
	// 	log.Fatal(err)

	// }

	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)

	// }

	// fmt.Println(string(body))
}
