package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"go.smartmachine.io/cumulus/pkg/client"
	"net/url"
	"strings"
)

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVarP(&verbose,"verbose", "v", false, "Show extra verbose check output.")
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "checks connections and permissions to our cloud provider",
	Long:  `Makes sure that we have valid credentials and roles to be able to provision cloudy stuff`,
	Run:   check,
}

var verbose bool

func check(cmd *cobra.Command, args []string) {

	passed := true

	profile, err := cmd.Flags().GetString("profile")

	if err != nil {
		panic(fmt.Sprintf("unable to retrieve profile flag: %+v", err))
	}

	c, err := client.New(profile, "")

	if err != nil {
		passed = false
		fmt.Printf("Unable to retrieve shared credentials: %+v\v", err)
	}

	var configs external.Configs
	configs = append(configs, external.WithSharedConfigProfile(profile))

	sharedConfig,err := external.LoadSharedConfig(configs)

	if err != nil {
		passed = false
		fmt.Printf("Unable to get assumed role: %+v\n", err)
	}

	shared,ok := sharedConfig.(external.SharedConfig)

	if !ok {
		passed = false
		fmt.Println("Unable to get shared config.")
	}

	fmt.Printf("%sProfile: %s\n", emoji.Sprint(":name_badge:"), shared.Profile)

	assumedRole := false

	if shared.AssumeRole.RoleARN != "" {
		fmt.Printf("%sAssumed Role: %s\n", emoji.Sprint(":tophat:"), shared.AssumeRole.RoleARN)
		assumedRole = true
	}

	stsSvc := sts.New(c.Config)

	req := stsSvc.GetCallerIdentityRequest(nil)
	res, err := req.Send(context.Background())

	if err != nil {
		passed = false
		fmt.Printf("Unable to get caller identity: %+v\n", err)
	}

	fmt.Printf("%sCaller identity: %s\n", emoji.Sprint(":sparkles:"), *res.Arn)

	iamSvc := iam.New(c.Config)

	if !assumedRole {

		userReq := iamSvc.GetUserRequest(nil)
		userRes, err := userReq.Send(context.Background())
		if err != nil {
			passed = false
			fmt.Printf("Unable to get current user: %+v\n", err)
		}

		fmt.Printf("%sUserName: %s\n", emoji.Sprint(":sparkles:"), *userRes.User.UserName)

		listGroupsReq := iamSvc.ListGroupsForUserRequest(&iam.ListGroupsForUserInput{UserName: userRes.User.UserName})
		listGroupsRes, err := listGroupsReq.Send(context.Background())

		if err != nil {
			passed = false
			fmt.Printf("Unable to get caller groups: %+v\n", err)
		}


		for _, group := range listGroupsRes.Groups {

			fmt.Println("     Group:")
			fmt.Printf("       Name: %s\n", *group.GroupName)

			polReq := iamSvc.ListAttachedGroupPoliciesRequest(&iam.ListAttachedGroupPoliciesInput{GroupName: group.GroupName})
			polRes, err := polReq.Send(context.Background())

			if err != nil {
				passed = false
				fmt.Printf("Unable to get group policies: %+v\n", err)
			}

			for _, policy := range polRes.AttachedPolicies {
				fmt.Println("       Policy:")

				polReq := iamSvc.GetPolicyRequest(&iam.GetPolicyInput{PolicyArn: policy.PolicyArn})
				polRes, err := polReq.Send(context.Background())

				if err != nil {
					passed = false
					fmt.Printf("Unable to get role policies: %+v\n", err)
				}

				fmt.Printf("         Name: %s\n", *polRes.Policy.PolicyName)

				if polRes.Policy.Description != nil {
					fmt.Printf("         Description: %s\n", *polRes.Policy.Description)
				}

				if verbose {

					polVerReq := iamSvc.GetPolicyVersionRequest(&iam.GetPolicyVersionInput{
						PolicyArn: polRes.Policy.Arn,
						VersionId: polRes.Policy.DefaultVersionId,
					})
					polVerRes, err := polVerReq.Send(context.Background())

					if err != nil {
						passed = false
						fmt.Printf("Unable to get policy document: %+v\n", err)
					}

					doc, err := url.QueryUnescape(*polVerRes.PolicyVersion.Document)
					if err != nil {
						passed = false
						fmt.Printf("Unable to parse document: %+v\n", err)
					}

					fmt.Println("         Document:")
					for _, line := range strings.Split(doc, "\n") {
						fmt.Printf("           %s\n", line)
					}

				}

			}

		}


	} else {

		roleARN, err := arn.Parse(shared.AssumeRole.RoleARN)


		if err != nil {
			passed = false
			fmt.Printf("Unable to get role policies: %+v\n", err)
		}

		roleName := strings.Replace(roleARN.Resource, "role/", "", 1)

		roleReq := iamSvc.GetRoleRequest(&iam.GetRoleInput{RoleName: &roleName})
		roleRes, err := roleReq.Send(context.Background())

		if err != nil {
			passed = false
			fmt.Printf("Unable to get role: %+v\n", err)
		}

		fmt.Printf("%sRole: %s\n", emoji.Sprint(":sparkles:"), *roleRes.Role.RoleName)

		polReq := iamSvc.ListAttachedRolePoliciesRequest(&iam.ListAttachedRolePoliciesInput{RoleName: &roleName})
		polRes, err := polReq.Send(context.Background())

		if err != nil {
			passed = false
			fmt.Printf("Unable to get role policies: %+v\n", err)
		}

		fmt.Println("     Policy:")

		for _, pol := range polRes.AttachedPolicies {

			polReq := iamSvc.GetPolicyRequest(&iam.GetPolicyInput{PolicyArn: pol.PolicyArn})
			polRes, err := polReq.Send(context.Background())

			if err != nil {
				passed = false
				fmt.Printf("Unable to get role policies: %+v\n", err)
			}

			fmt.Printf("       Name:        %s\n", *polRes.Policy.PolicyName)
			if polRes.Policy.Description != nil {
				fmt.Printf("       Description: %s\n", *polRes.Policy.Description)
			}

			if verbose {

				polVerReq := iamSvc.GetPolicyVersionRequest(&iam.GetPolicyVersionInput{
					PolicyArn: polRes.Policy.Arn,
					VersionId: polRes.Policy.DefaultVersionId,
				})
				polVerRes, err := polVerReq.Send(context.Background())

				if err != nil {
					passed = false
					fmt.Printf("Unable to get policy document: %+v\n", err)
				}

				doc, err := url.QueryUnescape(*polVerRes.PolicyVersion.Document)
				if err != nil {
					passed = false
					fmt.Printf("Unable to parse document: %+v\n", err)
				}

				fmt.Println("       Document:")
				for _, line := range strings.Split(doc, "\n") {
					fmt.Printf("         %s\n", line)
				}

			}
		}




	}

	if passed {
		fmt.Printf("%sAll checks passed.\n", emoji.Sprint(":sparkles:"))
	} else {
		fmt.Printf("%sOne or more checks failed.\n", emoji.Sprint(":poop:"))
	}
}
