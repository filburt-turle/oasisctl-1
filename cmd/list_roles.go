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

	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// listRolesCmd fetches roles of the given organization
	listRolesCmd = &cobra.Command{
		Use:   "roles",
		Short: "List all roles of the given organization",
		Run:   listRolesCmdRun,
	}
	listRolesArgs struct {
		organizationID string
	}
)

func init() {
	listCmd.AddCommand(listRolesCmd)
	f := listRolesCmd.Flags()
	f.StringVarP(&listRolesArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func listRolesCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	organizationID, argsUsed := optOption("organization-id", listRolesArgs.organizationID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch organization
	org := selection.MustSelectOrganization(ctx, cliLog, organizationID, rmc)

	// Fetch roles in organization
	list, err := iamc.ListRoles(ctx, &common.ListOptions{ContextId: org.GetId()})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list roles")
	}

	// Show result
	fmt.Println(format.RoleList(list.Items, rootArgs.format))
}
