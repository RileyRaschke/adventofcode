package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func AllUniqueChars(x string) bool {
	unique := true
	for idx, c := range x {
		lastIdx := strings.LastIndex(x, string(c))
		if lastIdx != idx {
			unique = false
			break
		}
	}
	return unique
}

func main() {
	var pos, ss, ms int = 0, -1, -1
	ssBuf, msBuf := "", ""
	reader := bufio.NewReader(os.Stdin)
	for {
		rune, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		pos++
		if ss < 0 {
			if len(ssBuf) < 4 {
				ssBuf += string(rune) // push
				continue
			}
			ssBuf = ssBuf[1:]     // pop front
			ssBuf += string(rune) // push
			if AllUniqueChars(ssBuf) {
				ss = pos
				msBuf = ssBuf
			}
		} else {
			if len(msBuf) < 14 {
				msBuf += string(rune) // push
				continue
			}
			msBuf = msBuf[1:]     // pop front
			msBuf += string(rune) // push
			if AllUniqueChars(msBuf) {
				ms = pos
				break
			}
		}
	}
	fmt.Printf("Packet found at %d\n", ss)
	fmt.Printf("Start-Of-Message found at %d\n", ms)
}
