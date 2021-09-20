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

type PingReturn struct {
	website string
	success bool
	latency float64
}

// MaxParallelism returns the max GOMAXPROCS value, by comparing runtime.GOMAXPROCS
// with number of CPUs, and returning larger value
func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

// Ping will execute the ping of the websites.
//Input: website, channel
//Output: Ping results are sent to channel
func ping(website string, resultschannel chan PingReturn) {
	pingCmd := exec.Command("ping", "-n", "1", "-a", website)
	pingOut := runCommand(pingCmd)
	var grepCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		grepCmd = exec.Command("findstr", "time=")
	} else {
		grepCmd = exec.Command("grep", "times=")
	}
	grepCmd.Stdin = strings.NewReader(pingOut)
	grepOut := runCommand(grepCmd)
	if grepOut == "" {
		resultschannel <- PingReturn{website, false, -1}
	} else {
		latencyString := grepOut[strings.Index(grepOut, "time=")+5 : strings.Index(grepOut, "ms")]
		latencyFloat, err := strconv.ParseFloat(latencyString, 64)
		checkError(err)
		resultschannel <- PingReturn{website, true, latencyFloat}
	}
}

// Pinger executes the ping command sends stats to current channel
func Pinger(gmp int, websites []string) {
	runtime.GOMAXPROCS(gmp)
	resultschannel := make(chan PingReturn)
	//var webinfocollection [][]string

	//For each valid website entered, ping the website, send/receive the data to/from the channel, save the data for later output
	for i := 0; i < RoutineCount; i++ {
		if websites[i%len(websites)] != "" {
			go ping(websites[i%len(websites)], resultschannel)
		}
	}
	for i := 0; i < RoutineCount; {
		select {
		case x, ok := <-resultschannel:
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

// Plot is a function to plot the statistics of GOMAXPROCS vs program run time
func plot(gmpToRuntime map[int]int64) {
	var keys []string
	// var values []opts.BarData
	var values []opts.ScatterData
	for keyvalue := 0; keyvalue < len(gmpToRuntime); keyvalue++ {
		keystring := strconv.Itoa(keyvalue)
		keys = append(keys, keystring)
		values = append(values, opts.ScatterData{Value: gmpToRuntime[keyvalue]})
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
			Name: "Number of CPU THREADS",
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
	// Where the magic happens
	f, _ := os.Create("bar.html")
	err := scatter.Render(f)
	checkError(err)

}

func runCommand(cmd *exec.Cmd) string {
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return out.String()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Sets the GOMAXPROCS value, and parallelizes ping command while increasing GOMAXPROCS value by one
// in each thread.
func main() {
	var err error
	if len(os.Args) > 1 {
		RoutineCount, err = strconv.Atoi(os.Args[1])
		if err != nil {
			RoutineCount = 100
		}
	} else {
		RoutineCount = 100
	}
	var gmpToRuntime = make(map[int]int64) // Create an int slice to map GOMAXPROCS values to runtime in nanoseconds.

	gmp := MaxParallelism() // Set the max GOMAXPROCS value

	//Initializing necessary variables
	var input string
	var inputsplice []string

	//Saving desired websites from a user to variables
	fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("Please enter the websites you would like to ping; each separated with a space. Otherwise, enter 'q' to quit! ")
	in := bufio.NewReader(os.Stdin)
	input, _ = in.ReadString('\n')
	input = strings.TrimSuffix(input, "\n")
	if input == "" || input == "q" {
		os.Exit(0)
	}
	inputsplice = strings.Split(input, " ")

	// Iterate through every possible value of GOMAXPROCS and run Pinger program for just the first entered website.
	// Return the runtime of each iteration.
	fmt.Printf("(1) We will now compare the runtime of the program against all possible values of GOMAXPROCS. \n")
	i := 1
	for i < gmp+1 {
		fmt.Printf("CPU CORES BEING CURRENTLY TESTED: %d \n", i)
		start := time.Now()
		Pinger(i, inputsplice)
		duration := time.Since(start)
		gmpToRuntime[i] = duration.Microseconds()
		i++
	}
	fmt.Printf("The output for the comparison between GOMAXPROCS and the program run time has completed. \n\n\n")
	//Plot gmpToRunTime -> Output Graph
	plot(gmpToRuntime)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 8, 3, '\t', 0)
	_, err = fmt.Fprintln(w, "\t")
	checkError(err)
	_, err = fmt.Fprintln(w, "WEBSITE\t MIN\t AVG\t MAX\t STDDEV\t TOTAL\t PERCENT\t")
	checkError(err)
	_, err = fmt.Fprintln(w, "-------\t -------\t -------\t -------\t -------\t -------\t -------\t")
	checkError(err)
	for i := range inputsplice {
		succCount := 0
		totCount := 0
		var pingSlice []float64
		for j := range pingRes {
			if pingRes[j].website == inputsplice[i] {
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
			_, err = fmt.Fprintln(w, inputsplice[i], "\t",
				strconv.FormatFloat(minimum, 'f', 2, 64), "\t",
				strconv.FormatFloat(means, 'f', 2, 64), "\t",
				strconv.FormatFloat(maximum, 'f', 2, 64), "\t",
				strconv.FormatFloat(stddev, 'f', 2, 64), "\t",
				totCount, "\t",
				strconv.FormatFloat(float64(succCount)/float64(totCount)*100.0, 'f', 2, 64)+"%", "\t")
			checkError(err)
		} else {
			_, err = fmt.Fprintln(w, inputsplice[i], "\t", "NaN", "\t", "NaN", "\t", "NaN", "\t", "NaN", "\t", totCount, "\t", "0.00%", "\t")
		}
	}
	_, err = fmt.Fprintln(w, "\t")
	checkError(err)
	err = w.Flush()
	checkError(err)
}
