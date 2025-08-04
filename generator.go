package main

import (
	"fmt"
	"strings"
)

func generate(document SyntaxDocument) string {
	var sb strings.Builder
	sb.WriteString("import * as z from \"zod\";\n")
	for _, enum := range document.enums {
		sb.WriteString(generateEnum(enum))
	}
	for _, message := range document.messages {
		sb.WriteString(generateMessage(message))
	}
	return sb.String()
}

func generateMessageField(field SyntaxMessageField) string {
	output := ""
	switch field.kind {
	// is a map
	// important that we check this first b/c TOKEN_TYPES includes map
	case "map":
		output = fmt.Sprintf("z.record(%s, %s)", kindToZod(*field.mapKey), kindToZod(*field.mapValue))
	default:
		output = kindToZod(field.kind)
	}

	// note: cannot be both optional and repeated
	// is optional
	if field.optional {
		return fmt.Sprintf("%s.optional()", output)
	} else if field.repeated {
		return fmt.Sprintf("z.array(%s)", output)
	}

	return output
}

func generateMessage(message SyntaxMessage) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("export const %s = z.object({\n", message.name))
	for _, field := range message.fields {
		sb.WriteString(fmt.Sprintf("  %s: %s,\n", field.name, generateMessageField(field)))
	}
	sb.WriteString("});\n")
	sb.WriteString(fmt.Sprintf("export type %sMessage = z.infer<typeof %s>;\n", message.name, message.name))
	// @todo generate child messages
	return sb.String()
}

func generateEnum(enum SyntaxEnum) string {
	var sb strings.Builder
	// apparently using actual enums is not recommended
	// https://zod.dev/api?id=enums
	sb.WriteString(fmt.Sprintf("export const %sLiteral = {\n", enum.name))
	for _, field := range enum.fields {
		sb.WriteString(fmt.Sprintf("  %s: %s,\n", field.name, field.id))
	}
	sb.WriteString("} as const;\n")
	sb.WriteString(fmt.Sprintf("export const %s = z.enum(%sLiteral);\n", enum.name, enum.name))
	return sb.String()
}

// utility to convert a type into its zod equivalent
// (I use kind since type is a reserved word)
func kindToZod(kind string) string {
	if TOKEN_TYPES.Contains(kind) {
		switch kind {
		case "bool":
			return "z.boolean()"
		case "string":
			return "z.string()"
		case "bytes":
			return "z.instanceof(Uint8Array)"
		}
		// since most types are numbers we use this as a default
		return "z.number()"
	}

	// just return identifier
	return kind 
}