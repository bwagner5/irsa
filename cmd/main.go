package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/bwagner5/irsa/pkg/irsa"
	"github.com/spf13/cobra"
)

// todos(bwagner5)
//   1. auto-discover cluster from current kube context
//   2. allow for updating if the role exists but something has changed
//        - policies have changed (different ones attached or policy body doesn't match)
//        - service account changed

var version string

type Options struct {
	clusterName    string
	serviceAccount string
	roleName       string
	policyARNs     []string
	policies       []string
	version        bool
	region         string
	profile        string
}

func main() {
	options := Options{}
	rootCmd := &cobra.Command{
		Use:   "irsa",
		Short: "irsa is a simple CLI tool that creates IAM Roles for K8s Service Accounts",
		Run: func(cmd *cobra.Command, _ []string) {
			if options.version {
				fmt.Println(version)
				os.Exit(0)
			}
			ctx := context.Background()
			cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(options.region), config.WithSharedConfigProfile(options.profile))
			if err != nil {
				fmt.Printf("❌ %s\n", err)
				os.Exit(1)
			}
			serviceAccountNamespacedName := strings.Split(options.serviceAccount, "/")
			if len(serviceAccountNamespacedName) != 2 {
				fmt.Printf("❌ Service account must be in the format: <namespace>/<service account>")
				os.Exit(1)
			}
			namespace := serviceAccountNamespacedName[0]
			serviceAccountName := serviceAccountNamespacedName[1]
			irsaGen := irsa.New(cfg)
			role, err := irsaGen.Create(
				context.Background(),
				options.clusterName,
				namespace,
				serviceAccountName,
				options.roleName,
				options.policyARNs,
				options.policies)
			if err != nil {
				fmt.Printf("❌ %s\n", err)
				os.Exit(1)
			}
			fmt.Println(role)
		},
	}
	rootCmd.PersistentFlags().StringVar(&options.roleName, "role-name", "", "the name of the IAM Role")
	rootCmd.PersistentFlags().StringSliceVar(&options.policies, "policies", []string{}, "policy from a file (file://<>) or a URL (http(s)://<>)")
	rootCmd.PersistentFlags().StringSliceVar(&options.policyARNs, "policy-arns", []string{}, "policy ARNs to add to the IAM Role")
	rootCmd.PersistentFlags().StringVar(&options.clusterName, "cluster-name", "", "the EKS cluster name")
	rootCmd.PersistentFlags().StringVar(&options.serviceAccount, "service-account", "", "the namespaced name of the service account (i.e. my-namespace/my-sa")
	rootCmd.PersistentFlags().BoolVarP(&options.version, "version", "v", false, "the version")
	rootCmd.PersistentFlags().StringVarP(&options.region, "region", "r", "", "the AWS Region")
	rootCmd.PersistentFlags().StringVarP(&options.profile, "profile", "p", "", "the AWS Profile")
	rootCmd.Execute()
}
