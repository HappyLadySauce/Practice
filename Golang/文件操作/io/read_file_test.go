package io_test

import (
	"fmt"
	"file/io"
	"testing"
)

func TestReadFile(t *testing.T) {
	info := io.ReadFile("../data/test_file.text")
	fmt.Println(info)
}

func TestReadFileWithBuffer(t *testing.T) {
	info := io.ReadFileWithBuffer("../data/test_file.text")
	fmt.Println(info)
}