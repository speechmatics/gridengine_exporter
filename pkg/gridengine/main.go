package gridengine

import (
	"encoding/xml"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	xmltypes "github.com/speechmatics/gridengine_exporter/pkg/xml"
)

type State string

const (
	Running               State = "running"
	Error                 State = "error"
	Disabled              State = "disabled"
	Suspended             State = "suspended"
	Orphaned              State = "orphaned"
	ConfigAmbiguous       State = "configuration ambiguous"
	LoadThresholdAlarm    State = "load threshold alarm"
	SuspendThresholdAlarm State = "suspend threshold alarm"
	CalenderDisabled      State = "disabled by calender"
	CalenderSuspend       State = "suspended by calender"
	Unknown               State = "unknown"
)

type JobsMap map[string][]Job

type Job struct {
	Number    int
	Name      string
	Owner     string
	State     string
	Slots     int
	Priority  float32
	StartTime time.Time
}

type Queue struct {
	Name          string
	Type          string
	Slots         int
	UsedSlots     int
	ReservedSlots int
	State         State
	Jobs          []Job
}

type Host struct {
	Hostname     string
	Architecture string
	Processors   int
	Sockets      int
	Cores        int
	Threads      int
	LoadAvg      float32
	TotalMemory  uint64
	UsedMemory   uint64
	TotalSwap    uint64
	UsedSwap     uint64
	Queues       map[string]Queue
}

func stateStrToState(state string) State {
	stateMap := map[string]State{
		"":  Running,
		"E": Error,
		"d": Disabled,
		"s": Suspended,
		"o": Orphaned,
		"c": ConfigAmbiguous,
		"a": LoadThresholdAlarm,
		"A": SuspendThresholdAlarm,
		"D": CalenderDisabled,
		"C": CalenderSuspend,
		"u": Unknown,
	}

	return stateMap[state]
}

func processQueues(qqueues []xmltypes.HostQueue, hostname string, jobs JobsMap) map[string]Queue {
	queues := make(map[string]Queue)

	for _, qqueue := range qqueues {
		qfullName := qqueue.Name + "@" + hostname

		queue := Queue{
			Name: qqueue.Name,
			Jobs: jobs[qfullName],
		}

		for _, property := range qqueue.Properties {
			switch property.Property {
			case "qtype_string":
				queue.Type = property.Value
			case "slots":
				slots, _ := strconv.Atoi(property.Value)
				queue.Slots = slots
			case "slots_used":
				slots_used, _ := strconv.Atoi(property.Value)
				queue.UsedSlots = slots_used
			case "slots_resv":
				slots_resv, _ := strconv.Atoi(property.Value)
				queue.ReservedSlots = slots_resv
			case "state_string":
				state := stateStrToState(property.Value)
				queue.State = state
			}
		}

		queues[queue.Name] = queue
	}

	return queues
}

func processJobs(qjobs []xmltypes.Job) []Job {
	jobs := make([]Job, len(qjobs))

	for i, qjob := range qjobs {

		job := Job{
			Number:   qjob.Number,
			Name:     qjob.Name,
			Owner:    qjob.Owner,
			State:    qjob.State,
			Slots:    qjob.Slots,
			Priority: qjob.Priority,
		}

		jobs[i] = job
	}

	return jobs
}

func GetHostQueuesJobs() (map[string]Host, []Job, error) {
	qhostRawXml, err := exec.Command("qhost", "-q", "-xml").Output()
	if err != nil {
		return nil, nil, fmt.Errorf("qhost returned non zero: %v", err)
	}

	var qhost xmltypes.Qhost
	err = xml.Unmarshal(qhostRawXml, &qhost)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing qhost xml: %v", err)
	}

	qstatRawXml, err := exec.Command("qstat", "-u", "*", "-q", "*", "-xml", "-f").Output()
	if err != nil {
		return nil, nil, fmt.Errorf("qstat returned non zero: %v", err)
	}

	var qstat xmltypes.Qstat
	err = xml.Unmarshal(qstatRawXml, &qstat)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing qstat xml: %v", err)
	}

	jobs := make(JobsMap)

	for _, qqueue := range qstat.Queues {
		jobs[qqueue.FullName] = processJobs(qqueue.Jobs)
	}

	pendingJobs := make([]Job, len(qstat.Jobs))

	hosts := make(map[string]Host)

	for _, qhost := range qhost.Hosts {
		host := Host{
			Hostname: qhost.Name,
		}

		for _, property := range qhost.Properties {
			switch property.Property {
			case "arch_string":
				host.Architecture = property.Value
			case "num_proc":
				num_proc, _ := strconv.Atoi(property.Value)
				host.Processors = num_proc
			case "m_socket":
				m_socket, _ := strconv.Atoi(property.Value)
				host.Sockets = m_socket
			case "m_core":
				m_core, _ := strconv.Atoi(property.Value)
				host.Cores = m_core
			case "m_thread":
				m_thread, _ := strconv.Atoi(property.Value)
				host.Threads = m_thread
			case "load_avg":
				load_avg, _ := strconv.ParseFloat(property.Value, 32)
				host.LoadAvg = float32(load_avg)
			}
		}

		host.Queues = processQueues(qhost.Queues, host.Hostname, jobs)
		hosts[host.Hostname] = host
	}

	return hosts, pendingJobs, nil
}
