package main

import (
	"project-go-warmup/cs4513_go_impl"
	"time"
	"fmt"
	"strconv"
)

func main() {
	//TODO: implement me
	// HINT: need to import the cs4513_go_impl package
	// HINT: use the time package for measurement

	//cs4513_go_impl.Sum(5, "q2_test1.txt")

	for i := 0; i < 5; i++ { // iterate through all of the 5 different tests 
		filename := "./cs4513_go_test/q2_test" + strconv.Itoa(i+1) +".txt" // get the current filename
		fmt.Println("Running workload file:", filename) // prinout information
		for j := 0; j < 10; j++ { // run the iterations for the 10 different threads
			start_time := time.Now() // set the start time
			sum_output := cs4513_go_impl.Sum(j+1, filename) // run the Sum() function and get the output
			end_time := time.Now() // set the end time
			delta_time_sec := end_time.Sub(start_time).Seconds() // get the delta time 
			fmt.Printf("result: %d, num of workers: %d, time: %vsec\n", sum_output, j+1, delta_time_sec) // printout results
		}
	}
}
