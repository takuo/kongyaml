// Package kongyaml is kong.Resolver for YAML configuration.

package kongyaml

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

//
// It reads YAML configuration from the provided io.Reader and returns a kong.Resolver.
// The resolver can be used to resolve command-line flags from the YAML configuration.
// It supports CamelCase keys only in the YAML file.
//
// The function returns an error if there is an issue decoding the YAML.
//
// Example usage:
//
//	r := strings.NewReader(`
//	SomeFlag: someValue
//	Nested:
//	  AnotherFlag: anotherValue
//	`)
//	resolver, err := YAML(r)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Parameters:
//
//	r (io.Reader): The reader from which to read the YAML configuration.
//
// Returns:
//
//	(kong.Resolver, error): A kong.Resolver to resolve flags from the YAML configuration, and an error if decoding fails.

//
// This code is copied from "kong.JSON" with modification to handle CamelCase YAML
//

// CamelCase creates a kong.Resolver for a CamelCase YAML configuration.
func CamelCase(r io.Reader) (kong.Resolver, error) {
	values := map[string]any{}
	err := yaml.NewDecoder(r).Decode(values)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("YAML config decode error: %w", err)
	}
	var f kong.ResolverFunc = func(_ *kong.Context, parent *kong.Path, flag *kong.Flag) (any, error) {
		name := strings.ReplaceAll(flag.Name, "-", "_")
		camelName := camelCaseName(name)
		raw, ok := values[camelName]
		if ok {
			return raw, nil
		} else if parent != nil && parent.Command != nil {
			if v, ok := values[camelCaseName(parent.Command.Name)]; ok {
				raw = v
			}
		} else {
			raw = values
		}
		for _, part := range strings.Split(camelName, ".") {
			if values, ok := raw.(map[string]any); ok {
				raw, ok = values[part]
				if !ok {
					return nil, nil
				}
			} else {
				return nil, nil
			}
		}
		return raw, nil
	}
	return f, nil
}

var (
	re       = regexp.MustCompile("(^[A-Za-z])|[_.-]([A-Za-z])")
	replacer = strings.NewReplacer("_", "", "-", "")
)

func camelCaseName(str string) string {
	return re.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(replacer.Replace(s))
	})
}
