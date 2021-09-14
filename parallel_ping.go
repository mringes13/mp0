package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	file, _ := os.Create("test.txt")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Problem with printing")
			log.Fatal(err)
		}
	}(file)
	c := make(chan int, 100)
	for i := 0; i < 100; i++ {
		go shellPing("google.com", c)
	}
	for i := 0; i < 100; i++ {
		_, err := file.WriteString(fmt.Sprintf("%d\n", <-c))
		if err != nil {
			fmt.Println("Problem with writing")
			log.Fatal(err)
		}
	}
}

func shellPing(url string, c chan int) {
	ping := exec.Command("ping", "-n", "1", "-a", url)
	pingOut := runCommand(ping)
	//grep := exec.Command("findstr", "time=")
	var grep *exec.Cmd
	if runtime.GOOS == "windows" {
		grep = exec.Command("findstr", "time=")
	} else {
		grep = exec.Command("grep", "times=")
	}
	grep.Stdin = strings.NewReader(pingOut)
	grepOut := runCommand(grep)
	latency := strings.TrimSpace(strings.Split(grepOut, "=")[1])
	actualPing, err := strconv.Atoi(latency[:len(latency)-2])
	if err != nil {
		fmt.Println("Problem with converting")
		log.Fatal(err)
	}
	c <- actualPing
}

func runCommand(cmd *exec.Cmd) string {
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}
