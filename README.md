# gridengine_exporter

Prometheus Exporter for Grid Engine.

## Options



## Metrics

| Metric Name          | Type  | Description                                                                                           |
| -------------------- | ----- | ----------------------------------------------------------------------------------------------------- |
| `grid_queue_up`      | Guage | Indicates wheather grid engine is running  on a host and is listening on a queue. It's either 1 or 0. |
| `grid_queue_state`   | Guage | Shows the state a queue is, see the `state` label.                                                    |
| `grid_slots_total`   | Guage | Slots available in a queue.                                                                           |
| `grid_slots_running` | Guage | Slots in a queue that is currently running a job.                                                     |
| `grid_slots_pending` | Guage | Slots that are required for a job that is pending.                                                    |
| `grid_jobs_running`  | Guage | Number of jobs that are currently running in a queue.                                                 |
| `grid_jobs_pending`  | Guage | Number of jobs that are currently pending and waiting to be scheduled on to a queue.                  |
