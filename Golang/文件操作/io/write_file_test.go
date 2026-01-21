package io_test

import (
	"file/io"
	"testing"
)

func TestWriteFile(t *testing.T) {
	io.Writefile("../data/test_file.text","你好世界")
}