package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
)

func generateKeyPair() error {
	// 生成密钥对
	privateKey, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	// 将私钥转换为字节切片
	privateKeyBytes := privateKey.D.Bytes()
	// 将私钥字节切片转换为十六进制字符串
	privateKeyStr := hex.EncodeToString(privateKeyBytes)
	// 获取公钥
	publicKey := privateKey.Public()
	// 将公钥转换为字节切片
	publicKeyBytes := publicKey.(*sm2.PublicKey).X.Bytes()
	// 将公钥字节切片转换为十六进制字符串
	publicKeyStr := hex.EncodeToString(publicKeyBytes)
	// 输出私钥和公钥
	fmt.Println("Private Key:", privateKeyStr)
	fmt.Println("Public Key:", publicKeyStr)
	return nil
}

func main() {
	err := generateKeyPair()
	if err != nil {
		fmt.Println("Error generating key pair:", err)
	}
}
