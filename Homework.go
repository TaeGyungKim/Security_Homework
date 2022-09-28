/*
정보보안 과제
SHA-512를 Github에서 찾아서 실행하고 결과를 보이시오
SHA-512를 Github에서 찾아서 실행하고, 음악이나 비디오 파일을 이용하여 메세지다이제스트를 만드는 과정을 보이고 설명하시오.

2018108254 김태경
2022-09-28 생성
 - sha512 사용

*/
/*
참고 사이트

https://pkg.go.dev/crypto/sha512


*/
package main

import (
	"crypto/sha512" //내장함수 crypto/sha512 패키지를 이용해 sha512 실행
	"encoding/hex"
	"log"
)

func main() {
	txt := "test"
	//
	hash := sha512.New() //해시 인스턴스 생성

	hash.Write([]byte(txt)) //데이터 입력

	result1 := hash.Sum(nil)              // 해시값 추출 []byte
	result2 := sha512.Sum512([]byte(txt)) // [64]byte
	log.Println(hex.EncodeToString(result1))
	log.Println(hex.EncodeToString(result2[:]))
}
