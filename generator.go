package main

import (
	"fmt"
	"strings"
)

func generate(document SyntaxDocument) string {
	var sb strings.Builder
	// import zod
	sb.WriteString("import * as z from \"zod\";\n")

	for _, enum := range document.enums {
		sb.WriteString(generateEnum(enum))
	}
	return sb.String()
}

func generateMessage(message SyntaxMessage) string {
	return ""
}

func generateEnum(enum SyntaxEnum) string {
	var sb strings.Builder
	// apparently using actual enums is not recommended
	// https://zod.dev/api?id=enums
	sb.WriteString(fmt.Sprintf("export const %s = {\n", enum.name))
	for _, field := range enum.fields {
		sb.WriteString(fmt.Sprintf("  %s: %s,\n", field.name, field.id))
	}
	sb.WriteString("} as const;\n")
	sb.WriteString(fmt.Sprintf("export const %sEnum = z.enum(%s)\n", enum.name, enum.name))
	return sb.String()
}