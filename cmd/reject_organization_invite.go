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

	rm "github.com/arangodb-managed/apis/resourcemanager/v1"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// rejectOrganizationInviteCmd rejects an organization invite that the user has access to
	rejectOrganizationInviteCmd = &cobra.Command{
		Use:   "invite",
		Short: "Reject an organization invite the authenticated user has access to",
		Run:   rejectOrganizationInviteCmdRun,
	}
	rejectOrganizationInviteArgs struct {
		organizationID string
		inviteID       string
	}
)

func init() {
	rejectOrganizationCmd.AddCommand(rejectOrganizationInviteCmd)
	f := rejectOrganizationInviteCmd.Flags()
	f.StringVarP(&rejectOrganizationInviteArgs.inviteID, "invite-id", "i", defaultOrganizationInvite(), "Identifier of the organization invite")
	f.StringVarP(&rejectOrganizationInviteArgs.organizationID, "organization-id", "o", defaultOrganization(), "Identifier of the organization")
}

func rejectOrganizationInviteCmdRun(cmd *cobra.Command, args []string) {
	// Validate arguments
	inviteID, argsUsed := optOption("invite-id", rejectOrganizationInviteArgs.inviteID, args, 0)
	mustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := mustDialAPI()
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := contextWithToken()

	// Fetch invite
	invite := selection.MustSelectOrganizationInvite(ctx, cliLog, inviteID, rejectOrganizationInviteArgs.organizationID, rmc)

	// Reject invite
	if _, err := rmc.RejectOrganizationInvite(ctx, invite); err != nil {
		cliLog.Fatal().Err(err).Msg("Failed to reject organization invite")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println("You have rejected the invite.")
}
