package key

import (
	"os"
)

func Isctrl(b byte) bool {
	return b <= 31 || b == 127
}

func Ctrl(key string) byte {
	c := []byte(key)
	return c[0] & 0x1f
}

func readKey() byte {
	var b []byte = make([]byte, 1)
	for {
		nread, err := os.Stdin.Read(b)
		if nread == 0 || err != nil {
			break
		}
		if len(b) > 0 {
			return b[0]
		}
	}
	return 0
}

func HandleKeypress(stop chan bool) {
	c := readKey()

	switch c {
	case Ctrl("x"):
		stop <- true
	}
}
