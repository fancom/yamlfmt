package features

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
	"unicode"

	"github.com/fancom/yamlfmt"
)

var keyValuePattern = regexp.MustCompile(`^\s*([\w\-]+)\s*:\s*(.*)\s*$`)
var listValuePattern = regexp.MustCompile(`^\s*-\s*(.*)\s*$`)

func MakeFeatureStripStringQuotes(linebreakStr string) yamlfmt.Feature {
	return yamlfmt.Feature{
		Name:        "Strip quotes in strings",
		AfterAction: stripStringQuotesFeature(linebreakStr),
	}
}

func stripStringQuotesFeature(linebreakStr string) yamlfmt.FeatureFunc {
	return func(content []byte) ([]byte, error) {
		var buf bytes.Buffer
		reader := bytes.NewReader(content)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			txt := scanner.Text()
			match := keyValuePattern.FindStringSubmatch(txt)
			if len(match) >= 3 {
				value := match[2]
				buf.WriteString(stripQuotes(txt, value))
				buf.WriteString(linebreakStr)
			} else {
				match := listValuePattern.FindStringSubmatch(txt)
				if len(match) >= 2 {
					value := match[1]
					buf.WriteString(stripQuotes(txt, value))
					buf.WriteString(linebreakStr)
				} else {
					buf.WriteString(txt)
					buf.WriteString(linebreakStr)
				}
			}
		}
		return buf.Bytes(), scanner.Err()
	}
}

func stripQuotes(txt string, value string) string {
	if !strings.HasPrefix(value, `"`) && !strings.HasSuffix(value, `"`) {
		return txt
	}
	if containsSpecialSymbols(value) {
		return txt
	}
	if isNumeric(value) {
		return txt
	}
	return strings.Replace(txt, value, value[1:len(value)-1], 1)
}
func containsSpecialSymbols(s string) bool {
	specialSymbols := map[rune]struct{}{
		'{': {}, '}': {}, '[': {}, ']': {}, ',': {}, '&': {},
		':': {}, '*': {}, '#': {}, '?': {}, '|': {}, '\\': {},
		'-': {}, '<': {}, '>': {}, '=': {}, '!': {}, '%': {},
		'@': {},
	}
	for _, char := range s {
		if _, exists := specialSymbols[char]; exists {
			return true
		}
	}
	return false
}
func isNumeric(s string) bool {
	s = strings.Trim(s, `"'`)
	if s == "" {
		return false
	}

	for _, char := range s {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
