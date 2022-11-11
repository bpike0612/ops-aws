package database_test

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/bpike0612/ops-aws/database"
	"github.com/bpike0612/ops-aws/mock"
)

func Test_awsRDS_DescribeEndpoint(t *testing.T) {
	mockEndpntRsp := rds.DescribeDBClusterEndpointsOutput{
		DBClusterEndpoints: []*rds.DBClusterEndpoint{
			{
				Endpoint:     aws.String("mock-db-test-1"),
				EndpointType: aws.String(database.WriterEndPointType),
			},
		},
	}
	mockClient := mock.RdsAPI{
		DescribeDBClusterEndpointsFn: func(input *rds.DescribeDBClusterEndpointsInput) (*rds.DescribeDBClusterEndpointsOutput, error) {
			return &mockEndpntRsp, nil
		},
	}
	want, err := json.Marshal(&mockEndpntRsp)
	mockApi := database.NewRds(&mockClient)
	assert.NotNil(t, mockApi)
	got, err := mockApi.DescribeEndpoint("mock-db-cluster", database.WriterEndPointType)
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func Test_awsRDS_UnableToDescribeEndpoint(t *testing.T) {
	mockClient := mock.RdsAPI{
		DescribeDBClusterEndpointsFn: func(input *rds.DescribeDBClusterEndpointsInput) (*rds.DescribeDBClusterEndpointsOutput, error) {
			return nil, errors.New("mock error")
		},
	}
	mockApi := database.NewRds(&mockClient)
	assert.NotNil(t, mockApi)
	got, err := mockApi.DescribeEndpoint("mock-db-cluster", database.WriterEndPointType)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "unable to describe endpoint, mock error")
	assert.Nil(t, got)
}

func Test_awsRDS_UnableToDescribeInstances(t *testing.T) {
	mockClient := mock.RdsAPI{
		DescribeDBInstancesFn: func(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
			return nil, errors.New("mock error")
		},
	}
	mockApi := database.NewRds(&mockClient)
	assert.NotNil(t, mockApi)
	got, err := mockApi.DescribeAllEndpoints(database.WriterEndPointType)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "unable to list instances, mock error")
	assert.Nil(t, got)
}

func Test_awsRDS_DescribeWriterInstances(t *testing.T) {
	mockInstRsp := rds.DescribeDBInstancesOutput{
		DBInstances: []*rds.DBInstance{
			{
				AvailabilityZone:    aws.String("us-east-1c"),
				DBClusterIdentifier: aws.String("mock-db"),
				DBInstanceStatus:    aws.String("available"),
				Engine:              aws.String("aurora-postgresql"),
			},
		},
	}
	mockEndpntRsp := rds.DescribeDBClusterEndpointsOutput{
		DBClusterEndpoints: []*rds.DBClusterEndpoint{
			{
				EndpointType: aws.String(database.WriterEndPointType),
				Endpoint:     aws.String("mock-db-test-1"),
			},
		},
	}
	mockClient := mock.RdsAPI{
		DescribeDBInstancesFn: func(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
			return &mockInstRsp, nil
		},
		DescribeDBClusterEndpointsFn: func(input *rds.DescribeDBClusterEndpointsInput) (*rds.DescribeDBClusterEndpointsOutput, error) {
			return &mockEndpntRsp, nil
		},
	}
	want, err := json.Marshal(&mockEndpntRsp.DBClusterEndpoints)
	mockApi := database.NewRds(&mockClient)
	assert.NotNil(t, mockApi)
	got, err := mockApi.DescribeAllEndpoints(database.WriterEndPointType)
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func Test_awsRDS_UnableToDescribeClusterEndpoints(t *testing.T) {
	mockInstRsp := rds.DescribeDBInstancesOutput{
		DBInstances: []*rds.DBInstance{
			{
				AvailabilityZone:    aws.String("us-east-1c"),
				DBClusterIdentifier: aws.String("mock-db"),
				DBInstanceStatus:    aws.String("available"),
				Engine:              aws.String("aurora-postgresql"),
			},
		},
	}
	mockClient := mock.RdsAPI{
		DescribeDBInstancesFn: func(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
			return &mockInstRsp, nil
		},
		DescribeDBClusterEndpointsFn: func(input *rds.DescribeDBClusterEndpointsInput) (*rds.DescribeDBClusterEndpointsOutput, error) {
			return nil, errors.New("mock error")
		},
	}
	mockApi := database.NewRds(&mockClient)
	assert.NotNil(t, mockApi)
	got, err := mockApi.DescribeAllEndpoints(database.WriterEndPointType)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "unable to describe endpoints, mock error")
	assert.Nil(t, got)
}

func Test_awsRDS_ListDBIdentifiers(t *testing.T) {
	mockCreateTime := time.Date(
		2020, 11, 17, 20, 34, 58, 651387237, time.UTC)
	mockInstRsp := rds.DescribeDBInstancesOutput{
		DBInstances: []*rds.DBInstance{
			{
				AvailabilityZone:     aws.String("us-east-1c"),
				DBClusterIdentifier:  aws.String("mock-db"),
				DBInstanceStatus:     aws.String("available"),
				Engine:               aws.String("aurora-postgresql"),
				DBInstanceIdentifier: aws.String("mock-db-test"),
				InstanceCreateTime:   &mockCreateTime,
			},
		},
	}
	mockClient := mock.RdsAPI{
		DescribeDBInstancesFn: func(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
			return &mockInstRsp, nil
		},
	}
	want := []byte(`[{"identifier":"mock-db-test","createdAt":"2020-11-17T20:34:58.651387237Z"}]`)
	mockApi := database.NewRds(&mockClient)
	assert.NotNil(t, mockApi)
	got, err := mockApi.ListDBIdentifiers()
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func Test_awsRDS_UnableToListDBIdentifiers(t *testing.T) {
	mockClient := mock.RdsAPI{
		DescribeDBInstancesFn: func(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error) {
			return nil, errors.New("mock error")
		},
	}
	mockApi := database.NewRds(&mockClient)
	assert.NotNil(t, mockApi)
	got, err := mockApi.ListDBIdentifiers()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "unable to list instances, mock error")
	assert.Nil(t, got)
}
