// Copyright  OpenTelemetry Authors
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

package k8seventsreceiver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer/consumertest"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig"
)

func TestDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()
	rCfg, ok := cfg.(*Config)
	require.True(t, ok)

	require.Equal(t, &Config{
		ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
		InitialLookback:  60 * time.Second,
		APIConfig: k8sconfig.APIConfig{
			AuthType: k8sconfig.AuthTypeServiceAccount,
		},
	}, rCfg)
}

func TestFactoryType(t *testing.T) {
	require.Equal(t, config.Type("k8s_events"), NewFactory().Type())
}

func TestCreateReceiver(t *testing.T) {
	rCfg := createDefaultConfig().(*Config)

	// Fails with bad K8s Config.
	r, err := createLogsReceiver(
		context.Background(), componenttest.NewNopReceiverCreateSettings(),
		rCfg, consumertest.NewNop(),
	)
	require.Error(t, err)
	require.Nil(t, r)

	// Override for test.
	rCfg.makeClient = func(apiConf k8sconfig.APIConfig) (k8s.Interface, error) {
		return fake.NewSimpleClientset(), nil
	}
	r, err = createLogsReceiver(
		context.Background(),
		componenttest.NewNopReceiverCreateSettings(),
		rCfg, consumertest.NewNop(),
	)
	require.NoError(t, err)
	require.NotNil(t, r)
}
