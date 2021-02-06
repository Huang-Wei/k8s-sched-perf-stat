# k8s-sched-perf-stat

A tool to analyze the result(s) of [Kubernetes Scheduler Integration Performance test](https://github.com/kubernetes/kubernetes/tree/master/test/integration/scheduler_perf).

## Preparation

Firstly, collect the results of a particular scheduler performance test. It's recommended to run it for 10+ times. (This tool will eliminate the outliers automatically.)

```bash
# Run the test for 10 times
for i in {1..10}; do
  make test-integration WHAT=./test/integration/scheduler_perf KUBE_TEST_VMODULE="''" KUBE_TEST_ARGS="-alsologtostderr=false -logtostderr=false -run=^$$ -benchtime=1ns -bench=BenchmarkPerfScheduling/Unschedulable/500Nodes/200InitPods/1000PodsToSchedule -data-items-dir /tmp"
done
```

Then concatenate the results into a file:

```bash
rm -f /tmp/old.txt
ls /tmp/BenchmarkPerfScheduling*.json | while read f; do
  cat $f >> /tmp/old.txt
  echo >> /tmp/old.txt # append newline
done
```

## Usage

You can use this tool to analyze the results against a fixed codebase, by providing one argument:

```shell
$ go run *.go testdata/old1.txt
+--------------------------+------------------------------------+----------+--------+---------+
|          METRIC          |               GROUP                | QUANTILE |  UNIT  | MEDIAN  |
+--------------------------+------------------------------------+----------+--------+---------+
| SchedulingThroughput     | Unschedulable/500Nodes/200InitPods | Average  | pods/s |  133.93 |
| SchedulingThroughput     | Unschedulable/500Nodes/200InitPods | Perc50   | pods/s |   30.00 |
| SchedulingThroughput     | Unschedulable/500Nodes/200InitPods | Perc90   | pods/s |  564.50 |
| SchedulingThroughput     | Unschedulable/500Nodes/200InitPods | Perc99   | pods/s |  564.50 |
| scheduler_e2e_scheduling | Unschedulable/500Nodes/200InitPods | Average  | ms     |  458.46 |
| scheduler_e2e_scheduling | Unschedulable/500Nodes/200InitPods | Perc50   | ms     |   15.53 |
| scheduler_e2e_scheduling | Unschedulable/500Nodes/200InitPods | Perc90   | ms     | 1288.26 |
| scheduler_e2e_scheduling | Unschedulable/500Nodes/200InitPods | Perc99   | ms     | 1749.03 |
| scheduler_pod_scheduling | Unschedulable/500Nodes/200InitPods | Average  | ms     | 3153.33 |
| scheduler_pod_scheduling | Unschedulable/500Nodes/200InitPods | Perc50   | ms     | 3074.85 |
| scheduler_pod_scheduling | Unschedulable/500Nodes/200InitPods | Perc90   | ms     | 4742.60 |
| scheduler_pod_scheduling | Unschedulable/500Nodes/200InitPods | Perc99   | ms     | 5783.46 |
+--------------------------+------------------------------------+----------+--------+---------+
```

Or, by providing two arguments, this tool will compare them and show the diff:

```shell
$ go run *.go testdata/old1.txt testdata/new1.txt
+--------------------------+------------------------------------+----------+--------+---------+---------+---------+
|          METRIC          |               GROUP                | QUANTILE |  UNIT  |   OLD   |   NEW   |  DIFF   |
+--------------------------+------------------------------------+----------+--------+---------+---------+---------+
| SchedulingThroughput     | Unschedulable/500Nodes/200InitPods | Average  | pods/s |  133.93 |  152.68 | +14.00% |
| SchedulingThroughput     | Unschedulable/500Nodes/200InitPods | Perc50   | pods/s |   30.00 |   49.00 | +63.33% |
| SchedulingThroughput     | Unschedulable/500Nodes/200InitPods | Perc90   | pods/s |  564.50 |  515.50 | -8.68%  |
| SchedulingThroughput     | Unschedulable/500Nodes/200InitPods | Perc99   | pods/s |  564.50 |  515.50 | -8.68%  |
| scheduler_e2e_scheduling | Unschedulable/500Nodes/200InitPods | Average  | ms     |  458.46 |  334.55 | -27.03% |
| scheduler_e2e_scheduling | Unschedulable/500Nodes/200InitPods | Perc50   | ms     |   15.53 |   13.64 | -12.14% |
| scheduler_e2e_scheduling | Unschedulable/500Nodes/200InitPods | Perc90   | ms     | 1288.26 | 1004.83 | -22.00% |
| scheduler_e2e_scheduling | Unschedulable/500Nodes/200InitPods | Perc99   | ms     | 1749.03 | 1472.26 | -15.82% |
| scheduler_pod_scheduling | Unschedulable/500Nodes/200InitPods | Average  | ms     | 3153.33 | 3339.24 | +5.90%  |
| scheduler_pod_scheduling | Unschedulable/500Nodes/200InitPods | Perc50   | ms     | 3074.85 | 3346.54 | +8.84%  |
| scheduler_pod_scheduling | Unschedulable/500Nodes/200InitPods | Perc90   | ms     | 4742.60 | 5580.41 | +17.67% |
| scheduler_pod_scheduling | Unschedulable/500Nodes/200InitPods | Perc99   | ms     | 5783.46 | 6009.02 | +3.90%  |
+--------------------------+------------------------------------+----------+--------+---------+---------+---------+
```