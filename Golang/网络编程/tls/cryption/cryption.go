package cryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"

	ecies "github.com/ecies/go/v2"
	"github.com/forgoer/openssl"
)

const (
	_ = iota
	AesCryptMode
	DesCryptMode
)

const (
	// 默认块大小
	AesBlockSize = 16 // AES块大小
	DesBlockSize = 8  // DES块大小
	
	// 最小密钥长度
	MinAesKeySize = 16 // AES最小密钥长度
	MinDesKeySize = 8  // DES最小密钥长度
)

// validateSymmetricKey 验证对称加密密钥
func validateSymmetricKey(key []byte, encryptMode int) error {
	if len(key) == 0 {
		return errors.New("密钥不能为空")
	}
	
	switch encryptMode {
	case AesCryptMode:
		if len(key) < MinAesKeySize {
			return fmt.Errorf("AES密钥长度不能小于%d字节", MinAesKeySize)
		}
	case DesCryptMode:
		if len(key) < MinDesKeySize {
			return fmt.Errorf("DES密钥长度不能小于%d字节", MinDesKeySize)
		}
	default:
		return fmt.Errorf("不支持的加密模式: %d", encryptMode)
	}
	
	return nil
}

// generateSecureIV 生成安全的随机IV
func generateSecureIV(blockSize int) ([]byte, error) {
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("生成IV失败: %w", err)
	}
	return iv, nil
}

func SymmetricEncrypt(plainText string, key []byte, blockSize int, encryptMode int) (string, error) {
	// 参数验证
	if err := validateSymmetricKey(key, encryptMode); err != nil {
		return "", err
	}
	
	if len(plainText) == 0 {
		return "", errors.New("明文不能为空")
	}
	
	var block cipher.Block
	var err error
	
	// 根据加密模式创建加密器
	switch encryptMode {
	case AesCryptMode:
		block, err = aes.NewCipher(key)
		if err != nil {
			return "", fmt.Errorf("创建AES加密器失败: %w", err)
		}
	case DesCryptMode:
		block, err = des.NewCipher(key)
		if err != nil {
			return "", fmt.Errorf("创建DES加密器失败: %w", err)
		}
	default:
		return "", fmt.Errorf("不支持的加密模式: %d", encryptMode)
	}

	// 数据填充（使用块大小）
	paddedText := openssl.PKCS7Padding([]byte(plainText), blockSize)
	
	// 创建加密结果缓冲区（与填充后数据相同大小）
	encrypted := make([]byte, len(paddedText))
	
	// 生成安全的随机IV
	iv, err := generateSecureIV(blockSize)
	if err != nil {
		return "", err
	}
	
	// 创建CBC加密器
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(encrypted, paddedText)
	
	// 将IV和密文一起返回，格式: hex(iv + ciphertext)
	result := append(iv, encrypted...)
	return hex.EncodeToString(result), nil
}

func SymmetricDecrypt(cipherText string, key []byte, blockSize int, decryptMode int) (string, error) {
	// 参数验证
	if err := validateSymmetricKey(key, decryptMode); err != nil {
		return "", err
	}
	
	if len(cipherText) == 0 {
		return "", errors.New("密文不能为空")
	}
	
	// 将十六进制字符串解码为字节数组
	data, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", fmt.Errorf("十六进制解码失败: %w", err)
	}
	
	// 检查数据长度是否足够（至少包含IV）
	if len(data) < blockSize {
		return "", errors.New("密文格式错误：长度不足")
	}
	
	// 提取IV和密文
	iv := data[:blockSize]
	encrypted := data[blockSize:]
	
	var block cipher.Block
	
	// 根据解密模式创建解密器
	switch decryptMode {
	case AesCryptMode:
		block, err = aes.NewCipher(key)
		if err != nil {
			return "", fmt.Errorf("创建AES解密器失败: %w", err)
		}
	case DesCryptMode:
		block, err = des.NewCipher(key)
		if err != nil {
			return "", fmt.Errorf("创建DES解密器失败: %w", err)
		}
	default:
		return "", fmt.Errorf("不支持的解密模式: %d", decryptMode)
	}
	
	// 创建解密结果缓冲区
	decrypted := make([]byte, len(encrypted))
	
	// 创建CBC解密器并解密
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(decrypted, encrypted)
	
	// 移除填充
	plainText, err := openssl.PKCS7UnPadding(decrypted)
	if err != nil {
		return "", fmt.Errorf("移除填充失败: %w", err)
	}
	
	// 返回原始明文字符串
	return string(plainText), nil
}

func RsaEncrypt(plainText []byte, publicKey []byte) (string, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", fmt.Errorf("PEM公钥块为空")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("解析公钥失败: %w", err)
	}
	pub := pubInterface.(*rsa.PublicKey)
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pub, plainText)
	if err != nil {
		return "", fmt.Errorf("RSA加密失败: %w", err)
	}
	return hex.EncodeToString(cipherText), nil
}

func RsaDecrypt(cipherText []byte, privateKey []byte) (string, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", fmt.Errorf("PEM私钥块为空")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %w", err)
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
	if err != nil {
		return "", fmt.Errorf("RSA解密失败: %w", err)
	}
	return string(plainText), nil
}

func GetPrivateKey() (*ecies.PrivateKey, error) {
	return ecies.GenerateKey()
}

func EccEncrypt(plainText string, publicKey *ecies.PublicKey) (string, error) {
	src := []byte(plainText)
	cipherText, err := ecies.Encrypt(publicKey, src)
	if err != nil {
		return "", fmt.Errorf("ECC加密失败: %w", err)
	}
	return hex.EncodeToString(cipherText), nil
}

func EccDecrypt(cipherText []byte, privateKey *ecies.PrivateKey) (string, error) {
	plainText, err := ecies.Decrypt(privateKey, cipherText)
	if err != nil {
		return "", fmt.Errorf("ECC解密失败: %w", err)
	}
	return string(plainText), nil
}

// ==================== 便利函数 ====================

// AesEncrypt AES加密便利函数
func AesEncrypt(plainText string, key []byte) (string, error) {
	return SymmetricEncrypt(plainText, key, AesBlockSize, AesCryptMode)
}

// AesDecrypt AES解密便利函数
func AesDecrypt(cipherText string, key []byte) (string, error) {
	return SymmetricDecrypt(cipherText, key, AesBlockSize, AesCryptMode)
}

// DesEncrypt DES加密便利函数
func DesEncrypt(plainText string, key []byte) (string, error) {
	return SymmetricEncrypt(plainText, key, DesBlockSize, DesCryptMode)
}

// DesDecrypt DES解密便利函数
func DesDecrypt(cipherText string, key []byte) (string, error) {
	return SymmetricDecrypt(cipherText, key, DesBlockSize, DesCryptMode)
}