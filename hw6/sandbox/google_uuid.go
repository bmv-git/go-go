package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
)

/*
Different versions of UUIDs are based on different information:

Version 1 UUIDs are generated from a time and the MAC address.
Version 2 UUIDs are generated from an identifier (usually user ID), time, and the MAC address.
Version 3 and 5 UUIDs are generated based on hashing a namespace identifier and name.
Version 4 UUIDs are based on a randomly generated number.

Version 3 or 5 UUID gives reproducible results - the same ID for a given string.
Because Version 5 uses the SHA-1 hashing algorithm, it is generally more secure
and recommended than Version 3 which uses MD5.
*/
func main() {
	// version 1 uuid
	v1, err := uuid.NewUUID()
	if err != nil {
		log.Fatal("cannot generate v1 uuid")
	}
	fmt.Printf("v1 uuid: %v\tType: %T\n", v1, v1)

	// version 2 uuid
	v2, err := uuid.NewDCEGroup()
	if err != nil {
		log.Fatal("cannot generate v2 uuid")
	}
	fmt.Printf("v2 uuid: %v\tType: %T\n", v2, v2)

	// version 3 uuid (MD5)
	v3 := uuid.NewMD5(uuid.NameSpaceURL, []byte("https://example.com"))
	fmt.Printf("v3 uuid: %v\tType: %T\n", v3, v3)

	// version 4 uuid (Random)
	v4, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("cannot generate v4 uuid")
	}
	fmt.Printf("v4 uuid.NewRandom(): %v\tType: %T\n", v4, v4)

	v4_ := uuid.NewString() // возвращает string или panic (нет ошибки!)
	fmt.Printf("v4 uuid.NewString(): %v\tType: %T\n", v4_, v4_)

	// version 5 uuid (SHA-1)
	v5 := uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://example.com"))
	v5_ := uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://example.com"))
	fmt.Printf("v5 uuid:  \t\t%v\tType: %T\n", v5, v5)
	fmt.Printf("v5 (копия) uuid: \t%v\tType: %T\n", v5_, v5_)
}
