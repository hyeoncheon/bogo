# Bogo, a looking glass for Hyeoncheon

[![Test](https://github.com/hyeoncheon/bogo/actions/workflows/test.yml/badge.svg)](https://github.com/hyeoncheon/bogo/actions/workflows/test.yml)
[![DeepSource](https://deepsource.io/gh/hyeoncheon/bogo.svg/?label=active+issues&token=suXiU-8eOt2HTIniLbcCLbq2)](https://deepsource.io/gh/hyeoncheon/bogo/?ref=repository-badge)
[![codecov](https://codecov.io/gh/hyeoncheon/bogo/branch/main/graph/badge.svg?token=TkcqVhww7F)](https://codecov.io/gh/hyeoncheon/bogo)
[![Coverage Status](https://coveralls.io/repos/github/hyeoncheon/bogo/badge.svg?branch=main)](https://coveralls.io/github/hyeoncheon/bogo?branch=main)
[![Maintainability](https://api.codeclimate.com/v1/badges/1a36b1292948783341d0/maintainability)](https://codeclimate.com/github/hyeoncheon/bogo/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/hyeoncheon/bogo)](https://goreportcard.com/report/github.com/hyeoncheon/bogo)
[![Go Reference](https://img.shields.io/badge/go-reference-blue)](https://pkg.go.dev/github.com/hyeoncheon/bogo)
[![GitHub license](https://img.shields.io/github/license/hyeoncheon/bogo)](https://github.com/hyeoncheon/bogo/blob/main/LICENSE.md)

Bogo, a part of the Hyeoncheon project, is a looking glass agent for
diagnosing network connectivity issue.

Bogo was originally designed as a general-purpose event collector and reporter
for the Hyeoncheon project, as a successor of Kyeong which was written in Ruby
programming language, but now the focusing area of Bogo is networking and also
it is written in Go programming language.

Note: when it reports to Google Cloud Monitoring(formerly known as Stackdriver
Monitoring) using `stackdriver` exporter, the feature is only supported when
Bogo is running on a GCE instance for now.


## Features

As a periodic connectivity data reporter:

* Reports ping latency and packet loss to targets periodically.
* Reports its own health status to check if the agent is running.
* Exports collected data to Google Cloud Monitoring.
* Exports collected data to the standard output (mostly for test purposes).
* Targets can be specified with GCE metadata.

The default exporter is the `stackdriver` exporter that reports to Google
Cloud Monitoring. For this, the agent must run on a GCE instance to get proper
information such as project, zone (which indicates the location of the prober)
it runs on, and the connection to Google Cloud Monitoring.
When it runs on a GCE instance, it will collect all necessary information from
the metadata server of the instance, including ping targets, and it makes the
configuration easier.

You can find the exported metrics on the project's monitoring explorer with
the following keys:

* `custom.googleapis.com/bogo/ping/packet_loss`
* `custom.googleapis.com/bogo/ping/rtt_average`

The best display mode is "Group by: mean".


As a looking glass server:

* Handles `/mtr` request to allow users to see the traceroute from the agent
  to the specified destination.

The response for the MTR request is the raw output of the `mtr` command. To
use this feature, you need to install the `mtr` command and make sure it is
executable from the agent process's execution environment.

Currently, Bogo does not support access control (to deny or to allow specific
sources of the requests) by itself so you need to configure firewall rules or
need to place the agent behind any reverse proxy which support access control.


## Architecture and Components

Bogo consists of checkers, exporters, request handlers, and the webserver to
run the handlers.

Checkers are data collectors. it can be expanded to collect any type of data,
however, it does not support external checkers so all checkers are compiled as
a single binary.
Currently, `healthcheck` checker and `ping` checker are supported.

Exporters are metric exporters. They are a kind of interface to the monitoring
systems such as Google Cloud Monitoring.  The technical/internal mechanism of
them is the same as checkers.
For now, `stackdriver` and `stdout` are supported.

Request handlers are basically HTTP request handlers that serve requests from
the users. Checkers and exporters run periodically, but handlers run request
basis.
Supported paths are `/mtr` and `/echo`.

The web server of Bogo is not a full-featured web server. It is just designed
to serve specific, pre-configured requests such as `/mtr` as described above.


## Install

This is a Go program and follows the standard go installation. You can use
`go install` command as the same as the others, also you can build your own
by cloning this repo manually and run `go build` command.


### Requirement

Supported Go versions are 1.16 or above.


### Get and Build

Just get it!

```console
$ go install github.com/hyeoncheon/bogo/cmd/bogo@latest
go: downloading github.com/hyeoncheon/bogo v0.4.0
$ which bogo
/home/sio4/go/bin/bogo
$ bogo --version
bogo v0.4.0
$ 
```


## Setup and Run

Desired running environment is GCE. This section assumes you are running Bogo
on a GCE environment. (This is not a limitation for whole feature of Bogo, but
`stackdriver` exporter, currently, only runs on GCE environment.)


### Runtime Options

You can see the available options by running it with `--help`.

```console
$ bogo --help
bogo v0.4.0
Usage: bogo [-dhv] [-a value] [-c value] [--copts value] [--eopts value] [-e value] [-l value] [parameters ...]
 -a, --address=value
                    webserver's listen address [127.0.0.1:6090]
 -c, --checkers=value
                    set checkers
     --copts=value  checker options
 -d, --debug        debugging mode
     --eopts=value  exporter options
 -e, --exporter=value
                    set exporter [stackdriver]
 -h, --help         show help message
 -l, --log=value    log level. (debug, info, warn, or error) [info]
 -v, --version      show version
$ 
```

Checkers and Exporters (`-c` and `-e`) can be configured with a comma-separated
list of checkers and exporters accordingly. By default, all checkers are
enabled if no `-c` (or `--checkers`) option is provided, and `stackdriver`
exporter will be enabled by default if no `-e` (or `--exporter`) option is
provided.

Note: Bogo does not support multiple exporters at the same time for now even
thought the option accepts list of exporters.


### Plugin Options

Additionally, Bogo plugins (checkers and exporters) support configurable
options depending on each but they are not shown in the help message. The
`--copts` and `--eopts` are used and the supported format is:

```
<plugin_name>:<plugin_option>:<option_value>;...
```

For example, the following command means:

* Enables all checkers (no `-c` option is provided)
* Enables `stackdriver` and `stdout` exporters. (`-e`)
* Configures `ping` exporter with
  * `www.example.com` and `master` as ping targets.
  * 30 seconds of checking interval.
  * 500 milliseconds of ping interval.

```console
$ bogo --copts "ping:targets:www.example.com,master;ping:check_interval:30;ping:ping_interval:500" -e stackdriver,stdout
```

Available plugin options are:

* `heartbeat` checker
  * `interval`: heartbeat interval. (default: 60 seconds)
* `ping` checker
  * `check_interval`: time interval in seconds for checks.
    (default: 30 seconds)
  * `ping_interval`: time interval in milliseconds for ping requests.
    (default: 1000 milliseconds)

Note: ping count is fixed to 10 times, so if you configured `ping_interval` as
1000 milliseconds, the check will spend 10 seconds for each which means you
should not configure `check_interval` less than 10 seconds. Otherwise, Bogo
will ping the target too much.

Note: ping targets could be configured with plugin option but the easiest
way to configure the targets is using GCE metadata. `ping` checker will use
the value provided with the option but if no option is provided for target,
it will check the metadata and uses it. Instance metadata will be checked,
but if there is no instance metadata, then it will check project metadata.
When you runs many probers, using project metadata will be the best option.


### Configure Prober Instance on GCE

To write the monitoring data from GCE instances, using `stackdriver` exporter,
your prober instances should have the following access scope. (This scope is
included in the default scope so no additional job is required if you already
use the default scope.)

* `https://www.googleapis.com/auth/monitoring.write`

Bogo uses https://github.com/go-ping/ping internally and it requires running
the following command on the prober instance. I would suggest you add the
following command to your startup script. (See also bundled `startup` script.
This script can be used as a startup script for prober instances.)

```console
sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"
```


### Configuration Example and Usecase


### Run

Currently, it does not support daemon mode. To run it in background, please
consider the other method such as `nohup` or run by `at` command.
(See also bundled `startup` script)

```console
$ ./bogo -e stdout -c ping --copts "ping:targets:www.google.com,www.youtube.com"
WARN[0000] hey, it seems like I am on a legacy server or unsupported cloud!
INFO[0000] starting checker ping...                      module=checker
INFO[0000] starting exporter stdout...                   module=exporter
INFO[0000] 1 checkers and 1 exporters started
INFO[0000] starting webserver on 127.0.0.1:6090...
INFO[0000] ping www.google.com every 30s                 checker=ping
INFO[0000] ping www.youtube.com every 30s                checker=ping
⇨ http server started on 127.0.0.1:6090
INFO[0039] ping: {www.youtube.com 172.217.175.14 10 0 33.175457ms 290.418808ms 109.783825ms 84.456195ms}  exporter=stdout
INFO[0039] ping: {www.google.com 172.217.26.228 10 0 31.665286ms 294.040755ms 111.926156ms 85.734723ms}  exporter=stdout
^CWARN[0040] signal caught: interrupt
INFO[0040] shutting down webserver...
INFO[0040] webserver closed
INFO[0040] ping checker for www.youtube.com exited       checker=ping
INFO[0040] stdout exporter exited                        exporter=stdout
INFO[0040] ping checker for www.google.com exited        checker=ping
$ 
```


### Run in debugging mode

Bogo supports `-d` flag for debugging purposes. The flag enables debugging
output and you can check some more information about its internal flows.
Also, you can use `stdout` exporter to see which data is collected from the
console.


## Author

Yonghwan SO https://github.com/sio4, http://www.sauru.so


## Copyright (GNU General Public License v3.0)

Copyright 2020-2022 Yonghwan SO

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT
ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program. If not, see <https://www.gnu.org/licenses/>.
