package mapreduce

import (
	//"hash/fnv"
	"fmt"
    "os"
    //"bufio"
	//"strconv"
	"encoding/json"
	//"math"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	// Remember that you've encoded the values in the intermediate files, so you
	// will need to decode them. If you chose to use JSON, you can read out
	// multiple decoded values by creating a decoder, and then repeatedly calling
	// .Decode() on it until Decode() returns an error.
	//
	// You should write the reduced output in as JSON encoded KeyValue
	// objects to a file named mergeName(jobName, reduceTaskNumber). We require
	// you to use JSON here because that is what the merger than combines the
	// output from all the reduce tasks expects. There is nothing "special" about
	// JSON -- it is just the marshalling format we chose to use. It will look
	// something like this:
	//
	// enc := json.NewEncoder(mergeFile)
	// for key in ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Use checkError to handle errors.

	// merge file
	merge_name := mergeName(jobName, reduceTaskNumber) // determine name of merge file using given function
	merge_file, merge_err := os.Create(merge_name) // create merge file
	if merge_err != nil { // catch errors
		fmt.Println(merge_err)
		return
	}
	merge_encoder := json.NewEncoder(merge_file) // create an encoder for the merge file

	// open all intermediate files, process values
	var keys []string // array for keys
	var vals [][]string // array for values
	for i := 0; i < nMap; i++ { // for each of the nMap files
		filename := reduceName(jobName, i, reduceTaskNumber) // grab file name for current i file and this reduceTask's number
		open_file, errFile := os.Open(filename) // open file given name
		if errFile != nil { // catch errors
            fmt.Println(errFile)
            return
        }
		new_decoder := json.NewDecoder(open_file) // create decoder using opened file

		for { // loop until broken
			var kv KeyValue; // keyValue holder
			err := new_decoder.Decode(&kv) // decode keyValue from file (each iteration of loop decodes the next line in intermediate file)
			if err != nil { // catch error, meaning no lines left to read
				break;
			} else { // if no error, a line was read
				// sort values by keys
				idx := get_index(kv.Key, keys) // get the index of the key if already encountered, otherwise is -1
				if idx == -1 { // if first time key is encountered
					// first instance of key
					// update keys
					appended_key := append(keys, kv.Key) // add first instance of key
					keys = appended_key
					// update vals
					var temp_arr []string // add first instance of value
					temp_append := append(temp_arr, kv.Value)
					temp_arr = temp_append

					appended_val := append(vals, temp_arr)
					vals = appended_val
				} else { // otherwise, key has been encountered before
					// key has other values
					// only update vals


					appended_val := append(vals[idx], kv.Value) // append value at location of key in value array
					vals[idx] = appended_val
				}
			}
		}
		
		open_file.Close() // close current intermediate file before next one is opened
	}

	// for each key, reduce and encode to the merge
	for i, k := range keys { // now that we have keys and a list of values matching each key, for each key
		enc_err := merge_encoder.Encode(&KeyValue{k, reduceF(k, vals[i])}) // encode key and reduceF output to a merge file
		if enc_err != nil { // catch error
			fmt.Println(enc_err)
			return
		}
	}

	merge_file.Close() // close merge file

}

func get_index(target string, arr []string) int { // finds whether key exists in array
	for i, str := range arr { // for loop through array to see if key matches anything in array
		if target == str { // if match found, return index of match
			return i
		}
	}
	return -1 // else return -1 (not found)
}
