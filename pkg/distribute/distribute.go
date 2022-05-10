// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package distribute

import v2 "github.com/sealerio/sealer/types/api/v2"

// distribute module is handle how to start a private registry, and how to send cloud image on all node
// the registry-admin.yaml specify the domain name or host about registry config
// CloudImage rootfs not contains registry data and config any more
// sealer save will contains registry config, Clusterfile, registry data and init script
/*
   --- /registry
   --- registry-admin.yaml
   --- registry-image.tar
   --- init.sh
   --- registry.yaml
   --- Clusterfile
*/
type Interface interface {
	// Run registry
	Run() error
	Pull(cluster *v2.Cluster) error
}
