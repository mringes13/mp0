package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/montanaflynn/stats"
)

var RoutineCount int
var pingRes []PingReturn

// PingReturn is the return instance of a shell ping command. 
// Its fields describe the website name, whether the ping was successful, and the latency in milliseconds.
type PingReturn struct {
	website string  // The site toward which ping command is executed
	success bool    // Success or failure of the ping
	latency float64 // The network latency of the website. If ping is unsuccessful, latency=-1
}

// MaxParallelism returns the max GOMAXPROCS value, by comparing runtime.GOMAXPROCS
// with number of CPU threads, and returning larger value
func MaxParallelism() int {
	maxProcess := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcess < numCPU {
		return maxProcess
	}
	return numCPU
}

/*
	@input website // A single instance of website to be pinged
		   resultsChannel // An unbuffered channel that stores a PingReturn instance
	The ping() function executes the shell ping command once on the given website 
	and sends the a string of the standard output into resultsChannel.
*/
func ping(website string, resultsChannel chan PingReturn) {
	var pingCmd *exec.Cmd
	var grepCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		pingCmd = exec.Command("ping", "-n", "1", website)
		grepCmd = exec.Command("findstr", "time=")
	} else {
		pingCmd = exec.Command("ping", "-c", "1", website)
		grepCmd = exec.Command("grep", "time=")
	}
	pingOut := runCommand(pingCmd)
	grepCmd.Stdin = strings.NewReader(pingOut)
	grepOut := runCommand(grepCmd)
	if grepOut == "" {
		resultsChannel <- PingReturn{website, false, -1}
	} else {
		latencyString := strings.TrimSpace(grepOut[strings.Index(grepOut, "time=")+5 : strings.Index(grepOut, "ms")])
		latencyFloat, err := strconv.ParseFloat(latencyString, 64)
		checkError(err)
		resultsChannel <- PingReturn{website, true, latencyFloat}
	}
}

/*
	@input gmp // Sets the environment variable GOMAXPROCS
		   websites // A list of websites to be pinged
	initiatePingRoutines() iteratively sets GOMAXPROCS values 
	and calls go-routines to ping websites, given in the string array,
	in parallel.
*/
func initiatePingRoutines(gmp int, websites []string) {
	runtime.GOMAXPROCS(gmp)
	resultsChannel := make(chan PingReturn)

	// For each valid website entered, ping the website, send/receive the data to/from the channel, save the data for later output
	for i := 0; i < RoutineCount; i++ {
		if websites[i%len(websites)] != "" {
			go ping(websites[i%len(websites)], resultsChannel)
		}
	}
	for i := 0; i < RoutineCount; {
		select {
		case x, ok := <-resultsChannel:
			if ok {
				pingRes = append(pingRes, x)
				i++
			} else {
				fmt.Printf("No value here %d\n", i)
				break
			}
		default:
			continue
		}
		if i+1 >= RoutineCount {
			break
		}
	}
}

/*
	@input gmpToRuntime // A map of coordinates where x=gmp (GOMAXPROCS) value and y=runtime in microseconds
	plot() writes the html file that constructs the gmp vs. runtime graph
*/
func plot(gmpToRuntime map[int]int64) {
	var keys []string
	var values []opts.ScatterData
	for keyValue := 1; keyValue <= len(gmpToRuntime); keyValue++ {
		keyString := strconv.Itoa(keyValue)
		keys = append(keys, keyString)
		values = append(values, opts.ScatterData{Value: gmpToRuntime[keyValue]})
	}

	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title:   "GOMAXPROCS vs. Program Run Time",
				SubLink: "https://github.com/mringes13/mp0#readme",
			},
		),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "#CPU Threads",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Time(μs)",
		}),
	)

	// Put data into instance
	scatter.SetXAxis(keys)
	scatter.AddSeries("Category A", values)
	scatter.SetSeriesOptions(charts.WithLabelOpts(
		opts.Label{
			Show:     true,
			Position: "right",
		}),
	)
	f, _ := os.Create("gomaxprocsvsruntime.html")
	err := scatter.Render(f)
	checkError(err)

}

/*
	@input: cmd // A shell command to be run
	@output: The command prompt output
	This function executes the input command and returns its output
*/
func runCommand(cmd *exec.Cmd) string {
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return out.String()
}

/*
	@input: err // An error instance to be checked
*/
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// main function
func main() {
	var err error
	if len(os.Args) > 1 {
		// Converts user-input arguments from string to int
		RoutineCount, err = strconv.Atoi(os.Args[1])
		if err != nil {
			RoutineCount = 100
		}
	} else {
		RoutineCount = 100
	}
	var gmpToRuntime = make(map[int]int64) // Create an int slice to map GOMAXPROCS values to runtime in nanoseconds.

	gmp := MaxParallelism() // Set the maximum GOMAXPROCS value possible on the client's system.

	// Initializing data structures to store user-input website names.
	var input string
	var inputSplice []string

	// Saving user-input websites to an string slice, to be fed into initiatePingRoutines() function later.
	fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("Hello! Please enter the websites you would like to ping; each separated with a space. Otherwise, enter 'q' to quit! ")
	in := bufio.NewReader(os.Stdin)
	input, _ = in.ReadString('\n')
	input = strings.TrimSuffix(input, "\n")
	for input == "" {
		fmt.Printf("Please enter the websites you would like to ping; each separated with a space. Otherwise, enter 'q' to quit! ")
		in = bufio.NewReader(os.Stdin)
		input, _ = in.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
	}
	if input == "q" {
		fmt.Println("You have quit the program by entering 'q'. Goodbye! \n")
		os.Exit(0)
	}
	inputSplice = strings.Split(input, " ")

	// Iterate through every possible value of GOMAXPROCS and run initiatePingRoutines() program for all the user-input websites.
	// Return the runtime of each iteration and store it in the gmpToRuntime map. 
	fmt.Printf("We will now compare the runtime of the program against all possible values of GOMAXPROCS. \n")
	i := 1
	for i < gmp+1 {
		fmt.Printf("# of CPU threads currently being tested: %d \n", i)
		start := time.Now() // Begin timing of ping program at current GOMAXPROCS setting
		initiatePingRoutines(i, inputSplice)
		duration := time.Since(start) // End timing of ping program at current GOMAXPROCS setting
		gmpToRuntime[i] = duration.Microseconds()
		i++
	}
	fmt.Printf("The output for the relation between GOMAXPROCS and the program run time has completed. \n")
	fmt.Printf("A file gomaxprocsvsruntime.html has been created in the current directory with a visual representation of CPU threads versus program run time. \n\n")
	
	// Plot gmpToRunTime -> Output HTML Graph
	plot(gmpToRuntime)

	// Print out a table of latency statistics for each website pinged.
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 7, 8, 3, '\t', 0)
	_, err = fmt.Fprintln(w, "\t")
	checkError(err)
	fmt.Println("The ping statistics for the input websites will now be displayed below.")
	_, err = fmt.Fprintln(w, "WEBSITE\t MIN\t AVG\t MAX\t STDDEV\t TOTAL\t PERCENT\t")
	checkError(err)
	_, err = fmt.Fprintln(w, "-------\t -----\t -----\t -----\t ------\t ------\t -------\t")
	checkError(err)
	for i := range inputSplice {
		succCount := 0
		totCount := 0
		var pingSlice []float64
		for j := range pingRes {
			if pingRes[j].website == inputSplice[i] {
				totCount += 1
				if pingRes[j].success {
					succCount += 1
					pingSlice = append(pingSlice, pingRes[j].latency)
				}
			}
		}
		if len(pingSlice) > 0 {
			minimum, _ := stats.Min(pingSlice)
			maximum, _ := stats.Max(pingSlice)
			means, _ := stats.Mean(pingSlice)
			stddev, _ := stats.StandardDeviation(pingSlice)
			_, err = fmt.Fprintln(w, inputSplice[i], "\t",
				strconv.FormatFloat(minimum, 'f', 2, 64), "\t",
				strconv.FormatFloat(means, 'f', 2, 64), "\t",
				strconv.FormatFloat(maximum, 'f', 2, 64), "\t",
				strconv.FormatFloat(stddev, 'f', 2, 64), "\t",
				totCount, "\t",
				strconv.FormatFloat(float64(succCount)/float64(totCount)*100.0, 'f', 2, 64)+"%", "\t")
			checkError(err)
		} else {
			_, err = fmt.Fprintln(w, inputSplice[i], "\t", "NaN", "\t", "NaN", "\t", "NaN", "\t", "NaN", "\t", totCount, "\t", "0.00%", "\t")
		}
	}
	_, err = fmt.Fprintln(w, "\t")
	checkError(err)
	err = w.Flush()
	checkError(err)
}
