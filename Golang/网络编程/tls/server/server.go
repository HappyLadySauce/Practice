package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"log"
	"net"
	"network/tls/cryption"
	"network/tls/common"
	"time"
)

// 初始化RSA密钥对
var (
	privateKey *rsa.PrivateKey
	publicKey  []byte
)

func init() {
	// 生成RSA密钥对用于密钥交换
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	common.DealError(err)
	
	// 将公钥编码为PEM格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	common.DealError(err)
	
	publicKey = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
}

func NewTcpServer(socketAddr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", socketAddr)
	common.DealError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	common.DealError(err)
	log.Println("wait connect for client...")
	conn, err := listener.Accept()
	common.DealError(err)
	log.Printf("connect success for client %v\n", conn.RemoteAddr().String())
	conn.SetDeadline(time.Now().Add(10 * time.Second))	// 设置连接超时时间
	defer conn.Close()

	// 首先发送公钥给客户端
	_, err = conn.Write(publicKey)
	common.DealError(err)

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	common.DealError(err)
	// 解密AES密钥
	encryptedAesKey, err := hex.DecodeString(string(buffer[:n]))
	common.DealError(err)
	aeskey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedAesKey)
	common.DealError(err)
	conn.Write([]byte("I received your aes key."))

	// 后续通信使用AES加密
	n, err = conn.Read(buffer)
	common.DealError(err)
	plainText, err := cryption.AesDecrypt(string(buffer[:n]), []byte(aeskey))
	common.DealError(err)
	log.Printf("received client message: %s\n", plainText)
}



