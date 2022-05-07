# irsa
`irsa` is a simple CLI that creates IAM Roles for K8s Service Accounts.


## About

IAM Roles for Service Accounts are a great way to operate a more secure infrastructure with scoped down permissions per pod. But crafting a trust policy for IRSA is a little difficult. This CLI wraps that process. 

## Usage

```bash
irsa is a simple CLI tool that creates IAM Roles for K8s Service Accounts

Usage:
  irsa [flags]

Flags:
      --cluster-name string      the EKS cluster name
  -h, --help                     help for irsa
      --policies strings         policy from a file (file://<>) or a URL (http(s)://<>)
      --policy-arns strings      policy ARNs to add to the IAM Role
  -p, --profile string           the AWS Profile
  -r, --region string            the AWS Region
      --role-name string         the name of the IAM Role
      --service-account string   the namespaced name of the service account (i.e. my-namespace/my-sa
  -v, --version                  the version
```

## Example:

```bash
go run cmd/main.go --cluster-name wagnerbm-karpenter-demo \
   --role-name ebs-csi-driver \
   --policies 'https://raw.githubusercontent.com/kubernetes-sigs/aws-ebs-csi-driver/master/docs/example-iam-policy.json' \
   --service-account kube-system/ebs-csi
arn:aws:iam::332273710158:role/ebs-csi-driver
```