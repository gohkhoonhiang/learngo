package greetings

import (
	"regexp"
	"testing"
)

func TestHelloName(t *testing.T) {
	name := "Gladys"
	want := regexp.MustCompile(`\b` + name + `\b`)
	msg, err := Hello(1, "Gladys")
	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Hello(1, "Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

func TestHelloEmpty(t *testing.T) {
	msg, err := Hello(0, "")
	if msg != "" || err == nil {
		t.Fatalf(`Hello(0, "") = %q, %v want "", error`, msg, err)
	}
}
