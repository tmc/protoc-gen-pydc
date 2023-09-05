package main

import (
	"fmt"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

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
var tmpl = template.Must(template.New("pydc").Funcs(template.FuncMap{
	"messageName": func(msg *protogen.Message) string {
		return msg.GoIdent.GoName
	},
	"comments": func(comments protogen.Comments) string {
		// strip "// " from the beginning of each line
		c := strings.TrimPrefix(string(comments), "// ")
		c = strings.TrimSpace(c)

		if c != "" {
			return fmt.Sprintf("# %s", c)
		}
		return ""
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
}).Parse(`
{{- $packageName := .Desc.Package -}}
# -*- coding: utf-8 -*-
"""
Python Dataclasses for {{.Desc.Package}}
"""
from dataclasses import dataclass
from collections import OrderedDict
from enum import Enum
from typing import Dict

{{ range .Enums}}
class {{.GoIdent.GoName}}(Enum):
    {{- range .Values}}
    {{.GoIdent.GoName}} = {{.Desc.Number}} {{comments .Comments.Trailing}}
    {{- end}}
{{- end}}

{{- range .Messages}}
{{ $message := . }}
@dataclass
class {{.Desc.Name}}:
    {{- range .Fields}}
    {{.Desc.TextName}}: {{fieldType .}} {{comments .Comments.Trailing}}
    {{- end}}

{{- end}}

{{- range .Services}}
{{- $service := . }}
{{- end}}
`))
