package main

import (
	"fmt"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-openvpn-check/openvpn"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"time"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	StatusFile        string
	StatusFileAgeCrit uint
	StatusFileAgeWarn uint
	MinClientsCrit    uint
	MinClientsWarn    uint
}

var parseFileFn = openvpn.ParseFile

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-openvpn-check",
			Short:    "OpenVPN server status check for Sensu",
			Keyspace: "sensu.io/plugins/sensu-openvpn-check/config",
			Timeout:  15000,
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "status-file",
			Env:       "OPENVPN_STATUS_FILE",
			Argument:  "status-file",
			Shorthand: "f",
			Default:   "",
			Usage:     "The OpenVPN status file",
			Value:     &plugin.StatusFile,
		},
		&sensu.PluginConfigOption[uint]{
			Path:      "status-file-age-crit",
			Env:       "OPENVPN_STATUS_FILE_AGE_CRIT",
			Argument:  "status-file-age-crit",
			Shorthand: "",
			Default:   180,
			Usage:     "The OpenVPN status file age threshold for critical",
			Value:     &plugin.StatusFileAgeCrit,
		},
		&sensu.PluginConfigOption[uint]{
			Path:      "status-file-age-warn",
			Env:       "OPENVPN_STATUS_FILE_AGE_WARN",
			Argument:  "status-file-age-warn",
			Shorthand: "",
			Default:   120,
			Usage:     "The OpenVPN status file age threshold for warning",
			Value:     &plugin.StatusFileAgeWarn,
		},
		&sensu.PluginConfigOption[uint]{
			Path:      "min-clients-crit",
			Env:       "OPENVPN_MIN_CLIENTS_CRIT",
			Argument:  "min-clients-crit",
			Shorthand: "",
			Default:   0,
			Usage:     "The OpenVPN minimum clients threshold for critical",
			Value:     &plugin.MinClientsCrit,
		},
		&sensu.PluginConfigOption[uint]{
			Path:      "min-clients-warn",
			Env:       "OPENVPN_MIN_CLIENTS_WARN",
			Argument:  "min-clients-warn",
			Shorthand: "",
			Default:   0,
			Usage:     "The OpenVPN minimum clients threshold for warning",
			Value:     &plugin.MinClientsWarn,
		},
	}
)

const (
	severityWarning  = "warning"
	severityCritical = "critical"
)

type clientCountThresholdError struct {
	count     uint
	threshold uint
	severity  string
}

func (c *clientCountThresholdError) Error() string {
	return fmt.Sprintf("number of connection lower than %s threshold (%d < %d)", c.severity, c.count, c.threshold)
}

type fileAgeThresholdError struct {
	age       float64
	threshold float64
	severity  string
}

func (f *fileAgeThresholdError) Error() string {
	return fmt.Sprintf("file older than %s threshold (%.2f > %.2f)", f.severity, f.age, f.threshold)
}

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(_ *corev2.Event) (int, error) {
	if len(plugin.StatusFile) == 0 {
		return sensu.CheckStateCritical, fmt.Errorf("--status-file or OPENVPN_STATUS_FILE environment variable is required")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(_ *corev2.Event) (int, error) {
	status, err := parseFileFn(plugin.StatusFile)
	if err != nil {
		return sensu.CheckStateCritical, err
	}

	fmt.Println(status.ClientCount)

	// Verify file age
	fileAgeSeconds := time.Since(status.LastModified).Seconds()
	if plugin.StatusFileAgeCrit > 0 && fileAgeSeconds > float64(plugin.StatusFileAgeCrit) {
		return sensu.CheckStateCritical, &fileAgeThresholdError{fileAgeSeconds, float64(plugin.StatusFileAgeCrit), severityCritical}
	} else if plugin.StatusFileAgeWarn > 0 && fileAgeSeconds > float64(plugin.StatusFileAgeWarn) {
		return sensu.CheckStateWarning, &fileAgeThresholdError{fileAgeSeconds, float64(plugin.StatusFileAgeWarn), severityWarning}
	}

	// Verify the min number of client connections
	if plugin.MinClientsCrit > 0 && status.ClientCount < plugin.MinClientsCrit {
		return sensu.CheckStateCritical, &clientCountThresholdError{status.ClientCount, plugin.MinClientsCrit, severityCritical}
	} else if plugin.MinClientsWarn > 0 && status.ClientCount < plugin.MinClientsWarn {
		return sensu.CheckStateWarning, &clientCountThresholdError{status.ClientCount, plugin.MinClientsWarn, severityWarning}
	}

	return sensu.CheckStateOK, nil
}
