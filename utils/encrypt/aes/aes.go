package aes

/*
 [wiki](https://zh.wikipedia.org/wiki/%E5%88%86%E7%BB%84%E5%AF%86%E7%A0%81%E5%B7%A5%E4%BD%9C%E6%A8%A1%E5%BC%8F)
（1）AES有5种加密模式，分别是：
 （a）电子密码本（Electronic Codebook Book，ECB）； 照块密码的块大小被分为数个块，并对每个块进行独立加密, 不能很好的隐藏数据模式
 （b）密码块链接（Cipher Block Chaining ，CBC），如果明文长度不是分组长度16字节的整数倍需要进行填充； CBC是最为常用的工作模式 加密串行，解密并行
 （c）计数器模式（Counter，CTR）；
 （d）密文反馈（Cipher FeedBack，CFB）； 模式类似于反向 CBC
 （e）输出反馈（Output FeedBack，OFB）。 然后将其与明文块进行异或，得到密文
（2）AES是对称分组加密算法，每组长度为128bits，即16字节。
（3）AES秘钥的长度只能是16、24或32字节，分别对应三种AES，即AES-128, AES-192和AES-256，三者的区别是加密的轮数不同；

数据填充：
ECB，CBC：块密码只能对确定长度的数据块进行处理，而消息的长度通常是可变的。因此部分模式（即ECB和CBC）需要最后一块在加密前进行填充。
CFB，OFB和CTR：不需要对长度不为密码块大小整数倍的消息进行特别的处理。因为这些模式是通过对块密码的输出与明文进行异或工作的。

PKCS Pkcs7    #5/7 padding strategy.
ANSI AnsiX923 X.923 padding strategy.
ISO Iso10126 padding strategy.
ISO/IEC Iso97971 9797-1 Padding Method 2.
ZeroPadding Zero padding strategy.
NoPadding: Padding;
*/

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"github.com/cute-angelia/go-utils/utils/encrypt/base"
	"io"
)

type Mode int

const (
	PaddingPkcs7 Mode = iota
	PaddingAnsiX923
	PaddingIso10126
	PaddingIso97971
	PaddingZeroPadding
	PaddingNoPadding
)

func GetAesKey() []byte {
	return []byte("passphrasewhichneedstobe32bytes.")
}

// return base64 string
func EncryptCFB(message []byte, key []byte) (string, error) {
	plainText := message

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//returns to base64 encoded string
	// encmess = base64.URLEncoding.EncodeToString(cipherText)
	return base.Base64Encode(cipherText), nil
}

func DecryptCFB(crypted string, key []byte) ([]byte, error) {
	cipherText, err := base.Base64Decode(crypted)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return nil, err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText, nil
}

// 加密 CBC 数据填充
func EncryptCBC(message, key []byte, mode Mode) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()

	switch mode {
	case PaddingPkcs7:
		message = PKCSPadding(message, blockSize)
	case PaddingZeroPadding:
		message = ZeroPadding(message, blockSize)
	default:
		return "", errors.New("mode类型不匹配")
	}

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	crypted := make([]byte, len(message))
	blockMode.CryptBlocks(crypted, message)
	return base.Base64Encode(crypted), nil
}

// 解密 CBC
func DecryptCBC(crypted string, key []byte, mode Mode) ([]byte, error) {
	cipherText, err := base.Base64Decode(crypted)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//AES分组长度为128位，所以blockSize=16，单位字节
	blockSize := block.BlockSize()
	iv := key[:blockSize]
	//blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) //初始向量的长度必须等于块block的长度16字节
	blockMode := cipher.NewCBCDecrypter(block, iv) //初始向量的长度必须等于块block的长度16字节
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, cipherText)

	switch mode {
	case PaddingPkcs7:
		origData = PKCSUnPadding(origData)
	case PaddingZeroPadding:
		origData = ZeroUnPadding(origData)
	default:
		return origData, errors.New("mode类型不匹配")
	}
	return origData, nil
}
