package cryption

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	// 测试密钥
	key := []byte("1234567890123456") // 16字节密钥
	
	// 测试明文
	originalText := "你好！我是HappyLadySauce"
	
	// 加密
	encryptedText, err := SymmetricEncrypt(originalText, key, 16, AesCryptMode)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}
	
	t.Logf("原始文本: %s", originalText)
	t.Logf("加密后: %s", encryptedText)
	
	// 解密
	decryptedText, err := SymmetricDecrypt(encryptedText, key, 16, AesCryptMode)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	
	t.Logf("解密后: %s", decryptedText)
	
	// 验证结果
	if decryptedText != originalText {
		t.Errorf("解密结果不匹配! 期望: %s, 实际: %s", originalText, decryptedText)
	}
}

func TestEmptyText(t *testing.T) {
	key := []byte("1234567890123456")
	
	// 测试空字符串 - 现在应该返回错误
	_, err := SymmetricEncrypt("", key, 16, AesCryptMode)
	if err == nil {
		t.Error("应该拒绝空明文")
	}
	
	// 测试只包含空格的字符串
	encrypted, err := SymmetricEncrypt("   ", key, 16, AesCryptMode)
	if err != nil {
		t.Fatalf("空格字符串加密失败: %v", err)
	}
	
	decrypted, err := SymmetricDecrypt(encrypted, key, 16, AesCryptMode)
	if err != nil {
		t.Fatalf("空格字符串解密失败: %v", err)
	}
	
	if decrypted != "   " {
		t.Errorf("空格字符串解密结果不匹配! 期望: '   ', 实际: '%s'", decrypted)
	}
}

func TestDesEncryptDecrypt(t *testing.T) {
	// 测试密钥 - DES需要8字节密钥
	key := []byte("12345678") // 8字节密钥
	
	// 测试明文
	originalText := "你好！我是HappyLadySauce"
	
	// 加密
	encryptedText, err := SymmetricEncrypt(originalText, key, 8, DesCryptMode)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}
	
	t.Logf("原始文本: %s", originalText)
	t.Logf("加密后: %s", encryptedText)
	
	// 解密
	decryptedText, err := SymmetricDecrypt(encryptedText, key, 8, DesCryptMode)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	
	t.Logf("解密后: %s", decryptedText)
	
	// 验证结果
	if decryptedText != originalText {
		t.Errorf("解密结果不匹配! 期望: %s, 实际: %s", originalText, decryptedText)
	}
}

func TestAesConvenienceFunctions(t *testing.T) {
	// 测试便利函数
	key := []byte("12345678901234567890123456789012") // 32字节密钥
	originalText := "测试便利函数"
	
	// 使用便利函数加密
	encryptedText, err := AesEncrypt(originalText, key)
	if err != nil {
		t.Fatalf("AES加密失败: %v", err)
	}
	
	t.Logf("AES加密后: %s", encryptedText)
	
	// 使用便利函数解密
	decryptedText, err := AesDecrypt(encryptedText, key)
	if err != nil {
		t.Fatalf("AES解密失败: %v", err)
	}
	
	t.Logf("AES解密后: %s", decryptedText)
	
	// 验证结果
	if decryptedText != originalText {
		t.Errorf("AES便利函数解密结果不匹配! 期望: %s, 实际: %s", originalText, decryptedText)
	}
}

func TestDesConvenienceFunctions(t *testing.T) {
	// 测试便利函数
	key := []byte("12345678") // 8字节密钥
	originalText := "测试DES便利函数"
	
	// 使用便利函数加密
	encryptedText, err := DesEncrypt(originalText, key)
	if err != nil {
		t.Fatalf("DES加密失败: %v", err)
	}
	
	t.Logf("DES加密后: %s", encryptedText)
	
	// 使用便利函数解密
	decryptedText, err := DesDecrypt(encryptedText, key)
	if err != nil {
		t.Fatalf("DES解密失败: %v", err)
	}
	
	t.Logf("DES解密后: %s", decryptedText)
	
	// 验证结果
	if decryptedText != originalText {
		t.Errorf("DES便利函数解密结果不匹配! 期望: %s, 实际: %s", originalText, decryptedText)
	}
}

func TestEncryptionSecurity(t *testing.T) {
	// 测试加密安全性 - 相同明文应该产生不同密文
	key := []byte("1234567890123456")
	originalText := "测试安全性"
	
	// 加密两次
	encrypted1, err := AesEncrypt(originalText, key)
	if err != nil {
		t.Fatalf("第一次加密失败: %v", err)
	}
	
	encrypted2, err := AesEncrypt(originalText, key)
	if err != nil {
		t.Fatalf("第二次加密失败: %v", err)
	}
	
	t.Logf("第一次加密: %s", encrypted1)
	t.Logf("第二次加密: %s", encrypted2)
	
	// 由于使用了随机IV，两次加密结果应该不同
	if encrypted1 == encrypted2 {
		t.Error("安全性问题：相同明文产生了相同的密文")
	}
	
	// 但都能正确解密
	decrypted1, err := AesDecrypt(encrypted1, key)
	if err != nil {
		t.Fatalf("第一次解密失败: %v", err)
	}
	
	decrypted2, err := AesDecrypt(encrypted2, key)
	if err != nil {
		t.Fatalf("第二次解密失败: %v", err)
	}
	
	if decrypted1 != originalText || decrypted2 != originalText {
		t.Error("解密结果不正确")
	}
}

func TestErrorHandling(t *testing.T) {
	// 测试错误处理
	key := []byte("1234567890123456")
	
	// 测试空明文
	_, err := AesEncrypt("", key)
	if err == nil {
		t.Error("应该拒绝空明文")
	}
	
	// 测试空密钥
	_, err = AesEncrypt("测试", []byte{})
	if err == nil {
		t.Error("应该拒绝空密钥")
	}
	
	// 测试短密钥
	_, err = AesEncrypt("测试", []byte("123"))
	if err == nil {
		t.Error("应该拒绝过短的密钥")
	}
	
	// 测试空密文
	_, err = AesDecrypt("", key)
	if err == nil {
		t.Error("应该拒绝空密文")
	}
	
	// 测试格式错误的密文
	_, err = AesDecrypt("invalid", key)
	if err == nil {
		t.Error("应该拒绝格式错误的密文")
	}
}

