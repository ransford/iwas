package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(policyCmd)
}

var policyCmd = &cobra.Command{
	Use:   "policy [arn]",
	Short: "Show an IAM policy by ARN",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		policyArn := args[0]

		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return err
		}
		client := iam.NewFromConfig(cfg)
		getPolicyInput := iam.GetPolicyInput{
			PolicyArn: &policyArn,
		}
		pol, err := client.GetPolicy(context.Background(), &getPolicyInput)
		if err != nil {
			return err
		}

		getPolicyVersionInput := iam.GetPolicyVersionInput{
			PolicyArn: pol.Policy.Arn,
			VersionId: pol.Policy.DefaultVersionId,
		}
		cur, err := client.GetPolicyVersion(context.Background(), &getPolicyVersionInput)
		if err != nil {
			return err
		}
		doc, err := url.QueryUnescape(aws.ToString(cur.PolicyVersion.Document))
		if err != nil {
			return err
		}

		fmt.Println(doc)
		return nil
	},
}
