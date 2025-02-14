package kongyaml

import (
	"net/netip"
	"strings"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYAML(t *testing.T) {
	type CLI struct {
		FlagName string
		Names    []string
		Command  struct {
			NestedFlag string
		} `cmd:""`
		Embedded struct {
			One string
			Two bool
		} `embed:"" prefix:"embed."`
		Dict          map[string]string
		NestedDict    map[string]map[string]bool
		NonStringDict map[netip.Addr][]string
		TypedDict     map[string]struct {
			Foo string
			Bar float64
		}
		TypedSlice []struct{ Foo string }
	}
	var cli CLI
	r := strings.NewReader(`
FlagName: "hello world"
Embed:
    One: "str"
    Two: true
Names:
    - "one"
    - "two"
    - "three"
Command:
    NestedFlag: "nested flag"
    Number: 1.0
    Int: 12342345234534
Dict:
    Foo: bar
NestedDict:
    Foo:
        Bar: true # also settable as --nested-dict=foo=bar=true
NonStringDict:
    "1.2.3.4": ["foo", "bar"]
TypedDict:
    Foo:
        Bar: 1.337
        Foo: bar
TypedSlice:
    - Foo: bar
    - Foo: baz
`)
	resolver, err := CamelCase(r)
	require.NoError(t, err)
	parser, err := kong.New(&cli, kong.Resolvers(resolver))
	require.NoError(t, err)
	_, err = parser.Parse([]string{"command"})
	require.NoError(t, err)
	expected := CLI{
		FlagName: "hello world",
		Names:    []string{"one", "two", "three"},
		Command: struct {
			NestedFlag string
		}{NestedFlag: "nested flag"},
		Embedded: struct {
			One string
			Two bool
		}{
			One: "str",
			Two: true,
		},
		Dict: map[string]string{
			"Foo": "bar",
		},
		NestedDict: map[string]map[string]bool{
			"Foo": {
				"Bar": true,
			},
		},
		NonStringDict: map[netip.Addr][]string{
			netip.MustParseAddr("1.2.3.4"): {"foo", "bar"},
		},
		TypedDict: map[string]struct {
			Foo string
			Bar float64
		}{
			"Foo": {
				Foo: "bar",
				Bar: 1.337,
			},
		},
		TypedSlice: []struct{ Foo string }{
			{Foo: "bar"},
			{Foo: "baz"},
		},
	}
	require.Equal(t, expected, cli)
}

func TestEmptyFile(t *testing.T) {
	type CLI struct {
		FlagName string
	}
	var cli CLI
	r := strings.NewReader("")
	resolver, err := CamelCase(r)
	require.NoError(t, err)
	parser, err := kong.New(&cli, kong.Resolvers(resolver))
	require.NoError(t, err)
	_, err = parser.Parse([]string{})
	require.NoError(t, err)
	expected := CLI{
		FlagName: "",
	}
	require.Equal(t, expected, cli)
}

func Test_camelCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "snake_case",
			args: args{str: "snake_case"},
			want: "SnakeCase",
		},
		{
			name: "dash-case",
			args: args{str: "dash-case"},
			want: "DashCase",
		},
		{
			// dot is used for embed prefix in kong, don't replace it
			name: "with.dot",
			args: args{str: "with.dot"},
			want: "With.Dot",
		},
		{
			// dot is used for embed prefix in kong, don't replace it
			name: "dot.dash-case",
			args: args{str: "dot.dash-case"},
			want: "Dot.DashCase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, camelCaseName(tt.args.str))
		})
	}
}
