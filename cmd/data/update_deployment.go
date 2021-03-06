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

package data

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	data "github.com/arangodb-managed/apis/data/v1"
	rm "github.com/arangodb-managed/apis/resourcemanager/v1"

	"github.com/arangodb-managed/oasisctl/cmd"
	"github.com/arangodb-managed/oasisctl/pkg/format"
	"github.com/arangodb-managed/oasisctl/pkg/selection"
)

func init() {
	cmd.InitCommand(
		cmd.UpdateCmd,
		&cobra.Command{
			Use:   "deployment",
			Short: "Update a deployment the authenticated user has access to",
		},
		func(c *cobra.Command, f *flag.FlagSet) {
			cargs := &struct {
				deploymentID   string
				organizationID string
				projectID      string
				name           string
				description    string
				ipwhitelistID  string

				version               string
				model                 string
				nodeSizeID            string
				nodeCount             int32
				nodeDiskSize          int32
				coordinators          int32
				coordinatorMemorySize int32
				dbservers             int32
				dbserverMemorySize    int32
				dbserverDiskSize      int32
			}{}
			f.StringVarP(&cargs.deploymentID, "deployment-id", "d", cmd.DefaultDeployment(), "Identifier of the deployment")
			f.StringVarP(&cargs.organizationID, "organization-id", "o", cmd.DefaultOrganization(), "Identifier of the organization")
			f.StringVarP(&cargs.projectID, "project-id", "p", cmd.DefaultProject(), "Identifier of the project")
			f.StringVar(&cargs.name, "name", "", "Name of the deployment")
			f.StringVar(&cargs.description, "description", "", "Description of the deployment")
			f.StringVarP(&cargs.ipwhitelistID, "ipwhitelist-id", "i", cmd.DefaultIPWhitelist(), "Identifier of the IP whitelist to use for the deployment")
			f.StringVar(&cargs.version, "version", "", "Version of ArangoDB to use for the deployment")
			f.StringVar(&cargs.model, "model", data.ModelOneShard, "Set model of the deployment")
			f.StringVar(&cargs.nodeSizeID, "node-size-id", "", "Set the node size to use for this deployment")
			f.Int32Var(&cargs.nodeCount, "node-count", 3, "Set the number of desired nodes")
			f.Int32Var(&cargs.nodeDiskSize, "node-disk-size", 0, "Set disk size for nodes (GB)")
			f.Int32Var(&cargs.coordinators, "coordinators", 3, "Set number of coordinators for flexible deployments")
			f.Int32Var(&cargs.coordinatorMemorySize, "coordinator-memory-size", 4, "Set memory size of coordinators for flexible deployments (GB)")
			f.Int32Var(&cargs.dbservers, "dbservers", 3, "Set number of dbservers for flexible deployments")
			f.Int32Var(&cargs.dbserverMemorySize, "dbserver-memory-size", 4, "Set memory size of dbservers for flexible deployments (GB)")
			f.Int32Var(&cargs.dbserverDiskSize, "dbserver-disk-size", 32, "Set disk size of dbservers for flexible deployments (GB)")

			c.Run = func(c *cobra.Command, args []string) {
				// Validate arguments
				log := cmd.CLILog
				deploymentID, argsUsed := cmd.OptOption("deployment-id", cargs.deploymentID, args, 0)
				cmd.MustCheckNumberOfArgs(args, argsUsed)

				// Connect
				conn := cmd.MustDialAPI()
				datac := data.NewDataServiceClient(conn)
				rmc := rm.NewResourceManagerServiceClient(conn)
				ctx := cmd.ContextWithToken()

				// Fetch deployment
				item := selection.MustSelectDeployment(ctx, log, deploymentID, cargs.projectID, cargs.organizationID, datac, rmc)
				ensureModel := func() *data.Deployment_ModelSpec {
					if item.Model == nil {
						item.Model = &data.Deployment_ModelSpec{}
					}
					return item.Model
				}
				ensureServers := func() *data.Deployment_ServersSpec {
					if item.Servers == nil {
						item.Servers = &data.Deployment_ServersSpec{}
					}
					return item.Servers
				}
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
				if f.Changed("ipwhitelist-id") {
					item.IpwhitelistId = cargs.ipwhitelistID
					hasChanges = true
				}
				if f.Changed("version") {
					item.Version = cargs.version
					hasChanges = true
				}
				if f.Changed("model") {
					ensureModel().Model = cargs.model
					hasChanges = true
				}
				if f.Changed("node-size-id") {
					ensureModel().NodeSizeId = cargs.nodeSizeID
					hasChanges = true
				}
				if f.Changed("node-count") {
					ensureModel().NodeCount = cargs.nodeCount
					hasChanges = true
				}
				if f.Changed("node-disk-size") {
					ensureModel().NodeDiskSize = cargs.nodeDiskSize
					hasChanges = true
				}
				if f.Changed("coordinators") {
					ensureServers().Coordinators = cargs.coordinators
					hasChanges = true
				}
				if f.Changed("coordinator-memory-size") {
					ensureServers().CoordinatorMemorySize = cargs.coordinatorMemorySize
					hasChanges = true
				}
				if f.Changed("dbservers") {
					ensureServers().Dbservers = cargs.dbservers
					hasChanges = true
				}
				if f.Changed("dbserver-memory-size") {
					ensureServers().DbserverMemorySize = cargs.dbserverMemorySize
					hasChanges = true
				}
				if f.Changed("dbserver-disk-size") {
					ensureServers().DbserverDiskSize = cargs.dbserverDiskSize
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
					fmt.Println(format.Deployment(updated, nil, cmd.RootArgs.Format, false))
				}
			}
		},
	)
}
