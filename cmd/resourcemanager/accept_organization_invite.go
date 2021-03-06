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
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// acceptOrganizationInviteCmd accepts an organization invite that the user has access to
	acceptOrganizationInviteCmd = &cobra.Command{
		Use:   "invite",
		Short: "Accept an organization invite the authenticated user has access to",
		Run:   acceptOrganizationInviteCmdRun,
	}
	acceptOrganizationInviteArgs struct {
		organizationID string
		inviteID       string
	}
)

func init() {
	acceptOrganizationCmd.AddCommand(acceptOrganizationInviteCmd)
	f := acceptOrganizationInviteCmd.Flags()
	f.StringVarP(&acceptOrganizationInviteArgs.inviteID, "invite-id", "i", cmd.DefaultOrganizationInvite(), "Identifier of the organization invite")
	f.StringVarP(&acceptOrganizationInviteArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func acceptOrganizationInviteCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := acceptOrganizationInviteArgs
	inviteID, argsUsed := cmd.OptOption("invite-id", cargs.inviteID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch invite
	invite := selection.MustSelectOrganizationInvite(ctx, log, inviteID, cargs.organizationID, rmc)

	// Accept invite
	if _, err := rmc.AcceptOrganizationInvite(ctx, &common.IDOptions{Id: invite.GetId()}); err != nil {
		log.Fatal().Err(err).Msg("Failed to accept organization invite")
	}

	// Fetch organization
	orgName := invite.GetOrganizationId()
	if org, err := rmc.GetOrganization(ctx, &common.IDOptions{Id: invite.GetOrganizationId()}); err != nil {
		log.Warn().Err(err).Msg("Failed to get organization")
	} else {
		orgName = org.GetName()
	}

	// Show result
	format.DisplaySuccess(cmd.RootArgs.Format)
	fmt.Printf("You are now a member of the '%s' organization.\n", orgName)
}
