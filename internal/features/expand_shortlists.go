package features

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
	"unicode"

	"github.com/fancom/yamlfmt"
)

func getPading(s string, ident int, indentlessArrays bool) string {
	var builder strings.Builder
	for _, r := range s {
		if unicode.IsSpace(r) {
			builder.WriteByte(' ')
		} else {
			break
		}
	}
	if !indentlessArrays {
		for i := 0; i < ident; i++ {
			builder.WriteByte(' ')
		}
	}
	return builder.String()
}

func MakeFeatureExpandShortlists(linebreakStr string, ident int, indentlessArrays bool) yamlfmt.Feature {
	return yamlfmt.Feature{
		Name:        "Expand Shortlists",
		AfterAction: expandShortlistsFeature(linebreakStr, ident, indentlessArrays),
	}
}

func expandShortlistsFeature(linebreakStr string, ident int, indentlessArrays bool) yamlfmt.FeatureFunc {
	return func(content []byte) ([]byte, error) {
		var buf bytes.Buffer
		reader := bytes.NewReader(content)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			txt := scanner.Text()
			trimTxt := strings.TrimSpace(txt)

			if regexp.MustCompile(`: \[(.*)\]$`).MatchString(trimTxt) {
				matches := regexp.MustCompile(`\[(.*)\]$`).FindAllStringSubmatch(trimTxt, -1)
				if len(matches) > 0 && len(matches[0]) > 1 {
					extracted := matches[0][1]
					padding := getPading(txt, ident, indentlessArrays)
					buf.WriteString(strings.Split(txt, ":")[0] + ":")
					buf.WriteString(linebreakStr)
					items := strings.Split(extracted, ",")
					for i := range items {
						items[i] = strings.TrimSpace(items[i])
						//TODO(AZ): move quotes removing into separate feature
						items[i] = strings.Trim(items[i], "\"")
						buf.WriteString(padding + "- " + items[i])
						buf.WriteString(linebreakStr)
					}
				}
			} else {
				buf.WriteString(txt)
				buf.WriteString(linebreakStr)
			}
		}
		return buf.Bytes(), scanner.Err()
	}
}
