// contains the logic for filtering strings
// - the main entry point is not clear between Filter and FilterByList currently.
package filter

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

func fuzzy(filter string) (*regexp.Regexp, error) {
	newFilter := "(?i)" // ?i makes is case insensitive

	regexpString := newFilter + strings.Join(strings.Split(filter, ""), ".*?")
	return regexp.Compile(regexpString)
}

type match struct {
	filePath string
	rank     int
}

type matches []match

// TODO: correct make errors
func (m matches) Len() int           { return len(m) }
func (m matches) Less(i, j int) bool { return m[i].rank < m[j].rank }
func (m matches) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func Filter(col []string, filter string) []string {
	matches := matches{}

	freg, err := fuzzy(strings.Replace(filter, ".", "\\.", -1))

	if err != nil {
		fmt.Errorf("bad regexp compilation %v", err)
	}

	for _, e := range col {
		fuzzyMatches := freg.FindAllString(e, -1)

		if len(fuzzyMatches) != 0 {
			fmatch := fuzzyMatches[len(fuzzyMatches)-1]

			if len(fmatch) >= len(filter) {
				matches = append(matches, match{e, len(fmatch)})
			}
		}
	}

	sort.Sort(matches)

	filePaths := []string{}

	for _, match := range matches {
		filePaths = append(filePaths, match.filePath)
	}

	return filePaths
}

func contains(filePath string, list []string) bool {
	for _, ignoreMatch := range list {
		if strings.HasPrefix(filePath, strings.TrimPrefix(ignoreMatch, "/")) {
			return true
		}
	}
	return false
}

func FilterByList(source []string, list []string) []string {
	filteredList := []string{}

	for _, filePath := range source {
		if !contains(filePath, list) {
			filteredList = append(filteredList, filePath)
		}
	}

	return filteredList
}
