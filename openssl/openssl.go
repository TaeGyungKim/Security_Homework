/*
고(Go language)에서 사용할 수 있도록 openssl에서 지원하는 패키지를 이용
openssl에서는 AES, DES뿐만이 아니라 다양한 암호화 알고리즘을 지원한다.

Security_Homework의 하위 패키지인 Openssl로서
고언어에서 지원하는 기능인 import를 이용하여 구현하였다.


2018108254 김태경
2022-09-23
*/

/*참고 사이트
https://pkg.go.dev/github.com/forgoer/openssl#section-readme
https://pkg.go.dev/github.com/forgoer/openssl

*/
package Openssl

import (
	"encoding/base64"
	"fmt"

	"github.com/forgoer/openssl" //openssl에서 지원하는 패키지
)

//AES에서의 ECD(전자코드북) 모드를 사용하는 함수
func AesECB(key, src []byte) {
	//src := []byte("123456") //ex

	dst, _ := openssl.AesECBEncrypt(src, key, openssl.PKCS7_PADDING)
	fmt.Printf(base64.StdEncoding.EncodeToString(dst)) // yXVUkR45PFz0UfpbDB8/ew==

	dst, _ = openssl.AesECBDecrypt(dst, key, openssl.PKCS7_PADDING)
	fmt.Println(string(dst)) // 123456
}

//AES에서의 CBC(암호 블록 체인)모드를 사용하는 함수
func AesCBC(key, src, iv []byte) {
	//src := []byte("123456") //ex
	//iv := []byte("1234567890123456") //ex

	dst, _ := openssl.AesCBCEncrypt(src, key, iv, openssl.PKCS7_PADDING)
	fmt.Println(base64.StdEncoding.EncodeToString(dst)) // 1jdzWuniG6UMtoa3T6uNLA==

	dst, _ = openssl.AesCBCDecrypt(dst, key, iv, openssl.PKCS7_PADDING)
	fmt.Println(string(dst)) // 123456

}
