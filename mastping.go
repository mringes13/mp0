package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Returns the max GOMAXPROCS value, by comparing runtime.GOMAXPROCS
// with number of CPUs, and returning larger value
func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

// TEAM EDIT
// Executes ping command and sends stats to current channel
func Pinger(gmp int, out chan string) {
	runtime.GOMAXPROCS(gmp)
	go func() {
		stdout, err := exec.Command("ping", "-c3", "google.com").Output() // Execute the ping command
		out <- string(stdout)  // Send output value to the channel
		// Error check for failed ping
		if err != nil {
			log.Fatal(err)
			fmt.Println(err.Error())
			return
		}
		return
	}()
}

// TEAM EDIT
// Gets substring of relevant ping statistics
func get_statistics(s string) string {
	index_start := strings.Index(s, "stddev")
	index_stop := strings.Index(s, "ms")
	s = s[index_start+9:index_stop+3]
	stats := strings.Split(s, "/")
	return stats[1]
}

// Sets the GOMAXPROCS value, and parallelizes ping command while increasing GOMAXPROCS value by one
// in each thread.
func main() {
	var gmpToRuntime = make(map[int]int) // Create an int slice to map GOMAXPROCS values to runtime in nanoseconds.

	gmp := MaxParallelism() // Set the max GOMAXPROCS value
	out := make(chan string) // Create an output channel

	// Iterate through every possible value of GOMAXPROCS and run Pinger program.
	// Return the runtime of each iteration.
	for i := 0; i < gmp + 1; i++ {
		fmt.Printf("Iteration is %d\n", i)
		start := time.Now()
		go Pinger(i, out)
		duration := time.Since(start)
		gmpToRuntime[i] = int(duration)
	}
	close(out) // can you still iterate over this channel after closing it?

	// Iterate through and print all values from the channel
	for c := range out {
		fmt.Println(c)
	}

	// Print out the final map of GOMAXPROCS values to program runtimes
	fmt.Printf("Time map is: \n")
	fmt.Println(gmpToRuntime)
	fmt.Println("Close")

}
