package cryption_test

import (
	"network/cryption"
	"testing"
)


func TestCryptFile(t *testing.T) {
	var key [16]byte = [16]byte{'1','2','3','4','5','6','7','8','1','2','3','4','5','6','7','8'}
	
	// 加密文件
	if err := cryption.EncryptFile("./data/test.txt", key); err != nil {
		t.Fatalf("加密失败: %v", err)
	}
	
	// 解密文件
	if err := cryption.DecryptFile("./data/test.txt.encrypt", key); err != nil {
		t.Fatalf("解密失败: %v", err)
	}
}

