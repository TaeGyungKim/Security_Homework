/*
정보보안 과제
깃허브에서 공개 AES 알고리즘을 찾아서 실행하고  암호화 및 복호화되는 과정 보이시오

2018108254 김태경
2022-09-23 생성
2022-09-24 수정
 - 16byte 난수 생성기 추가
*/
/*참고사이트

쉽게 배우는 소프트웨어 공학

//builtin.go
https://pkg.go.dev/crypto
https://pkg.go.dev/crypto/aes
https://pkg.go.dev/crypto/cipher
https://pkg.go.dev/crypto/cipher#NewCBCDecrypter //예제 이용
https://pkg.go.dev/github.com/golang-module/dongle

//openssl.go
https://pkg.go.dev/github.com/forgoer/openssl#section-readme
https://pkg.go.dev/github.com/forgoer/openssl //예제 이용

//Homework.go
https://pkg.go.dev/image
https://go.dev/play/p/WUgHQ3pRla //예제 이용
https://go.dev/blog/image
https://pkg.go.dev/image/png
https://stackoverflow.com/questions/71141811/convert-png-image-to-raw-byte-golang



https://pkg.go.dev/image
https://go.dev/play/p/WUgHQ3pRla
https://go.dev/blog/image
https://pkg.go.dev/image/png
https://stackoverflow.com/questions/71141811/convert-png-image-to-raw-byte-golang

*/

package main

import (
	openssl "Security_Homework/Openssl" //openssl를 이용하는 패키지
	//builtin "Security_Homework/builtin" //내장 함수를 이용하는 패키지
	"bufio"
	"crypto/rand"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
)

func main() {
	key := []byte("1234567890123456") //16바이트
	src := []byte("1234567890987655") //test값, 16바이트
	//iv := []byte("1234567890123456") //ex

	//key := make([]byte, 16)
	//randomGenerater16(key, "key") //키값에 임의의 16바이트값 부여

	//src := make([]byte, 16)
	//randomGenerater16(src, "plainText") //test값에 임의의 16바이트값 부여

	//iv := make([]byte, 16)
	//randomGenerater16(iv, "Init vector") // iv에 임의 16바이트 값 부여

	log.Println(string(src)) // 암호화하기 전의 평문 출력

	openssl.AesECB(key, src) // openssl에서 지원하는 AES의 ECB 모드
	//builtin.AesECB(key, src)

	//openssl.AesCBC(key, src, iv)     // CBC 모드

	//text 사용하는 경우
	//txt := textFile("test.txt")
	//builtin.AesECB(key, txt)
	//openssl.AesECB(key, txt)
	//builtin.AesECB(key, img)

	//image 사용하는 경우
	//img := imageFile("lena.png")      //이미지(png)를 []byte 형식으로 변환
	//openssl.AesECB(key, img)
}

// 입력값에 임의의 16바이트값 부여하는 함수
func randomGenerater16(value []byte, name string) []byte {
	_, err := rand.Read(value) //랜덤 부여
	if err != nil {
		log.Fatal("error:", err)
		return nil
	}
	log.Println(name, " value : ", value)

	return value
}

// 텍스트 디코딩
func textFile(filePath string) []byte {
	file, err := os.Open(filePath) //Read Only
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()

	//text
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) //1줄씩
	}

	return []byte(scanner.Text())
}

// 이미지 디코딩
func imageFile(filePath string) []byte {
	file, err := os.Open(filePath) //Read Only
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()

	m, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	//bounds := m.Bounds()

	imageData := imageToRGBA(m)
	return imageData
}

// png 이미지를 raw 파일로 변환시키는 함수
func imageToRGBA(img image.Image) []uint8 {
	sz := img.Bounds()
	raw := make([]uint8, (sz.Max.X-sz.Min.X)*(sz.Max.Y-sz.Min.Y)*4)
	idx := 0
	for y := sz.Min.Y; y < sz.Max.Y; y++ {
		for x := sz.Min.X; x < sz.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			raw[idx], raw[idx+1], raw[idx+2], raw[idx+3] = uint8(r), uint8(g), uint8(b), uint8(a)
			idx += 4
		}
	}
	return raw
}
