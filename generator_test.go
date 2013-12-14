package testgen

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

type Test struct {}

func (*Test) Package() string {
	return "test"
}

func (*Test) Prefix() string {
	return "Test"
}

func (*Test) Imports() map[string]string {
	return map[string]string{
		"importa": "",
		"importb": "b",
	}
}

func (*Test) Dimensions() []Dimension {
	return []Dimension{
		{ {"A", 1}, {"B", 2}, {"C", 3} },
		{ {"D", 4}, {"E", 5} },
		{ {"F", 6 } },
	}
}

func (*Test) Comment(w io.Writer, elts ...Element) error {
	if _, err := w.Write([]byte("//")); err != nil {
		return err
	}
	for _, elt := range elts {
		if _, err := w.Write([]byte(" " + elt.Name)); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}

func (*Test) Body(w io.Writer, elts ...Element) error {
	if _, err := w.Write([]byte("\t_ = 0")); err != nil {
		return err
	}
	for _, elt := range elts {
		if _, err := w.Write([]byte(fmt.Sprintf(" + %v", elt.Value))); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}

func TestGenerate(t *testing.T) {
	var buf bytes.Buffer
	Generate(&Test{}, &buf)

	src :=
`package test

import (
	"testing"
	"importa"
	"importb" b
)

// A D F
func TestTestADF(t *Testing) {
	_ = 0 + 1 + 4 + 6
}

// A E F
func TestTestAEF(t *Testing) {
	_ = 0 + 1 + 5 + 6
}

// B D F
func TestTestBDF(t *Testing) {
	_ = 0 + 2 + 4 + 6
}

// B E F
func TestTestBEF(t *Testing) {
	_ = 0 + 2 + 5 + 6
}

// C D F
func TestTestCDF(t *Testing) {
	_ = 0 + 3 + 4 + 6
}

// C E F
func TestTestCEF(t *Testing) {
	_ = 0 + 3 + 5 + 6
}
`
	actual := strings.Split(buf.String(), "\n")
	expected := strings.Split(src, "\n")
	i := 0
	diff := ""
	s40 := "                                         "
	for ; i < len(actual) && i < len(expected); i += 1 {
		// Rough hack, won't work for unicode
		actual[i] = strings.Replace(actual[i], "\t", "........", -1)
		expected[i] = strings.Replace(expected[i], "\t", "........", -1)
		padding := 40 - len(expected[i])
		paddings := s40[:padding]
		if actual[i] == expected[i] {
			diff += fmt.Sprintf("%s%s == %s\n", expected[i], paddings, actual[i])
		} else {
			diff += fmt.Sprintf("%s%s != %s\n", expected[i], paddings, actual[i])
		}
	}
	for ;i < len(expected); i += 1 {
		expected[i] = strings.Replace(expected[i], "\t", "........", -1)
		diff += fmt.Sprintf("%s\n", expected[i])
	}
	for ;i < len(actual); i += 1 {
		expected[i] = strings.Replace(expected[i], "\t", "........", -1)
		diff += fmt.Sprintf("%s    %s\n", s40, actual[i])
	}
	if buf.String() != src {
		t.Fatalf("Expected, Unexpected:\n" + diff)
	}
}
