//
// DISCLAIMER
//
// Copyright 2020 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//
// Author Ewout Prangsma
//

package rm

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// deleteOrganizationInviteCmd deletes an organization invite that the user has access to
	deleteOrganizationInviteCmd = &cobra.Command{
		Use:   "invite",
		Short: "Delete an organization invite the authenticated user has access to",
		Run:   deleteOrganizationInviteCmdRun,
	}
	deleteOrganizationInviteArgs struct {
		organizationID string
		inviteID       string
	}
)

func init() {
	deleteOrganizationCmd.AddCommand(deleteOrganizationInviteCmd)
	f := deleteOrganizationInviteCmd.Flags()
	f.StringVarP(&deleteOrganizationInviteArgs.inviteID, "invite-id", "i", cmd.DefaultOrganizationInvite(), "Identifier of the organization invite")
	f.StringVarP(&deleteOrganizationInviteArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func deleteOrganizationInviteCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := deleteOrganizationInviteArgs
	inviteID, argsUsed := cmd.OptOption("invite-id", cargs.inviteID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch invite
	item := selection.MustSelectOrganizationInvite(ctx, log, inviteID, cargs.organizationID, rmc)

	// Delete invite
	if _, err := rmc.DeleteOrganizationInvite(ctx, &common.IDOptions{Id: item.GetId()}); err != nil {
		log.Fatal().Err(err).Msg("Failed to delete organization invite")
	}

	// Show result
	fmt.Println("Deleted organization invite!")
}
