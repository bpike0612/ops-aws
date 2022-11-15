package ssm

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type ssmAPI interface {
	DescribeParameters(input *ssm.DescribeParametersInput) (*ssm.DescribeParametersOutput, error)
	GetParametersByPath(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error)
}

type awsSSM struct {
	client ssmAPI
	filter []*ssm.ParameterStringFilter
}

// NewSsmClient returns awsSSM.
func NewSsmClient(client ssmAPI) *awsSSM { //nolint:revive
	if client != nil {
		return &awsSSM{client: client}
	}
	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials //todo pull from env
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		panic(err)
	}

	// Create SSM service client
	svc := ssm.New(sess)

	return &awsSSM{client: svc}
}

func (s *awsSSM) Get(filters []*ssm.ParameterStringFilter) error {
	input := ssm.DescribeParametersInput{
		Filters:          nil,
		MaxResults:       aws.Int64(10), //nolint:gomnd
		NextToken:        nil,
		ParameterFilters: filters,
	}
	output, err := s.client.DescribeParameters(&input)
	if err != nil {
		err := fmt.Errorf("unable to desribe parameters, %w", err)

		return err
	}
	fmt.Print(output) //nolint

	return nil
}

func (s *awsSSM) GetParameters() error {
	input := ssm.DescribeParametersInput{
		Filters:          nil,
		MaxResults:       aws.Int64(10), //nolint:gomnd
		NextToken:        nil,
		ParameterFilters: s.filter,
	}
	output, err := s.client.DescribeParameters(&input)
	if err != nil {
		err := fmt.Errorf("unable to describe parameters, %w", err)

		return err
	}
	fmt.Print(output) //nolint

	return nil
}

func (s *awsSSM) WithParameterFilters(key string, values []string) *awsSSM {
	var filter ssm.ParameterStringFilter
	var filters []*ssm.ParameterStringFilter
	filter.SetKey(key)
	filter.SetValues(aws.StringSlice(values))
	filters = append(filters, &filter) //nolint:ineffassign,staticcheck

	return s
}

func With(key string, values []string) []*ssm.ParameterStringFilter {
	var filter ssm.ParameterStringFilter
	var filters []*ssm.ParameterStringFilter
	// filter.SetKey(key)
	// filter.SetValues(aws.StringSlice(values))
	for i := range values {
		filter.SetKey(key)
		filter.SetValues(aws.StringSlice(values))
		filters[i] = &filter
	}
	// filters = append(filters, &filter)

	return filters
}

// GetParameterByPath will return the value for a given path parameter or an error.
// Note: The hierarchy for the parameter. Hierarchies start with a forward slash (/). The hierarchy is the
// parameter name except the last part of the parameter. For the API call to succeed, the last part of the parameter
// name can't be in the path. A parameter name hierarchy can have a maximum of 15 levels.
func (s *awsSSM) GetParameterByPath(path string) ([]byte, error) {
	input := ssm.GetParametersByPathInput{
		MaxResults: aws.Int64(10), //nolint:gomnd
		Path:       aws.String(path),
		Recursive:  aws.Bool(true),
	}
	out, err := s.client.GetParametersByPath(&input)
	if err != nil {
		return nil, err
	}
	if len(out.Parameters) == 0 {
		err := fmt.Sprintf("no results found for path, %s", path)

		return nil, errors.New(err)
	}
	b, err := json.Marshal(&out)
	if err != nil {
		err := fmt.Errorf("error marshalling json, %w", err)

		return nil, err
	}

	return b, nil
}
