package main

import (
	"testing"
)

func Test_Decode_20_u(t *testing.T) {
	Ans, _ := Decode("u")
    if Ans != 20 {
		t.Error("wrong result")
	}
}

func Test_Encode_u_20(t *testing.T) {
	Ans := Encode(20)
    if Ans != "u" {
		t.Error("wrong result")
	}
}