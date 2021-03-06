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

package iam

import (
	"fmt"

	"github.com/spf13/cobra"

	common "github.com/arangodb-managed/apis/common/v1"
	iam "github.com/arangodb-managed/apis/iam/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

var (
	// listGroupMembersCmd fetches the members of a group the user is a part of
	listGroupMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "List members of a group the authenticated user is a member of",
		Run:   listGroupMembersCmdRun,
	}
	listGroupMembersArgs struct {
		groupID        string
		organizationID string
	}
)

func init() {
	listGroupCmd.AddCommand(listGroupMembersCmd)
	f := listGroupMembersCmd.Flags()
	f.StringVarP(&listGroupMembersArgs.groupID, "group-id", "g", cmd.DefaultGroup(), "Identifier of the group")
	f.StringVarP(&listGroupMembersArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
}

func listGroupMembersCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := listGroupMembersArgs
	groupID, argsUsed := cmd.ReqOption("group-id", cargs.groupID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	iamc := iam.NewIAMServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch group
	group := selection.MustSelectGroup(ctx, log, groupID, cargs.organizationID, iamc, rmc)

	list, err := iamc.ListGroupMembers(ctx, &common.ListOptions{ContextId: group.GetId()})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list group members")
	}

	// Show result
	fmt.Println(format.GroupMemberList(ctx, list.GetItems(), iamc, cmd.RootArgs.Format))
}
