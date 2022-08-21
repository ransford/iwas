package internal

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func ParseVersion(ver string) (i int, err error) {
	v := strings.TrimPrefix(ver, "v")
	i, err = strconv.Atoi(v)
	if i < 1 {
		return -1, errors.New("version number too low")
	}
	return
}

// IAM is global, implying us-east-1
const IAM_REGION = "us-east-1"

// PolicyNameToArn converts a policy name to an ARN.
func PolicyNameToArn(pol string) (arn.ARN, error) {
	if a, err := arn.Parse(pol); err == nil {
		if a.Service == "iam" && strings.HasPrefix(a.Resource, "policy/") {
			return a, nil
		}
		return arn.ARN{}, errors.New("not an IAM policy ARN")
	}

	// Error was non-nil, so not a valid ARN; construct one by interpreting pol as the name of a
	// policy resource within the current account
	cfg, err := config.LoadDefaultConfig(
		context.TODO())
	if err != nil {
		return arn.ARN{}, err
	}
	cli := sts.NewFromConfig(cfg)
	out, err := cli.GetCallerIdentity(
		context.TODO(),
		&sts.GetCallerIdentityInput{})
	if err != nil {
		return arn.ARN{}, err
	}

	return arn.ARN{
		Partition: "aws",
		Service:   "iam",
		Region:    "",
		AccountID: *out.Account,
		Resource:  fmt.Sprintf("policy/%s", pol),
	}, nil
}
