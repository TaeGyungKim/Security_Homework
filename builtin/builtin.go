/*
고언어(Go language)에서 지원하는 내장패키지인 crypto를 이용하는 패키지


2018108254 김태경
2022-09-23
*/

/*참고 사이트
https://pkg.go.dev/crypto
https://pkg.go.dev/crypto/aes
https://pkg.go.dev/crypto/cipher
https://pkg.go.dev/crypto/cipher#NewCBCDecrypter //예제 이용
https://pkg.go.dev/github.com/golang-module/dongle //코드 이용하지 않음.
*/
package builtin

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
)

//고 언어에 내장된 패키지 aes를 이용하는 함수
//키와 평문을 받아와서 실행
//이 함수는 ECB 방식으로 16바이트씩 잘라서 실행
//모드를 설정하지 않을 경우 패딩 등을 지원하지 않음.
func AesECB(key, src []byte) {
	block, err := aes.NewCipher(key) //AES 대칭키 암호화 블록 생성
	if err != nil {
		log.Fatal(err)
		return
	}

	ciphertext := make([]byte, len(src))
	block.Encrypt(ciphertext, []byte(src)) // 평문을 암호화
	log.Printf("%x\n", ciphertext)

	plaintext := make([]byte, len(src))
	block.Decrypt(plaintext, ciphertext) //암호화문를 평문으로 복호화
	log.Println(string(plaintext))

}

//AES CBC 암호화
func AesCBCEncrypt(key, plaintext []byte) []byte {
	block, err := aes.NewCipher(key) //AES 대칭키 암호화 블록 생성
	if err != nil {
		log.Fatal(err)
		return nil
	}

	if len(plaintext)%aes.BlockSize != 0 { //길이가 16비트단위로 끊어지지 않을 경우
		panic("plaintext is not a multiple of the block size")
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Printf("%x\n", ciphertext)

	return ciphertext
}

//Aes CBC 복호화
func AesCBCDecrypt(key, ciphertext []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
		return
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		log.Fatal(err)
		return
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		log.Fatal(err)
		return
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	log.Printf("%s\n", ciphertext)
}
