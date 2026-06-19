package main

import "fmt"

type Reporter interface {
	Report() string
}

func printReport(r Reporter) {
	fmt.Println(r.Report())
}
