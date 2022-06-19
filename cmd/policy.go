package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/spf13/cobra"
)

var iamClient *iam.Client

func init() {
	rootCmd.AddCommand(policyCmd)

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}
	iamClient = iam.NewFromConfig(cfg)
}

func getPolicyVersion(arn, version string) (string, error) {
	if version == "" {
		getPolicyInput := iam.GetPolicyInput{
			PolicyArn: &arn,
		}
		pol, err := iamClient.GetPolicy(context.Background(), &getPolicyInput)
		if err != nil {
			return "", err
		}
		version = *pol.Policy.DefaultVersionId
	}

	getPolicyVersionInput := iam.GetPolicyVersionInput{
		PolicyArn: &arn,
		VersionId: &version,
	}
	cur, err := iamClient.GetPolicyVersion(context.Background(), &getPolicyVersionInput)
	if err != nil {
		return "", err
	}
	doc, err := url.QueryUnescape(aws.ToString(cur.PolicyVersion.Document))
	if err != nil {
		return "", err
	}
	fmt.Printf("// %s\n", version)
	return doc, nil
}

func validVersion(v string) bool {
	_, err := regexp.MatchString(`v[0-9]+`, v)
	return err == nil
}

var policyCmd = &cobra.Command{
	Use:     "policy <arn> [version]",
	Aliases: []string{"get"},
	Short:   "Show an IAM policy by ARN",
	Args: func(cmd *cobra.Command, args []string) error {
		// ARN only
		if len(args) == 1 {
			if _, err := arn.Parse(args[0]); err != nil {
				return err
			}
			return nil
		}

		// ARN and version
		if len(args) == 2 {
			if _, err := arn.Parse(args[0]); err != nil {
				return err
			}
			if !validVersion(args[1]) {
				return errors.New("invalid version string")
			}
			return nil
		}

		return errors.New("wrong number of arguments")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		arn := args[0]
		version := ""
		if len(args) == 2 {
			version = args[1]
		}
		pol, err := getPolicyVersion(arn, version)
		if err != nil {
			return err
		}
		fmt.Println(pol)
		return nil
	},
}
