package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/ransford/iwas/internal"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/spf13/cobra"
)

func (p *Policy) PrintHistory() error {
	ars := p.arn.String()
	getPolicyInput := iam.GetPolicyInput{
		PolicyArn: &ars,
	}
	pol, err := iamClient.GetPolicy(context.Background(), &getPolicyInput)
	if err != nil {
		return err
	}
	ver, err := internal.ParseVersion(*pol.Policy.DefaultVersionId)
	if err != nil {
		return err
	}
	for i := ver; i >= 1; i-- {
		log.Debugf("Getting policy version %d", i)
		doc, err := p.GetVersion(fmt.Sprintf("v%d", i))
		if err != nil {
			return err
		}
		fmt.Println(doc)
	}
	return nil
}

var historyCmd = &cobra.Command{
	Use:     "history <arn>",
	Aliases: []string{"get"},
	Short:   "Show all versions of an IAM policy",
	Args: func(cmd *cobra.Command, args []string) error {
		// ARN only
		if len(args) == 1 {
			if _, err := internal.PolicyNameToArn(args[0]); err != nil {
				return err
			}
			return nil
		}
		return errors.New("wrong number of arguments")
	},
	PreRun: setLogLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		parn, _ := internal.PolicyNameToArn(args[0])
		pol := &Policy{parn}
		return pol.PrintHistory()
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
