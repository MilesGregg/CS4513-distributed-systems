package mapreduce

import (
	"hash/fnv"
	"fmt"
    "os"
    "bufio"
	//"strconv"
	"encoding/json"
)

// doMap does the job of a map worker: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.
func doMap(
	jobName string, // the name of the MapReduce job
	mapTaskNumber int, // which map task this is
	inFile string,
	nReduce int, // the number of reduce task that will be run ("R" in the paper)
	mapF func(file string, contents string) []KeyValue,
) {
	// TODO:
	// You will need to write this function.
	// You can find the filename for this map task's input to reduce task number
	// r using reduceName(jobName, mapTaskNumber, r). The ihash function (given
	// below doMap) should be used to decide which file a given key belongs into.
	//
	// The intermediate output of a map task is stored in the file
	// system as multiple files whose name indicates which map task produced
	// them, as well as which reduce task they are for. Coming up with a
	// scheme for how to store the key/value pairs on disk can be tricky,
	// especially when taking into account that both keys and values could
	// contain newlines, quotes, and any other character you can think of.
	//
	// One format often used for serializing data to a byte stream that the
	// other end can correctly reconstruct is JSON. You are not required to
	// use JSON, but as the output of the reduce tasks *must* be JSON,
	// familiarizing yourself with it here may prove useful. You can write
	// out a data structure as a JSON string to a file using the commented
	// code below. The corresponding decoding functions can be found in
	// common_reduce.go.
	//
	//   enc := json.NewEncoder(file)
	//   for _, kv := ... {
	//     err := enc.Encode(&kv)
	//
	// Remember to close the file after you have written all the values!
	// Use checkError to handle errors.

	// read in file content
	// open up the incoming file
    // read in file content
    file, errFile := os.Open(inFile) // open given file
    if errFile != nil { // check for file error
        fmt.Println(errFile)
        return
    }
    defer file.Close() // defer closing the file

    it := bufio.NewScanner(file) // scan in file text
    it.Split(bufio.ScanWords) // split text

    file_content := "" // placeholder for scanned in words
    it.Scan() // first scan to handle proper whitespace
    file_content += string(it.Text()) // first process
    for it.Scan() { // scan in the rest
        file_content += "\n" + string(it.Text()) // process the rest
    }

	var all_files = make([]*os.File, nReduce) // array for intermediate files
	var all_encoders = make([]*json.Encoder, nReduce) // array for intermediate encoders
    // create nReduce files
    for i := 0; i < nReduce; i++ { // for each of the intermediate files
        filename := reduceName(jobName, mapTaskNumber, i) // create the file name using given function
        new_file, errFile := os.Create(filename) // create file for mapping
        if errFile != nil { // catch errors
            fmt.Println(errFile)
            return
        }
		new_encoder := json.NewEncoder(new_file) // create encoder for new file
		all_files[i] = new_file // add file to array
		all_encoders[i] = new_encoder // add encoder to array
    }

    // get KeyValue pairs and ihash into nReduce files using .json encoders
    var transport []KeyValue = mapF(inFile, file_content) // create keyValue pairs using mapF function
    for _, keyval := range transport { // for each keyValue pair
		var hash uint32 = ihash(keyval.Key) % uint32(nReduce) // hash the key into one of the intermediate files
		
		err := all_encoders[hash].Encode(&keyval) // encode into file determined by hash
		if err != nil { // catch errors
			fmt.Println(err)
			return
		}
		//fmt.Println("key " + keyval.Key)
		//fmt.Println("val " + keyval.Value)
    }

	// close all files in all_files array
	for _, f := range all_files { // close all intermediate files once finished updating
		f.Close()
	}

}

func ihash(s string) uint32 { // starter code provided hash function
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}