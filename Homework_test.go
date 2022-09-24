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

	"log"
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
	//src := []byte("1234567890987655") //test값, 16바이트
	key := []byte("1234567890123456")

	src := make([]byte, 16)             //uint src[16]; 선언
	randomGenerater16(src, "plainText") //바이트 배열 16바이트로 난수 생성

	//iv := make([]byte, 16)
	//randomGenerater16(iv, "Init vector") // iv에 임의 16바이트 값 부여

	log.Println(string(src)) // 암호화하기 전의 평문 출력

	openssl.AesECB(key, src)
	//openssl.AesCBC(key, src, iv)
	//builtin.AesECB(key, src)
}
