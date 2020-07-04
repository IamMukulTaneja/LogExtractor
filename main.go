package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	layout = "2006-01-02T15:04:05.0000Z"
)

//It contains Flags
type data struct {
	fromTime        time.Time
	toTime          time.Time
	logFileLocation string
}

func setFlags(fromTime *string, toTime *string, logFileLoc *string) {
	flag.StringVar(fromTime, "f", "", "From Time [Required]")
	flag.StringVar(toTime, "t", "", "To Time [Required]")
	flag.StringVar(logFileLoc, "i", "", "Log File Name [Required]")

	flag.Parse()

	if *fromTime == "" || *toTime == "" || *logFileLoc == "" {
		flag.PrintDefaults()
		log.Fatal("All the flags are mandatory")

	}
}

// This function prints the content of the file which contains fromTime.
func (d *data) onlyFromTime(fileLoc string) {

	fileLocation := "/" + d.logFileLocation + "/" + fileLoc

	file, err := os.Open(fileLocation)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		timeString := strings.Split(line, ",")[0]
		isoTime, _ := time.Parse(layout, timeString)
		//Condition to print the log line.
		if isoTime.After(d.fromTime) || isoTime.Equal(d.fromTime) {
			fmt.Println(line)
		}
	}

}

// This function prints the content of the file which contains toTime.
func (d *data) onlyToTime(fileLoc string) {

	fileLocation := "/" + d.logFileLocation + "/" + fileLoc

	file, err := os.Open(fileLocation)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		timeString := strings.Split(line, ",")[0]
		isoTime, _ := time.Parse(layout, timeString)
		//Condition to print
		if isoTime.Before(d.toTime) || isoTime.Equal(d.toTime) {
			fmt.Println(line)
		}
	}
}

// This function prints the content of the file which contains both the times.
func (d *data) bothFromAndTo(fileLoc string) {

	fileLocation := "/" + d.logFileLocation + "/" + fileLoc

	file, err := os.Open(fileLocation)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		timeString := strings.Split(line, ",")[0]
		isoTime, _ := time.Parse(layout, timeString)
		//Condition to print
		if (isoTime.Equal(d.fromTime) || isoTime.After(d.fromTime)) && (isoTime.Equal(d.toTime) || isoTime.Before(d.toTime)) {
			fmt.Println(line)
		}
	}
}

func (d *data) findLogFilesRange(logFileLocation string) []string {

	//Open Metadata File
	file, err := os.Open(logFileLocation)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	defer file.Close()

	var files []string

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		lineArray := strings.Split(line, ",")

		fileStartTime, _ := time.Parse(layout, lineArray[1])
		fileEndTime, _ := time.Parse(layout, lineArray[2])

		// If the file contains from time
		if (d.fromTime.Equal(fileStartTime) || d.fromTime.After(fileStartTime)) && (d.fromTime.Equal(fileEndTime) || d.fromTime.Before(fileEndTime)) {
			files = append(files, lineArray[0])

		}

		// If the file contains to time
		if (d.toTime.Equal(fileStartTime) || d.toTime.After(fileStartTime)) && (d.toTime.Equal(fileEndTime) || d.toTime.Before(fileEndTime)) {
			files = append(files, lineArray[0])
		}

		// If the file falls between from and to time
		if d.fromTime.Before(fileStartTime) && d.toTime.After(fileEndTime) {
			files = append(files, lineArray[0])

		}

	}

	return files

}

// It prints the content of the files which falls between to and from time
func (d *data) printAll(logFileName string) {

	fileLocation := "/" + d.logFileLocation + "/" + logFileName
	file, err := os.Open(fileLocation)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func main() {

	var fromTimeString string
	var toTimeString string
	var logFileLocation string

	//Set and Get Flags
	setFlags(&fromTimeString, &toTimeString, &logFileLocation)

	fromTime, err := time.Parse(layout, fromTimeString)
	if err != nil {
		log.Fatal(err)
	}
	toTime, err := time.Parse(layout, toTimeString)
	if err != nil {
		log.Fatal(err)
	}

	//Initializing  struct of the flags
	d := data{
		fromTime:        fromTime,
		toTime:          toTime,
		logFileLocation: logFileLocation,
	}

	//Fetch the fileName which needs to be searched for log printing
	files := d.findLogFilesRange("LogfilesMetadata.txt")

	startFile := files[0]
	endFile := files[len(files)-1]

	if startFile == endFile {
		d.bothFromAndTo(startFile)
	} else {
		for i, val := range files {
			if i == 0 {
				d.onlyFromTime(val)
			} else if i == len(files)-1 {
				d.onlyToTime(val)
			} else {
				d.printAll(val)
			}
		}
	}

}
