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

package cmd

import (
	"github.com/arangodb-managed/oasisctl/pkg/format"
)

// DefaultFormat returns the default value for the output format
func DefaultFormat() string { return envOrDefault("FORMAT", format.DefaultFormat) }

// DefaultOrganization returns the default value for an organization identifier
func DefaultOrganization() string { return envOrDefault("ORGANIZATION", "") }

// DefaultOrganizationInvite returns the default value for an organization-invite identifier
func DefaultOrganizationInvite() string { return envOrDefault("ORGANIZATION_INVITE", "") }

// DefaultProject returns the default value for a project identifier
func DefaultProject() string { return envOrDefault("PROJECT", "") }

// DefaultGroup returns the default value for a group identifier
func DefaultGroup() string { return envOrDefault("GROUP", "") }

// DefaultRole returns the default value for a role identifier
func DefaultRole() string { return envOrDefault("ROLE", "") }

// DefaultCACertificate returns the default value for a CA certificate identifier
func DefaultCACertificate() string { return envOrDefault("CACERTIFICATE", "") }

// DefaultIPWhitelist returns the default value for an IP whitelist identifier
func DefaultIPWhitelist() string { return envOrDefault("IPWHITELIST", "") }

// DefaultProvider returns the default value for a provider identifier
func DefaultProvider() string { return envOrDefault("PROVIDER", "") }

// DefaultRegion returns the default value for a region identifier
func DefaultRegion() string { return envOrDefault("REGION", "") }

// DefaultDeployment returns the default value for a deployment identifier
func DefaultDeployment() string { return envOrDefault("DEPLOYMENT", "") }

// DefaultURL returns the default value for a URL
func DefaultURL() string { return envOrDefault("URL", "") }
