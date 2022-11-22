package opsaws

import "github.com/aws/aws-sdk-go/service/ssm"

type SSM interface {
	Get(filters []*ssm.ParameterStringFilter) error
	GetParameters() error
	GetParameterByPath(path string) ([]byte, error)
}
