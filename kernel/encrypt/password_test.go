package encrypt_test

import (
	"encoding/base64"
	"fmt"

	k "github.com/aileron-gateway/aileron-gateway/apis/kernel"
	"github.com/aileron-gateway/aileron-gateway/kernel/encrypt"
)

func ExampleBCrypt_Hash() {
	password := []byte("password")

	spec := &k.BCryptSpec{
		Cost: 10,
	}

	crypt, err := encrypt.NewBCrypt(spec)
	if err != nil {
		panic(err)
	}
	ciphertext, err := crypt.Hash(password)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(string(ciphertext))
	// Example Output:
	// $2a$10$rwnLBPvWrS95E7ERqayoI.41ve0xJrCA0C3bxoGBFgCWHUB34ZJhO
}

func ExampleBCrypt_Compare() {
	password := []byte("password")

	spec := &k.BCryptSpec{
		Cost: 10,
	}

	crypt, err := encrypt.NewBCrypt(spec)
	if err != nil {
		panic(err)
	}
	ciphertext, err := crypt.Hash(password)
	if err != nil {
		panic(err)
	}

	err = crypt.Compare(ciphertext, password)
	if err != nil {
		panic(err)
	}

	fmt.Println(err == nil)
	// Output:
	// true
}

func ExampleSCrypt_Hash() {
	password := []byte("password")

	spec := &k.SCryptSpec{
		SaltLen: 10,
		N:       32768,
		R:       8,
		P:       1,
		KeyLen:  32,
	}

	crypt, err := encrypt.NewSCrypt(spec)
	if err != nil {
		panic(err)
	}
	ciphertext, err := crypt.Hash(password)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(base64.StdEncoding.EncodeToString(ciphertext))
	// Example Output:
	// AUZzuD9KYkxXnUoExalWcpGd4L9H7Flxf3dZrpXeEHRsL3ET+cncjvIw
}

func ExampleSCrypt_Compare() {
	password := []byte("password")

	spec := &k.SCryptSpec{
		SaltLen: 10,
		N:       32768,
		R:       8,
		P:       1,
		KeyLen:  32,
	}

	crypt, err := encrypt.NewSCrypt(spec)
	if err != nil {
		panic(err)
	}
	ciphertext, err := crypt.Hash(password)
	if err != nil {
		panic("handle error here")
	}

	err = crypt.Compare(ciphertext, password)
	if err != nil {
		panic("password not matched")
	}

	fmt.Println(err == nil)
	// Output:
	// true
}

func ExamplePBKDF2_Hash() {
	password := []byte("password")

	spec := &k.PBKDF2Spec{
		SaltLen: 10,
		Iter:    4096,
		KeyLen:  32,
		HashAlg: k.HashAlg_SHA256,
	}

	crypt, err := encrypt.NewPBKDF2(spec)
	if err != nil {
		panic(err)
	}
	ciphertext, err := crypt.Hash(password)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(base64.StdEncoding.EncodeToString(ciphertext))
	// Example Output:
	// Gi/yF1Aqfk6dIP13jOcjz8tHL7o13Zdxt6krgtEa0LYcFFE6pS4qvtdG
}

func ExamplePBKDF2_Compare() {
	password := []byte("password")

	spec := &k.PBKDF2Spec{
		SaltLen: 10,
		Iter:    4096,
		KeyLen:  32,
		HashAlg: k.HashAlg_SHA256,
	}

	crypt, err := encrypt.NewPBKDF2(spec)
	if err != nil {
		panic(err)
	}
	ciphertext, err := crypt.Hash(password)
	if err != nil {
		panic("handle error here")
	}

	err = crypt.Compare(ciphertext, password)
	if err != nil {
		panic("password not matched")
	}

	fmt.Println(err == nil)
	// Output:
	// true
}

func ExampleArgon2i_Hash() {
	password := []byte("password")

	spec := &k.Argon2Spec{
		SaltLen: 10,
		Time:    3,
		Memory:  32 * 1024,
		Threads: 4,
		KeyLen:  32,
	}

	crypt, err := encrypt.NewArgon2i(spec)
	if err != nil {
		panic(err)
	}
	ciphertext, err := crypt.Hash(password)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(base64.StdEncoding.EncodeToString(ciphertext))
	// Example Output:
	// aB3NQtOu3qfiML9uhkIN8PFHxpDcQkdfgwtE80keLhH7aSZoU1BwhLmY
}

func ExampleArgon2i_Compare() {
	password := []byte("password")

	spec := &k.Argon2Spec{
		SaltLen: 10,
		Time:    3,
		Memory:  32 * 1024,
		Threads: 4,
		KeyLen:  32,
	}

	crypt, err := encrypt.NewArgon2i(spec)
	if err != nil {
		panic(err)
	}
	ciphertext, err := crypt.Hash(password)
	if err != nil {
		panic("handle error here")
	}

	err = crypt.Compare(ciphertext, password)
	if err != nil {
		panic("password not matched")
	}

	fmt.Println(err == nil)
	// Output:
	// true
}

func ExampleArgon2id_Hash() {
	password := []byte("password")

	spec := &k.Argon2Spec{
		SaltLen: 10,
		Time:    1,
		Memory:  63 * 1024,
		Threads: 4,
		KeyLen:  32,
	}

	crypt, err := encrypt.NewArgon2id(spec)
	if err != nil {
		panic(err)
	}
	ciphertext, err := crypt.Hash(password)
	if err != nil {
		panic("handle error here")
	}

	fmt.Println(base64.StdEncoding.EncodeToString(ciphertext))
	// Example Output:
	// E6/M0U9snWH026tJMIm22mcNqg5otqPoFFL9TV+X69URs0ybGyI8vIb7
}

func ExampleArgon2id_Compare() {
	password := []byte("password")

	spec := &k.Argon2Spec{
		SaltLen: 10,
		Time:    1,
		Memory:  63 * 1024,
		Threads: 4,
		KeyLen:  32,
	}

	crypt, err := encrypt.NewArgon2id(spec)
	if err != nil {
		panic(err)
	}
	ciphertext, err := crypt.Hash(password)
	if err != nil {
		panic("handle error here")
	}

	err = crypt.Compare(ciphertext, password)
	if err != nil {
		panic("password not matched")
	}

	fmt.Println(err == nil)
	// Output:
	// true
}
