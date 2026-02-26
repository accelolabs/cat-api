package meow

import (
	"math/rand/v2"
	"slices"
	"strings"
)

var meows = []string{
	"meow",
	"maow",
	"miaw",
	"nya",
	"mew",
	"purr",
	"prrp",
	"mrp",
	"mrow",
}

var stretchRunes = []rune{'e', 'o', 'a', 'i', 'y', 'r', 'w'}

// 3: meeooow
func getMeow(maxStretch int) string {
    if maxStretch < 1 {
        maxStretch = 1
    }

	meow := []rune(meows[rand.IntN(len(meows))])

	var builder strings.Builder
	builder.Grow(len(meow) * maxStretch)

	for _, m := range meow {
		if !slices.Contains(stretchRunes, m) {
			builder.WriteRune(m)
			continue
		}

		for range rand.IntN(int(maxStretch)) + 1 {
			builder.WriteRune(m)
		}
	}

	return builder.String()
}

// (5, 4): meeooow-miiaaaww-meeeewwww-nyyyyaaa-purrr
func Meow(length int, maxStretch int) string {
	var builder strings.Builder

	for range length - 1 {
		builder.WriteString(getMeow(maxStretch))
		builder.WriteRune('-')
	}

	builder.WriteString(getMeow(maxStretch))

	return builder.String()
}
