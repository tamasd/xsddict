package xsddict

import (
	"io"
	"strings"

	"github.com/antchfx/xmlquery"
)

const (
	SpaceLength   = 8
	TabLength     = 4
	NewLineLength = 4
)

func GenerateDict(dict io.Writer, input io.Reader) error {
	n, err := xmlquery.Parse(input)
	if err != nil {
		return err
	}

	if err = generateXmlParts(dict, n, "//xsd:element", "name", func(name string) string { return "/" + name + ">" }); err != nil {
		return err
	}

	if err = generateXmlParts(dict, n, "//xsd:attribute", "name", func(name string) string { return name + "=" }); err != nil {
		return err
	}

	return nil
}

func generateXmlParts(dict io.Writer, root *xmlquery.Node, query, attrName string, processName func(name string) string) error {
	m := make(map[string]struct{})

	results, err := xmlquery.QueryAll(root, query)
	if err != nil {
		return err
	}

	for _, result := range results {
		name := result.SelectAttr(attrName)
		if name == "" {
			continue
		}

		if _, ok := m[name]; ok {
			continue
		}

		m[name] = struct{}{}

		if processName != nil {
			name = processName(name)
		}

		if _, err := io.WriteString(dict, name); err != nil {
			return err
		}
	}

	return nil
}

func WhiteSpaces(dict io.Writer) error {
	if _, err := io.WriteString(dict, strings.Repeat(" ", SpaceLength)); err != nil {
		return err
	}

	if _, err := io.WriteString(dict, strings.Repeat("\t", TabLength)); err != nil {
		return err
	}

	if _, err := io.WriteString(dict, strings.Repeat("\n", NewLineLength)); err != nil {
		return err
	}

	return nil
}
