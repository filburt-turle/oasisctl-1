//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package iam

import (
	"fmt"

	"github.com/spf13/cobra"

	iam "github.com/arangodb-managed/apis/iam/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// updatePolicyDeleteBindingCmd deleted a role binding from a policy
	updatePolicyDeleteBindingCmd = &cobra.Command{
		Use:   "binding",
		Short: "Delete a role binding from a policy",
		Run:   updatePolicyDeleteBindingCmdRun,
	}
	updatePolicyDeleteBindingArgs struct {
		url      string
		roleID   string
		userIDs  []string
		groupIDs []string
	}
)

func init() {
	updatePolicyDeleteCmd.AddCommand(updatePolicyDeleteBindingCmd)
	f := updatePolicyDeleteBindingCmd.Flags()
	f.StringVarP(&updatePolicyDeleteBindingArgs.url, "url", "u", cmd.DefaultURL(), "URL of the resource to update the policy for")
	f.StringVarP(&updatePolicyDeleteBindingArgs.roleID, "role-id", "r", cmd.DefaultRole(), "Identifier of the role to delete bind for")
	f.StringSliceVar(&updatePolicyDeleteBindingArgs.userIDs, "user-id", nil, "Identifiers of the users to delete bindings for")
	f.StringSliceVar(&updatePolicyDeleteBindingArgs.groupIDs, "group-id", nil, "Identifiers of the groups to delete bindings for")
}

func updatePolicyDeleteBindingCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := updatePolicyDeleteBindingArgs
	url, argsUsed := cmd.OptOption("url", cargs.url, args, 0)
	roleID, _ := cmd.ReqOption("role-id", cargs.roleID, nil, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)
	if len(cargs.userIDs) == 0 &&
		len(cargs.groupIDs) == 0 {
		log.Fatal().Msg("Provide at least one --user-id or --group-id")
	}

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Add role binding
	req := &iam.RoleBindingsRequest{
		ResourceUrl: url,
	}
	for _, uid := range cargs.userIDs {
		req.Bindings = append(req.Bindings, &iam.RoleBinding{
			MemberId: iam.CreateMemberIDFromUserID(uid),
			RoleId:   roleID,
		})
	}
	for _, gid := range cargs.groupIDs {
		req.Bindings = append(req.Bindings, &iam.RoleBinding{
			MemberId: iam.CreateMemberIDFromGroupID(gid),
			RoleId:   roleID,
		})
	}
	updated, err := iamc.DeleteRoleBindings(ctx, req)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to update policy")
	}

	// Show result
	fmt.Println("Updated policy!")
	fmt.Println(format.Policy(ctx, updated, iamc, cmd.RootArgs.Format))
}