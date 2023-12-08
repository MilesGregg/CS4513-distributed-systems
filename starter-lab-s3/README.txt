CS4513: A Simple AWS S3 
==================================

Note, this document includes a number of design questions that can help your implementation. We highly recommend that you answer each design question **before** attempting the corresponding implementation.
These questions will help you design and plan your implementation and guide you towards the resources you need.
Finally, if you are unsure how to start the project, we recommend you visit office hours for some guidance on these questions before attempting to implement this project.


Team members
-----------------

1. Joshua Malcarne (jrmalcarne@wpi.edu)
2. Miles Gregg (mgregg@wpi.edu)

Design Questions
------------------

1. When implementing the `CreateBucket(bucketName string,)` function, you need to create the bucket in each node.
   1.1 How will you determine the number of nodes in the S3 system?
   We will determine the number of nodes in the S3 system using the provided helper function: numNodes := getNumberNodes()

   1.2 How will you create the directories?
   We will create the directories by passing a call to os.MkdirAll("nodes/" + strconv.Itoa(i) + "/" + bucketName, os.ModePerm).
   This will create the bucket with the destired name if it does not exist, and if it does already exist it DOES NOT truncate the existing directory

Brief response please!
---------------------

1. When implementing the `WriteNodeFile(nodeIndex int, bucketName string, fileName string, contents []byte, version time.Time,)` function, you need to write the file to the specified location.
   1.1 How will you find the relative path to the file?
   We can find the relative path to the existing file by concatenating the passed parameters into a filepath as so:
   filename_absolute_path := "nodes/" + strconv.Itoa(nodeIndex) + "/" + bucketName + "/" + fileName

   1.2 How will you handle if the file exists?
   If the file exists, we open the file, and write both the file and the version file as normal

   1.3 How will you handle if the file doesn't exist?
   If the file doesn't exist (the open fails), then we create the file before writing the file and the version file as normal

   1.4 How will you keep track of file versions?
   We take the version file (whose path is the filepath + .version as such: filename_absolute_path + ".version") and write time information to it using the
   conversion: version_file.Write([]byte(version.Format(time.RFC3339)))

Brief response please!
---------------------

1. When implementing the `ReadNodeFile(nodeIndex int, bucketName string, fileName string)` function, you need to read the file at the specified location.
   1.1 How will you find the relative path to the file?
   Similar to WriteNodeFile, we can find the relative path to the existing file by concatenating the passed parameters into a filepath as so:
   filename_absolute_path := "nodes/" + strconv.Itoa(nodeIndex) + "/" + bucketName + "/" + fileName

   1.2 How will you get the version of that file?
   We can use the following process and conversion to get the version of that file:
   Scan the version file into an info variable of type []byte: _, err = fmt.Fscan(file_version, &info)
   Convert the []byte variable to a time variable with the following conversion: version, err := time.Parse(time.RFC3339, string(info))

Brief response please!
---------------------

1. When implementing the `RequestWriteFile(bucketName string, fileName string, fileContents []byte)` function, you need to write the file to the specified bucket.
   1.1 How will you ensure that two clients accessing the same file at the same time will behave predictably? In other words, how will you handle scheduling requests?
   We can do this by using a lock variable of type sync.Map and then locking with the following commands:
   fileLock, _ := fileLock.LoadOrStore(fileName, &sync.Mutex{})
	fileLock.(*sync.Mutex).Lock()
   defer fileLock.(*sync.Mutex).Unlock()
   This locks a mutex around each specific fileName, eliminating parallel writes/race conditions by ensuring the clients go in the order of who gets to the mutex
   first. First come, first serve. The file is unlocked after the RequestWriteFile has ended for the client holding the lock (defer statement)

   1.2 How will you choose how many nodes to write to (hint: `src/replication/common.go`)?
   To choose how many nodes to write to, we find 
   numNodes := getNumberNodes()
	numQ := getWriteQuorum()
   numNodes is the total number of nodes that exist and numQ is the number of nodes that the quorum is able to utilize. 
   We then find the minimum of these two numbers min(numNodes, numQ), and that is the number of nodes to write to.

   1.3 How will you choose which nodes to write to?
   If the number of nodes we are writing to is the total number of nodes, then we simply write to all of the nodes. If not, then it is up to us how we would
   like to determine which nodes are written to of min(numNodes, numQ). We could choose randomly, for example, or we could write to whichever nodes contain the
   smallest number of buckets/files.

   1.4 How will you ensure that all copies of the file written in a single request have the same version?
   To ensure that all copies of the file written in a single request have the same version, we save the value returned that time.Now() to a variable t before
   creating all the files. Then, we pass t to each of the files such that t is the same. Otherwise, if we called time.Now() for each file, the versions would be
   different.

Brief response please!
---------------------

1. When implementing the `RequestReadFile(bucketName string, fileName string)` function, you need to read the file at the specified bucket.
   1.1 How will you ensure that two clients accessing the same file at the same time will behave predictably? In other words, how will you handle scheduling requests?
   We will ensure that two clients accessing the same file at the same time will behave predictably/handle scheduling request using the same method as for RequestWriteFile.
   We can do this by using a lock variable of type sync.Map and then locking with the following commands:
   fileLock, _ := fileLock.LoadOrStore(fileName, &sync.Mutex{})
	fileLock.(*sync.Mutex).Lock()
   defer fileLock.(*sync.Mutex).Unlock()
   This locks a mutex around each specific fileName, eliminating parallel reads/race conditions by ensuring the clients go in the order of who gets to the mutex
   first. First come, first serve. The file is unlocked after the RequestReadFile has ended for the client holding the lock (defer statement)
   
   1.2 How will you choose how many nodes to read from (hint: `src/replication/common.go`)?
   Similar to RequestWriteFile, to choose how many nodes to read from, we find 
   numNodes := getNumberNodes()
	numQ := getReadQuorum()
   numNodes is the total number of nodes that exist and numQ is the number of nodes that the quorum is able to utilize. 
   We then find the minimum of these two numbers min(numNodes, numQ), and that is the number of nodes to read from.

   1.3 How will you choose which nodes to read from?
   We can choose which nodes to read from in a couple ways. Assuming the nodes are able to be read from/valid, we can choose at random, or keep a queue of which
   have been read and update the queue as we read nodes that have been in the queue longer (FIFO). The implementation is up to us.

   1.4 How will you determine which version of the file is the newest?
   We can determine which verison of file is newest by comparing the second of the two values returned by ReadNodeFile(checkedNode, bucketName, fileName).
   This value is the version of the file read, and we can compare versions over the loop to find the most recent/up to date version.

   1.5 How will you handle if a node cannot be read from (is faulty)?
   We can handle this by concatenating passed values into the filePath and checking if the file exists quickly using os.Stat:
   nodePath := fmt.Sprintf("nodes/%d/%s/%s", i, bucketName, fileName)
	_, err := os.Stat(nodePath); !os.IsNotExist(err)
   If an error is encountered, the node is faulty and we skip over it. If we went above and beyond, we could also choose to restore this node by creating the file
   based on an existing version, but that is beyond the scope of this homework (and possibly a good addition for the future).

Brief response please!
---------------------

Errata
------
You can earn up to ten points for bugs, errors, and additional test cases. Two points each. 

6. Bug bounty: describe any known errors or bugs of the project writeup and starter code. 
We encountered a bug with the final test case, TestQuorumThreeClientsFiveNodes, where the test incorrectly identified the number of files created.
In our directory, 3 files were properly created, but the test case only identified 2 and failed.
We tried to recreate this bug over the course of ~50-60 runs and cache clears, but we were not able to.
While I'm not entirely sure why this happened, my best guess is that there is a very small race condition somewhere in the test case, and it triggered.
I wish I could provide more detail, but as I said we were unable to replicate the bug.

7. Additional Test Cases: describe any new test cases that you developed and explain why they are needed. 
There are a couple new test cases we can propose for this project:
- Starvation Test: a test to check whether a node is beind starved by a non-terminal race condition. This is needed because it prevents students from simply abusing
   mutex to deal with parallel clients. To do this, we can run many, many clients at once and keep track of time over them. If a client's requests are fulfilled
   multiple times before another file's request is fulfilled once, the requests are not being properly scheduled and the race condition exists.

- Node Utilization Test: this test checks for whether all nodes are being used (i.e. if nodes = 5 and writeQuorum = 3), or if some nodes are not being utilized
   properly by the code. This test could be achieved by having a client write many many files (say, 100) using a write quorum less than the total number of nodes, 
   and then checking the node directories to see if any nodes were no utilized upon the completion of creating all those files.

- Version Check Tests: this test would check to make sure the code correctly identifies the file with the latest version when passed a read request from a client.
   This could be achieved simply by updating a version file at random and then passing a read request to see if this "newest" version is correctly identified.

- Node Repair Test: this test is somewhat beyond the current scope of the project, but could be added as a nice extension. In the real world, when a faulty
   server is encoutnered, the system attempts to recover that server if at all possible (though this is handled by a separate process). We can replicate this
   during the Faulty Node Test by checking whether students repair a node that their read request function has identified as faulty (i.e. whether the student
   recreates the file using existing copies in order to restore the node). This would be a cool addition to this project.

---------------------

Misc 
-------
8. Describe any deviations, if any, from the requirement.
