package mock

import (
	"github.com/aws/aws-sdk-go/service/kms"
	opsaws "github.com/bpike0612/ops-aws"
)

var _ opsaws.Decrypter = (*KMS)(nil)

type KMS struct {
	DecodeDataFn func(data string) ([]byte, error)
}

func (m *KMS) DecodeData(data string) ([]byte, error) {
	return m.DecodeDataFn(data)
}

type KmsAPI struct {
	DecryptFn func(input *kms.DecryptInput) (*kms.DecryptOutput, error)
}

func (k *KmsAPI) Decrypt(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	return k.DecryptFn(input)
}
