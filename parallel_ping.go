package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const PacketCount int = 500

var db *sql.DB

func main() {
	var e error
	db, e = sql.Open("sqlite3", "./ping_results.db")
	if e != nil {
		panic(e)
	}
	file, _ := os.Create("test.txt")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Problem with printing")
			log.Fatal(err)
		}
	}(file)
	c := make(chan int)
	//concurrent shell ping
	start := time.Now()
	for i := 0; i < PacketCount; i++ {
		go shellPing("google.com", c)
	}
	for i := 0; i < PacketCount; {
		select {
		case x, ok := <-c:
			if ok {
				fmt.Printf("%d\n", x)
				stmt, err := db.Prepare("INSERT INTO ping_results (time, website, ping) values (DATETIME('now', 'localtime'),?,?)")
				if err != nil {
					log.Fatal(err)
				}
				_, err = stmt.Exec("google.com", x)
				if err != nil {
					log.Fatal(err)
				}
				i++
			} else {
				fmt.Printf("No value here %d\n", i)
				break
			}
		default:
			continue
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("Concurrent shell ping: %s\n", elapsed)
}

func shellPing(url string, c chan int) {
	pinger := exec.Command("ping", "-n", "1", "-a", url)
	pingOut := runCommand(pinger)
	var grep *exec.Cmd
	if runtime.GOOS == "windows" {
		grep = exec.Command("findstr", "time=")
	} else {
		grep = exec.Command("grep", "times=")
	}
	grep.Stdin = strings.NewReader(pingOut)
	grepOut := runCommand(grep)
	latency := grepOut[strings.Index(grepOut, "time=")+5 : strings.Index(grepOut, "ms")]
	actualPing, err := strconv.Atoi(latency)
	if err != nil {
		c <- -1
		//fmt.Println("Problem with converting")
		//log.Fatal(err)
	} else {
		c <- actualPing
	}
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
