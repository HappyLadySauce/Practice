package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"net"
	"network/tls/cryption"
	"network/tls/common"
	"time"
)

// 初始化RSA密钥对
var (
	publicKey  []byte
)

func init() {
	// 生成RSA密钥对用于密钥交换
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	common.DealError(err)
	
	// 将公钥编码为PEM格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	common.DealError(err)
	
	publicKey = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
}

func NewTcpClient(socketAddr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", socketAddr)
	common.DealError(err)
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	common.DealError(err)
	log.Printf("connect success for server %v\n", conn.RemoteAddr().String())
	conn.SetDeadline(time.Now().Add(10 * time.Second))	// 设置连接超时时间
	defer conn.Close()

	// 接收服务器的公钥
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	common.DealError(err)
	serverPublicKey := buffer[:n]

	// 生成AES密钥
	aeskey := []byte("1234567890123456")
	// 使用服务器的公钥加密AES密钥
	encryptedAesKey, err := cryption.RsaEncrypt(aeskey, serverPublicKey)
	common.DealError(err)
	_, err = conn.Write([]byte(encryptedAesKey))
	common.DealError(err)

	// 接收服务器确认
	conn.Read(buffer)

	// 后续通信使用AES加密
	plain := "你好！Server"
	encryptedMsg, err := cryption.AesEncrypt(plain, []byte(aeskey))
	common.DealError(err)
	_, err = conn.Write([]byte(encryptedMsg))
	common.DealError(err)
}
