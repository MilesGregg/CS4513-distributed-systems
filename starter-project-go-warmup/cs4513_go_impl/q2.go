package cs4513_go_impl

import (
	"bufio"
	"io"
	"strconv"
	"math"
	"os"
	"fmt"
)

/*
Do NOT modify function signature.

Sum numbers from channel `nums` and output sum to `out`.
You should only output to `out` once.
*/
func sumWorker(nums chan int, out chan int) {
	// TODO: implement me
	// HINT: use for loop over `nums`

	//fmt.Println("starting up thread")
	// this threads individual sum
	individual_sum := 0
	// sum everything that is inside of the chan 
	for val := range nums {
		//fmt.Println(val)
		individual_sum += val // append the current val 
	}
	//fmt.Println(individual_sum)
	out <- individual_sum // set the individual sum to the out channel
}

/*
Do NOT modify function signature.

Read integers from the file `fileName` and return sum of all values.
This function must launch `num` go routines running `sumWorker` to find the sum of the values concurrently.

You should use `checkError` to handle potential errors.
*/
func Sum(num int, fileName string) int {
	// TODO: implement me
	// HINT: use `readInts` and `sumWorkers`
	// HINT: used buffered channels for splitting numbers between workers

	// final output sum for this function
	final_sum_output := 0

	// read in file from path
	file, errFile := os.Open(fileName) //"./cs4513_go_test/" +     --> for main testing not tests cases!!
	if errFile != nil { // check for error
		fmt.Println(errFile) // printout error
		return -1
	}
	defer file.Close() // close file at the very end of everything

	// read ints now from the file
	array_of_ints, errInts := readInts(file)
	if errInts != nil { // check for error
		fmt.Println(errInts) // printout error
		return -1
	}

	// METHOD 1: currently commented out because method 2 is more efficient
/*
	// put input channel directy into thread first or create array of channels....??????
	for i := 0; i < num; i++ {
		// total size to read
		size := int(math.Ceil(float64(len(array_of_ints)) / float64(num))) // size of the partition
		channel_nums := make(chan int, size) // make channel nums buffer
		channel_sum := make(chan int, 1)
		for j := 0; j < size; j++ {
			//fmt.Println(i*size+j)
			if i*size+j < len(array_of_ints) {
				channel_nums <- array_of_ints[i*size+j] // input current array index into the channel
			} else {
				break
			}
		}
		//fmt.Printf("PARTITION ENDS HERE")

		//fmt.Println("starting worker")
		go sumWorker(channel_nums, channel_sum)
		close(channel_nums)
		//fmt.Println(<-output)
		output_from_channel := <-channel_sum
		//fmt.Println(output_from_channel)
		final_sum_output += output_from_channel
		close(channel_sum)
	}
	//fmt.Println(final_sum_output)
*/

	// METHOD 2:
	// make array of channels for the inputs into the threads
	input_channels_num := make([]chan int, num)

    // put input channel directy into thread first or create array of channels....??????
    for i := 0; i < num; i++ { // iterate through the num of threads
        // total size to read for each thread
        size := int(math.Ceil(float64(len(array_of_ints)) / float64(num))) // size of the partition
        channel_nums := make(chan int, size) // make channel nums buffer
        //channel_sum := make(chan int)
        for j := 0; j < size; j++ {
            //fmt.Println(i*size+j)
            if i*size+j < len(array_of_ints) {
                channel_nums <- array_of_ints[i*size+j] // input current array index into the channel
            } else {
                break // else break
            }
        }

        input_channels_num[i] = channel_nums // append the channel onto the array channels

        //fmt.Printf("END HERE")

        //fmt.Println("starting worker")
        /*go sumWorker(channel_nums, channel_sum)
        close(channel_nums)
        //fmt.Println(<-output)
        output_from_channel := <-channel_sum
        //fmt.Println(output_from_channel)
        final_sum_output += output_from_channel
        close(channel_sum)*/
    }
    //fmt.Println(final_sum_output)

    for i := 0; i < num; i++ { // iterate through the num of threads
        channel_sum := make(chan int, 1) // channel output

        go sumWorker(input_channels_num[i], channel_sum) // run the worker (aka the thread)
        close(input_channels_num[i]) // close the current input channel num

        final_sum_output += <-channel_sum // add onto the final sum output

        //close(channel_sum)
    }

	return final_sum_output // return the final otuput
}

/*
Do NOT modify this function.
Read a list of integers separated by whitespace from `r`.
Return the integers successfully read with no error, or
an empty slice of integers and the error that occurred.
*/
func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var elems []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return elems, err
		}
		elems = append(elems, val)
	}
	return elems, nil
}
