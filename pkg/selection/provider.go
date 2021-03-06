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

package selection

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	common "github.com/arangodb-managed/apis/common/v1"
	platform "github.com/arangodb-managed/apis/platform/v1"
)

// MustSelectProvider fetches the provider with given ID, or name and fails if no provider is found.
// If no ID is specified, all providers are fetched and if the user
// is member of exactly 1, that provider is returned.
func MustSelectProvider(ctx context.Context, log zerolog.Logger, id string, organizationID string, platformc platform.PlatformServiceClient) *platform.Provider {
	provider, err := SelectProvider(ctx, log, id, organizationID, platformc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to list providers")
	}
	return provider
}

// SelectProvider fetches the provider with given ID, or name or returns an error if not found.
// If no ID is specified, all providers are fetched and if the user
// is member of exactly 1, that provider is returned.
func SelectProvider(ctx context.Context, log zerolog.Logger, id string, organizationID string, platformc platform.PlatformServiceClient) (*platform.Provider, error) {
	if id == "" {
		list, err := platformc.ListProviders(ctx, &platform.ListProvidersRequest{OrganizationId: organizationID, Options: &common.ListOptions{}})
		if err != nil {
			log.Debug().Err(err).Msg("Failed to list providers")
			return nil, err
		}
		if len(list.Items) != 1 {
			log.Debug().Err(err).Msgf("You're member of %d providers. Please specify one explicitly.", len(list.Items))
			return nil, fmt.Errorf("You're member of %d providers. Please specify one explicitly.", len(list.Items))
		}
		return list.Items[0], nil
	}
	result, err := platformc.GetProvider(ctx, &common.IDOptions{Id: id})
	if err != nil {
		if common.IsNotFound(err) {
			// Try to lookup provider by name
			list, err := platformc.ListProviders(ctx, &platform.ListProvidersRequest{OrganizationId: organizationID, Options: &common.ListOptions{}})
			if err == nil {
				for _, x := range list.Items {
					if x.GetName() == id {
						return x, nil
					}
				}
			}
		}
		log.Debug().Err(err).Str("provider", id).Msg("Failed to get provider")
		return nil, err
	}
	return result, nil
}
