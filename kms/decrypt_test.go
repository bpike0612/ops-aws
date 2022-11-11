package kms_test

import (
	"encoding/base64"
	"errors"
	"github.com/aws/aws-sdk-go/service/kms"
	opsKms "github.com/bpike0612/ops-aws/kms"
	"github.com/bpike0612/ops-aws/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_kmsDecryptUnableToDecrypt(t *testing.T) {
	mockClient := mock.KmsAPI{
		DecryptFn: func(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
			return nil, errors.New("mock error")
		},
	}
	mockApi := opsKms.NewKmsClient(&mockClient)
	assert.NotNil(t, mockApi)
	got, err := mockApi.DecodeData("")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "mock error")
	assert.Nil(t, got)
}

func Test_kmsDecrypt(t *testing.T) {
	data := base64.StdEncoding.EncodeToString([]byte(`this is a test`))
	output := kms.DecryptOutput{Plaintext: []byte(`this is a test`)}
	mockClient := mock.KmsAPI{
		DecryptFn: func(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
			return &output, nil
		},
	}
	mockApi := opsKms.NewKmsClient(&mockClient)
	assert.NotNil(t, mockApi)
	got, err := mockApi.DecodeData(data)
	assert.Nil(t, err)
	assert.Equal(t, got, []byte(`this is a test`))
}
