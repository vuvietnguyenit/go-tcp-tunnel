package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var scriptCheckPort = `
	#!/bin/sh
	port=%s
	cmd=$(ss -tlnp | awk '{print $4}' | grep :$port)
	if [ -n "$cmd" ]; then
		# Port is running
		echo "PORT_IS_RUNNING"
	else
		# Port is not running:
		echo "PORT_IS_NOT_RUNNING"
	fi
	`

var scriptForwardPort = `
	#!/bin/sh
	source=%s
	dest=%s
	desc=%s
	socat TCP4-LISTEN:$source,fork,reuseaddr TCP4:$dest > /dev/null 2>&1 &
`

func init() {
	// Check socat exists
	var out strings.Builder

	cmdCheckSocatExists := exec.Command("socat", "-V")
	cmdCheckSocatExists.Stdout = &out
	err := cmdCheckSocatExists.Run()
	if err != nil {
		log.Fatal("socat doesn't seem to be installed on your machine, please install socat first")
	}

}

type ForwardData struct {
	Name       string `json:"name"`
	SourcePort string `json:"source_port"`
	Dest       string `json:"dest"`
}

func procesFile(filename string) *os.File {
	inputFileJson, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return inputFileJson
}

// Parser dataset to struct
func datasetParser(fileData io.Reader) []ForwardData {
	byteValue, _ := ioutil.ReadAll(fileData)
	var listForwardData []ForwardData
	err := json.Unmarshal(byteValue, &listForwardData)
	if err != nil {
		log.Fatal(err)
	}
	return listForwardData
}

// Function to get TCP forward data input from JSON file
func getInput() string {
	if len(os.Args) < 2 {
		log.Fatal("forwarding input is not provide")
	}
	filename := os.Args[1]
	// Validate dataset path
	_, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatalf("file input: %v is not exists \n", filename)
	}

	if err != nil {
		log.Fatal(err)
	}
	return filename
}

// Function to check port number is running is Linux system
func executeCommandCheckRunningPort(port string) error {
	// Execute command bash
	var outputCode strings.Builder
	commandCreated := fmt.Sprintf(scriptCheckPort, port)
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = strings.NewReader(commandCreated)
	cmd.Stdout = &outputCode
	err := cmd.Run()
	if err != nil {
		return err
	}
	result := outputCode.String()
	if result == "PORT_IS_RUNNING\n" {
		errmsg := fmt.Sprintf("Port %s is already used. \n", port)
		return errors.New(errmsg)
	}
	return nil
}

// Function to check port is valid and using socat to
// execute command line to action TCP forward
func startForward(forwardData ForwardData) {
	log.Printf("Start creating forward from %s -> :%s [%s]",
		forwardData.Dest,
		forwardData.SourcePort,
		forwardData.Name,
	)
	err := executeCommandCheckRunningPort(forwardData.SourcePort)
	if err != nil {
		log.Printf("Error: %v \n", err)
		return
	}
	// Create shell command
	commandCreated := fmt.Sprintf(scriptForwardPort, forwardData.SourcePort, forwardData.Dest, forwardData.Name)

	cmd := exec.Command("/bin/sh")
	cmd.Stdin = strings.NewReader(commandCreated)
	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}
	log.Printf("-> Forward [%s] done.", forwardData.Name)

}

func main() {
	// Get input args
	fileInput := getInput()

	// open and read content file
	f := procesFile(fileInput)
	defer f.Close()

	// Start parsing
	data := datasetParser(f)
	for _, el := range data {
		startForward(el)
	}
}
