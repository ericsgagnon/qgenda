package main

import (
	"log"
	"os"

	"github.com/ericsgagnon/qgenda/app"
)

func main() {

	a, err := app.NewCLIApp(nil)
	if err != nil {
		log.Println(err)
	}
	if err := a.Run(os.Args); err != nil {
		log.Fatalln(err)
	}

	// 	src := `
	// ---
	// apiVersion: batch/v1beta1
	// kind: CronJob
	// name: ${PG_CONNECTION_STRING}
	// qname: "${PG_CONNECTION_STRING}"
	// sname: $(PG_CONNECTION_STRING)
	// cname: $PG_CONNECTION_STRING
	// # metadata:
	// #	name: resourcecleanup
	// # spec:
	// # 10:00 UTC == 1200 CET
	// # schedule: '0 10 * * 1-5'
	// `

	// 	re := regexp.MustCompile(`\$\{.+\}`)
	// 	envvars := map[string]string{}
	// 	for _, m := range re.FindAllString(src, -1) {
	// 		mre := regexp.MustCompile(`[${}]`)
	// 		mtrimmed := mre.ReplaceAllString(m, "")
	// 		// fmt.Printf("%s:\t%s\n", mtrimmed, os.Getenv(mtrimmed))
	// 		envvars[m] = os.Getenv(mtrimmed)
	// 	}
	// 	s := src
	// 	for k, v := range envvars {
	// 		s = strings.ReplaceAll(s, k, `"`+v+`"`)
	// 	}
	// fmt.Println(s)
	// fmt.Println(app.ExpandEnvVars(src))
	// dcfg := app.DefaultConfig(nil)
	// b, err := yaml.Marshal(dcfg)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(string(b))
	// fmt.Println(dcfg.DBClients["odbc"].String())
	// ex := qgenda.ExampleDBClientConfig()
	// ex.ExpandEnvVars = false
	// ex.ExpandFileContents = false
	// fmt.Println(ex)
	// fmt.Println(qgenda.ExpandFileContents("{file:/home/liveware/.wget-hsts}"))
	// out := qgenda.ExpandFileContents(qgenda.ExpandEnvVars("{file:${HOME}/.wget-hsts} and then there's: {file:${HOME}/.bashrc}"))
	// outa := strings.Split(out, "\n")
	// fmt.Println(outa[0])
	// fmt.Println(out)
	// a: "123123"
	// b: "${PG_CONNECTION_STRING}"
	// c:
	// - "firstEnstry"
	// - "secondEntry"

	// fmt.Println(src)
	// node := yaml.Node{}
	// if err := yaml.Unmarshal([]byte(src), &node); err != nil {
	// 	log.Println(err)
	// }
	// fmt.Printf("%#v\n", node)
	// fmt.Println("-------------------------------------------------------")
	// for _, nc := range node.Content {
	// 	fmt.Printf("%#v\n", *nc)
	// 	for _, ncnc := range nc.Content {
	// 		fmt.Printf("%#v\n", *ncnc)

	// 	}
	// }
	// fmt.Println("-------------------------------------------------------")
	// b, err := yaml.Marshal(node)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(string(b))

	// a.Command.Execute()
	// cmd, err := app.NewCommand()
	// if err != nil {
	// 	log.Println(err)
	// }
	// cmd.Execute()
}
