package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	//check for the argument given
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . sample.txt result.txt")
		return
		//if os.Args != 3 that is the amount of word you put in the terminal as in the example above then the program returns
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	//reading input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file")
		return
	}

	// Converting the slices of numbers to strings we can read
	text := string(data)

	// HEXADECIMAL CONVERSION

	text = regexp.MustCompile(`\b([0-9A-Fa-f]+)\s*\(hex\)`).ReplaceAllStringFunc(text,
		func(m string) string {
			sub := regexp.MustCompile(`\b([0-9A-Fa-f]+)`).FindStringSubmatch(m)
			if len(sub) < 2 {
				return m
			}

			val, _ := strconv.ParseInt(sub[1], 16, 64)
			return fmt.Sprintf("%d", val)
		})

	// BINARY CONVERSION

	text = regexp.MustCompile(`\b([01]+)\s*\(bin\)`).ReplaceAllStringFunc(text,
		func(m string) string {
			sub := regexp.MustCompile(`\b([01]+)`).FindStringSubmatch(m)
			if len(sub) < 2 {
				return m
			}

			val, _ := strconv.ParseInt(sub[1], 2, 64)
			return fmt.Sprintf("%d", val)
		})

	// CASE TRANSFORMATIONS

	transformRe := regexp.MustCompile(`\((up|low|cap)(?:,\s*(\d+))?\)`)

	for {
		loc := transformRe.FindStringIndex(text)
		if loc == nil {
			break
		}

		sub := transformRe.FindStringSubmatch(text[loc[0]:loc[1]])

		action := sub[1]

		count := 1
		if sub[2] != "" {
			count, _ = strconv.Atoi(sub[2])
		}

		// Remove token
		text = text[:loc[0]] + text[loc[1]:]

		segment := text[:loc[0]]

		wordRe := regexp.MustCompile(`\S+`)
		words := wordRe.FindAllStringIndex(segment, -1)

		if len(words) >= count {
			start := len(words) - count

			for i := start; i < len(words); i++ {
				word := segment[words[i][0]:words[i][1]]

				switch action {
				case "up":
					word = strings.ToUpper(word)
				case "low":
					word = strings.ToLower(word)
				case "cap":
					word = strings.Title(strings.ToLower(word))
				}

				text = text[:words[i][0]] + word + text[words[i][1]:]
			}
		}
	}

	//ARTICLE CORRECTION

	text = regexp.MustCompile(`\b([Aa])\s+([aeiouhAEIOUH])`).
		ReplaceAllString(text, "${1}n $2")

	// PUNCTUATION NORMALIZATION

	// 1. Remove spaces BEFORE punctuation
	text = regexp.MustCompile(`\s+([,.!?:;])`).ReplaceAllString(text, "$1")

	// 2. Ensure exactly ONE space AFTER punctuation (unless it's the end of the string)
	text = regexp.MustCompile(`([,.!?:;])([^\s])`).ReplaceAllString(text, "$1 $2")

	// 3. Special fix for the ' (quote) issue
	text = regexp.MustCompile(`'\s*(.+?)\s*'`).ReplaceAllString(text, "'$1'")

	// 4. Clean up any double spaces created by removing tags
	text = regexp.MustCompile(` +`).ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	// Replace 2 or more spaces with a single space
	text = regexp.MustCompile(` {2,}`).ReplaceAllString(text, " ")

	err = os.WriteFile(outputFile, []byte(text+"\n"), 0644)
	if err != nil {
		fmt.Println("Error writing file")
	}
}
