package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

// Global vars
var (
	wg         sync.WaitGroup
	callsCount = 0       // Number of calls to the TPM function - used to write at log file
	pcrIndex   = int(13) // PCR index - better use a clear one
)

// Reads the input file and sends lines to the lineChan channel
func read(lineChan chan string, done chan bool) {
	// Open input.txt file
	inputFile, err := os.Open("./input.txt")
	if err != nil {
		fmt.Println("Error opening input file:", err)
		close(lineChan)
		done <- true
		return
	}
	defer inputFile.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Reading line: " + line + "\n" + "________________________________")
		lineChan <- line
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	close(lineChan)
	done <- true
}

// Process input - receives input, extends the PCR value calling TPM through Go-TPM and writes the output at outputFile and logFile.
func processInput(rwc io.ReadWriter, outWriter, logWriter *bufio.Writer, lineChan chan string) {
	defer wg.Done()

	for line := range lineChan {
		callsCount++
		pcrValue, err := hex.DecodeString(line)
		if err != nil {
			log.Fatalf("Failed to decode line: %v", err)
		}

		oldPCRValue, err := tpm2.ReadPCR(rwc, pcrIndex, tpm2.AlgSHA256)
		if err != nil {
			log.Fatalf("Can't read PCR: %v", err)
		}
		fmt.Println("OLD PCR: ", oldPCRValue)

		if err = tpm2.PCRExtend(rwc, tpmutil.Handle(pcrIndex), tpm2.AlgSHA256, pcrValue[:], ""); err != nil {
			log.Fatalf("Erro ao executar PCR extend: %v", err)
		}

		newPCRValue, err := tpm2.ReadPCR(rwc, pcrIndex, tpm2.AlgSHA256)
		if err != nil {
			log.Fatalf("Erro ao ler o valor do PCR: %v", err)
		}
		fmt.Println("NEW PCR: ", newPCRValue)
		finalPCR := sha256.Sum256(append(oldPCRValue, pcrValue[:]...))

		fmt.Println("FINAL PCR: ", strings.ToUpper(hex.EncodeToString(finalPCR[:])))

		outWriter.WriteString(strings.ToUpper(hex.EncodeToString(finalPCR[:])) + "\n")
		logWriter.WriteString(fmt.Sprintf("%d PCR: %s | %s\n", callsCount, strings.ToUpper(hex.EncodeToString(finalPCR[:])), time.Now().Format("2006-01-02 15:04:05")))
		fmt.Println("________________________________")
		outWriter.Flush()
		logWriter.Flush()
	}
}

func main() {
	startTime := time.Now().UnixNano() // used to calculate the total execution time

	// Create output file and write at it
	outputFile, err := os.Create("./output.txt")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Create log file and writer to it
	logFile, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	// Create a writer to both files
	outWriter := bufio.NewWriter(outputFile)
	logWriter := bufio.NewWriter(logFile)

	lineChan := make(chan string)
	done := make(chan bool)

	// Open a connection with the TPM.
	rwc, err := tpmutil.OpenTPM("/dev/tpm0")
	if err != nil {
		log.Fatalf("Failed to open TPM: %v", err)
	}
	defer func() {
		if err := rwc.Close(); err != nil {
			log.Fatalf("Failed to close TPM: %v", err)
		}
	}()

	// Start the goroutine to read lines
	go read(lineChan, done)

	wg.Add(1) // Increment the WaitGroup counter

	// Start the goroutine to process lines
	go processInput(rwc, outWriter, logWriter, lineChan)

	<-done // Wait for reading to complete

	wg.Wait() // Wait for all goroutines to finish

	// Flush the writers
	outWriter.Flush()
	logWriter.Flush()

	endTime := time.Now().UnixNano()
	totalTime := endTime - startTime
	fmt.Println("Total execution time (nanoseconds):", totalTime)
}
