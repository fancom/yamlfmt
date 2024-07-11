package features

import (
	"bufio"
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/fancom/yamlfmt"
)

var keyValuePattern = regexp.MustCompile(`^\s*(- )?([\w\-]+)\s*:\s*(.*)\s*$`)
var listValuePattern = regexp.MustCompile(`^\s*-\s*(.*)\s*$`)

var boolValues = []string{"true", "false", "on", "off", "yes", "no"}

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
			if len(match) >= 4 {
				value := match[3]
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
	if isBoolean(value) {
		return txt
	}
	if len(value) < 3 {
		return txt
	}
	return strings.Replace(txt, value, value[1:len(value)-1], 1)
}

func isBoolean(s string) bool {
	s = strings.Trim(s, `"'`)
	if s == "" {
		return false
	}
	if !contains(s, boolValues) {
		return false
	}
	return true
}
func contains(value string, array []string) bool {
	for _, v := range array {
		if strings.EqualFold(value, v) {
			return true
		}
	}
	return false
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
	_, err := strconv.ParseFloat(s, 64)

	return err == nil
}
