/*
정보보안 과제
깃허브에서 공개 AES 알고리즘을 찾아서 실행하고  암호화 및 복호화되는 과정 보이시오


2018108254 김태경
2022-09-23
*/

/*참고사이트

https://pkg.go.dev/image
https://go.dev/play/p/WUgHQ3pRla
https://go.dev/blog/image
https://pkg.go.dev/image/png
https://stackoverflow.com/questions/71141811/convert-png-image-to-raw-byte-golang

*/

package main

import (
	//openssl "Security_Homework/Openssl" //openssl를 이용하는 패키지
	builtin "Security_Homework/builtin" //내장 함수를 이용하는 패키지
	"bufio"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
)

func main() {
	key := []byte("1234567890123456") //16바이트
	src := []byte("1234567890987655") //test값, 16바이트
	builtin.AesECB(key, src)
	//openssl.AesECB(key, src)

	//text
	//txt := textFile("test.txt")
	//builtin.AesECB(key, txt)
	//openssl.AesECB(key, txt)

	//image
	//img := imageFile("lena.png")
	//builtin.AesECB(key, img)
	//openssl.AesECB(key, img)

}

//텍스트 디코딩
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

//이미지 디코딩
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

//png 이미지를 raw 파일로 변환시키는 함수
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
