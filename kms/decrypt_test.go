package kms_test

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/kms"
	opsKms "github.com/bpike0612/ops-aws/kms"
	"github.com/bpike0612/ops-aws/mock"
	"github.com/stretchr/testify/assert"
)

func Test_kmsDecryptUnableToDecrypt(t *testing.T) {
	t.Parallel()
	mockClient := mock.KmsAPI{
		DecryptFn: func(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
			return nil, errors.New("mock error")
		},
	}
	mockAPI := opsKms.NewKmsClient(&mockClient)
	assert.NotNil(t, mockAPI)
	got, err := mockAPI.DecodeData("")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "unable to decrypt blob, mock error")
	assert.Nil(t, got)
}

func Test_kmsDecrypt(t *testing.T) {
	t.Parallel()
	data := base64.StdEncoding.EncodeToString([]byte(`this is a test`))
	output := kms.DecryptOutput{Plaintext: []byte(`this is a test`)}
	mockClient := mock.KmsAPI{
		DecryptFn: func(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
			return &output, nil
		},
	}
	mockAPI := opsKms.NewKmsClient(&mockClient)
	assert.NotNil(t, mockAPI)
	got, err := mockAPI.DecodeData(data)
	assert.Nil(t, err)
	assert.Equal(t, got, []byte(`this is a test`))
}
