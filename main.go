package main

import (
	"errors"
	"os"

	"github.com/lusingander/topi/internal/openapi"
	"github.com/lusingander/topi/internal/ui"
)

func pathFromArgs(args []string) (string, error) {
	if len(args) <= 1 {
		return "", errors.New("must set OpenAPI spec json/yaml filepath as argument")
	}
	return args[1], nil
}

func run(args []string) error {
	path, err := pathFromArgs(args)
	if err != nil {
		return err
	}
	doc, err := openapi.Load(path)
	if err != nil {
		return err
	}
	return ui.Start(doc)
}

func main() {
	if err := run(os.Args); err != nil {
		panic(err)
	}
}
