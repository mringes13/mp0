package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
)

func ping(website string, count string, channel chan []string) {

	//Check to make sure website is valid
	if website == "" {
		return
	}

	//Communicate with user
	fmt.Printf("Thank you! ~~%s~~ is now being sent %s packets! Please wait for results.\n", website, count)

	//Assign the arguments
	arg1 := "ping"
	arg2 := website
	arg3 := "-c"
	arg4 := count
	arg5 := "-q"

	//Ping the websites using the command line
	cmd := exec.Command(arg1, arg2, arg3, arg4, arg5)
	stdout, err := cmd.Output()

	//Check for request timeout
	if strings.Contains(string(stdout), "100.0%") {
		channel <- []string{website + " received no packets.", "--", "--", "--", "--"}
		return
	}
	//Check for errors
	if err != nil {
		if err.Error() == "exit status 68" {
			channel <- []string{website + " not available.", "--", "--", "--", "--"}
		}
		return
	}
	//Understand the output of the ping
	understand(string(stdout), website, channel)
}

//Understand & displaying ping outputs
func understand(stdout string, website string, channel chan []string) {
	stdoutnew := string(stdout)
	valuesum := strings.Split(stdoutnew, "= ")[1]
	splitvalues := strings.Split(valuesum, "/")
	newstddev := strings.Split(splitvalues[3], " ms")[0]
	min, avg, max, stddev := splitvalues[0], splitvalues[1], splitvalues[2], newstddev
	webinfo := []string{website, min, avg, max, stddev}
	channel <- webinfo
}

func main() {
	//Initializing necessary variables
	var web1, web2, web3, web4, web5, web6, web7, web8, web9, web10 string
	var packetcount string
	var statvalues = make(chan []string)
	defer close(statvalues)
	var webinfo []string
	var webinfocollection [][]string

	//Saving desired websites from a user to variables
	fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
	fmt.Printf("Please enter up to 10 websites you would like to ping; each separated with a space. Otherwise, enter 'q' to quit! ")
	fmt.Scanf("%s %s %s %s %s %s %s %s %s %s", &web1, &web2, &web3, &web4, &web5, &web6, &web7, &web8, &web9, &web10)
	websites := []string{web1, web2, web3, web4, web5, web6, web7, web8, web9, web10}

	//Understanding input from the user.
	//Does the user want to quit? Yes? End Program. No? Ask how many packets they want to send to each website.
	if web1 != "q" {
		fmt.Printf("Enter the # of packets you would like to send to each valid website: ")
		fmt.Scanln(&packetcount)
		fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")
		fmt.Printf("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -\n")

		//For each valid website entered, ping the website, send/receive the data to/from the channel, save the data for later output
		for i := range websites {
			if websites[i] != "" {
				go ping(websites[i], packetcount, statvalues)
				webinfo = <-statvalues
				webinfocollection = append(webinfocollection, webinfo)
			}
		}

		//Create a table with the output data
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
	} else {
		fmt.Println("You have quit the program. See you soon!")
	}

}
