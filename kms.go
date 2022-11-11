package ops_aws

type Decrypter interface {
	DecodeData(data string) ([]byte, error)
}
