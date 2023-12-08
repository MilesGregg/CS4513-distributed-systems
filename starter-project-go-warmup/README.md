CS4513: Warmup The Go Basic
==========================

Team members
-----------------

1. Joshua Malcarne (jrmalcarne@wpi.edu)
2. Miles Gregg (mgregg@wpi.edu)

Two Problems
------------------

1. What is the difference between unbuffered and buffered channels? And why do you choose one over the other for this assignment?
An unbuffered channel has no space to store elements, and thus an unbuffer channel is blocked until a receive is encountered and vice versa
A buffered channel has a fixed capacity for storing elements, and thus a buffered channel will fill to capacity and then be blocked until a receive frees some space within the buffer and vice versa.
In this assignment, we choose a buffered channeled for the input (since we send in multiple values) and unbuffered for the channel output (since we send out a single value).

2. Briefly explain how you approached the two problems.

For q1, we approached the problem by scanning in the file lne by line, and for each line (word) checking if it was greater than the required length, and then updating that word's count in a wordcount array. However, we ran into an issue with appending the first wordcounts of each unique word, which we fixed by introducing an if-else for first time cases. Finally, we take our wordcount array and pass it into the provided sorting function before truncating it to the desired return size and returning the truncated array.

For q2, we began by implementing a simple loop-sum for the sumWorker helper function. Then, we handled readins from the file using the helper function readInts to acquire an array of ints, and partitioned the array into num channels. Then, we passed the channels to parallel calls of sumWorker.


Measurement Part
------------------

Your Reference output:

```
root@367e59db242a:~/host-share/starter-project-go-warmup# go build main.go && ./main
Running workload file: ./cs4513_go_test/q2_test1.txt
result: 499500, num of workers: 1, time: 0.007253252sec
result: 499500, num of workers: 2, time: 0.002445625sec
result: 499500, num of workers: 3, time: 0.002443166sec
result: 499500, num of workers: 4, time: 0.002478382sec
result: 499500, num of workers: 5, time: 0.002559388sec
result: 499500, num of workers: 6, time: 0.002304935sec
result: 499500, num of workers: 7, time: 0.002354592sec
result: 499500, num of workers: 8, time: 0.002619868sec
result: 499500, num of workers: 9, time: 0.00238198sec
result: 499500, num of workers: 10, time: 0.002748516sec
Running workload file: ./cs4513_go_test/q2_test2.txt
result: 117652, num of workers: 1, time: 0.007424048sec
result: 117652, num of workers: 2, time: 0.002456796sec
result: 117652, num of workers: 3, time: 0.002548875sec
result: 117652, num of workers: 4, time: 0.00261805sec
result: 117652, num of workers: 5, time: 0.002752651sec
result: 117652, num of workers: 6, time: 0.002439293sec
result: 117652, num of workers: 7, time: 0.00241409sec
result: 117652, num of workers: 8, time: 0.002510693sec
result: 117652, num of workers: 9, time: 0.0023536sec
result: 117652, num of workers: 10, time: 0.002971609sec
Running workload file: ./cs4513_go_test/q2_test3.txt
result: 617152, num of workers: 1, time: 0.007880242sec
result: 617152, num of workers: 2, time: 0.002534228sec
result: 617152, num of workers: 3, time: 0.00274097sec
result: 617152, num of workers: 4, time: 0.002976415sec
result: 617152, num of workers: 5, time: 0.002797853sec
result: 617152, num of workers: 6, time: 0.002433141sec
result: 617152, num of workers: 7, time: 0.002621099sec
result: 617152, num of workers: 8, time: 0.002471347sec
result: 617152, num of workers: 9, time: 0.00288009sec
result: 617152, num of workers: 10, time: 0.002739429sec
Running workload file: ./cs4513_go_test/q2_test4.txt
result: 4995000, num of workers: 1, time: 0.010053213sec
result: 4995000, num of workers: 2, time: 0.004334628sec
result: 4995000, num of workers: 3, time: 0.003849578sec
result: 4995000, num of workers: 4, time: 0.005028523sec
result: 4995000, num of workers: 5, time: 0.004160686sec
result: 4995000, num of workers: 6, time: 0.004024934sec
result: 4995000, num of workers: 7, time: 0.004228636sec
result: 4995000, num of workers: 8, time: 0.003960798sec
result: 4995000, num of workers: 9, time: 0.004187664sec
result: 4995000, num of workers: 10, time: 0.004030889sec
Running workload file: ./cs4513_go_test/q2_test5.txt
result: 49950000, num of workers: 1, time: 0.028454648sec
result: 49950000, num of workers: 2, time: 0.016452426sec
result: 49950000, num of workers: 3, time: 0.016979557sec
result: 49950000, num of workers: 4, time: 0.017543344sec
result: 49950000, num of workers: 5, time: 0.018333911sec
result: 49950000, num of workers: 6, time: 0.017253085sec
result: 49950000, num of workers: 7, time: 0.01809353sec
result: 49950000, num of workers: 8, time: 0.017740952sec
result: 49950000, num of workers: 9, time: 0.017596218sec
result: 49950000, num of workers: 10, time: 0.017696843sec
```

Observations and Explanation:
While the result always remains the same, as the number of workers increases, the total time required for computation decreases. While our measurement doesn't represent this very well, this is the relationship that they hint at as well as the relationship that the assignment example output hints at.
This relationship occurs because as we increase the number of threads, we have more workers for the same number of problems. This means that it we have 100 problems and 1 worker, that 1 worker is doing all 100 problems, whereas 100 problems and 10 workers means that each worker is doing 10 problems, so the work is split significantly between the parallel processes/threads.
I believe that our processers are fast enough that for this small amount of data, the only significant increase in performance is between 1 and  threads, and past that gains plateau. However, for very large amounts of data, the relationship would manifest past 2 threads rather than plateauing.


Errata
------

Describe any known errors, bugs, or deviations from the requirements.