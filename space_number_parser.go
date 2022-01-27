package main

import (
	"os"
	"bufio"
	"fmt"
	"strconv"
//	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	lineBreak = &unicode.RangeTable{
		R16: []unicode.Range16{
			{
				Lo:     0x0A,
				Hi:     0x0D,
				Stride: 3,
			},
		},
		R32:         nil,
		LatinOffset: 0,
	}
)

func main() {
	var numbers []int64

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(scanWordsUntil(lineBreak))
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			panic(err)
		}

		i, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			panic(err)
		}
		numbers = append(numbers, i)
	}

	fmt.Println("Numbers:", numbers)
}

func scanWordsUntil(rangeTable *unicode.RangeTable) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		var pos int

		for width := 0; pos < len(data); pos += width {
			var r rune
			r, width = utf8.DecodeRune(data[pos:])
			if unicode.Is(rangeTable, r) {
				return 0, nil, bufio.ErrFinalToken
			}
			if !unicode.IsSpace(r) {
				break
			}
		}

		for width, i := 0, pos; i < len(data); i += width {
			var r rune
			r, width = utf8.DecodeRune(data[i:])
			if unicode.Is(rangeTable, r) {
				return i + width, data[pos:i], bufio.ErrFinalToken
			}
			if unicode.IsSpace(r) {

				return i + width, data[pos:i], nil
			}
		}

		if atEOF && len(data) > pos {
			return len(data), data[pos:], nil
		}

		return pos, nil, nil
	}
}

