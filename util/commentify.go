package util

import (
	"bytes"
)

// TODO: add lineCharacterLimit to configurable options
var lineCharacterLimit int = 80
var commentSpace string = "  "

var commentChars = map[string][]string{
	"javascript": []string{"/*", "*/"},
	"typescript": []string{"/*", "*/"},
	"python":     []string{`"""`, `"""`},
	"ruby":       []string{"=begin", "=end"},
	"golang":     []string{"/*", "*/"},
	"scala":      []string{"/*", "*/"},
	"cpp":        []string{"/*", "*/"},
	"python3":    []string{`"""`, `"""`},
	"c":          []string{"/*", "*/"},
	"csharp":     []string{"/*", "*/"},
	"swift":      []string{"/*", "*/"},
	"java":       []string{"/*", "*/"},
	"rust":       []string{"/*", "*/"},
	"php":        []string{"/*", "*/"},
	"kotlin":     []string{"/*", "*/"},
}

func Commentify(s string, language string) []byte {
	var buf bytes.Buffer
	var lineLength int
	var currentWord = []rune(commentSpace)

	buf.WriteString(commentChars[language][0])
	buf.WriteByte('\n')
	for _, ch := range s {
		lineLength++

		currentWord = append(currentWord, ch)

		if ch == ' ' {
			buf.WriteString(string(currentWord))
			currentWord = []rune{}
		}

		if ch == '\n' {
			buf.WriteString(string(currentWord))
			currentWord = []rune(commentSpace)
			lineLength = len(currentWord)
		}

		if lineLength >= lineCharacterLimit-1 {
			buf.WriteByte('\n')
			currentWord = append([]rune(commentSpace), currentWord...)
			lineLength = len(currentWord)
		}
	}
	buf.WriteString(string(currentWord))
	buf.WriteByte('\n')
	buf.WriteString(commentChars[language][1])

	return buf.Bytes()
}
