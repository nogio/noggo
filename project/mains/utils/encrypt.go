package utils

import (
	"encoding/pem"
	"errors"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"encoding/base64"
	"crypto/md5"
	"encoding/hex"
	"os"
	"io"
	"fmt"
	"crypto/sha1"
)




//sha1加密
func Sha1(str string) string {
	md5Ctx := sha1.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}


//sha1加密文件
func Sha1File(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := sha1.New()
		if _, e := io.Copy(h, f); e == nil {
			return fmt.Sprintf("%x", h.Sum(nil))
		}
	}
	return ""
}



//md5加密
func Md5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//md5加密文件
func Md5File(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := md5.New()
		if _, e := io.Copy(h, f); e == nil {
			return fmt.Sprintf("%x", h.Sum(nil))
		}
	}
	return ""
}


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



func Encode64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}
func Decode64(value string) string {
	d, e := base64.StdEncoding.DecodeString(value)
	if e == nil {
		return string(d)
	}
	return value
}









const (
	hash64Table = "01234AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz56789+/"
)

func Enhash64(value string, tables ...string) string {
	table := hash64Table
	if len(tables) > 0 {
		table =tables[0]
	}

	var base64Coder = base64.NewEncoding(table)
	return base64Coder.EncodeToString([]byte(value))
}
func Dehash64(value string, tables ...string) string {
	table := hash64Table
	if len(tables) > 0 {
		table =tables[0]
	}

	var base64Coder = base64.NewEncoding(table)
	d, e := base64Coder.DecodeString(value)
	if e == nil {
		return string(d)
	}
	return value
}




const (
	hashIdAlphabet  = "1ab2cd3ef4hk5mn6rs7tu8vw9xz"
	hashIdAlphabetEasy  = "123456789abcdefghjkmnpqrstuvwxyz"    //此表去除难以识别的字母
	hashIdSalt      = "iys.hashId.Salt"
)
