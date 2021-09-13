package controllers2_test

import (
	"os"
	"testing"
)

var tt = new(TestData)

type TestData struct {
	Name   string
	Method string
}

func TestMain(m *testing.M) {
	tt.Name = "hello"
	tt.Method = "bomb"
	os.Exit(m.Run())
}
