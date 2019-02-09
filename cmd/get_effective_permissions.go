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

	"github.com/arangodb-managed/oasis/pkg/format"
)

var (
	// getEffectivePermissionsCmd fetches the effective permissions of the user for a given URL.
	getEffectivePermissionsCmd = &cobra.Command{
		Use:   "permissions",
		Short: "Get the effective permissions, the authenticated user has for a given URL",
		Run:   getEffectivePermissionsCmdRun,
	}
	getEffectivePermissionsArgs struct {
		url string
	}
)

func init() {
	getEffectiveCmd.AddCommand(getEffectivePermissionsCmd)
	f := getEffectivePermissionsCmd.Flags()
	f.StringVarP(&getEffectivePermissionsArgs.url, "url", "u", defaultURL(), "URL of resource to get effective permissions for")
}

func getEffectivePermissionsCmdRun(cmd *cobra.Command, args []string) {
	// Connect
	conn := mustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	ctx := contextWithToken()

	// Fetch permissions
	list, err := iamc.GetEffectivePermissions(ctx, &common.URLOptions{Url: getEffectivePermissionsArgs.url})
	if err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to list organizations")
	}

	// Show result
	fmt.Println(format.PermissionList(list.Items, rootArgs.format))
}
