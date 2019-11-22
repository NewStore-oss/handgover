package main

import "testing"

func TestMessage(t *testing.T) {
	if message() != "Hello World" {
		t.FailNow()
	}
}
