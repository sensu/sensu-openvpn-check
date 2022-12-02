{{- /* Ignore this text, until templating is ran via [sensu-plugin-tool](https://github.com/sensu/sensu-plugin-tool) the
below badge links will not render */ -}}

[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/sensu/sensu-openvpn-check)
![Go Test](https://github.com/sensu/sensu-openvpn-check/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/sensu/sensu-openvpn-check/workflows/goreleaser/badge.svg)

# Sensu OpenVPN Check

## Table of Contents

<!-- TOC -->
* [Sensu OpenVPN Check](#sensu-openvpn-check)
  * [Table of Contents](#table-of-contents)
  * [Overview](#overview)
  * [Files](#files)
  * [Usage examples](#usage-examples)
    * [Help Output](#help-output)
    * [Environment Variables](#environment-variables)
    * [Output](#output)
    * [Number of Clients Thresholds](#number-of-clients-thresholds)
    * [Status File Age Thresholds](#status-file-age-thresholds)
  * [Configuration](#configuration)
    * [Asset registration](#asset-registration)
    * [Check definition](#check-definition)
  * [Installation from source](#installation-from-source)
  * [Additional notes](#additional-notes)
  * [Contributing](#contributing)

## Overview

The Sensu OpenVPN Check is a [Sensu Check][6] that verifies the status of an OpenVPN server using its status files.
The check outputs the number of active sessions and can optionally evaluate session count thresholds and file age
thresholds to make sure the server is active.

## Files

## Usage examples

### Help Output

```
OpenVPN server status check for Sensu

Usage:
  sensu-openvpn-check [flags]
  sensu-openvpn-check [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -h, --help                        help for sensu-openvpn-check
      --min-clients-crit uint       The OpenVPN minimum clients threshold for critical
      --min-clients-warn uint       The OpenVPN minimum clients threshold for warning
  -f, --status-file string          The OpenVPN status file
      --status-file-age-crit uint   The OpenVPN status file age threshold for critical (default 180)
      --status-file-age-warn uint   The OpenVPN status file age threshold for warning (default 120)

Use "sensu-openvpn-check [command] --help" for more information about a command.
```

### Environment Variables

| Argument               | Environment Variable         |
|------------------------|------------------------------|
| --min-clients-crit     | OPENVPN_MIN_CLIENTS_CRIT     |
| --min-clients-warn     | OPENVPN_MIN_CLIENTS_WARN     |
| --status-file          | OPENVPN_STATUS_FILE          |
| --status-file-age-crit | OPENVPN_STATUS_FILE_AGE_CRIT |
| --status-file-age-warn | OPENVPN_STATUS_FILE_AGE_WARN |

### Output

Unless there's an error the check will output the number of active clients connected to the OpenVPN server.
Additionally,
a warning or error message based on threshold evaluation can be printed on the next line.

### Number of Clients Thresholds

The check can make sure there are a minimum number of clients connected to the OpenVPN server. The `--min-clients-crit`
and `--min-clients-warn`
options can be used to accomplish that. If the actual number of clients is lower than the critical threshold a status of
2 (critical) is returned by the check.
If the actual number of clients is less than the warning threshold a status of 1 (warning) is returned by the check.
In both cases a similar line is also printed on the terminal:

```
Error executing sensu-openvpn-check: error executing check: number of connection lower than critical threshold (13 < 200)
```

By default, the number of clients thresholds are not set.

### Status File Age Thresholds

The check can evaluate the OpenVPN status file age and generate a warning or error status. An older status file
indicates the file
is not getting updated, potentially exposing a problem with the server.

If the file is older than the critical value set with `--status-file-age-crit` a status of 2 (critical) is returned by
the check.
If the file is older than the warning value set with `--status-file-age-warn` a status of 1 (warning) is returned by the
check.
In both cases a line is printed on the terminal:

```
Error executing sensu-openvpn-check: error executing check: file older than critical threshold (211.18 > 180.00)
```

By default, the warning threshold is set to 120 seconds and the critical threshold is set to 180 seconds.

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add sensu/sensu-openvpn-check
```

If you're using an earlier version of sensuctl, you can find the asset on
the [Bonsai Asset Index][https://bonsai.sensu.io/assets/sensu/sensu-openvpn-check].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-openvpn-check
  namespace: default
spec:
  command: sensu-openvpn-check --status-file /var/run/openvpn-status.log
  subscriptions:
    - system
  runtime_assets:
    - sensu/sensu-openvpn-check
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-openvpn-check repository:

```
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[2]: https://github.com/sensu/sensu-plugin-sdk

[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md

[4]: https://github.com/sensu/sensu-openvpn-check/blob/master/.github/workflows/release.yml

[5]: https://github.com/sensu/sensu-openvpn-check/actions

[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/

[7]: https://github.com/sensu/check-plugin-template/blob/master/main.go

[8]: https://bonsai.sensu.io/

[9]: https://github.com/sensu/sensu-plugin-tool

[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
