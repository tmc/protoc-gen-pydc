package main

import (
	"bytes"

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

	gf := gen.NewGeneratedFile(f.GeneratedFilenamePrefix+".py", f.GoImportPath)
	_, err := gf.Write(buf.Bytes())
	if err != nil {
		gen.Error(err)
		return
	}
	// // write out python file:
	// if err := ioutil.WriteFile(f.GeneratedFilenamePrefix+".py", buf.Bytes(), 0600); err != nil {
	// 	gen.Error(err)
	// 	return
	// }
}
