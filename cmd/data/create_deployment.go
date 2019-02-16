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
	// createDeploymentCmd creates a new deployment
	createDeploymentCmd = &cobra.Command{
		Use:   "deployment",
		Short: "Create a new deployment",
		Run:   createDeploymentCmdRun,
	}
	createDeploymentArgs struct {
		name           string
		description    string
		organizationID string
		projectID      string
		regionID       string
		// TODO add other fields
	}
)

func init() {
	cmd.CreateCmd.AddCommand(createDeploymentCmd)

	f := createDeploymentCmd.Flags()
	f.StringVar(&createDeploymentArgs.name, "name", "", "Name of the deployment")
	f.StringVar(&createDeploymentArgs.description, "description", "", "Description of the deployment")
	f.StringVarP(&createDeploymentArgs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization to create the deployment in")
	f.StringVarP(&createDeploymentArgs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project to create the deployment in")
	f.StringVarP(&createDeploymentArgs.regionID, "region-id", "r", cmd.DefaultRegion(), "Identifier of the region to create the deployment in")
}

func createDeploymentCmdRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := cmd.CLILog
	cargs := createDeploymentArgs
	name, argsUsed := cmd.ReqOption("name", cargs.name, args, 0)
	description := cargs.description
	cmd.MustCheckNumberOfArgs(args, argsUsed)

	// Connect
	conn := cmd.MustDialAPI()
	datac := data.NewDataServiceClient(conn)
	rmc := rm.NewResourceManagerServiceClient(conn)
	ctx := cmd.ContextWithToken()

	// Fetch project
	project := selection.MustSelectProject(ctx, log, cargs.projectID, cargs.organizationID, rmc)

	// Create ca certificate
	result, err := datac.CreateDeployment(ctx, &data.Deployment{
		ProjectId:   project.GetId(),
		Name:        name,
		Description: description,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create deployment")
	}

	// Show result
	fmt.Println("Success!")
	fmt.Println(format.Deployment(result, cmd.RootArgs.Format))
}
