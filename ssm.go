package ops_aws

type SSM interface {
	GetParameters() error
	GetParameterByPath() error
}
