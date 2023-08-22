package main

import (
	"fmt"
	"testing"
)

func mapFromStrings(input ...string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, v := range input {
		result[v] = struct{}{}
	}
	return result
}

type tc struct {
	pattern      string
	examples     []tcrow
	placeholders map[string]struct{}
}

type tcrow struct {
	name         string
	input        string
	shouldMatch  bool
	matchContext map[string]string
}

func TestMatcher(t *testing.T) {
	testcases := []tc{
		{
			pattern:      "log_$1_another_$2",
			placeholders: mapFromStrings("$1", "$2"),
			examples: []tcrow{
				{
					name:        "good case",
					input:       "log_firstMatch_another_secondMatch",
					shouldMatch: true,
					matchContext: map[string]string{
						"$1": "firstMatch",
						"$2": "secondMatch",
					},
				},
				{
					name:        "bad case. input is shorter than pattern",
					input:       "log_firstMatch_another",
					shouldMatch: false,
				},
				{
					name:        "bad case. no match from the beginning",
					input:       "some input",
					shouldMatch: false,
				},
			},
		},
		{
			pattern:      "$1__$2",
			placeholders: mapFromStrings("$1", "$2"),
			examples: []tcrow{
				{
					name:        "good case",
					input:       "first__second",
					shouldMatch: true,
					matchContext: map[string]string{
						"$1": "first",
						"$2": "second",
					},
				},
				{
					name:        "good case 2",
					input:       "____",
					shouldMatch: true,
					matchContext: map[string]string{
						"$1": "",
						"$2": "__",
					},
				},
			},
		},
		{
			pattern:      "$1_value_$1",
			placeholders: mapFromStrings("$1"),
			examples: []tcrow{
				{
					name:        "good case",
					input:       "one_value_one",
					shouldMatch: true,
					matchContext: map[string]string{
						"$1": "one",
					},
				},
				{
					name:        "bad case. more than one match for the placeholder",
					input:       "one_value_two",
					shouldMatch: false,
				},
				{
					name:        "bad case. other match not found",
					input:       "one_value_",
					shouldMatch: false,
				},
			},
		},
	}

	for _, tc := range testcases {
		for _, example := range tc.examples {
			testname := fmt.Sprintf("(Pattern: %s) - %s", tc.pattern, example.name)

			t.Run(testname, func(t *testing.T) {

				resultContext, match := GetMatch(example.input, tc.pattern, tc.placeholders)
				if match != example.shouldMatch {
					t.Errorf("got match=%v. expected=%v\n", match, example.shouldMatch)
				}

				if example.shouldMatch {

					if len(resultContext) != len(example.matchContext) {
						t.Errorf("len of match contexts not the same\n")
					}

					for placeholder, matchedStr := range example.matchContext {
						if resultContext[placeholder] != matchedStr {
							t.Errorf("expected %s=%s got %s=%s\n", placeholder, matchedStr, placeholder, resultContext[placeholder])
						}
					}
				}
			})
		}
	}
}
