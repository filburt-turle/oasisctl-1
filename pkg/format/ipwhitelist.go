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

package format

import (
	"strings"

	security "github.com/arangodb-managed/apis/security/v1"
)

// IPWhitelist returns a single IP whitelist formatted for humans.
func IPWhitelist(x *security.IPWhitelist, opts Options) string {
	return formatObject(opts,
		kv{"id", x.GetId()},
		kv{"name", x.GetName()},
		kv{"description", x.GetDescription()},
		kv{"cidr-ranges", strings.Join(x.GetCidrRanges(), ", ")},
		kv{"url", x.GetUrl()},
		kv{"created-at", formatTime(opts, x.GetCreatedAt())},
	)
}

// IPWhitelistList returns a list of IP whitelists formatted for humans.
func IPWhitelistList(list []*security.IPWhitelist, opts Options) string {
	return formatList(opts, list, func(i int) []kv {
		x := list[i]
		return []kv{
			kv{"id", x.GetId()},
			kv{"name", x.GetName()},
			kv{"description", x.GetDescription()},
			kv{"cidr-ranges", strings.Join(x.GetCidrRanges(), ", ")},
			kv{"url", x.GetUrl()},
			kv{"created-at", formatTime(opts, x.GetCreatedAt())},
		}
	}, false)
}
