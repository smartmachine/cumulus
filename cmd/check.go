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
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "checks connections and permissions to our cloud provider",
	Long:  `Makes sure that we have valid credentials and roles to be able to provision cloudy stuff`,
	Run:   check,
}

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
		fmt.Printf("Unable to get shared config.\n", err)
	}

	fmt.Printf("%sProfile: %s\n", emoji.Sprint(":sparkles:"), shared.Profile)

	assumedRole := false

	if shared.AssumeRole.RoleARN != "" {
		fmt.Printf("%sAssumed Role: %s\n", emoji.Sprint(":sparkles:"), shared.AssumeRole.RoleARN)
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

		fmt.Printf("%sUser: %+v\n", emoji.Sprint(":sparkles:"), userRes.User)

		listGroupsReq := iamSvc.ListGroupsForUserRequest(&iam.ListGroupsForUserInput{UserName: userRes.User.UserName})
		listGroupsRes, err := listGroupsReq.Send(context.Background())

		if err != nil {
			passed = false
			fmt.Printf("Unable to get caller groups: %+v\n", err)
		}

		fmt.Printf("%sCaller groups: %+v\n", emoji.Sprint(":sparkles:"), listGroupsRes.Groups)

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

		fmt.Printf("%sRole: %+v\n", emoji.Sprint(":sparkles:"), roleRes.Role)

		polReq := iamSvc.ListAttachedRolePoliciesRequest(&iam.ListAttachedRolePoliciesInput{RoleName: &roleName})
		polRes, err := polReq.Send(context.Background())

		if err != nil {
			passed = false
			fmt.Printf("Unable to get role policies: %+v\n", err)
		}

		fmt.Printf("%sPolicies: %+v\n", emoji.Sprint(":sparkles:"), polRes.AttachedPolicies)

		for _, pol := range polRes.AttachedPolicies {

			polReq := iamSvc.GetPolicyRequest(&iam.GetPolicyInput{PolicyArn: pol.PolicyArn})
			polRes, err := polReq.Send(context.Background())

			if err != nil {
				passed = false
				fmt.Printf("Unable to get role policies: %+v\n", err)
			}

			fmt.Printf("%sPolicy: %+v\n", emoji.Sprint(":sparkles:"), polRes.Policy)

			polVerReq := iamSvc.GetPolicyVersionRequest(&iam.GetPolicyVersionInput{
				PolicyArn: polRes.Policy.Arn,
				VersionId: polRes.Policy.DefaultVersionId,
			})
			polVerRes,err := polVerReq.Send(context.Background())

			if err != nil {
				passed = false
				fmt.Printf("Unable to get policy document: %+v\n", err)
			}

			doc, err := url.QueryUnescape(*polVerRes.PolicyVersion.Document)
			if err != nil {
				passed = false
				fmt.Printf("Unable to parse document: %+v\n", err)
			}

			fmt.Printf("%sPolicy Document: %s\n", emoji.Sprint(":sparkles:"), doc)
		}




	}

	if passed {
		fmt.Printf("%sAll checks passed.\n", emoji.Sprint(":sparkles:"))
	} else {
		fmt.Printf("%sOne or more checks failed.\n", emoji.Sprint(":poop:"))
	}
}
