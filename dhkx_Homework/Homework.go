/*
정보보안 과제
Diffie-Hellman 알고리즘 구현
Diffie-Hellman 알고리즘을 구현한 소스코드를 찾아, 이를 실행하고 키 공유가 되는 예를 보이시오.
실행하는 화면을 캡쳐해서 실행결과를 보이고, 설명하시오.

2018108254 김태경
2022. 10. 13. 생성
 - 깃허브에서 코드 가져옴
 - Send(), Recv()의 데이터 주고 받는 함수 추가
 - go 키워드로 동시성 프로그래밍 추가
 - 동시성 프로그래밍 제어를 위한 sync.WaitGroup 추가

2022. 10. 15. 수정
 - chan 채널 추가
 - deadlock 발생

2022. 10. 17. 수정
 - chan 추가
 - sync.WaitGroup 제거
 - 반복문의 활성화를 담당하는 pk 구조체에 active 플래그 추가
 - key 생성을 알리는 check 플래그 추가

2022. 10. 18. 수정
 - 구조체와 변수 명칭 수정
 - 주석 수정

*/
/*참고 사이트
RFC 2409: The Internet Key Exchange (IKE) - https://www.rfc-editor.org/rfc/rfc2409.html - diffie-hellman
RFC 3526: https://www.rfc-editor.org/rfc/rfc3526.html - diffie-hellman
https://ko.wikipedia.org/wiki/%EB%94%94%ED%94%BC-%ED%97%AC%EB%A8%BC_%ED%82%A4_%EA%B5%90%ED%99%98 - diffie-hellman 정의
https://pkg.go.dev/github.com/monnand/dhkx#section-readme - 예제설명
https://github.com/monnand/dhkx - 예제 이용

https://etloveguitar.tistory.com/m/40 - 동시성프로그래밍 방식

*/

package main

import (
	"log"
	"runtime"

	"github.com/monnand/dhkx"
)

// 키 교환 채널 구조체 생성
type kx struct {
	bob      chan []byte
	alice    chan []byte
	isActive chan bool
}

var isCheck [2]bool

// kx 구조체로 선언된 pubKey 객체
// 동시성 프로그래밍을 위한 키교환 채널(버퍼) 객체
var keyChange = kx{}

func main() {
	//Processor의 모든 코어 사용
	runtime.GOMAXPROCS(runtime.NumCPU())

	//키 채널(버퍼) 초기화
	initKeyChannel()

	//활성화
	keyChange.isActive <- true
	//isActive가 활성화되어 있다면 순환
	//go 키워드로 동시성 프로그래밍 처리
	for key := range keyChange.isActive {
		log.Println("Activation Status : ", key)
		go keyChange.AliceSide()
		go keyChange.BobSide()
	}
}

// bob's key, alice's key의 채널을 []byte로 버퍼 1개를 할당
// isActive process information communication's activation
func initKeyChannel() {
	keyChange.bob = make(chan []byte, 1)
	keyChange.alice = make(chan []byte, 1)
	keyChange.isActive = make(chan bool, 1)
}

// Alice의 경우의 키교환
func (keyMessage kx) AliceSide() {
	// Get a group. Use the default one would be enough.
	//RFC2409, RFC3526에 정의된 ID로 DHGroup 가져옴
	g, err := dhkx.GetGroup(0)
	if err != nil {
		log.Panic("dont get group")
	}

	// Generate a private key from the group.
	// Use the default random number generator.
	// nil값으로 할당시 random 처리
	// 비밀키 생성
	priv, err := g.GeneratePrivateKey(nil)
	if err != nil {
		log.Panic("dont create private key")
	}

	// 비밀키로부터 공개키를 얻는다.
	pub := priv.Bytes()

	// Bob에게 공개키 보냄.
	keyMessage.Send("Bob", pub)

	// Receive a slice of bytes from Bob, which contains Bob's public key
	b := keyMessage.Recv("Bob")

	// Recover Bob's public key
	bobPubKey := dhkx.NewPublicKey(b)

	// Compute the key
	k, _ := g.ComputeKey(bobPubKey, priv)

	// Get the key in the form of []byte
	key := k.Bytes()

	//shared key 있으면 check
	if key != nil {
		isCheck[0] = true
	}

	//Alice의 생성 결과 키(shared key)
	log.Printf("Alice side Result: %x\n", key)

	//check가 모두 활성화시 함수 간 통신의 종료를 알린다.
	if isCheck[0] && isCheck[1] {
		close(keyChange.isActive)
	}
}

// Bob의 경우의 키 교환
func (keyMessage kx) BobSide() {
	// Get a group. Use the default one would be enough.
	g, _ := dhkx.GetGroup(0)

	// Generate a private key from the group.
	// Use the default random number generator.
	priv, _ := g.GeneratePrivateKey(nil)

	// Get the public key from the private key.
	pub := priv.Bytes()

	// Receive a slice of bytes from Alice, which contains Alice's public key
	a := keyMessage.Recv("Alice")

	// Send the public key to Alice.
	keyMessage.Send("Alice", pub)

	// Recover Alice's public key
	alicePubKey := dhkx.NewPublicKey(a)

	// Compute the key
	k, _ := g.ComputeKey(alicePubKey, priv)

	// Get the key in the form of []byte
	key := k.Bytes()

	//shared key 있으면 check
	if key != nil {
		isCheck[1] = true
	}

	//Bob의 생성 결과 키 (shared key)
	log.Printf("Bob side Result: %x\n", key)
}

// 키를 송신하는 함수
func (kx) Send(name string, key []byte) {
	//송신자 이름이 alice 경우 bob 채널 버퍼에 키를 보냄
	if name == "Alice" {
		keyChange.bob <- key
	} else if name == "Bob" {
		keyChange.alice <- key
	}

	//log.Printf("Bob's key : %x, Alice's key : %x \n", keyChange.bob, keyChange.alice)
	log.Printf("Send %s'key: %x\n ", name, key)
}

// key를 수신하는 함수
func (kx) Recv(name string) []byte {
	var key []byte

	//수신자가 alice인 경우 키 버퍼에서 키를 받아오고 버퍼를 닫는다.
	if name == "Alice" {
		key = <-keyChange.alice
		close(keyChange.alice)
		//상대 이름이 Bob인 경우
	} else if name == "Bob" {
		key = <-keyChange.bob
		close(keyChange.bob)
	}

	log.Printf("Recv %s'key: %x\n ", name, key)
	return key
}
