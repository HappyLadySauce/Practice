package cryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/forgoer/openssl"
)


func EncryptFile(filePath string, key [16]byte) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file %v failed: %w", filePath, err)
	}
	
	encryptContent, err := encryptLine(string(content), key)
	if err != nil {
		return fmt.Errorf("encrypt file failed: %w", err)
	}

	src := filePath + ".encrypt"
	err = os.WriteFile(src, []byte(encryptContent), 0644)
	if err != nil {
		return fmt.Errorf("write encrypted file %v failed: %w", src, err)
	}
	
	return nil
}

func DecryptFile(filePath string, key [16]byte) error {
	encryptedContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read encrypted file %v failed: %w", filePath, err)
	}
	
	decryptContent, err := decryptLine(string(encryptedContent), key)
	if err != nil {
		return fmt.Errorf("decrypt file failed: %w", err)
	}

	src := filePath + ".decrypt"
	err = os.WriteFile(src, []byte(decryptContent), 0644)
	if err != nil {
		return fmt.Errorf("write decrypted file %v failed: %w", src, err)
	}
	
	return nil
}


func encryptLine(line string, key [16]byte) (string, error) {
	src := []byte(line)
	// 进行数据填充以满足 AES 分组大小
	src = openssl.PKCS7Padding(src, 16)
	block, err := aes.NewCipher(key[:])	// 创建一个加密器
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCBCEncrypter(block, key[:])	// 采用CBC分组模式加密
	encrypted := make([]byte, len(src))	// 给密文申请内存空间
	encrypter.CryptBlocks(encrypted, src)	//进行加密
	return hex.EncodeToString(encrypted), nil
}

func decryptLine(cryptLine string, key [16]byte) (string, error) {
	// 先将十六进制字符串解码为字节数组
	src, err := hex.DecodeString(cryptLine)
	if err != nil {
		return "", fmt.Errorf("hex decode failed: %v", err)
	}
	
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	
	// 使用相同的key作为IV进行解密
	decrypter := cipher.NewCBCDecrypter(block, key[:])
	decrypted := make([]byte, len(src))
	decrypter.CryptBlocks(decrypted, src)
	
	out, err := openssl.PKCS7UnPadding(decrypted)	// 进行反填充
	if err != nil {
		return "", err
	}
	
	return string(out), nil
}


















