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
 -

*/
/*참고 사이트

RFC 2409: The Internet Key Exchange (IKE) - https://www.rfc-editor.org/rfc/rfc2409.html
RFC 3526:  - https://www.rfc-editor.org/rfc/rfc3526.html

https://pkg.go.dev/github.com/monnand/dhkx#section-readme
https://github.com/monnand/dhkx

*/

package main

import (
	"log"
	"runtime"
	"sync"

	"github.com/monnand/dhkx"
)

type pk struct {
	bob   chan []byte
	alice chan []byte
}

var pubKey = &pk{}
var wg sync.WaitGroup

func main() {
	//Processor의 모든 코어 사용
	runtime.GOMAXPROCS(runtime.NumCPU())

	wg = sync.WaitGroup{}
	pkCh := make(chan *pk)
	//bobCh := make(chan *pk)
	//aliceCh := make(chan *pk)

	wg.Add(2)
	go AliceSide(pkCh)
	go BobSide(pkCh)
	wg.Wait()
	close(pkCh)

}

func AliceSide(pkCh chan *pk) {
	for publicKey := range pkCh {
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
		// Bob으로부터 공개키를 받는다.
		b := publicKey.Recv("Bob")

		// Recover Bob's public key
		bobPubKey := dhkx.NewPublicKey(b)

		// Compute the key
		k, _ := g.ComputeKey(bobPubKey, priv)

		// Get the key in the form of []byte
		key := k.Bytes()

		log.Println("Alice side : ", key)
		if key != nil {
			wg.Done()
			return
		}
	}

}

func BobSide(pkCh chan *pk) {
	for publicKey := range pkCh {

		// Get a group. Use the default one would be enough.
		g, _ := dhkx.GetGroup(0)

		// Generate a private key from the group.
		// Use the default random number generator.
		priv, _ := g.GeneratePrivateKey(nil)

		// Get the public key from the private key.
		pub := priv.Bytes()

		// Receive a slice of bytes from Alice, which contains Alice's public key
		a := publicKey.Recv("Alice")

		// Send the public key to Alice.
		publicKey.Send("Alice", pub)

		// Recover Alice's public key
		alicePubKey := dhkx.NewPublicKey(a)

		// Compute the key
		k, _ := g.ComputeKey(alicePubKey, priv)

		// Get the key in the form of []byte
		key := k.Bytes()

		log.Println("Bob side : ", key)
		if key != nil {
			wg.Done()
			return
		}
	}

}

func (*pk) Send(name string, key []byte) {
	//송신자 이름이 alice 경우 bob에게 키를 보냄

	if name == "Alice" {
		pubKey.bob <- key
	} else if name == "Bob" {
		pubKey.alice <- key
	}

	log.Println(name, "에게 보내는 공개키 : ", key)
}

func (*pk) Recv(name string) []byte {
	var key []byte
	//상대 이름이 alice인 경우
	if name == "Alice" {
		key = <-pubKey.alice
	} else if name == "Bob" {
		key = <-pubKey.bob
	}
	//상대 이름이 Bob인 경우
	log.Println(name, "에게 받는 공개키 : ", key)

	return key
}
