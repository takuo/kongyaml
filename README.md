
## Description

*kongyaml.CamelCase*

It reads YAML configuration from the provided io.Reader and returns a kong.Resolver.
The resolver can be used to resolve command-line flags from the YAML configuration.
It supports CamelCase keys only in the YAML file.

## Usage

config.yaml
```yaml
Flag: "flag"
Timeout: 30
Verbose: true
LogLevel: "info"
Debug: false
Nested:
    Key: "key"
    Value: "value"
```

main.go
```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/takuo/kongyaml"
)

func main() {
	var cli struct {
		ConfigFile kong.ConfigFlag `help:"Configuration file." type:"existingfile" default:"config.yaml" short:"c"`
		Flag       string          `help:"Flag description"`
		Timeout    int             `help:"Request timeout in seconds" default:"10"`
		Verbose    bool            `help:"Enable verbose output"`
		LogLevel   string          `help:"Set the logging level" enum:"debug,info,warn,error" default:"info"`
		Debug      bool            `help:"Enable debug mode"`
		Nested     struct {
			Key   string `help:"Nested key"`
			Value string `help:"Nested value"`
		} `embed:"" prefix:"nested."`
	}
	kong.Parse(&cli, kong.Configuration(kongyaml.CamelCase))
	b, _ := json.MarshalIndent(cli, "", "  ")
	fmt.Printf("%v\n", string(b))
}
```

result
```console
$ go run main.go -c config.yaml
{
  "ConfigFile": "/path/to/src/kongyaml/config.yaml",
  "Flag": "flag",
  "Timeout": 30,
  "Verbose": true,
  "LogLevel": "info",
  "Debug": false,
  "Nested": {
    "Key": "key",
    "Value": "value"
  }
}
```