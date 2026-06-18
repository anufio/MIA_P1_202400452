package parser

// parser.go

import (
	"strings"
	"unicode"
)

type ParsedCommand struct {
	Name   string
	Params map[string]string
}

func ParseLine(line string) *ParsedCommand {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil
	}

	tokens := tokenize(line)
	if len(tokens) == 0 {
		return nil
	}

	cmd := &ParsedCommand{
		Name:   strings.ToLower(tokens[0]),
		Params: make(map[string]string),
	}

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		if strings.HasPrefix(token, "-") {
			eqIdx := strings.Index(token, "=")
			if eqIdx == -1 {
				key := strings.ToLower(token[1:])
				cmd.Params[key] = ""
			} else {
				key := strings.ToLower(token[1:eqIdx])
				val := token[eqIdx+1:]
				cmd.Params[key] = val
			}
		}
	}

	return cmd
}

func tokenize(line string) []string {
	var tokens []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(line); i++ {
		ch := rune(line[i])
		if ch == '"' {
			inQuotes = !inQuotes
		} else if unicode.IsSpace(ch) && !inQuotes {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}

func ParseScript(content string) []string {
	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		result = append(result, trimmed)
	}
	return result
}
