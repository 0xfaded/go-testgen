package testgen

import (
	"io"
	"strings"
)

type Generator interface {
	Package() string
	Prefix() string
	Imports() map[string] string
	Dimensions() []Dimension
	Comment(io.Writer, ...Element) error
	Body(io.Writer, ...Element) error
}

func Generate(gen Generator, w io.Writer) error {
	if _, err := w.Write([]byte("package ")); err != nil {
		return err
	}
	if _, err := w.Write([]byte(gen.Package())); err != nil {
		return err
	}

	imports := gen.Imports()
	if _, err := w.Write([]byte("\n\nimport (\n\t\"testing\"\n")); err != nil {
		return err
	}
	for im, name := range(imports) {
		s := "\t\""+im+"\""
		if name != "" {
			s += " "+name
		}
		s += "\n"
		if _, err := w.Write([]byte(s)); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte(")\n")); err != nil {
		return err
	}

	prefix := gen.Prefix()
	dims := gen.Dimensions()

	end := uint64(1)
	for _, dim := range(dims) {
		end *= uint64(len(dim))
	}

	elts := make([]Element, len(dims))
	names := make([]string, len(dims))
	for test := uint64(0); test < end; test += 1 {
		tmp := test
		for i := len(dims)-1; i >= 0; i -= 1 {
			j := tmp % uint64(len(dims[i]))
			elts[i] = dims[i][j]
			names[i] = dims[i][j].Name
			tmp /= uint64(len(dims[i]))
		}

		testName := strings.Join(names, "")

		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
		if err := gen.Comment(w, elts...); err != nil {
			return err
		}
		if _, err := w.Write([]byte("func Test" + prefix + testName)); err != nil {
			return err
		}
		if _, err := w.Write([]byte("(t *Testing) {\n")); err != nil {
			return err
		}
		if err := gen.Body(w, elts...); err != nil {
			return err
		}
		if _, err := w.Write([]byte("}\n")); err != nil {
			return err
		}
	}
	return nil
}
