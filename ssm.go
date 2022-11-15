package opsaws

type SSM interface {
	GetParameters() error
	GetParameterByPath() error
}
