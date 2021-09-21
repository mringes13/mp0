## MP0
MP0 is a simple application that pings desired websites, as defined by the user, using Go-Routines and Go-Channels. 
This project was assigned to practice the newly learned Go-Routines and Go-Channels.

## How to Run
### Step 1: Initialize ping of websites
Start the analytical ping with `go run shellPing.go`

###### If an error of the following form is triggered: "cannot find package 'github.com/go-echarts/go-echarts/v2/charts'", change module settings with `export GO111MODULE=on`. Then restart the program with `go run shellPing.go`.
### Step 2: Interact with Command Line
Enter the website domains you wish you ping. You can input as many inputs as you want, but they must be separated by a white space. 
Press ENTER to generate analysis of parallel pings. 
### Step 3: Open the created HTML file
If input from command line was valid, an HTML file will be created. Open the created HTML file to see the GOMAXPROCS vs. Program Run Time plot with `open gomaxprocsvsruntime.html`.


#### User Ping Websites Process
- To quit the program when being asked to input websites, return 'q'.
- To quit the program while running, interact with the system command line (i.e. OS Terminal - Control + C)
- Enter the desired websites to be pinged.
- To specify the number of GoRoutines, start the program with 'go run github.com/mast/shellPing 10'

## Screenshots
1. Command Line Interface - Valid User Input
![final](https://user-images.githubusercontent.com/60116121/133951208-c88dff0c-a7da-4ef5-9df5-ac0a7542c0db.png)

2. Command Line Interface - User Quit Program
<img width="900" alt="quit" src="https://user-images.githubusercontent.com/60116121/133951221-30d0ffb2-a05d-4ab3-88f7-097d80ee6ac5.png">

3. Output GOMAXPROCS vs Runtime Plot
<img width="856" alt="Screen Shot 2021-09-19 at 9 52 51 PM" src="https://user-images.githubusercontent.com/60116121/133951225-f3efb8e8-1721-4d0d-8f80-73c0cb4aca60.png">


## Workflow
![MP0 Workflow Diagram-4](https://user-images.githubusercontent.com/60116121/133952461-c621afac-5cc9-4e80-a71a-e42e0318dbb5.png)

## Custom Data Structures
1. PingReturn Struct 
`type PingReturn struct {
	website string
	success bool
	latency float64
}`

2. GOMAXPROCS-to-Runtime Map
`var gmpToRuntime = make(map[int]int64)`

## References
- The plotting function is a modified version of sample code from [Go E-Charts Examples](https://github.com/go-echarts/examples/blob/master/examples/scatter.go "Go E-Charts Examples").
- The data collected from the output table invokes functions from [Go/montanaflynn/stats](https://github.com/montanaflynn/stats)
- Our exact implementation of the `MaxParallelism()` function, which identifies and sets GOMAXPROCS to the maximum possible number of CPU Threads, was taken from this [github repository](https://gist.github.com/peterhellberg/5848304).
