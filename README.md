## MP0
MP0 is a simple application that pings desired websites, as defined by the user, using Go-Routines and Go-Channels. 
This project was assigned to practice the newly learned Go-Routines and Go-Channels.

## How to Run
### Step 1: Initialize ping of websites
Start the analytical ping with `go run mastping.go`
### Step 2: Interact with Command Line
### Step 3: Open the created HTML file
If input from command line was valid, an HTML file will be created. Open the created HTML file to see the GOMAXPROCS vs. Program Run Time plot with `open bar.html`.


#### User Ping Websites Process
- To quit the program when being asked to input websites, return 'q'.
- To quit the program while running, interact with the system command line (i.e. OS Terminal - Control + C)
- Enter the desired websites to be pinged.

## Screenshots
1. Command Line Interface
<img width="954" alt="Screen Shot 2021-09-19 at 11 13 15 AM" src="https://user-images.githubusercontent.com/60116121/133932785-92ff1f81-7a14-4b6c-8635-5d6292c97616.png">

2. Output GOMAXPROCS vs Runtime Plot
<img width="1051" alt="Screen Shot 2021-09-19 at 11 13 41 AM" src="https://user-images.githubusercontent.com/60116121/133932810-8bb6d754-c2ee-410e-99d4-ea42a23a58fc.png">

## Workflows
### Ping Process
![MP0 Workflow Diagram](https://user-images.githubusercontent.com/60116121/133932682-9a37ebe8-20af-487f-95b2-b4035317fc1b.png)


## References
- The plotting function is a modified version of sample code from [Go E-Charts Examples](https://github.com/go-echarts/examples/blob/master/examples/scatter.go "Go E-Charts Examples").
