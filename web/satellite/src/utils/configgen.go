// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package main

//go:generate go run $GOFILE -o=config.gen.ts

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"storj.io/storj/private/apigen"
	"storj.io/storj/satellite/console/consoleweb"
)

func main() {
	flag.CommandLine = flag.NewFlagSet("configgen", flag.ExitOnError)
	outPath := flag.String("o", "", "path to the output file")
	flag.Parse()

	if *outPath == "" {
		fmt.Fprintln(os.Stderr, "missing required argument -o")
		os.Exit(1)
	}

	var result apigen.StringBuilder
	pf := result.Writelnf

	pf("// AUTOGENERATED BY configgen.go")
	pf("// DO NOT EDIT.")
	pf("")

	types := apigen.NewTypes()
	types.Register(reflect.TypeOf(consoleweb.FrontendConfig{}))
	result.WriteString(types.GenerateTypescriptDefinitions())

	content := strings.ReplaceAll(result.String(), "\t", "    ")
	err := os.WriteFile("config.gen.ts", []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}
