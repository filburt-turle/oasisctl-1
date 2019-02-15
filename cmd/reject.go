//
// DISCLAIMER
//
// Copyright 2019 ArangoDB Inc, Cologne, Germany
//
// Author Ewout Prangsma
//

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// RejectCmd is root for various `reject ...` commands
	RejectCmd = &cobra.Command{
		Use:   "reject",
		Short: "Reject invites",
		Run:   ShowUsage,
	}
)

func init() {
	RootCmd.AddCommand(RejectCmd)
}
