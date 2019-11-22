package main

import "fmt"

func main() {
	fmt.Println(message())
}

const m = "Hello World"

func message() string {
	return m
}
