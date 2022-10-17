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

// 공유키 채널 구조체 생성
type pk struct {
	bob      chan []byte
	alice    chan []byte
	isActive chan bool
}

var check [2]bool

// pk 구조체로 선언된 pubKey 객체
// 동시성 프로그래밍을 위한 공유키 채널(버퍼) 객체
var pubKey = pk{}

func main() {
	//Processor의 모든 코어 사용
	runtime.GOMAXPROCS(runtime.NumCPU())

	//키 채널(버퍼) 초기화
	initKeyChannel()

	//활성화
	pubKey.isActive <- true
	//isActive가 활성화되어 있다면 순환
	//go 키워드로 동시성 프로그래밍 처리
	for key := range pubKey.isActive {
		log.Println("Activation Status : ", key)
		go pubKey.AliceSide()
		go pubKey.BobSide()
	}

}

// bob's key, alice's key의 채널을 []byte로 버퍼 1개를 할당
// isActive process information communication's activation
func initKeyChannel() {
	pubKey.bob = make(chan []byte, 1)
	pubKey.alice = make(chan []byte, 1)
	pubKey.isActive = make(chan bool, 1)
}

// Alice의 경우
func (publicKey pk) AliceSide() {
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
	publicKey.Send("Bob", pub)

	// Receive a slice of bytes from Bob, which contains Bob's public key
	b := publicKey.Recv("Bob")

	// Recover Bob's public key
	bobPubKey := dhkx.NewPublicKey(b)

	// Compute the key
	k, _ := g.ComputeKey(bobPubKey, priv)

	// Get the key in the form of []byte
	key := k.Bytes()

	//shared key 있으면 check
	if key != nil {
		check[0] = true
	}

	//Alice의 생성 결과 키
	log.Printf("Alice side Result: %x\n", key)

	//check가 모두 활성화시 함수 간 통신의 종료를 알린다.
	if check[0] && check[1] {
		close(pubKey.isActive)
	}

}

func (publicKey pk) BobSide() {
	// Get a group. Use the default one would be enough.
	g, _ := dhkx.GetGroup(0)

	// Generate a private key from the group.
	// Use the default random number generator.
	priv, _ := g.GeneratePrivateKey(nil)

	// Get the public key from the private key.
	pub := priv.Bytes()

	// Receive a slice of bytes from Alice, which contains Alice's public key
	a := publicKey.Recv("Alice") //, pkCh

	// Send the public key to Alice.
	publicKey.Send("Alice", pub) //, pkCh

	// Recover Alice's public key
	alicePubKey := dhkx.NewPublicKey(a)

	// Compute the key
	k, _ := g.ComputeKey(alicePubKey, priv)

	// Get the key in the form of []byte
	key := k.Bytes()

	//shared key 있으면 check
	if key != nil {
		check[1] = true
	}

	//Bob의 생성 결과 키
	log.Printf("Bob side Result: %x\n", key)
}

// 키를 송신하는 함수
func (pk) Send(name string, key []byte) {
	//송신자 이름이 alice 경우 bob 채널 버퍼에 키를 보냄
	if name == "Alice" {
		pubKey.bob <- key
	} else if name == "Bob" {
		pubKey.alice <- key
	}

	log.Printf("Bob's key : %x, Alice's key : %x \n", pubKey.bob, pubKey.alice)
	log.Printf("Send %s'key: %x\n ", name, key)
}

// key를 수신하는 함수
func (pk) Recv(name string) []byte {
	var key []byte

	//수신자가 alice인 경우
	if name == "Alice" {
		key = <-pubKey.alice
		close(pubKey.alice)
		//상대 이름이 Bob인 경우
	} else if name == "Bob" {
		key = <-pubKey.bob
		close(pubKey.bob)
	}

	log.Printf("Recv %s'key: %x\n ", name, key)
	return key
}
