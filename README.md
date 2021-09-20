## MP0
MP0 is a simple application that pings desired websites, as defined by the user, using Go-Routines and Go-Channels. 
This project was assigned to practice the newly learned Go-Routines and Go-Channels.

## How to Run
### Step 1: Initialize ping of websites
Start the analytical ping with `go run mastping.go`
### Step 2: Interact with Command Line
Enter the website domains you wish you ping. You can input as many inputs as you want, but they must be separated by a white space. 
Press ENTER to generate analysis of parallel pings. 
### Step 3: Open the created HTML file
If input from command line was valid, an HTML file will be created. Open the created HTML file to see the GOMAXPROCS vs. Program Run Time plot with `open gomaxprocsvsruntime.html`.


#### User Ping Websites Process
- To quit the program when being asked to input websites, return 'q'.
- To quit the program while running, interact with the system command line (i.e. OS Terminal - Control + C)
- Enter the desired websites to be pinged.

## Screenshots
1. Command Line Interface - Valid User Input
<img width="967" alt="Screen Shot 2021-09-19 at 12 23 08 PM" src="https://user-images.githubusercontent.com/60116121/133935098-32439c47-6344-44c0-929f-3d67fc6c6ab0.png">

2. Command Line Interface - User Quit Program
<img width="958" alt="Screen Shot 2021-09-19 at 12 24 18 PM" src="https://user-images.githubusercontent.com/60116121/133935141-89a10fcf-fd0d-41a5-882c-a6337fd2431e.png">

3. Output GOMAXPROCS vs Runtime Plot
<img width="1073" alt="Screen Shot 2021-09-19 at 12 23 46 PM" src="https://user-images.githubusercontent.com/60116121/133935120-2ae8518f-7659-4b5e-8099-60f4877932c3.png">


## Workflow
![MP0 Workflow Diagram](https://user-images.githubusercontent.com/60116121/133932682-9a37ebe8-20af-487f-95b2-b4035317fc1b.png)

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
