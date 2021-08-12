package xml

import "encoding/xml"

type Job struct {
	XMLName      xml.Name `xml:"job_list"`
	Number       int      `xml:"JB_job_number"`
	Name         string   `xml:"JB_name"`
	Owner        string   `xml:"JB_owner"`
	State        string   `xml:"state"`
	Slots        int      `xml:"slots"`
	Priority     float32  `xml:"JAT_prio"`
	StartTimeStr string   `xml:"JAT_start_time"`
}

type Queue struct {
	XMLName       xml.Name `xml:"Queue-List"`
	FullName      string   `xml:"name"`
	QType         string   `xml:"qtype"`
	SlotsUsed     int      `xml:"slots_used"`
	SlotsReserved int      `xml:"slots_resv"`
	SlotsTotal    int      `xml:"slots_total"`
	LoadAvg       float32  `xml:"load_avg"`
	Architecture  string   `xml:"arch"`
	Jobs          []Job    `xml:"job_list,omitempty"`
}

type Qstat struct {
	XMLName xml.Name `xml:"job_info"`
	Queues  []Queue  `xml:"queue_info>Queue-List"`
	Jobs    []Job    `xml:"job_info>job_list"`
}

type BaseProperty struct {
	Property string `xml:"name,attr"`
	Value    string `xml:",innerxml"`
}

type HostQueueProperty struct {
	XMLName xml.Name `xml:"queuevalue"`
	BaseProperty
}

type HostQueue struct {
	XMLName    xml.Name            `xml:"queue"`
	Name       string              `xml:"name,attr"`
	Properties []HostQueueProperty `xml:"queuevalue"`
}

type HostProperty struct {
	XMLName xml.Name `xml:"hostvalue"`
	BaseProperty
}

type Host struct {
	XMLName    xml.Name       `xml:"host"`
	Name       string         `xml:"name,attr"`
	Properties []HostProperty `xml:"hostvalue"`
	Queues     []HostQueue    `xml:"queue"`
}

type Qhost struct {
	XMLName xml.Name `xml:"qhost"`
	Hosts   []Host   `xml:"host"`
}
