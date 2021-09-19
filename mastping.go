package main

import (
	"bufio"
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
)

var displayresults bool = false

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

// Ping will executes the ping of the websites.
//Input: website, channel
//Output: Ping results are sent to channel
func ping(website string, resultschannel chan []string) {
	//Check to make sure website is valid
	if website == "" {
		return
	}

	//Communicate with user
	if displayresults == true {
		fmt.Printf("Thank you! ~~%s~~ is now being sent packets! Please wait for results.\n", website)
	}

	//Assign the arguments
	arg1 := "ping"
	arg2 := website
	arg3 := "-c 5"
	arg4 := "-q"

	//Ping the websites using the command line
	cmd := exec.Command(arg1, arg2, arg3, arg4)
	stdout, err := cmd.Output()

	//Check for request timeout
	if strings.Contains(string(stdout), "100.0%") {
		resultschannel <- []string{website + " received no packets.", "--", "--", "--", "--"}
		return
	}
	//Check for errors
	if err != nil {
		fmt.Println(err)
		if err.Error() == "exit status 68" {
			resultschannel <- []string{website + " is not available.", "--", "--", "--", "--"}
		}
		return
	}

	stdoutnew := string(stdout)
	valuesum := strings.Split(stdoutnew, "= ")[1]
	splitvalues := strings.Split(valuesum, "/")
	newstddev := strings.Split(splitvalues[3], " ms")[0]
	min, avg, max, stddev := splitvalues[0], splitvalues[1], splitvalues[2], newstddev
	webinfo := []string{website, min, avg, max, stddev}
	resultschannel <- webinfo

}

// Pinger executes the ping command sends stats to current channel
func Pinger(gmp int, websites []string) {
	runtime.GOMAXPROCS(gmp)
	resultschannel := make(chan []string)
	webinfocollection := [][]string{}

	//For each valid website entered, ping the website, send/receive the data to/from the channel, save the data for later output
	for i := range websites {
		if websites[i] != "" {
			go ping(websites[i], resultschannel)
			webinfo := <-resultschannel
			webinfocollection = append(webinfocollection, webinfo)
		}
	}

	//Create a table with the output data
	if displayresults == true {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 8, 8, 3, '\t', 0)
		fmt.Fprintln(w, "\t")
		fmt.Fprintln(w, "WEBSITE\t MIN\t AVG\t MAX\t STDDEV\t")
		fmt.Fprintln(w, "-------\t -------\t -------\t -------\t -------\t")
		for i := range webinfocollection {
			fmt.Fprintln(w, webinfocollection[i][0], "\t", webinfocollection[i][1], "\t", webinfocollection[i][2], "\t", webinfocollection[i][3], "\t", webinfocollection[i][4], "\t")
		}
		fmt.Fprintln(w, "\t")
		w.Flush()
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
			Name: "CPU Threads",
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
	f, _ := os.Create("gomaxprocsvsruntime.html")
	scatter.Render(f)

}

// Sets the GOMAXPROCS value, and parallelizes ping command while increasing GOMAXPROCS value by one
// in each thread.
func main() {
	//Initializing necessary variables
	var input string
	var inputsplice []string
	var inputgomaxprocs []string
	var i int = 0
	var gmpToRuntime = make(map[int]int64) // Create an int slice to map GOMAXPROCS values to runtime in microseconds.
	gmp := MaxParallelism()                // Set the max GOMAXPROCS value

	//Saving desired websites from a user to variables
	fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("Please enter the websites you would like to ping; each separated with a space. Otherwise, enter 'q' to quit! ")
	in := bufio.NewReader(os.Stdin)
	input, _ = in.ReadString('\n')
	input = strings.TrimSuffix(input, "\n")
	inputsplice = strings.Split(input, " ")
	if inputsplice[0] != "q" {
		inputgomaxprocs = inputsplice
		// Iterate through every possible value of GOMAXPROCS and run Pinger program for just the first entered website.
		// Return the runtime of each iteration.
		fmt.Printf("(1) We will now compare the runtime of the program against all possible values of GOMAXPROCS. \n")
		for i < gmp+1 {
			fmt.Printf("CPU threads currently beig tested: %d \n", i)
			start := time.Now()
			Pinger(i, inputgomaxprocs)
			duration := time.Since(start)
			gmpToRuntime[i] = duration.Microseconds()
			i++
		}
		fmt.Printf("The output for the comparison between GOMAXPROCS and the program run time has completed. An HTML file with the results has been created. \n\n\n")

		//Plot gmpToRunTime -> Output Graph
		plot(gmpToRuntime)

		//Call Pinger for each website the user has called
		fmt.Println("(2) Statistics for each of the input websites will now be completed.")
		displayresults = true
		Pinger(i, inputsplice)
	} else {
		fmt.Println("You have quit the program. If you wish to rereun, enter 'go run mastping.go' into the command line. Thank you!")
		fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
		fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	}
}
