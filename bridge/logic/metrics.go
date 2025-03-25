package logic

import "github.com/prometheus/client_golang/prometheus"

const (
	TypeFunc = "func"
	TypeTri  = "tri"
)

var (
	L2WatchCounterCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "reddio",
			Subsystem: "logic_l2_watch",
			Name:      "counter",
			Help:      "Total Counter for l2 watcher",
		},
		[]string{TypeFunc, TypeTri},
	)
)

func init() {
	prometheus.MustRegister(L2WatchCounterCounter)
}
