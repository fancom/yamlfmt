// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package features

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/fancom/yamlfmt"
)

func MakeFeatureTrimTrailingWhitespace(linebreakStr string) yamlfmt.Feature {
	return yamlfmt.Feature{
		Name:         "Trim Trailing Whitespace",
		BeforeAction: trimTrailingWhitespaceFeature(linebreakStr),
	}
}

func trimTrailingWhitespaceFeature(linebreakStr string) yamlfmt.FeatureFunc {
	return func(content []byte) ([]byte, error) {
		buf := bytes.NewBuffer(content)
		s := bufio.NewScanner(buf)
		newLines := []string{}
		for s.Scan() {
			newLines = append(newLines, strings.TrimRight(s.Text(), " "))
		}
		return []byte(strings.Join(newLines, linebreakStr)), nil
	}
}
