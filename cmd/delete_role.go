//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
)

var (
	// deleteRoleCmd deletes a role that the user has access to
	deleteRoleCmd = &cobra.Command{
		Use:   "role",
		Short: "Delete a role the authenticated user has access to",
		Run:   deleteRoleCmdRun,
	}
	deleteRoleArgs struct {
		organizationID string
		roleID         string
	}
)

func init() {
	deleteCmd.AddCommand(deleteRoleCmd)
	f := deleteRoleCmd.Flags()
	f.StringVarP(&deleteRoleArgs.roleID, "role-id", "r", defaultRole(), "Identifier of the role")
	f.StringVarP(&deleteRoleArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func deleteRoleCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	roleID, argsUsed := optOption("role-id", deleteRoleArgs.roleID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch role
	item := mustSelectRole(ctx, roleID, deleteRoleArgs.organizationID, iamc, rmc)

	// Delete role
	if _, err := iamc.DeleteRole(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to delete role")
	}

	// Show result
	fmt.Println("Deleted role!")
}
