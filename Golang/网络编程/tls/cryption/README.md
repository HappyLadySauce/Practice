# 加密库 (Cryption Package)

一个功能完整的Go语言加密库，支持对称加密(AES/DES)、非对称加密(RSA)和椭圆曲线加密(ECC)。

## 🚀 功能特性

### ✅ 已优化功能
- **安全性增强**: 使用随机IV代替固定IV
- **参数验证**: 完善的输入参数检查
- **错误处理**: 详细的错误信息和异常处理
- **便利函数**: 简化的AES/DES加密接口
- **代码结构**: 清晰的代码组织和常量定义

### 🔐 支持的加密算法
- **AES**: 高级加密标准 (支持128/192/256位密钥)
- **DES**: 数据加密标准 (56位密钥)
- **RSA**: 非对称加密 (支持公钥加密)
- **ECC**: 椭圆曲线加密

## 📦 安装使用

```go
import "network/tls/cryption"
```

## 💡 快速开始

### AES加密/解密
```go
// 使用便利函数
key := []byte("1234567890123456") // 16字节密钥
plainText := "你好，世界！"

// 加密
encrypted, err := cryption.AesEncrypt(plainText, key)
if err != nil {
    log.Fatal(err)
}

// 解密
decrypted, err := cryption.AesDecrypt(encrypted, key)
if err != nil {
    log.Fatal(err)
}
```

### DES加密/解密
```go
key := []byte("12345678") // 8字节密钥
plainText := "Hello World"

// 加密
encrypted, err := cryption.DesEncrypt(plainText, key)
if err != nil {
    log.Fatal(err)
}

// 解密
decrypted, err := cryption.DesDecrypt(encrypted, key)
if err != nil {
    log.Fatal(err)
}
```

### 高级用法 - 指定参数
```go
// 使用底层函数，可以指定块大小和模式
key := []byte("1234567890123456")
plainText := "需要加密的内容"

// AES加密，16字节块大小
encrypted, err := cryption.SymmetricEncrypt(plainText, key, 16, cryption.AesCryptMode)
if err != nil {
    log.Fatal(err)
}

// 解密
decrypted, err := cryption.SymmetricDecrypt(encrypted, key, 16, cryption.AesCryptMode)
if err != nil {
    log.Fatal(err)
}
```

## 🔑 密钥要求

| 算法 | 密钥长度 | 说明 |
|------|----------|------|
| AES  | 16/24/32 字节 | 分别对应128/192/256位 |
| DES  | 8 字节 | 固定长度 |

## 🛡️ 安全特性

### ✅ 已实现
- **随机IV**: 每次加密生成不同的随机初始化向量
- **参数验证**: 严格的密钥长度和内容检查
- **错误包装**: 详细的错误信息，便于调试
- **安全编码**: 使用十六进制编码存储密文

### 🔍 加密格式
优化后的加密格式：
```
密文 = hex(IV + 加密数据)
```
- IV长度: AES=16字节, DES=8字节
- 每次加密产生不同的密文（即使明文相同）
- 兼容解密函数自动提取IV

## 🧪 测试

运行所有测试：
```bash
go test ./tls/cryption/ -v
```

测试覆盖：
- ✅ 基本加密/解密功能
- ✅ 空字符串处理
- ✅ 错误处理
- ✅ 安全性测试（随机IV）
- ✅ 便利函数测试

## 📋 错误处理

库函数会返回详细的错误信息：
```go
// 空明文错误
_, err := cryption.AesEncrypt("", key)
// 返回: 明文不能为空

// 短密钥错误  
_, err := cryption.AesEncrypt("测试", []byte("123"))
// 返回: AES密钥长度不能小于16字节

// 格式错误
_, err := cryption.AesDecrypt("invalid", key)
// 返回: 十六进制解码失败/密文格式错误
```

## 🔄 向后兼容性

⚠️ **注意**: 由于优化了IV生成方式，新版本的加密结果与旧版本不兼容。如果需要兼容旧版本数据，请使用原来的固定IV方式。

## 📁 文件结构
```
tls/cryption/
├── cryption.go      # 主要加密实现
├── cryption_test.go # 测试文件
├── data/           # 示例和测试数据
│   ├── example.go  # 使用示例
│   └── README.md   # 本文档
└── README.md       # 项目文档
```

## 🎯 使用建议

1. **选择合适的算法**: AES比DES更安全，推荐使用AES
2. **密钥管理**: 使用足够长的随机密钥
3. **错误处理**: 始终检查和处理返回的错误
4. **性能考虑**: 对于大量数据，考虑分块处理

## 🔧 未来优化方向

- [ ] 支持更多加密模式（GCM、CTR等）
- [ ] 添加密钥派生功能（PBKDF2）
- [ ] 支持大文件流式加密
- [ ] 添加数字签名功能
- [ ] 支持更多密钥格式