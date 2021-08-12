package collector

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/speechmatics/gridengine_exporter/pkg/gridengine"
)

type GridengineCollector struct {
	queueState   *prometheus.Desc
	slotsTotal   *prometheus.Desc
	slotsFree    *prometheus.Desc
	slotsUsed    *prometheus.Desc
	slotsRunning *prometheus.Desc
	slotsPending *prometheus.Desc
	jobsRunning  *prometheus.Desc
	jobsPending  *prometheus.Desc
}

func NewGridengineCollector() *GridengineCollector {
	queueLabels := []string{"queue", "host"}
	userLabels := []string{"user"}
	jobLabels := append(queueLabels, userLabels...)

	return &GridengineCollector{
		queueState: prometheus.NewDesc("grid_queue_state",
			"Shows state a queue is in", queueLabels, nil,
		),
		slotsTotal: prometheus.NewDesc("grid_slots_total",
			"Total slots present in queue", queueLabels, nil,
		),
		slotsFree: prometheus.NewDesc("grid_slots_free",
			"Free slots available in queue", queueLabels, nil,
		),
		slotsUsed: prometheus.NewDesc("grid_slots_used",
			"Used slots in queue", queueLabels, nil,
		),
		slotsRunning: prometheus.NewDesc("grid_slots_running",
			"Slots with a running job in queue", queueLabels, nil,
		),
		slotsPending: prometheus.NewDesc("grid_slots_pending",
			"Slots required by jobs that are pending", nil, nil,
		),
		jobsRunning: prometheus.NewDesc("grid_jobs_running",
			"Jobs running in queue", jobLabels, nil,
		),
		jobsPending: prometheus.NewDesc("grid_jobs_pending",
			"Jobs pending to be run on a queue", userLabels, nil,
		),
	}
}

func (collector *GridengineCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.queueState
	ch <- collector.slotsTotal
	ch <- collector.slotsFree
	ch <- collector.slotsUsed
	ch <- collector.slotsRunning
	ch <- collector.slotsPending
	ch <- collector.jobsRunning
	ch <- collector.jobsPending
}

func (collector *GridengineCollector) Collect(ch chan<- prometheus.Metric) {
	_, pendingJobs, err := gridengine.GetHostQueuesJobs()
	if err != nil {
		log.Printf("failed to get gridengine data: %v", err)
		return
	}

	pendingJobsPerUser := make(map[string]int)
	for _, pendingJob := range pendingJobs {
		pendingJobsPerUser[pendingJob.Owner]++
	}
	for user, pendingJobs := range pendingJobsPerUser {
		ch <- prometheus.MustNewConstMetric(collector.jobsPending, prometheus.GaugeValue, float64(pendingJobs), user)
	}

	// ch <- prometheus.MustNewConstMetric(collector.barMetric, prometheus.CounterValue, metricValue)
}
