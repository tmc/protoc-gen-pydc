package main

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"google.golang.org/protobuf/compiler/protogen"
)

// Use embed to include the templates directory in the binary:
//
//go:embed templates
var templates embed.FS

var protoToPythonTypeMap = map[string]string{
	"string": "str",
}

func toPyType(t string) string {
	if pyType, ok := protoToPythonTypeMap[t]; ok {
		return pyType
	}
	return t
}

// tmpl is the template used for generating python dataclasses from proto files.
func getTemplate() *template.Template {
	tmpl := template.New("main.tmpl").Funcs(sprig.TxtFuncMap())
	fm := template.FuncMap{
		"include": func(name string, data interface{}) (string, error) {
			buf := bytes.NewBuffer(nil)
			if err := tmpl.ExecuteTemplate(buf, name, data); err != nil {
				return "", err
			}
			return buf.String(), nil
		},
		"messageName": func(msg *protogen.Message) string {
			return msg.GoIdent.GoName
		},
		"comments": func(comments protogen.Comments) string {
			// strip "// " from the beginning of each line
			c := strings.TrimPrefix(string(comments), "// ")
			c = strings.TrimSpace(c)
			if c == "" {
				return ""

			}
			return fmt.Sprintf(" # %s", c)
		},
		"fieldType": func(field *protogen.Field) string {
			// Check if the field is a map
			if field.Desc.IsMap() {
				keyType := field.Desc.MapKey().Kind().String()
				valueType := field.Desc.MapValue().Kind().String()
				return "Dict[" + toPyType(keyType) + ", " + toPyType(valueType) + "]"
			}

			// Handle other field types
			switch {
			case field.Enum != nil:
				return field.Enum.GoIdent.GoName
			case field.Message != nil:
				return field.Message.GoIdent.GoName
			default:
				ts := field.Desc.Kind().String()
				return toPyType(ts)
			}
		},
	}
	return template.Must(tmpl.Funcs(fm).ParseFS(templates, "**/*.tmpl"))
}
