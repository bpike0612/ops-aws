package mock

import (
	"github.com/aws/aws-sdk-go/service/ssm"
	opsaws "github.com/bpike0612/ops-aws"
)

var _ opsaws.SSM = (*SSM)(nil)

type SSM struct {
	GetFn                func(filters []*ssm.ParameterStringFilter) error
	GetParametersFn      func() error
	GetParameterByPathFn func(path string) ([]byte, error)
}

func (s *SSM) Get(filters []*ssm.ParameterStringFilter) error {
	return s.GetFn(filters)
}

func (s *SSM) GetParameters() error {
	return s.GetParametersFn()
}

func (s *SSM) GetParameterByPath(path string) ([]byte, error) {
	return s.GetParameterByPathFn(path)
}
