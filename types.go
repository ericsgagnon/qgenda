package main

import "time"

type Nullable struct {
	Valid bool
}

type Time struct {
	Value time.Time
	Nullable
}

type String struct {
	Value string
	Nullable
}

type Int struct {
	Value int
	Nullable
}

type Byte struct {
	Value byte
	Nullable
}

type Bool struct {
	Value bool
	Nullable
}


