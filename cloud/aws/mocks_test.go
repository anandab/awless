package aws

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/wallix/awless/graph"
)

type mockEc2 struct {
	ec2iface.EC2API
	vpcs             []*ec2.Vpc
	subnets          []*ec2.Subnet
	instances        []*ec2.Instance
	securityGroups   []*ec2.SecurityGroup
	keyPairs         []*ec2.KeyPairInfo
	internetGateways []*ec2.InternetGateway
	routeTables      []*ec2.RouteTable
}

func (m *mockEc2) DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return &ec2.DescribeVpcsOutput{Vpcs: m.vpcs}, nil
}

func (m *mockEc2) DescribeSubnets(input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return &ec2.DescribeSubnetsOutput{Subnets: m.subnets}, nil
}

func (m *mockEc2) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{Reservations: []*ec2.Reservation{{Instances: m.instances}}}, nil
}

func (m *mockEc2) DescribeInstancesPages(input *ec2.DescribeInstancesInput, fn func(p *ec2.DescribeInstancesOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&ec2.DescribeInstancesOutput{Reservations: []*ec2.Reservation{{Instances: m.instances}}}, true)
	return nil
}

func (m *mockEc2) DescribeSecurityGroups(input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return &ec2.DescribeSecurityGroupsOutput{SecurityGroups: m.securityGroups}, nil
}

func (m *mockEc2) DescribeKeyPairs(input *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	return &ec2.DescribeKeyPairsOutput{KeyPairs: m.keyPairs}, nil
}

func (m *mockEc2) DescribeInternetGateways(input *ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
	return &ec2.DescribeInternetGatewaysOutput{InternetGateways: m.internetGateways}, nil
}

func (m *mockEc2) DescribeRouteTables(input *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
	return &ec2.DescribeRouteTablesOutput{RouteTables: m.routeTables}, nil
}

// Not tested
func (m *mockEc2) DescribeVolumes(input *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	return &ec2.DescribeVolumesOutput{}, nil
}
func (m *mockEc2) DescribeVolumesPages(input *ec2.DescribeVolumesInput, fn func(p *ec2.DescribeVolumesOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&ec2.DescribeVolumesOutput{}, true)
	return nil
}
func (m *mockEc2) DescribeAvailabilityZones(input *ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
	return &ec2.DescribeAvailabilityZonesOutput{}, nil
}

type mockIam struct {
	iamiface.IAMAPI
	groups          []*iam.GroupDetail
	managedPolicies []*iam.ManagedPolicyDetail
	roles           []*iam.RoleDetail
	users           []*iam.User
	usersDetails    []*iam.UserDetail
}

func (m *mockIam) ListUsers(input *iam.ListUsersInput) (*iam.ListUsersOutput, error) {
	return &iam.ListUsersOutput{Users: m.users}, nil
}

func (m *mockIam) ListUsersPages(input *iam.ListUsersInput, fn func(p *iam.ListUsersOutput, lastPage bool) (shouldContinue bool)) error {
	fn(&iam.ListUsersOutput{Users: m.users}, true)
	return nil
}

func (m *mockIam) ListPolicies(input *iam.ListPoliciesInput) (*iam.ListPoliciesOutput, error) {
	var policies []*iam.Policy
	for _, p := range m.managedPolicies {
		policy := &iam.Policy{PolicyId: p.PolicyId, PolicyName: p.PolicyName}
		policies = append(policies, policy)
	}
	return &iam.ListPoliciesOutput{Policies: policies}, nil
}

func (m *mockIam) GetAccountAuthorizationDetails(input *iam.GetAccountAuthorizationDetailsInput) (*iam.GetAccountAuthorizationDetailsOutput, error) {
	return &iam.GetAccountAuthorizationDetailsOutput{GroupDetailList: m.groups, Policies: m.managedPolicies, RoleDetailList: m.roles, UserDetailList: m.usersDetails}, nil
}

func stringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

type mockS3 struct {
	s3iface.S3API
	bucketsACL       map[string][]*s3.Grant
	bucketsPerRegion map[string][]*s3.Bucket
	objectsPerBucket map[string][]*s3.Object
}

func (m *mockS3) GetBucketAcl(input *s3.GetBucketAclInput) (*s3.GetBucketAclOutput, error) {
	return &s3.GetBucketAclOutput{Grants: m.bucketsACL[awssdk.StringValue(input.Bucket)]}, nil
}
func (m *mockS3) Name() string {
	return ""
}
func (m *mockS3) Provider() string {
	return ""
}
func (m *mockS3) ProviderAPI() string {
	return ""
}
func (m *mockS3) ProviderRunnableAPI() interface{} {
	return m
}
func (m *mockS3) ResourceTypes() []string {
	return []string{}
}
func (m *mockS3) FetchResources() (*graph.Graph, error) {
	return nil, nil
}
func (m *mockS3) FetchByType(t string) (*graph.Graph, error) {
	return nil, nil
}
func (m *mockS3) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	var buckets []*s3.Bucket
	for _, b := range m.bucketsPerRegion {
		buckets = append(buckets, b...)
	}
	return &s3.ListBucketsOutput{Buckets: buckets}, nil
}
func (m *mockS3) ListObjects(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	return &s3.ListObjectsOutput{Contents: m.objectsPerBucket[awssdk.StringValue(input.Bucket)]}, nil
}
func (m *mockS3) GetBucketLocation(input *s3.GetBucketLocationInput) (*s3.GetBucketLocationOutput, error) {
	for region, buckets := range m.bucketsPerRegion {
		for _, bucket := range buckets {
			if awssdk.StringValue(bucket.Name) == awssdk.StringValue(input.Bucket) {
				return &s3.GetBucketLocationOutput{LocationConstraint: awssdk.String(region)}, nil
			}
		}
	}
	return nil, fmt.Errorf("bucket location mock: bucket %s not found", awssdk.StringValue(input.Bucket))
}
