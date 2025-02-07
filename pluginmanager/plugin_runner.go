// Copyright 2022 iLogtail Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pluginmanager

import (
	"github.com/alibaba/ilogtail"
)

type pluginCategory string

const (
	pluginMetricInput  pluginCategory = "MetricInput"
	pluginServiceInput pluginCategory = "ServiceInput"
	pluginProcessor    pluginCategory = "Porcessor"
	pluginAggregator   pluginCategory = "Aggregator"
	pluginFlusher      pluginCategory = "Flusher"
)

type PluginRunner interface {
	Init(inputQueueSize int, aggrQueueSize int) error

	Initialized() error

	ReceiveRawLog(log *ilogtail.LogWithContext)

	AddPlugin(pluginName string, category pluginCategory, plugin interface{}, config map[string]interface{}) error

	Run()

	RunPlugins(category pluginCategory, control *ilogtail.AsyncControl)

	Merge(p PluginRunner)

	Stop(exit bool) error
}
