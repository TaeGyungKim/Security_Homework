/*
고 언어(Go Language)에서 지원하는 테스트 코드


2018108254 김태경
2022-09-23 생성
2022-09-24 수정
 - 키값을 난수로 생성하여 벤치마크 테스트하는 기능 추가
*/

package main

import (
	openssl "Security_Homework/Openssl" //openssl를 이용하는 패키지
	//builtin "Security_Homework/builtin" //내장 함수를 이용하는 패키지

	"crypto/rand"
	"testing"
	//"github.com/stretchr/testify/assert"
)

/*
func Test(t *testing.T) {
	assert := assert.New(t)

}
*/
/*
func Benchmark_Image(b *testing.B) {
	key := []byte("1234567890123456")
	img := imageFile("lena.png")

	//builtin.AesECB(key, img)
	openssl.AesECB(key, img)
}*/

func Benchmark_Text(b *testing.B) {
	key := make([]byte, 16)
	_, err := rand.Read(key) //임의의 16바이트값 부여
	if err != nil {
		return
	}

	src := make([]byte, 16)
	_, err = rand.Read(src) //임의의 16바이트값 부여
	if err != nil {
		return
	}

	//src := []byte("1234567890987655") //test값, 16바이트
	//builtin.AesECB(key, src)
	openssl.AesECB(key, src)
}
