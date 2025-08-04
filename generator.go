package main

import (
	"fmt"
	"strings"
)

// @follow-up idea: instead of flattened structure we have going on separated by underscores (ie: Message__NestedMessage)
// what if we build a nested object like this:
// const nested = { Message: { NestedMessage: {} }  }

const DELIMITER = "__"

func generate(document SyntaxDocument) string {
	var sb strings.Builder
	sb.WriteString("import * as z from \"zod\";\n")
	for _, enum := range document.enums {
		sb.WriteString(generateEnum(enum))
	}
	for _, message := range document.messages {
		sb.WriteString(generateMessage(message, []string{}))
	}
	return sb.String()
}

func generateMessage(message SyntaxMessage, messagePath []string) string {
	messagePath = append(messagePath, message.name)
	messageName := strings.Join(messagePath, DELIMITER)

	var sb strings.Builder

	// generate child messages first b/c we declare w/ const (so child needs to be accessible first)
	for _, childMessage := range message.messages {
		sb.WriteString(generateMessage(childMessage, messagePath))
	}

	sb.WriteString(fmt.Sprintf("export const %s = z.object({\n", messageName))
	for _, field := range message.fields {
		output := ""
		switch field.kind {
		// is a map - important that we check this first b/c TOKEN_TYPES includes map
		case "map":
			output = fmt.Sprintf("z.record(%s, %s)", kindToZod(*field.mapKey, &message, &messagePath), kindToZod(*field.mapValue, &message, &messagePath))
		default:
			output = kindToZod(field.kind, &message, &messagePath)
		}

		// note: cannot be both optional and repeated
		if field.optional {
			output = fmt.Sprintf("%s.optional()", output)
		} else if field.repeated {
			output = fmt.Sprintf("z.array(%s)", output)
		}

		sb.WriteString(fmt.Sprintf("  %s: %s,\n", field.name, output))
	}
	sb.WriteString("});\n")
	sb.WriteString(fmt.Sprintf("export type %sMessage = z.infer<typeof %s>;\n", messageName, messageName))

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

// utility to convert a type into its zod equivalent; I use kind since type is a reserved word
func kindToZod(kind string, message *SyntaxMessage, messagePath *[]string) string {
	if TOKEN_TYPES.Contains(kind) {
		switch kind {
		case "bool":
			return "z.boolean()"
		case "string":
			return "z.string()"
		case "bytes":
			return "z.instanceof(Uint8Array)"
		}
		// since most scalar types are numbers we use this as a default
		return "z.number()"
	}
	
	// @todo actually read imports
	// this is a pretty shit mechanism but oh well
	if kind == "google.protobuf.Any" {
		return "z.any()"
	}

	// @audit we should probably verify types against the document & throw an error otherwise
	// case 1: check that external type exists
	// case 2: check that internal type exists (we only check first part to verify scope & assume the rest is in there)

	/*
		resolve identifier in nested messages
		for example:
		message Message {
			message Inner {} <- generated as Message__Inner
			Inner field = 1;
			Outer field = 2;
		}
		so if we're given "Inner" and a message of "Message", we should resolve to Message__Inner
		if we're given "Outer" and a message of "Message", we just resolve to Outer and assume it exists outside
	*/
	output := strings.Split(kind, ".")
	for _, childMessage := range message.messages {
		if childMessage.name == output[0] {
			output = append(*messagePath, output...)
			break
		}
	}

	return strings.Join(output, DELIMITER)
}
