package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	PatternRegex = regexp.MustCompile(`\$\d`)
)

type mvCmd struct {
	from string
	to   string
}

const usage string = `Usage: mvp [pattern] [new-pattern]
Examples:
$ ls
log_1.json  log_2.json warn_3.json error_3.json README.md

$ mvp \$1_\$2.json \$2_\$1.json 
$ ls
1_log.json  2_log.json 3_warn.json 3_error.json README.md`

func Usage() {
	fmt.Println(usage)

}

func main() {
	if len(os.Args) == 1 {
		Usage()
		return
	}
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 1 && args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		Usage()
		return 0
	}
	if len(args) != 2 {
		fmt.Printf("Wrong usage\n")
		Usage()
		return 1
	}

	pattern, newPattern := args[0], args[1]

	// Matches
	placeholders := getMatches(pattern)
	if len(placeholders) == 0 {
		fmt.Printf("Bad pattern\n")
		Usage()
		return 1
	}

	// Matches in the new pattern should be the same as in the original pattern
	newPatternPlaceholders := getMatches(newPattern)
	if len(newPatternPlaceholders) == 0 {
		fmt.Printf("new-pattern %q doesn't contain exactly the same placeholders used in the pattern %q\n",
			newPattern, pattern)
		return 1
	}

	if !areMapsEqual(toMap(placeholders), toMap(newPatternPlaceholders)) {
		fmt.Printf("Placeholders in the two patterns don't match\n")
		return 1
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	placeHolderMap := make(map[string]struct{})
	for _, match := range placeholders {
		if _, ok := placeHolderMap[match]; !ok {
			placeHolderMap[match] = struct{}{}
		}
	}

	mvCmds := make([]mvCmd, 0)

	for _, entry := range entries {
		matchContext, isMatch := GetMatch(entry.Name(), pattern, placeHolderMap)
		if isMatch {
			newName := newPattern
			for placeholder, value := range matchContext {
				newName = strings.ReplaceAll(newName, placeholder, value)
			}
			mvCmds = append(mvCmds, mvCmd{
				from: entry.Name(),
				to:   newName,
			})
		}
	}

	if len(mvCmds) == 0 {
		fmt.Printf("Nothing found that matches the given pattern\n")
		return 1
	}

	// Show an overview about the things that will be moved
	var overview strings.Builder
	for _, cmd := range mvCmds {
		overview.WriteString(fmt.Sprintf("%s -> %s\n", cmd.from, cmd.to))
	}
	fmt.Printf("%s", &overview)

	var prompt string
	fmt.Print("Should i mv these? [y/n]: ")
	fmt.Scan(&prompt)
	prompt = strings.ToLower(prompt)
	var shouldMv bool
	if prompt == "y" || prompt == "yes" {
		shouldMv = true
	} else if prompt == "n" || prompt == "no" {
		shouldMv = false
	} else {
		fmt.Printf("Bad value. Aborting...\n")
		return 1
	}

	if !shouldMv {
		return 0
	}

	for _, cmd := range mvCmds {
		err := os.Rename(cmd.from, cmd.to)
		if err != nil {
			fmt.Printf("%s", err)
			return 1
		}
	}
	return 0
}

func GetMatch(entry string, pattern string, placeholders map[string]struct{}) (map[string]string, bool) {

	isPlaceHolder := func(index int) bool {
		_, ok := placeholders[pattern[index:index+2]]
		return ok
	}

	placeholderMatches := make(map[string]string)

	canAddMatch := func(placeholder string, match string) bool {
		current, ok := placeholderMatches[placeholder]
		if ok {
			return current == match
		}
		return true
	}

	entryIdx := 0
	patternIdx := 0

	// go over the pattern, trying to match characters with placeholders
	for patternIdx < len(pattern) && entryIdx < len(entry) {

		// exact match
		if pattern[patternIdx] == entry[entryIdx] {
			patternIdx += 1
			entryIdx += 1
			continue
		}

		isCurrentPlaceHolder := patternIdx+1 < len(pattern) && isPlaceHolder(patternIdx)

		// pattern not matching
		if !isCurrentPlaceHolder {
			return placeholderMatches, false
		}

		currentPlaceHolder := pattern[patternIdx : patternIdx+2]
		patternIdx += 2 // we will process the placeholder

		hasCharAfterPlaceHolder := patternIdx < len(pattern)
		if !hasCharAfterPlaceHolder {
			// nothing more in the pattern, the current placeholder should match the remaining of the entry if it exists
			// if not then there is a pattern with no corresponding chars in the entry, so it's not a match
			hasMoreInEntry := entryIdx < len(entry)
			if !hasMoreInEntry {
				return placeholderMatches, false
			}

			if canAddMatch(currentPlaceHolder, entry[entryIdx:]) {
				placeholderMatches[currentPlaceHolder] = entry[entryIdx:]
			} else {
				return placeholderMatches, false
			}

			return placeholderMatches, true
		}

		// we stop expanding chars to the placeholder when meeting this char
		// try to match as many chars as possible into the placeholder
		charAfterPlaceHolder := pattern[patternIdx]
		matchingStr := ""
		for entryIdx < len(entry) && (entry[entryIdx] != charAfterPlaceHolder) {
			matchingStr += string(entry[entryIdx])
			entryIdx += 1
		}

		if entryIdx >= len(entry) {
			// no more chars in the entry, but we know we have a char after the placeholder. not a match
			return placeholderMatches, false
		}

		if canAddMatch(currentPlaceHolder, matchingStr) {
			placeholderMatches[currentPlaceHolder] = matchingStr
		} else {
			return placeholderMatches, false
		}

	}

	return placeholderMatches, patternIdx >= len(pattern)
}

func toMap(input []string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, v := range input {
		result[v] = struct{}{}
	}
	return result
}

func areMapsEqual(map1, map2 map[string]struct{}) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key := range map1 {
		_, existsIn2 := map2[key]
		if !existsIn2 {
			return false
		}
	}
	return true
}

func getMatches(pattern string) []string {
	indexes := PatternRegex.FindAllStringIndex(pattern, -1)
	if len(indexes) == 0 {
		return nil
	}
	result := make([]string, 0)
	for _, indexEntry := range indexes {
		result = append(result, pattern[indexEntry[0]:indexEntry[1]])
	}
	return result
}
