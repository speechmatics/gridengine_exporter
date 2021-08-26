package collector

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/speechmatics/gridengine_exporter/pkg/gridengine"
)

type GridengineCollector struct {
	Filter       string
	queueUp      *prometheus.Desc
	queueState   *prometheus.Desc
	slotsTotal   *prometheus.Desc
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
		queueUp: prometheus.NewDesc("grid_queue_up",
			"Indicates wheather grid queue is up", queueLabels, nil,
		),
		queueState: prometheus.NewDesc("grid_queue_state",
			"Shows state a queue is in", append(queueLabels, "state"), nil,
		),
		slotsTotal: prometheus.NewDesc("grid_slots_total",
			"Total slots present in queue", queueLabels, nil,
		),
		slotsRunning: prometheus.NewDesc("grid_slots_running",
			"Slots with a running job in queue", jobLabels, nil,
		),
		slotsPending: prometheus.NewDesc("grid_slots_pending",
			"Slots required by jobs that are pending", userLabels, nil,
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
	ch <- collector.queueUp
	ch <- collector.queueState
	ch <- collector.slotsTotal
	ch <- collector.slotsRunning
	ch <- collector.slotsPending
	ch <- collector.jobsRunning
	ch <- collector.jobsPending
}

func (collector *GridengineCollector) Collect(ch chan<- prometheus.Metric) {
	hosts, pendingJobs, err := gridengine.GetHostQueuesJobs(collector.Filter)
	if err != nil {
		log.Printf("failed to get gridengine data: %v", err)
		return
	}

	for hostName, host := range hosts {
		for queueName, queue := range host.Queues {
			ch <- prometheus.MustNewConstMetric(collector.queueUp, prometheus.GaugeValue, 1, queueName, hostName)
			ch <- prometheus.MustNewConstMetric(collector.queueState, prometheus.GaugeValue, 1, queueName, hostName, string(queue.State))
			ch <- prometheus.MustNewConstMetric(collector.slotsTotal, prometheus.GaugeValue, float64(queue.Slots), queueName, hostName)

			runningJobsPerUser := make(map[string]int)
			runningSlotsPerUser := make(map[string]int)
			for _, job := range queue.Jobs {
				runningJobsPerUser[job.Owner]++
				runningSlotsPerUser[job.Owner] += job.Slots
			}
			for user, jobs := range runningJobsPerUser {
				ch <- prometheus.MustNewConstMetric(collector.jobsRunning, prometheus.GaugeValue, float64(jobs), queueName, hostName, user)
				ch <- prometheus.MustNewConstMetric(collector.slotsRunning, prometheus.GaugeValue, float64(runningSlotsPerUser[user]), queueName, hostName, user)
			}
		}
	}

	pendingJobsPerUser := make(map[string]int)
	pendingSlotsPerUser := make(map[string]int)
	for _, pendingJob := range pendingJobs {
		pendingJobsPerUser[pendingJob.Owner]++
		pendingSlotsPerUser[pendingJob.Owner] += pendingJob.Slots
	}
	for user, pendingJobs := range pendingJobsPerUser {
		ch <- prometheus.MustNewConstMetric(collector.jobsPending, prometheus.GaugeValue, float64(pendingJobs), user)
		ch <- prometheus.MustNewConstMetric(collector.slotsPending, prometheus.GaugeValue, float64(pendingSlotsPerUser[user]), user)
	}
}
