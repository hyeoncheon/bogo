# Bogo, an event collector/reporter for Hyeoncheon.

[![Go Reference](https://pkg.go.dev/badge/github.com/hyeoncheon/bogo.svg)](https://pkg.go.dev/github.com/hyeoncheon/bogo)

Bogo is an event collector and reporter for the Hyeoncheon project.cw
It will be a successor of Kyeong which was ruby based event collector.
(Basically it was designed for Hyeoncheon project but currently works as
standalone collector)

Currently, it only runs on a GCE(Google Compute Engine) VM when you want
to use integrated exporter for Google Cloud Monitoring. Also, only ping
checker is supported now.



## Feature

* reports ping latency and packet loss from prober vm to specified targets.
  (powered by https://github.com/go-ping/ping)
* exports to standard out
* exports to Google Cloud Monitoring (fka Stackdriver Monitoring)

Currently, the default (not fallback) exporter is the Stackdriver exporter.
For using this exporter, the prober VM should runs on the GCE vm to interact
with Google Cloud Monitoring automatically without any additional exporter
configuration.

You can find the exported metrics on the project's monitoring explorer:

* `custom.googleapis.com/bogo/ping/packet_loss`
* `custom.googleapis.com/bogo/ping/rtt_average`

The best display mode is "Group by: mean".



## Architecture and Components



## Data Flow



## Install

This is a Go program and follows the standard go installation. You can use
`go install` command as the same as the others, also you can build your own
by cloning this repo manually and run `go build` command.


### Requirement

As of now, supported Go version is 1.16. (well tested, but not fully tested)
However, it works fine with go 1.17, and I guess 1.18 also fine.


### Get and Build

Just get it!

```console
$ go install github.com/hyeoncheon/bogo/cmd/bogo@latest
go: downloading github.com/hyeoncheon/bogo v0.1.0
$ which bogo
/home/sio4/go/bin/bogo
$ bogo --version
bogo v0.1.0
$ 
```



## Setup and Run

Since Bogo runs on GCP VM, you need to configure a GCP project and VMs.


### Configure Prober VM on GCP

To write the monitoring data from VM instances, your VMs should have the
following access scope. (This scope is included in the default scope so
no additional job required if you already use the default scope.)

* `https://www.googleapis.com/auth/monitoring.write`

Bogo uses https://github.com/go-ping/ping and it requires running following
command on the prober VM. I would suggest you to add the following command
on your startup script. (See also bundled `startup` script. This script can
be used as startup script for prober vms.)

```console
sudo sysctl -w net.ipv4.ping_group_range="0   2147483647"
```


### Configuration Example and Usecase


### Run

Currently, it does not support daemon mode. To run it in background,
please consider the other method such as `nohup` or run by `at` command.
(See also bundled `startup` script)

```console
$ ./bogo www.google.com www.youtube.com
starting bogo for [www.google.com www.youtube.com]
stackdriver exporter: initialize exporter...
stat: 74.125.70.104 10 10 0 1.004398ms 1.444576ms 3.124352ms 582.238µs
stat: 142.251.120.93 10 10 0 1.562428ms 2.002906ms 3.66068ms 566.932µs
stat: 142.250.152.106 10 10 0 598.31µs 782.081µs 959.154µs 91.505µs
stat: 142.250.1.190 10 10 0 613.655µs 692.17µs 766.469µs 45.128µs
^Csignal caught: interrupt
interrupted!
stackdriver: bye
$ 
```


### Run in debugging mode

Bogo has `-d` flag but currently not working (since debuggins is enabled by
default :-) You can also use standard out exporter for testing purpose.
(`bogo -e stdout`)



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
