package main

import "testing"

func TestStub(t *testing.T) {
	if 1 != 1 { //nolint
		t.Fatal("Math is a lie. We should never have taught computers to think.")
	}
}
