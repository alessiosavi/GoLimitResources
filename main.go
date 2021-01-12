package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	stringutils "github.com/alessiosavi/GoGPUtils/string"
)

func main() {

	programName := flag.String("name", "", "Name of the program that resources will be limited")
	stopTime := flag.Int("stop", 300, "Number of milliseconds to stop the program")
	waitTime := flag.Int("wait", 5000, "Number of milliseconds to wait before a new stop")

	if stringutils.IsBlank(*programName) {
		flag.Usage()
		panic("No program name provided")
	}

	flag.Parse()
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			pause := exec.Command("cmd", "/C", "pssuspend64.exe", "-nobanner", *programName)
			_ = pause.Run()
			time.Sleep(time.Duration(*stopTime) * time.Millisecond)
			resume := exec.Command("cmd", "/C", "pssuspend64.exe", "-r", "-nobanner", *programName)
			_ = resume.Run()
			time.Sleep(time.Duration(*waitTime) * time.Millisecond)
		}
	}()
	sig := <-cancelChan
	log.Printf("Caught SIGTERM %v", sig)
	resume := exec.Command("cmd", "/C", "pssuspend64.exe", "-r", "-nobanner", *programName)
	_ = resume.Run()
}
