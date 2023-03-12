package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

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
func executeCommandCheckRunningPort(port string) bool {
	var out strings.Builder
	cmdCheckPortRunning := exec.Command("./sh/check_port_running.sh", port)
	cmdCheckPortRunning.Stdout = &out
	err := cmdCheckPortRunning.Run()
	if err != nil {
		log.Println(err)
		// Port is running
		return true
	}
	return false
}

// Function to check port is valid and using socat to
// execute command line to action TCP forward
func startForward(forwardData ForwardData) {
	log.Printf("Start creating forward from %s -> :%s [%s]",
		forwardData.Dest,
		forwardData.SourcePort,
		forwardData.Name,
	)
	portIsRunning := executeCommandCheckRunningPort(forwardData.SourcePort)
	if portIsRunning {
		log.Printf("port %s is already use", forwardData.SourcePort)
		return
	}
	cmd := exec.Command("./sh/socat_forward.sh", forwardData.SourcePort, forwardData.Dest)
	cmd.Stdout = os.Stdout
	err := cmd.Start()
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
