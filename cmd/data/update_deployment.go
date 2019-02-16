//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package data

import (
	"fmt"

	"github.com/spf13/cobra"

	data "github.com/arangodb-managed/apis/data/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasis/cmd"
	"github.com/arangodb-managed/oasis/pkg/format"
	"github.com/arangodb-managed/oasis/pkg/selection"
)

var (
	// updateDeploymentCmd updates a CA certificate that the user has access to
	updateDeploymentCmd = &cobra.Command{
		Use:   "cacertificate",
		Short: "Update a CA certificate the authenticated user has access to",
		Run:   updateDeploymentCmdRun,
	}
	updateDeploymentArgs struct {
		deploymentID   string
		organizationID string
		projectID      string
		name           string
		description    string
	}
)

func init() {
	cmd.UpdateCmd.AddCommand(updateDeploymentCmd)
	f := updateDeploymentCmd.Flags()
	f.StringVarP(&updateDeploymentArgs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment")
	f.StringVarP(&updateDeploymentArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
	f.StringVarP(&updateDeploymentArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
	f.StringVar(&updateDeploymentArgs.name, "name", "", "Name of the deployment")
	f.StringVar(&updateDeploymentArgs.description, "description", "", "Description of the deployment")
}

func updateDeploymentCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := updateDeploymentArgs
	deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	datac := data.NewDataServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch deployment
	item := selection.MustSelectDeployment(ctx, log, deploymentID, cargs.projectID, cargs.organizationID, datac, rmc)

	// Set changes
	f := c.Flags()
	hasChanges := false
	if f.Changed("name") {
		item.Name = cargs.name
		hasChanges = true
	}
	if f.Changed("description") {
		item.Description = cargs.description
		hasChanges = true
	}
	if !hasChanges {
		fmt.Println("No changes")
	} else {
		// Update deployment
		updated, err := datac.UpdateDeployment(ctx, item)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to update deployment")
		}

		// Show result
		fmt.Println("Updated deployment!")
		fmt.Println(format.Deployment(updated, cmd.RootArgs.Format))
	}
}
