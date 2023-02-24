// Copyright 2021 iLogtail Authors
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

package decoder

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alibaba/ilogtail/helper/decoder/common"
	"github.com/alibaba/ilogtail/helper/decoder/influxdb"
	"github.com/alibaba/ilogtail/helper/decoder/opentelemetry"
	"github.com/alibaba/ilogtail/helper/decoder/prometheus"
	"github.com/alibaba/ilogtail/helper/decoder/pyroscope"
	"github.com/alibaba/ilogtail/helper/decoder/raw"
	"github.com/alibaba/ilogtail/helper/decoder/sls"
	"github.com/alibaba/ilogtail/helper/decoder/statsd"
	"github.com/alibaba/ilogtail/pkg/pipeline/extensions"
)

type Option struct {
	FieldsExtend      bool
	DisableUncompress bool
}

var errDecoderNotFound = errors.New("no such decoder")

// GetDecoder return a new decoder for specific format
func GetDecoder(format string) (extensions.Decoder, error) {
	return GetDecoderWithOptions(format, Option{})
}

func GetDecoderWithOptions(format string, option Option) (extensions.Decoder, error) {
	switch strings.TrimSpace(strings.ToLower(format)) {
	case common.ProtocolSLS:
		return &sls.Decoder{}, nil
	case common.ProtocolPrometheus:
		return &prometheus.Decoder{}, nil
	case common.ProtocolInflux, common.ProtocolInfluxdb:
		return &influxdb.Decoder{FieldsExtend: option.FieldsExtend}, nil
	case common.ProtocolStatsd:
		return &statsd.Decoder{
			Time: time.Now(),
		}, nil
	case common.ProtocolOTLPLogV1:
		return &opentelemetry.Decoder{Format: common.ProtocolOTLPLogV1}, nil
	case common.ProtocolOTLPMetricV1:
		return &opentelemetry.Decoder{Format: common.ProtocolOTLPMetricV1}, nil
	case common.ProtocolOTLPTraceV1:
		return &opentelemetry.Decoder{Format: common.ProtocolOTLPTraceV1}, nil
	case common.ProtocolRaw:
		return &raw.Decoder{DisableUncompress: option.DisableUncompress}, nil

	case common.ProtocolPyroscope:
		return &pyroscope.Decoder{}, nil
	default:
		return extensions.CreateDecoder(format, option)
	}
}

//RegisterDecodersAsExtension register builtin decoders as extension, to allow them available in external plugins
func RegisterDecodersAsExtension() {
	creator := func(protocol string) extensions.Decoder {
		d, err := GetDecoder(protocol)
		if err != nil {
			panic(fmt.Sprintf("failed create decoder for protocol: %s", protocol))
		}
		return d
	}
	extensions.AddDecoderCreator(common.ProtocolSLS, func() extensions.Decoder {
		return creator(common.ProtocolSLS)
	})
	extensions.AddDecoderCreator(common.ProtocolPrometheus, func() extensions.Decoder {
		return creator(common.ProtocolPrometheus)
	})
	extensions.AddDecoderCreator(common.ProtocolInfluxdb, func() extensions.Decoder {
		return creator(common.ProtocolInfluxdb)
	})
	extensions.AddDecoderCreator(common.ProtocolInflux, func() extensions.Decoder {
		return creator(common.ProtocolInflux)
	})
	extensions.AddDecoderCreator(common.ProtocolStatsd, func() extensions.Decoder {
		return creator(common.ProtocolStatsd)
	})
	extensions.AddDecoderCreator(common.ProtocolOTLPLogV1, func() extensions.Decoder {
		return creator(common.ProtocolOTLPLogV1)
	})
	extensions.AddDecoderCreator(common.ProtocolOTLPMetricV1, func() extensions.Decoder {
		return creator(common.ProtocolOTLPMetricV1)
	})
	extensions.AddDecoderCreator(common.ProtocolOTLPTraceV1, func() extensions.Decoder {
		return creator(common.ProtocolOTLPTraceV1)
	})
	extensions.AddDecoderCreator(common.ProtocolRaw, func() extensions.Decoder {
		return creator(common.ProtocolRaw)
	})
	extensions.AddDecoderCreator(common.ProtocolPyroscope, func() extensions.Decoder {
		return creator(common.ProtocolPyroscope)
	})
}
