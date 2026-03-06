/*package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . sample.txt result.txt")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	text := string(content)

	// Convert hex numbers
	hexRe := regexp.MustCompile(`\b([0-9A-Fa-f]+)\s*\(hex\)`)
	text = hexRe.ReplaceAllStringFunc(text, func(match string) string {
		sub := hexRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}

		val, _ := strconv.ParseInt(sub[1], 16, 64)
		return fmt.Sprintf("%d", val)
	})

	// Convert binary numbers
	binRe := regexp.MustCompile(`\b([01]+)\s*\(bin\)`)
	text = binRe.ReplaceAllStringFunc(text, func(match string) string {
		sub := binRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}

		val, _ := strconv.ParseInt(sub[1], 2, 64)
		return fmt.Sprintf("%d", val)
	})

	// Case transformations
	// Case transformations using regex
	transformRe := regexp.MustCompile(`\((up|low|cap)(?:,\s*(\d+))?\)`)

	for {
		loc := transformRe.FindStringIndex(text)
		if loc == nil {
			break
		}

		match := text[loc[0]:loc[1]]

		sub := transformRe.FindStringSubmatch(match)

		action := sub[1]

		count := 1
		if sub[2] != "" {
			count, _ = strconv.Atoi(sub[2])
		}

		// Remove transformation token
		text = text[:loc[0]] + text[loc[1]:]

		// Split words safely
		words := strings.Fields(text[:loc[0]])
		if len(words) < count {
			continue
		}

		startIndex := len(words) - count

		for i := startIndex; i < len(words); i++ {
			switch action {
			case "up":
				words[i] = strings.ToUpper(words[i])
			case "low":
				words[i] = strings.ToLower(words[i])
			case "cap":
				words[i] = strings.Title(strings.ToLower(words[i]))
			}
		}

		// Rebuild text
		text = strings.Join(words, " ") + text[loc[0]:]
	}

	// Fix punctuation spacing
	punctRe := regexp.MustCompile(`\s+([,\.!?:;])`)
	text = punctRe.ReplaceAllString(text, "$1")

	// Normalize ellipsis
	text = regexp.MustCompile(`\s*\.\s*\.\s*\.`).ReplaceAllString(text, "...")

	// Fix quotes spacing
	text = regexp.MustCompile(`'\s+`).ReplaceAllString(text, "'")
	text = regexp.MustCompile(`\s+'`).ReplaceAllString(text, "'")

	// Fix article a → an
	text = fixArticles(text)

	// Write output file
	err = os.WriteFile(outputFile, []byte(text), 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
	}
}

func filterEmpty(arr []string) []string {
	res := []string{}
	for _, v := range arr {
		if v != "" {
			res = append(res, v)
		}
	}
	return res
}

func fixArticles(text string) string {
	re := regexp.MustCompile(`\b([Aa])\s+([aeiouhAEIOUH])`)
	return re.ReplaceAllString(text, "${1}n $2")
}
*/

/*
package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . sample.txt result.txt")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	text := string(content)

	// ---------- Hex conversion ----------
	hexRe := regexp.MustCompile(`\b([0-9A-Fa-f]+)\s*\(hex\)`)
	text = hexRe.ReplaceAllStringFunc(text, func(match string) string {
		sub := hexRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}

		val, _ := strconv.ParseInt(sub[1], 16, 64)
		return fmt.Sprintf("%d", val)
	})

	// ---------- Binary conversion ----------
	binRe := regexp.MustCompile(`\b([01]+)\s*\(bin\)`)
	text = binRe.ReplaceAllStringFunc(text, func(match string) string {
		sub := binRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}

		val, _ := strconv.ParseInt(sub[1], 2, 64)
		return fmt.Sprintf("%d", val)
	})

	// ---------- Case transformation ----------
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

		// Remove token first
		text = text[:loc[0]] + text[loc[1]:]

		// Work only on segment before token (preserves newlines)
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
	// ---------- Article correction ----------
	text = fixArticles(text)

	// ---------- Normalize punctuation ----------
	text = normalizeText(text)

	// ---------- Write output with newline ----------
	err = os.WriteFile(outputFile, []byte(text+"\n"), 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
	}
}

func fixArticles(text string) string {
	re := regexp.MustCompile(`\b([Aa])\s+([aeiouhAEIOUH])`)
	return re.ReplaceAllString(text, "${1}n $2")
}

func normalizeText(text string) string {
	text = regexp.MustCompile(`\s+([,\.!?:;])`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`([,\.!?:;])([^\s])`).ReplaceAllString(text, "$1 $2")
	text = regexp.MustCompile(`\s*\.\s*\.\s*\.`).ReplaceAllString(text, "...")
	text = regexp.MustCompile(`'\s+`).ReplaceAllString(text, "'")
	text = regexp.MustCompile(`\s+'`).ReplaceAllString(text, "'")

	return text
}
*/

/*package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . sample.txt result.txt")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	text := string(content)

	// ---------- Hex conversion ----------
	hexRe := regexp.MustCompile(`\b([0-9A-Fa-f]+)\s*\(hex\)`)
	text = hexRe.ReplaceAllStringFunc(text, func(match string) string {
		sub := hexRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}

		val, _ := strconv.ParseInt(sub[1], 16, 64)
		return fmt.Sprintf("%d", val)
	})

	// ---------- Binary conversion ----------
	binRe := regexp.MustCompile(`\b([01]+)\s*\(bin\)`)
	text = binRe.ReplaceAllStringFunc(text, func(match string) string {
		sub := binRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}

		val, _ := strconv.ParseInt(sub[1], 2, 64)
		return fmt.Sprintf("%d", val)
	})

	// ---------- Case transformation ----------
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

		// Remove token first
		text = text[:loc[0]] + text[loc[1]:]

		// Work only on segment before token (preserves newlines)
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
	// ---------- Article correction ----------
	text = fixArticles(text)

	// ---------- Normalize punctuation ----------
	text = normalizeText(text)

	// ---------- Write output with newline ----------
	err = os.WriteFile(outputFile, []byte(text+"\n"), 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
	}
}

func fixArticles(text string) string {
	re := regexp.MustCompile(`\b([Aa])\s+([aeiouhAEIOUH])`)
	return re.ReplaceAllString(text, "${1}n $2")
}

func normalizeText(text string) string {

	// Remove space before punctuation
	text = regexp.MustCompile(`\s+([,.!?:;])`).ReplaceAllString(text, "$1")

	// Remove space after punctuation if it is misplaced
	text = regexp.MustCompile(`([,.!?:;])\s+`).ReplaceAllString(text, "$1 ")

	// Fix double punctuation like ! !
	text = regexp.MustCompile(`([!?.])\s*\1`).ReplaceAllStringFunc(text, func(m string) string {
	r := regexp.MustCompile(`([!?.])`)
	sub := r.FindStringSubmatch(m)
	if len(sub) > 1 {
		return sub[1] + sub[1]
	}
	return m
})

	// Normalize ellipsis
	text = regexp.MustCompile(`\s*\.\s*\.\s*\.`).ReplaceAllString(text, "...")

	// Fix quotes spacing
	text = regexp.MustCompile(`'\s*`).ReplaceAllString(text, "'")
	text = regexp.MustCompile(`\s*'`).ReplaceAllString(text, "'")

	// Fix colon quote pattern: : 'word'
	text = regexp.MustCompile(`:\s*'`).ReplaceAllString(text, ": '")

	return text
}
*/

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . sample.txt result.txt")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file")
		return
	}

	text := string(data)

	/* ================= HEX CONVERSION ================= */

	text = regexp.MustCompile(`\b([0-9A-Fa-f]+)\s*\(hex\)`).ReplaceAllStringFunc(text,
		func(m string) string {
			sub := regexp.MustCompile(`\b([0-9A-Fa-f]+)`).FindStringSubmatch(m)
			if len(sub) < 2 {
				return m
			}

			val, _ := strconv.ParseInt(sub[1], 16, 64)
			return fmt.Sprintf("%d", val)
		})

	/* ================= BINARY CONVERSION ================= */

	text = regexp.MustCompile(`\b([01]+)\s*\(bin\)`).ReplaceAllStringFunc(text,
		func(m string) string {
			sub := regexp.MustCompile(`\b([01]+)`).FindStringSubmatch(m)
			if len(sub) < 2 {
				return m
			}

			val, _ := strconv.ParseInt(sub[1], 2, 64)
			return fmt.Sprintf("%d", val)
		})

	/* ================= CASE TRANSFORMATIONS ================= */

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

	/* ================= ARTICLE CORRECTION ================= */

	text = regexp.MustCompile(`\b([Aa])\s+([aeiouhAEIOUH])`).
		ReplaceAllString(text, "${1}n $2")

	/* ================= PUNCTUATION NORMALIZATION ================= */

	text = regexp.MustCompile(`\s+([,.!?:;])`).ReplaceAllString(text, "$1")
	text = regexp.MustCompile(`([,.!?:;])([^\s])`).ReplaceAllString(text, "$1 $2")

	text = regexp.MustCompile(`([!?.])\s+([!?.])`).ReplaceAllString(text, "$1$2")

	text = regexp.MustCompile(`\s*\.\s*\.\s*\.`).ReplaceAllString(text, "...")

	text = regexp.MustCompile(`'\s*`).ReplaceAllString(text, "'")
	text = regexp.MustCompile(`\s*'`).ReplaceAllString(text, "'")

	text = regexp.MustCompile(`:\s*'`).ReplaceAllString(text, ": '")

	/* ================= LINE SEPARATION FIX ================= */

	// Normalize multiple newlines
	text = regexp.MustCompile(`\n{2,}`).ReplaceAllString(text, "\n")

	// Insert newline between sentences starting with capital letter after punctuation
	text = regexp.MustCompile(`([.?!])([A-Z])`).ReplaceAllString(text, "$1\n$2")

	/* ================= FINAL CLEANUP ================= */

	text = strings.TrimRight(text, " \t\r\n")

	err = os.WriteFile(outputFile, []byte(text+"\n"), 0644)
	if err != nil {
		fmt.Println("Error writing file")
	}
}
