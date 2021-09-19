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

const RoutineCount = 500

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
	////Check to make sure website is valid
	//if website == "" {
	//	return
	//}
	//
	////Communicate with user
	//fmt.Printf("Thank you! ~~%s~~ is now being sent packets! Please wait for results.\n", website)
	//
	////Assign the arguments
	//arg1 := "ping"
	//arg2 := website
	//arg3 := "-c 5"
	//arg4 := "-q"
	//
	////Ping the websites using the command line
	//cmd := exec.Command(arg1, arg2, arg3, arg4)
	//stdout, err := cmd.Output()
	//
	////Check for request timeout
	//if strings.Contains(string(stdout), "100.0%") {
	//	resultschannel <- []string{website + " received no packets.", "--", "--", "--", "--"}
	//	return
	//}
	////Check for errors
	//if err != nil {
	//	fmt.Println(err)
	//	if err.Error() == "exit status 68" {
	//		resultschannel <- []string{website + " is not available.", "--", "--", "--", "--"}
	//	}
	//	return
	//}
	//
	//stdoutnew := string(stdout)
	//valuesum := strings.Split(stdoutnew, "= ")[1]
	//splitvalues := strings.Split(valuesum, "/")
	//newstddev := strings.Split(splitvalues[3], " ms")[0]
	//min, avg, max, stddev := splitvalues[0], splitvalues[1], splitvalues[2], newstddev
	//webinfo := []string{website, min, avg, max, stddev}
	//resultschannel <- webinfo

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
			//webinfo := <-resultschannel
			//webinfocollection = append(webinfocollection, webinfo)
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
	//for ping := range resultschannel {
	//	fmt.Println(ping)
	//}

	//Create a table with the output data
	//w := new(tabwriter.Writer)
	//w.Init(os.Stdout, 8, 8, 3, '\t', 0)
	//_, err := fmt.Fprintln(w, "\t")
	//checkError(err)
	//_, err = fmt.Fprintln(w, "WEBSITE\t MIN\t AVG\t MAX\t STDDEV\t")
	//checkError(err)
	//_, err = fmt.Fprintln(w, "-------\t -------\t -------\t -------\t -------\t")
	//checkError(err)
	//for i := range webinfocollection {
	//	_, err = fmt.Fprintln(w, webinfocollection[i][0], "\t", webinfocollection[i][1], "\t", webinfocollection[i][2], "\t", webinfocollection[i][3], "\t", webinfocollection[i][4], "\t")
	//	checkError(err)
	//}
	//_, err = fmt.Fprintln(w, "\t")
	//checkError(err)
	//err = w.Flush()
	//checkError(err)
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
			Name: "CPU THREADS",
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
	var gmpToRuntime = make(map[int]int64) // Create an int slice to map GOMAXPROCS values to runtime in nanoseconds.

	gmp := MaxParallelism() // Set the max GOMAXPROCS value

	//Initializing necessary variables
	var input string
	var inputsplice []string

	//Saving desired websites from a user to variables
	fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("Please enter the websites you would like to ping; each separated with a space. Otherwise, enter 'q' to quit! ")
	fmt.Println("Websites entered, pinging.")
	in := bufio.NewReader(os.Stdin)
	input, _ = in.ReadString('\n')
	input = strings.TrimSuffix(input, "\n")
	//input = os.Args[1]
	if input == "" || input == "q" {
		os.Exit(0)
	}
	inputsplice = strings.Split(input, " ")

	// Iterate through every possible value of GOMAXPROCS and run Pinger program for just the first entered website.
	// Return the runtime of each iteration.
	fmt.Printf("(1) We will now compare the runtime of the program against all possible values of GOMAXPROCS. \n")
	i := 0
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
	_, err := fmt.Fprintln(w, "\t")
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
			_, err = fmt.Fprintln(w, inputsplice[i], "\t", "NaN", "\t", "NaN", "\t", "NaN", "\t", "NaN", "\t", "0.00%", "\t")
		}
	}
	_, err = fmt.Fprintln(w, "\t")
	checkError(err)
	err = w.Flush()
	checkError(err)

	//Call Pinger for each website the user has called
	//fmt.Println("(2) Statistics for each of the input websites will now be completed.")
	//Pinger(i, inputsplice)
}
