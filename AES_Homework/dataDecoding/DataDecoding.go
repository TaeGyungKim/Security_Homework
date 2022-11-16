package data

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
)

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
