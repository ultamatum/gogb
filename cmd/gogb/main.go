package main

import (
	"zandergray.dev/gogb/internal"
)

func main() {
	var tetrisPath string = "/home/alexander/Downloads/tetris.gb"

	if internal.CartLoad(tetrisPath) {
		internal.CartPrintInfo()
	}
}
