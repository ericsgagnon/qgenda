package meta

import "strings"

type Namespace string

func (n Namespace) Separator() string {
	return "."
}

func (n Namespace) Join(ns ...string) Namespace {

	out := strings.Join(ns, ".")
	// x := append([]string{string(n)}, ns...)
	return Namespace(out)
}

// func (n Namespace) Test() {

// 	i := 0

// }
