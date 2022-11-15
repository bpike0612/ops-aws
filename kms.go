package opsaws

type Decrypter interface {
	DecodeData(data string) ([]byte, error)
}
