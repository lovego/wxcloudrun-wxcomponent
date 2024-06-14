package custom

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// DecryptMsg 消息解密
func DecryptMsg(appID, encryptedMsg, aesKey string) (random, rawMsgXMLBytes []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic error: err=%v", e)
			return
		}
	}()
	var encryptedMsgBytes, key, getAppIDBytes []byte
	encryptedMsgBytes, err = base64.StdEncoding.DecodeString(encryptedMsg)
	if err != nil {
		return
	}
	key, err = aesKeyDecode(aesKey)
	if err != nil {
		panic(err)
	}
	random, rawMsgXMLBytes, getAppIDBytes, err = AESDecryptMsg(encryptedMsgBytes, key)
	if err != nil {
		err = fmt.Errorf("消息解密失败,%v", err)
		return
	}
	if appID != string(getAppIDBytes) {
		err = fmt.Errorf("消息解密校验APPID失败")
		return
	}
	return
}

func aesKeyDecode(encodedAESKey string) (key []byte, err error) {
	if len(encodedAESKey) != 43 {
		err = fmt.Errorf("the length of encodedAESKey must be equal to 43")
		return
	}
	key, err = base64.StdEncoding.DecodeString(encodedAESKey + "=")
	if err != nil {
		return
	}
	if len(key) != 32 {
		err = fmt.Errorf("encodingAESKey invalid")
		return
	}
	return
}

func AESDecryptMsg(ciphertext []byte, aesKey []byte) (random, rawXMLMsg, appID []byte, err error) {
	const (
		BlockSize = 32
		BlockMask = BlockSize - 1
	)

	if len(ciphertext) < BlockSize {
		err = fmt.Errorf("the length of ciphertext too short: %d", len(ciphertext))
		return
	}
	if len(ciphertext)&BlockMask != 0 {
		err = fmt.Errorf("ciphertext is not a multiple of the block size, the length is %d", len(ciphertext))
		return
	}

	plaintext := make([]byte, len(ciphertext))

	// 解密
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, aesKey[:16])
	mode.CryptBlocks(plaintext, ciphertext)

	// PKCS#7 去除补位
	amountToPad := int(plaintext[len(plaintext)-1])
	if amountToPad < 1 || amountToPad > BlockSize {
		err = fmt.Errorf("the amount to pad is incorrect: %d", amountToPad)
		return
	}
	plaintext = plaintext[:len(plaintext)-amountToPad]

	// 反拼接
	// len(plaintext) == 16+4+len(rawXMLMsg)+len(appId)
	if len(plaintext) <= 20 {
		err = fmt.Errorf("plaintext too short, the length is %d", len(plaintext))
		return
	}
	rawXMLMsgLen := int(decodeNetworkByteOrder(plaintext[16:20]))
	if rawXMLMsgLen < 0 {
		err = fmt.Errorf("incorrect msg length: %d", rawXMLMsgLen)
		return
	}
	appIDOffset := 20 + rawXMLMsgLen
	if len(plaintext) <= appIDOffset {
		err = fmt.Errorf("msg length too large: %d", rawXMLMsgLen)
		return
	}

	random = plaintext[:16:20]
	rawXMLMsg = plaintext[20:appIDOffset:appIDOffset]
	appID = plaintext[appIDOffset:]
	return
}

func decodeNetworkByteOrder(orderBytes []byte) (n uint32) {
	return uint32(orderBytes[0])<<24 |
		uint32(orderBytes[1])<<16 |
		uint32(orderBytes[2])<<8 |
		uint32(orderBytes[3])
}
