/*
정보보안 과제
SHA-512를 Github에서 찾아서 실행하고 결과를 보이시오
SHA-512를 Github에서 찾아서 실행하고, 음악이나 비디오 파일을 이용하여 메세지다이제스트를 만드는 과정을 보이고 설명하시오.

2018108254 김태경
2022-09-28 생성
 - sha512 사용
 - 오디오 깃허브 패키지 추가

*/
/*
참고 사이트
https://github.com/avelino/awesome-go
https://github.com/hajimehoshi/oto - 오디오
https://pkg.go.dev/crypto/sha512


*/
package main

import (
	"bytes"
	"crypto/sha512" //내장함수 crypto/sha512 패키지를 이용해 sha512 실행
	"encoding/hex"
	"log"
	"os"

	"github.com/hajimehoshi/go-mp3"
    "github.com/hajimehoshi/oto/v2"
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

	//연결x
	// Read the mp3 file into memory
	fileBytes, err := os.ReadFile("./my-file.mp3")
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	   // Prepare an Oto context (this will use your default audio device) that will
    // play all our sounds. Its configuration can't be changed later.

    // Usually 44100 or 48000. Other values might cause distortions in Oto
    samplingRate := 44100

    // Number of channels (aka locations) to play sounds from. Either 1 or 2.
    // 1 is mono sound, and 2 is stereo (most speakers are stereo). 
    numOfChannels := 2

    // Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
    audioBitDepth := 2

    // Remember that you should **not** create more than one context
    otoCtx, readyChan, err := oto.NewContext(samplingRate, numOfChannels, audioBitDepth)
    if err != nil {
        panic("oto.NewContext failed: " + err.Error())
    }
    // It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
    <-readyChan

    // Create a new 'player' that will handle our sound. Paused by default.
    player := otoCtx.NewPlayer(decodedMp3)
    
    // Play starts playing the sound and returns without waiting for it (Play() is async).
    player.Play()

    // We can wait for the sound to finish playing using something like this
    for player.IsPlaying() {
        time.Sleep(time.Millisecond)
    }

    // Now that the sound finished playing, we can restart from the beginning (or go to any location in the sound) using seek
    // newPos, err := player.(io.Seeker).Seek(0, io.SeekStart)
    // if err != nil{
    //     panic("player.Seek failed: " + err.Error())
    // }
    // println("Player is now at position:", newPos)
    // player.Play()

    // If you don't want the player/sound anymore simply close
    err = player.Close()
    if err != nil {
        panic("player.Close failed: " + err.Error())
    }
}


}
