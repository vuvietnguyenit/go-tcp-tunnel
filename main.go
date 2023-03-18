package main

import (
	"embed"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

//go:embed sh/*
var bash embed.FS

//go:embed sh/socat_forward.sh
var scriptSocatForward string

func openBash(filepath string) (fs.File, error) {
	f, err := bash.Open(filepath)
	if err != nil {
		return nil, err
	}
	return f, nil
}

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
	bashCheckPortRunningFile, err := openBash("sh/check_port_running.sh")
	if err != nil {
		return err
	}
	cmd := exec.Command("bash")
	cmd.Stdin = bashCheckPortRunningFile
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
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
		log.Printf("Error %v \n", err)
		return
	}
	cmd := exec.Command(scriptSocatForward, forwardData.SourcePort, forwardData.Dest)
	cmd.Stdout = os.Stdout
	err = cmd.Start()
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
