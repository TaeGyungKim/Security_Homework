package main

import (
	//"github.com/stretchr/testify/assert"
	openssl "Security_Homework/Openssl" //openssl를 이용하는 패키지
	builtin "Security_Homework/builtin" //내장 함수를 이용하는 패키지

	"testing"
)

/*
func Test(t *testing.T) {
	assert := assert.New(t)

}
*/
func Benchmark_Image(b *testing.B) {
	key := []byte("1234567890123456")
	img := imageFile("lena.png")

	//builtin.AesECB(key, img)
	openssl.AesECB(key, img)
}

func Benchmark_Text(b *testing.B) {
	key := []byte("1234567890123456")
	src := []byte("1234567890987655") //test값, 16바이트
	builtin.AesECB(key, src)
	//openssl.AesECB(key, src)
}
