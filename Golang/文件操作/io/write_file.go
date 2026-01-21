package io

import (
	"fmt"
	"os"
)

func Writefile(file, text string) {
	hand, err := os.OpenFile(file, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0o666)
	if err != nil {
		fmt.Printf("open file failed. err: %s\n", err)
	}else {
		defer hand.Close()
		hand.WriteString(text)
	}
}

