package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/ransford/iwas/internal"
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

func getPolicyVersion(a arn.ARN, version string) (string, error) {
	ars := a.String()
	if version == "" {
		getPolicyInput := iam.GetPolicyInput{
			PolicyArn: &ars,
		}
		pol, err := iamClient.GetPolicy(context.Background(), &getPolicyInput)
		if err != nil {
			return "", err
		}
		version = *pol.Policy.DefaultVersionId
	}

	getPolicyVersionInput := iam.GetPolicyVersionInput{
		PolicyArn: &ars,
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
	m, err := regexp.MatchString(`v[1-9][0-9]*`, v)
	return err == nil && m
}

func prettyPrintPolicy(pol string) {
	var pretty bytes.Buffer
	json.Indent(&pretty, []byte(pol), "", "  ")
	fmt.Println(pretty.String())
}

var policyCmd = &cobra.Command{
	Use:     "policy <arn> [version]",
	Aliases: []string{"get"},
	Short:   "Show an IAM policy by ARN",
	Args: func(cmd *cobra.Command, args []string) error {
		// ARN only
		if len(args) == 1 {
			if _, err := internal.PolicyNameToArn(args[0]); err != nil {
				return err
			}
			return nil
		}

		// ARN and version
		if len(args) == 2 {
			if _, err := internal.PolicyNameToArn(args[0]); err != nil {
				return err
			}
			if !validVersion(args[1]) {
				return errors.New("invalid version string")
			}
			return nil
		}

		return errors.New("wrong number of arguments")
	},
	PreRun: setLogLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		var a arn.ARN
		var err error
		if a, err = internal.PolicyNameToArn(args[0]); err != nil {
			return err
		}
		version := ""
		if len(args) == 2 {
			version = args[1]
		}
		pol, err := getPolicyVersion(a, version)
		if err != nil {
			return err
		}
		prettyPrintPolicy(pol)
		return nil
	},
}
