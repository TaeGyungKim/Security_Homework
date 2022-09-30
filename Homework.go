/*
정보보안 과제
SHA-512를 Github에서 찾아서 실행하고 결과를 보이시오
SHA-512를 Github에서 찾아서 실행하고, 음악이나 비디오 파일을 이용하여 메세지다이제스트를 만드는 과정을 보이고 설명하시오.

2018108254 김태경
2022-09-28 생성
 - sha512 사용
 - 오디오 깃허브 패키지 추가

2022-09-30 수정
 - 오디오 깃허브 패키지 제거
 - fileRead 함수 생성
 - Path 추가

*/
/*
참고 사이트
https://github.com/avelino/awesome-go
https://github.com/hajimehoshi/oto - 오디오 패키지 (사용x)
https://pkg.go.dev/crypto/sha512

https://lostark.game.onstove.com/OST - audio
https://store.steampowered.com/app/1057090/Ori_and_the_Will_of_the_Wisps/ - video
*/

package main

import (
	"crypto/sha512" //내장함수 crypto/sha512 패키지를 이용해 sha512 실행
	"encoding/hex"
	"log"
	"os"
)

func main() {

	hash := sha512.New() //해시 인스턴스 생성

	AudioPath := "vol1_02_Leonhart.wav"
	hash.Write(fileRead(AudioPath)) //음악 파일 입력

	VideoPath := "OriAndTheWilloftheWisps 2021-11-27.mp4"
	hash.Write(fileRead(VideoPath)) //비디오 파일 입력

	result := hash.Sum(nil)                 // write한 해시값 추출
	log.Println(hex.EncodeToString(result)) //][]byte(char)를 16진수로 출력

}

//파일 경로 얻어와서 파일 읽어내는 함수
func fileRead(filePath string) []byte {
	// Read the file into memory
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		panic("reading file failed: " + err.Error())
	}

	return fileBytes
}
