package parser

import (
	"regexp"
)

func GetCCppHeaders(cCode string) []string {
	re := regexp.MustCompile(`#include\s+<[^>]*>`)
	return re.FindAllString(cCode, -1)
}

func GetCCppFunction(cCode, funcName string) string {
	funcPattern := regexp.MustCompile(`\b\w+\s+\**\b` + regexp.QuoteMeta(funcName) + `\s*\([^)]*\)\s*\{`)
	match := funcPattern.FindStringIndex(cCode)
	if match == nil {
		return ""
	}

	start := match[0]
	end := start

	braceCount := 0
	inString := false
	escaped := false

	for i := start; i < len(cCode); i++ {
		if cCode[i] == '"' && !escaped {
			inString = !inString
		}
		if cCode[i] == '\\' && !escaped {
			escaped = true
			continue
		}
		if !inString {
			if cCode[i] == '{' {
				braceCount++
			} else if cCode[i] == '}' {
				braceCount--
				if braceCount == 0 {
					end = i + 1
					break
				}
			}
		}
		escaped = false
	}

	return cCode[start:end]
}
