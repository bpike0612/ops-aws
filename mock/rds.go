package mock

import (
	"github.com/aws/aws-sdk-go/service/rds"
	opsaws "github.com/bpike0612/ops-aws"
)

// This ensures the MockDb type implements the rdsAPI interface via a compiler check,
// even if it is not used elsewhere. You can read more on this pattern on the effective go site.
// https://golang.org/doc/effective_go#blank_implements
var _ opsaws.RDS = (*MockDb)(nil)

type MockDb struct {
	DescribeAllEndpointsFn func(epType string) ([]byte, error)
	ListDBIdentifiersFn    func() ([]byte, error)
	DescribeEndpointFn     func(clusterName string, epType string) ([]byte, error)
}

func (mdb *MockDb) DescribeAllEndpoints(epType string) ([]byte, error) {
	return mdb.DescribeAllEndpointsFn(epType)
}

func (mdb *MockDb) ListDBIdentifiers() ([]byte, error) {
	return mdb.ListDBIdentifiersFn()
}

func (mdb *MockDb) DescribeEndpoint(clusterName string, epType string) ([]byte, error) {
	return mdb.DescribeEndpointFn(clusterName, epType)
}

type RdsAPI struct {
	DescribeDBClusterEndpointsFn func(input *rds.DescribeDBClusterEndpointsInput) (*rds.DescribeDBClusterEndpointsOutput, error)
	DescribeDBInstancesFn        func(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error)
}

func (rds *RdsAPI) DescribeDBClusterEndpoints(input *rds.DescribeDBClusterEndpointsInput) (*rds.DescribeDBClusterEndpointsOutput, error) {
	return rds.DescribeDBClusterEndpointsFn(input)
}

func (rds *RdsAPI) DescribeDBInstances(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
	return rds.DescribeDBInstancesFn(input)
}
