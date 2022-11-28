package main

import (
	"github.com/sensu/sensu-openvpn-check/openvpn"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func getPluginConfig(minClientsWarn, minClientsCrit, statusFileAgeWarn, statusFileAgeCrit uint) *Config {
	return &Config{
		StatusFileAgeCrit: statusFileAgeCrit,
		StatusFileAgeWarn: statusFileAgeWarn,
		MinClientsCrit:    minClientsCrit,
		MinClientsWarn:    minClientsWarn,
	}
}

func getMockParseFileFn(status *openvpn.Status, err error) func(string) (*openvpn.Status, error) {
	return func(s string) (*openvpn.Status, error) {
		return status, err
	}
}

func TestExecuteCheck(t *testing.T) {
	now := time.Now()
	oneMinuteAgo := now.Add(time.Minute * -1)
	parseError := openvpn.NewParseError("Unable to Parse Status file")

	tests := []struct {
		name           string
		status         openvpn.Status
		parseFileErr   error
		config         *Config
		expectedStatus int
		expectedError  error
	}{
		{
			name: "all good",
			status: openvpn.Status{
				ClientCount:  25,
				RouteCount:   25,
				LastModified: now,
				IsUp:         true,
			},
			config:         getPluginConfig(0, 0, 0, 0),
			expectedStatus: sensu.CheckStateOK,
			expectedError:  nil,
		}, {
			name: "num client warning",
			status: openvpn.Status{
				ClientCount:  25,
				RouteCount:   25,
				LastModified: now,
				IsUp:         true,
			},
			config:         getPluginConfig(30, 20, 0, 0),
			expectedStatus: sensu.CheckStateWarning,
			expectedError:  &clientCountThresholdError{25, 30, severityWarning},
		}, {
			name: "num client critical",
			status: openvpn.Status{
				ClientCount:  25,
				RouteCount:   25,
				LastModified: now,
				IsUp:         true,
			},
			config:         getPluginConfig(40, 30, 0, 0),
			expectedStatus: sensu.CheckStateCritical,
			expectedError:  &clientCountThresholdError{25, 30, severityCritical},
		}, {
			name: "file time warning",
			status: openvpn.Status{
				ClientCount:  25,
				RouteCount:   25,
				LastModified: oneMinuteAgo,
				IsUp:         true,
			},
			config:         getPluginConfig(40, 30, 30, 120),
			expectedStatus: sensu.CheckStateWarning,
			expectedError:  &fileAgeThresholdError{60, 30, severityWarning},
		}, {
			name: "file time critical",
			status: openvpn.Status{
				ClientCount:  25,
				RouteCount:   25,
				LastModified: oneMinuteAgo,
				IsUp:         true,
			},
			config:         getPluginConfig(40, 30, 30, 45),
			expectedStatus: sensu.CheckStateCritical,
			expectedError:  &fileAgeThresholdError{60, 45, severityCritical},
		}, {
			name: "is down",
			status: openvpn.Status{
				ClientCount:  0,
				RouteCount:   0,
				LastModified: time.Time{},
				IsUp:         false,
			},
			parseFileErr:   parseError,
			config:         getPluginConfig(40, 30, 30, 45),
			expectedStatus: sensu.CheckStateCritical,
			expectedError:  parseError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// replace defaults from main
			plugin = *test.config
			parseFileFn = getMockParseFileFn(&test.status, test.parseFileErr)

			status, err := executeCheck(nil)

			if test.expectedError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.expectedStatus, status)
		})
	}
}
