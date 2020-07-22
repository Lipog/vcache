package onehash

import (
	"log"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := NewOneHash(3, nil)

	//for i := 0; i < 9; i++ {
	//	value := hash.hash([]byte(strconv.Itoa(i) + "adfc"))
	//	log.Println(value)
	//}
	value := hash.hash([]byte("0"))
	log.Println("0", value)
	value = hash.hash([]byte("1"))
	log.Println("1", value)
	value = hash.hash([]byte("http://119.3.101.129"))
	log.Println("http://119.3.101.129", value)
	value = hash.hash([]byte("Tom"))
	log.Println("Tom", value)
	value = hash.hash([]byte("Jack"))
	log.Println("Jack", value)
	value = hash.hash([]byte("Lip"))
	log.Println("Lip", value)
}
