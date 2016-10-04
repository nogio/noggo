package utils

import (
	"encoding/base64"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"errors"
)



//加密, 是用的公钥
func RsaEncrypt(data, key []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// 解密, 用的是私钥
func RsaDecrypt(data, key []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, data)
}



//rsa加密
func EncodeRsa(data,key string) (string) {
	result, err := RsaEncrypt([]byte(data), []byte(key))
	if err == nil {
		//可以直接使用base64。 这样省一次类型转换
		return Encode64(string(result))
	}
	return ""
}

//rsa解密
func DecodeRsa(data,key string) (string) {
	//先反转
	datas, err := base64.StdEncoding.DecodeString(data)
	if err == nil {
		result, err := RsaDecrypt(datas, []byte(key))
		if err == nil {
			return string(result)
		}
	}

	return ""
}

