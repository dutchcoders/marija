// Package mention provides function for parsing twitter like mentions and hashtags
package mention

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// GetTags returns a slice of tags, that is all characters after rune char up to occurance of space
// or another occurance of rune char. Additionally you can provide a coma separated unicode characters to
// be used as terminating sequence.
func GetTags(char rune, src io.Reader, terminator ...rune) []string {
	t := getTag(char, src, terminator...)
	var out []string
	for v := range t {
		out = append(out, v)
	}
	return out
}

// getTag sends matched tags to the output channel, the output channel is closed
// when scanning is complete. it is safe to use range on the output channel.
func getTag(char rune, src io.Reader, terminator ...rune) <-chan string {
	out := make(chan string, 2)
	go func() {
		scan := bufio.NewScanner(src)
		scan.Split(splitTag(char, terminator...))
		for scan.Scan() {
			txt := scan.Text()
			if txt != "" {
				out <- txt
			}
		}
		if err := scan.Err(); err != nil {
			fmt.Println(err)
		}
		close(out)
	}()
	return out
}

// splitTag splits the input at any occurence of rune char up to eather space of
// another occurence of rune char.
//
// Example:
// 	@gernest will be split when char is @ resulting in gernest
func splitTag(char rune, terminator ...rune) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		start := getRuneBytes(char)
		if x := bytes.Index(data, start); x >= 0 {
			begin := x + len(start)
			for n := begin; n < len(data); n++ {
				xFirst, width := utf8.DecodeRune(data[n:])
				if n == x+len(start) {
					if unicode.IsSpace(xFirst) {
						return n + width, data[begin:n], nil
					}
				} else {
					// we stop when we encounter another char instance
					// this can be a case like @gernest@gernest
					if xFirst == char {
						// in the case we have mulitple @ (ie `@@@@`), we return no mention
						if n-begin == 1 {
							return n + width, data[begin : n-1], nil
						}
						return n - width, data[begin:n], nil
					}

					for _, term := range terminator {
						if xFirst == term {
							// If when reaching our terminator, its the only
							// character in the data, ignore the mention. (ie "@,")
							if n-begin == 1 {
								return n + width, data[begin : n-1], nil
							}
							return n + width, data[begin:n], nil
						}
					}

					// the end of our tag
					if unicode.IsSpace(xFirst) {
						// make sure our result isn't just a single terminator (ie "@@")
						for _, term := range terminator {
							if n-begin == 1 && rune(data[begin:n][0]) == term {
								return n + width, data[begin : n-1], nil
							}
						}
						return n + width, data[begin:n], nil
					}
				}

			}
			if atEOF && len(data) > begin {
				return len(data), data[begin:], nil
			}
			return x, nil, nil
		}
		return len(data), nil, nil
	}
}

// getRuneBytes encodes rune r into bytes and return the byte slice.
func getRuneBytes(r rune) []byte {
	rst := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(rst, r)
	return rst
}
