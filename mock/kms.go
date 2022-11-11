package mock

import (
	"github.com/aws/aws-sdk-go/service/kms"
	ops_aws_cli "ops-aws-cli"
)

var _ ops_aws_cli.Decrypter = (*MockKMS)(nil)

type MockKMS struct {
	DecodeDataFn func(data string) ([]byte, error)
}

func (m *MockKMS) DecodeData(data string) ([]byte, error) {
	return m.DecodeDataFn(data)
}

type KmsAPI struct {
	DecryptFn func(input *kms.DecryptInput) (*kms.DecryptOutput, error)
}

func (k *KmsAPI) Decrypt(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	return k.DecryptFn(input)
}
