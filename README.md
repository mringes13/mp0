## MP0
MP0 is a simple application that pings desired websites, as defined by the user, using Go-Routines and Go-Channels. 
This project was assigned to practice the newly learned Go-Routines and Go-Channels.

## How to Run

### Step 1: Clone Git Repository
Clone the following git repository with `git clone https://github.com/mringes13/mp0.git`.

### Step 2: Initialize ping of websites
Change the current directory to be within the recently cloned folder. Start the analytical ping with `go run shellPing.go`

##### If an error of the following form is triggered: "cannot find package 'github.com/go-echarts/go-echarts/v2/charts'", change module settings with `export GO111MODULE=on`. Then restart the program with `go run shellPing.go`.

##### If the error is not solved, install the following dependencies with the following: 
`go get github.com/go-echarts/go-echarts`

`go get github.com/montanaflynn/stats`

### Step 3: Interact with Command Line
Enter the website domains you wish you ping. You can input as many inputs as you want, but they must be separated by a white space. 
Press ENTER to generate analysis of parallel pings. 

If you wish to quit the program, enter `q`.
### Step 4: Open the created HTML file
If input from command line was valid, an HTML file will be created. Open the created HTML file to see the GOMAXPROCS vs. Program Run Time plot with `open gomaxprocsvsruntime.html`.

## Screenshots
1. Command Line Interface - Valid User Input
![final](https://user-images.githubusercontent.com/60116121/133951208-c88dff0c-a7da-4ef5-9df5-ac0a7542c0db.png)

2. Command Line Interface - User Quit Program
<img width="900" alt="quit" src="https://user-images.githubusercontent.com/60116121/133951221-30d0ffb2-a05d-4ab3-88f7-097d80ee6ac5.png">

3. Output GOMAXPROCS vs Runtime Plot
![runtime_analysis](https://user-images.githubusercontent.com/60116121/134258272-1b832d70-6672-4375-8629-d3ea5661b3be.png)


## Workflow
![workflow](https://user-images.githubusercontent.com/60116121/134258300-e17f0a46-1ab0-4b75-80de-d51604cddbab.png)


## Custom Data Structures
1. PingReturn Struct 
`type PingReturn struct {
	website string
	success bool
	latency float64
}`

2. GOMAXPROCS-to-Runtime Map
`var gmpToRuntime = make(map[int]int64)`

## Exit Code
- `0`: Successful
- `1`: Incorrect command line input format
- `2`: External package function error

## References
- The plotting function is a modified version of sample code from [Go E-Charts Examples](https://github.com/go-echarts/examples/blob/master/examples/scatter.go "Go E-Charts Examples").
- The data collected from the output table invokes functions from [Go/montanaflynn/stats](https://github.com/montanaflynn/stats)
- Our exact implementation of the `MaxParallelism()` function, which identifies and sets GOMAXPROCS to the maximum possible number of CPU Threads, was taken from this [github repository](https://gist.github.com/peterhellberg/5848304).
