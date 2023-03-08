package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ransford/iwas/internal"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/spf13/cobra"
)

func (p *Policy) PrintHistory(since time.Time) error {
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
		if len(args) != 1 {
			return errors.New("wrong number of arguments")
		}
		if _, err := internal.PolicyNameToArn(args[0]); err != nil {
			return err
		}

		if sinceSpec != "" {
			var err error
			since, err = time.Parse("2023-12-01", sinceSpec)
			if err != nil {
				return err
			}
		}

		return nil
	},
	PreRun: setLogLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		parn, _ := internal.PolicyNameToArn(args[0])
		pol := &Policy{parn}
		return pol.PrintHistory(since)
	},
}

var sinceSpec string
var since time.Time

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().StringVar(&sinceSpec, "since", "", "Date in YYYY-MM-DD")
}
