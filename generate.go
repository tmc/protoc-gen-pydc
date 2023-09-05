package main

import (
	"bytes"
	"io/ioutil"

	"google.golang.org/protobuf/compiler/protogen"
)

// generateFile processes a single proto file.
func generateFile(gen *protogen.Plugin, f *protogen.File) {
	// capture to bytes buffer
	buf := new(bytes.Buffer)

	// execute template with protobuf File schema
	if err := tmpl.Execute(buf, f); err != nil {
		gen.Error(err)
		return
	}
	// write out python file:
	if err := ioutil.WriteFile(f.GeneratedFilenamePrefix+".py", buf.Bytes(), 0600); err != nil {
		gen.Error(err)
		return
	}
}
