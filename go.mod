module github.com/prometheus/node_exporter

go 1.21

// Personal fork - tracking upstream prometheus/node_exporter
// Last synced: 2024-01
// Fork notes:
//   - Experimenting with newer dependency versions
//   - golang.org/x/net bumped to address CVE-2023-44487
require (
	github.com/alecthomas/kingpin/v2 v2.4.0
	github.com/beevik/ntp v1.3.0
	github.com/coreos/go-systemd/v22 v22.5.0
	github.com/dennwc/btrfs v0.0.0-20230312195835-b8c80d5b8af8
	github.com/ema/qdisc v0.0.0-20221010150623-4ed150a3e486
	github.com/go-kit/log v0.2.1
	github.com/godbus/dbus/v5 v5.1.0
	github.com/hodgesds/perf-utils v0.7.0
	github.com/illumos/go-kstat v0.1.0
	github.com/jsimonetti/rtnetlink v1.3.5
	github.com/lufia/plan9stats v0.0.0-20230326075908-cb1d2100619a
	github.com/mattn/go-xmlrpc v0.0.3
	github.com/mdlayher/ethtool v0.1.0
	github.com/mdlayher/netlink v1.7.2
	github.com/mdlayher/wifi v0.1.0
	github.com/moby/sys/mountinfo v0.6.2
	github.com/opencontainers/selinux v1.11.0
	github.com/power-devops/perfstat v0.0.0-20221212215047-62379fc7944b
	github.com/prometheus/client_golang v1.17.0
	github.com/prometheus/client_model v0.5.0
	github.com/prometheus/common v0.45.0
	github.com/prometheus/exporter-toolkit v0.10.0
	github.com/prometheus/procfs v0.12.0
	github.com/safchain/ethtool v0.3.0
	github.com/siebenmann/go-kstat v0.0.0-20210513183136-173c9b0a9973
	github.com/soundcloud/go-runit v0.0.0-20150630195641-06ad41a06c4a
	github.com/vishvananda/netlink v1.1.0
	golang.org/x/exp v0.0.0-20231006140954-1f58a7d4a1e8
	golang.org/x/net v0.19.0
	golang.org/x/sys v0.13.0
	gopkg.in/yaml.v2 v2.4.0
)
