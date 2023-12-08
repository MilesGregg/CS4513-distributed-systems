package mapreduce

import (
    "sync"
    "fmt"
)

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
    var ntasks int
    var nios int // number of inputs (for reduce) or outputs (for map)
    switch phase {
    case mapPhase:
        ntasks = len(mr.files)
        nios = mr.nReduce
    case reducePhase:
        ntasks = mr.nReduce
        nios = len(mr.files)
    }

    debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

    // All ntasks tasks have to be scheduled on workers, and only once all of
    // them have been completed successfully should the function return.
    // Remember that workers may fail, and that any given worker may finish
    // multiple tasks.
    //
    // TODO:

    // create wait groups in order to know when specific tasks are done on the scheduler
    var tasks_wait_group sync.WaitGroup
    tasks_wait_group.Add(ntasks) // add the number of taks that are 

    // create all of the tasks before the scheduler
    var tasks []DoTaskArgs // create array of task argumetns 
    for i := 0; i < ntasks; i++ { // iterate through number of tasks and make
		// only need file name for the mapping phase
		var file string
		switch phase {
		case mapPhase: // when phase is equal to the mapPhase 
			file = mr.files[i] // set file to current task filename
		default:
			file = "" // deault set filename to nothing
		}

		//fmt.Println("setup current arguments for task")
		// setup task arguments for the specific task
		var task DoTaskArgs
		task.JobName = mr.jobName
		task.File = file
		task.Phase = phase
		task.TaskNumber = i
		task.NumOtherPhase = nios
		
		appened_task := append(tasks, task) // append task arguments onto the array of tasks
		tasks = appened_task
		//fmt.Println(tasks[i])
    }

    // send all of the created tasks into scheduler in order to complete them
    for _, task := range tasks {
		//fmt.Println("starting up task in scheduler")
		go schedule_current_task(task, &tasks_wait_group, mr.registerChannel)
		// startup thread to schedule the current task
    }

    // after all tasks are sceduled wait for all of them to be done
    tasks_wait_group.Wait()

    debug("Schedule: %v phase done\n", phase)
}

func schedule_current_task(current_task DoTaskArgs,
						   tasks_wait_group *sync.WaitGroup,
						   registerChannel chan string) {
    defer tasks_wait_group.Done() // once task is fully done then mark task as complete
	// run task forever until the DoTask function is completed
    for {
        curr, err := <-registerChannel // get the current registered channel

        if !err { // if there is an error with the channel then stop
            fmt.Printf("error with channel!")
            return // exit
        } else if call(curr, "Worker.DoTask", current_task, new(struct{})) {
            go receiver(registerChannel, curr) // startup the receiver in a go routine
            return // exit
        }
    }
	//return // don't need this
}

func receiver(registerChannel chan string, curr string) {
    registerChannel <- curr // pass the result back to the register channel, must be done in separate function call to avoid race condition
}