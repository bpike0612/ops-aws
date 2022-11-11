package kms

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type kmsAPI interface {
	Decrypt(input *kms.DecryptInput) (*kms.DecryptOutput, error)
}

type awsKMS struct {
	client kmsAPI
}

func NewKmsClient(client kmsAPI) *awsKMS {
	if client != nil {
		return &awsKMS{client: client}
	}
	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials //todo pull from env
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		panic(err)
	}

	// Create KMS service client
	svc := kms.New(sess)
	return &awsKMS{client: svc}
}

func (k *awsKMS) DecodeData(data string) ([]byte, error) {
	blob, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Printf("issue decoding, %v\n", err)
		return nil, err
	}
	input := kms.DecryptInput{
		CiphertextBlob: blob,
	}
	output, err := k.client.Decrypt(&input)
	if err != nil {
		fmt.Printf("unable to decrypt blob, %v\n", err)
		return nil, err
	}
	return output.Plaintext, nil
}
