package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

const (
	dbClusterEndpointType = "db-cluster-endpoint-type"
	WriterEndPointType    = "writer"
	ReaderEndpointType    = "reader"
)

type rdsAPI interface {
	DescribeDBClusterEndpoints(input *rds.DescribeDBClusterEndpointsInput) (*rds.DescribeDBClusterEndpointsOutput, error)
	DescribeDBInstances(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error)
}

type awsRDS struct {
	client rdsAPI
}

// NewRds accepts an rdsAPI interface and returns an awsRDS.
func NewRds(client rdsAPI) *awsRDS { //nolint:revive
	if client != nil {
		return &awsRDS{client: client}
	}
	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials //todo pull from env
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		panic(err)
	}

	// Create RDS service client
	svc := rds.New(sess)

	return &awsRDS{client: svc}
}

// DescribeEndpoint will return information for a provisioned RDS instance using the cluster name
// and type.
func (r *awsRDS) DescribeEndpoint(clusterName string, epType string) ([]byte, error) {
	filters := []*rds.Filter{
		{
			Name:   aws.String(dbClusterEndpointType),
			Values: aws.StringSlice([]string{epType}),
		},
	}
	output, err := r.client.DescribeDBClusterEndpoints(&rds.DescribeDBClusterEndpointsInput{
		DBClusterIdentifier: aws.String(clusterName),
		Filters:             filters,
	})
	if err != nil {
		err := fmt.Errorf("unable to describe endpoint, %w", err)

		return nil, err
	}
	b, err := json.Marshal(&output)
	if err != nil {
		err := fmt.Errorf("error marshalling json, %w", err)

		return nil, err
	}

	return b, nil
}

// DescribeAllEndpoints returns the URL for all provisioned RDS instances and filters on type
// epType can be either a "writer" or a "reader".
func (r *awsRDS) DescribeAllEndpoints(epType string) ([]byte, error) {
	endpoints := make([]*rds.DBClusterEndpoint, 0)
	result, err := r.client.DescribeDBInstances(nil)
	if err != nil {
		err := fmt.Sprintf("unable to list instances, %v", err)

		return nil, errors.New(err)
	}
	filters := []*rds.Filter{
		{
			Name:   aws.String(dbClusterEndpointType),
			Values: aws.StringSlice([]string{epType}),
		},
	}
	for _, d := range result.DBInstances {
		output, err := r.client.DescribeDBClusterEndpoints(&rds.DescribeDBClusterEndpointsInput{
			DBClusterIdentifier: d.DBClusterIdentifier,
			Filters:             filters,
		})
		if err != nil {
			err := fmt.Sprintf("unable to describe endpoints, %v", err)

			return nil, errors.New(err)
		}
		endpoints = append(endpoints, output.DBClusterEndpoints[0])
	}
	b, err := json.Marshal(&endpoints)
	if err != nil {
		err := fmt.Errorf("error marshalling json, %w", err)

		return nil, err
	}

	return b, nil
}

type instance struct {
	Identifier string    `json:"identifier"`
	CreatedAt  time.Time `json:"createdAt"`
}

// ListDBIdentifiers Returns the identifier and created date/time for provisioned RDS instances.
// This method does not currently support pagination. However, the AWS API does.
func (r *awsRDS) ListDBIdentifiers() ([]byte, error) {
	/*
		to := make([]string, len(s.To))
		for i, t := range s.To {
		    to[i] = t.String()
		}
	*/
	var inst instance
	instances := make([]instance, 0)
	result, err := r.client.DescribeDBInstances(nil)
	if err != nil {
		err := fmt.Sprintf("unable to list instances, %v", err)

		return nil, errors.New(err)
	}
	for _, d := range result.DBInstances {
		if err != nil {
			err := fmt.Sprintf("unable to describe endpoints, %v", err)

			return nil, errors.New(err)
		}
		inst.Identifier = aws.StringValue(d.DBInstanceIdentifier)
		inst.CreatedAt = aws.TimeValue(d.InstanceCreateTime)
		instances = append(instances, inst)
	}
	b, err := json.Marshal(&instances)
	if err != nil {
		err := fmt.Errorf("error marshalling json, %w", err)

		return nil, err
	}

	return b, nil
}
