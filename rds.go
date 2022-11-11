package ops_aws

type RDS interface {
	DescribeAllEndpoints(epType string) ([]byte, error)
	ListDBIdentifiers() ([]byte, error)
	DescribeEndpoint(clusterName string, epType string) ([]byte, error)
}
