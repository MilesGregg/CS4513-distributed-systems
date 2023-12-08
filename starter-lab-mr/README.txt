CS4513: The MapReduce Library
=======================================

Note, this document includes a number of design questions that can help your implementation. We highly recommend that you answer each design question **before** attempting the corresponding implementation.
These questions will help you design and plan your implementation and guide you towards the resources you need.
Finally, if you are unsure how to start the project, we recommend you visit office hours for some guidance on these questions before attempting to implement this project.


Team members
-----------------

1. Joshua Malcarne (jrmalcarne@wpi.edu)
2. Miles Gregg (mgregg@wpi.edu)

Design Questions
------------------

(2 point) 1. If there are n input files, and nReduce number of reduce tasks , how does the the MapReduce Library uniquely name the intermediate files?

With n input files, the following function will be called to generate a name for each of the nReduce tasks per input file. That is, we end up with
n*nReduce intermediate files, named according to the job name, the current mapTask (n mapTasks), and the reduceTask allocated to the file (nReduce reduceTasks).
This naming follows the function defined below, giving us a name in the form "mrtmp.jobName-mapTask-reduceTask", "mrtmp.jobName-{0,n}-{0,nReduce}"
func reduceName(jobName string, mapTask int, reduceTask int) string {
	return "mrtmp." + jobName + "-" + strconv.Itoa(mapTask) + "-" + strconv.Itoa(reduceTask)
}


(1 point) 2. Following the previous question, for the reduce task r, what are the names of files will it work on?

The reduceTask r will work on all files for which the name contains the reduceTask's number in the name. That is, all files named "mrtmp.jobName-{0,n}-r"
according to the naming conventions described above. Then, the reduceTask will produce files with the names "mrtmp.jobName-res-r".


(1 point) 3. If the submitted mapreduce job name is "test", what will be the final output file's name?

The final output file's name will be mrtmp.test


(2 point) 4. Based on `mapreduce/test_test.go`, when you run the `TestBasic()` function, how many master and workers will be started? And what are their respective addresses and their naming schemes?

----When we run the 'TestBasic()' function (pasted below), 1 master and 2 workers will be started----

Master:
The master's naming scheme is: port("master") for "" //see port function below
The master's address created using the naming scheme is:
Master name: /var/tmp/824-{user_id}/mr{process_id}-master

Workers:
The worker's naming schemes are port("worker"+strconv.Itoa(i)) 
The worker's addresses created using the naming scheme is:
Worker1 name: /var/tmp/824-{user_id}/mr{process_id}-worker1
Worker2 name: /var/tmp/824-{user_id}/mr{process_id}-worker2//see port function below

Port:
For both master and workers, the port naming scheme uses the port function pasted below to generate a unique-ish UNIX-domain socket name in /var/tmp. 
They can't use the current directory since AFS doesn't support UNIX-domain sockets. The output is the name/address generated for the master/worker (whichever it was called for)

----Functions----
func TestBasic(t *testing.T) {
	mr := setup()
	for i := 0; i < 2; i++ {
		go RunWorker(mr.address, port("worker"+strconv.Itoa(i)),
			MapFunc, ReduceFunc, -1)
	}
	mr.Wait()
	check(t, mr.files)
	checkWorker(t, mr.stats)
	cleanup(mr)
}

func port(suffix string) string {
	s := "/var/tmp/824-"
	s += strconv.Itoa(os.Getuid()) + "/"
	os.Mkdir(s, 0777)
	s += "mr"
	s += strconv.Itoa(os.Getpid()) + "-"
	s += suffix
	return s
}


(4 point) 5. In real-world deployments, when giving a mapreduce job, we often start master and workers on different machines (physical or virtual). Describe briefly the protocol that allows master and workers be aware of each other's existence, and subsequently start working together on completing the mapreduce job. Your description should be grounded on the RPC communications.

In real-world deployments, although the master node and worker nodes often start on separate machines of virtual instances, 
we can achieve communication/make allow the master and workers to be aware of each other's existences through RPC communication (Remote Procedure Calls).
This procedure works by first starting the master node which listens on the RPC for incoming connections from worker nodes. As each worker node is
started, it contacts the master over the RPC, providing its IP and additional information such as the channel along which it listens for instructions.
The master node maintains a list of each of the worker nodes and their status (idle, busy, crashed/erred, etc...). 

When a MR job is encountered, the master partitions the input data to the workker nodes for the map task and sends it via the channels. Then, the workers 
receive the map task and process. The internmediate outputs are returned to the master node via channels, which then partitions reduce tasks and sends the 
necessary data to workers via the channels. The workers the process the reduce task and send final output data back to the master over the channel. 
The master combines the final output data and returns the result to the user. 

Throughout this process, the RPCs (can be synchronous or asynchronous) provide reliable and efficient communication between the master and workers,
allowing the master to send requests and data as needed, as well as track the status of workers and receive processed outputs from workers. RPC could use
TCP or UDP network communication protocols, however for this application we would TCP communication becuase of connection based RPC with master and 
workers on different machines.


(2 point) 6. The current design and implementation uses a number of RPC methods. Can you find out all the RPCs and list their signatures? Briefly describe the criteria a method needs to satisfy to be considered a RPC method. (Hint: you can look up at: https://golang.org/pkg/net/rpc/)

----All RPCs and their signatures:----
1. "Worker.doTask"
2. "Worker.Shutdown"

We're pretty sure these are the RPCs since they're exported signatures for function calls, but it's possible that we're incorrect in which case
Call ```func call(srv string, rpcname string, args interface{}, reply interface{}) bool {}``` and Master.Shutdown ```func (mr *Master) Shutdown(_, _ *struct{}) error {}``` are
the correct answers

----Criteria a method needs to satisy to be considered an RPC method:----
In Go, a method needs to satisfy the following criteria to be an RPC method:
1. The method must be exported
2. The method must take at least two exported arguments
3. One of the arguments must be a pointer (passed args often include one pointer to input, one pointer to output)
4. One of the returned values must be an error
5. The method must belong to a struct type that's registered with an RPC instance


Errata
------

Describe any known errors, bugs, or deviations from the requirements.
Runtimes and Notations:
1. Phase 1 takes ~100 seconds to run
2. Phase 2 takes ~13-50 minutes to run depending on the laptop. Phase 2 will pass, but we were only able to do so on a mac (we couldn't figure out why the whitespace wasn't matching on windows, see email with alias/Professor for details)
3. Phase 3 takes ~50 seconds to run
4. Phase 4 takes ~120 seconds to run
---
