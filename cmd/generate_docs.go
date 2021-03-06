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
// Author Gergely Brautigam
//

package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	// GenerateCmd is root for various `generate ...` commands
	GenerateCmd = &cobra.Command{
		Use:                "generate-docs",
		Short:              "Generate output",
		Run:                generateMarkdownRun,
		DisableAutoGenTag:  true,
		DisableSuggestions: true,
	}
	generateArgs struct {
		outputDir             string
		linkFileExt           string
		replaceUnderscoreWith string
	}
)

const fmTemplate = `---
layout: default
description: %s
title: %s
---
`

func init() {
	RootCmd.AddCommand(GenerateCmd)
	f := GenerateCmd.Flags()
	f.StringVarP(&generateArgs.outputDir, "output-dir", "o", "./docs", "Output directory")
	f.StringVarP(&generateArgs.linkFileExt, "link-file-ext", "l", "", "What file extensions the links should point to")
	f.StringVarP(&generateArgs.replaceUnderscoreWith, "replace-underscore-with", "r", "", "Replace the underscore in links with the given character")
}

func generateMarkdownRun(c *cobra.Command, args []string) {
	// Validate arguments
	log := CLILog
	cargs := generateArgs

	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		command := strings.Replace(base, "_", " ", -1)
		title := strings.Title(command)
		description := "Description of the " + command + " command"
		return fmt.Sprintf(fmTemplate, description, title)
	}
	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		extension := ".md"
		if cargs.linkFileExt != "" {
			extension = cargs.linkFileExt
		}
		baseName := base
		if cargs.replaceUnderscoreWith != "" {
			baseName = strings.ReplaceAll(base, "_", cargs.replaceUnderscoreWith)
		}
		return baseName + extension
	}

	if _, err := os.Stat(cargs.outputDir); os.IsNotExist(err) {
		log.Fatal().Err(err).Str("output", cargs.outputDir).Msg("Output directory does not exist.")
	}

	err := doc.GenMarkdownTreeCustom(RootCmd, cargs.outputDir, filePrepender, linkHandler)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to generate document")
	}
}
