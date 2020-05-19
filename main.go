package hcl2json

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func hcl2Json() {
	logger := log.New(os.Stderr, "", 0)
	flag.Parse()
	var err error

	var filename = flag.Arg(0)
	var bytes []byte
	if filename == "" || filename == "-" {
		bytes, err = ioutil.ReadAll(os.Stdin)
	} else {
		bytes, err = ioutil.ReadFile(filename)
	}

	if err != nil {
		logger.Fatalf("Failed to read file: %s\n", err)
	}

	var content interface{}
	content, err = GetHclJSON(bytes, filename)

	if err != nil {
		logger.Fatalf("Failed to convert file: %v", err)
	}

	jb, err := json.MarshalIndent(content, "", "    ")

	if err != nil {
		logger.Fatalf("Failed to generate JSON: %v", err)
	}

	os.Stdout.Write(jb)
}

// GetHclJSON is the main exported function
func GetHclJSON(bytes []byte, filename string) (interface{}, error) {
	file, diags := hclsyntax.ParseConfig(bytes, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}
	return convertFile(file)
}
