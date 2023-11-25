package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseFile   = "commands.db"
	maxConcurrency = 2
)

// CommandInfo represents information about a command.
type CommandInfo struct {
	ID      int
	Command string
	Status  string
	Log     string
}

func isDockerInstalled() bool {
	cmd := exec.Command("docker", "--version")

	err := cmd.Run()
	return err == nil
}

func main() {
	if !isDockerInstalled() {
		fmt.Println("Docker is not installed. Please install Docker and try again.")
		os.Exit(1)
	}
	// Create a SQLite database
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the commands table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS commands (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			command TEXT NOT NULL,
			status TEXT,
			log TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Create channels for communication
	commandQueue := make(chan string)
	resultQueue := make(chan CommandInfo)
	stop := make(chan struct{})
	var wg sync.WaitGroup

	// Start background process
	go backgroundProcess(commandQueue, resultQueue, stop, &wg, DoneWithItAll)

	// Enqueue some commands for demonstration
	enqueueCommands(commandQueue, []string{"echo 'Command 1'", "echo 'Command 2'", "echo 'Command 3'"})

	// Monitor the result queue for completed commands
	go func() {
		for result := range resultQueue {
			fmt.Printf("Command %d completed with status: %s\n", result.ID, result.Status)
		}
	}()

	// Wait for user to exit
	fmt.Println("Press Enter to stop...")
	fmt.Scanln()

	// Stop the background process
	close(stop)
	wg.Wait()
	fmt.Println("Program terminated.")
	printAllCommands()
}

func printAllCommands() {
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, command, status, log FROM commands")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("All Commands:")
	for rows.Next() {
		var cmdInfo CommandInfo
		if err := rows.Scan(&cmdInfo.ID, &cmdInfo.Command, &cmdInfo.Status, &cmdInfo.Log); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Command: %s, Status: %s, Log: %s\n", cmdInfo.ID, cmdInfo.Command, cmdInfo.Status, cmdInfo.Log)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func backgroundProcess(commandQueue <-chan string, resultQueue chan<- CommandInfo, stop <-chan struct{}, wg *sync.WaitGroup, doneFunc func()) {
	var mu sync.Mutex
	runningCommands := make(map[int]struct{})
	totalCommands := 0
	completedCommands := 0
	semaphore := make(chan struct{}, maxConcurrency)

	for command := range commandQueue {
		semaphore <- struct{}{} // Acquire a semaphore slot
		totalCommands++
		wg.Add(1)
		go func(cmd string) {
			defer func() {
				<-semaphore // Release the semaphore slot
				wg.Done()
				fmt.Println("Finished", cmd)
			}()

			// Execute the command
			output, err := executeCommand(cmd)

			fmt.Println("Starting", cmd)

			// Update the database with the result
			mu.Lock()
			defer mu.Unlock()
			commandID := insertCommand(cmd, output, err)
			delete(runningCommands, commandID)
			completedCommands++

			// Send the result to the result queue
			resultQueue <- CommandInfo{
				ID:      commandID,
				Command: cmd,
				Status:  getStatus(err),
				Log:     output,
			}

			// Check if all commands are completed
			if completedCommands == totalCommands {
				doneFunc()
			}
		}(command)
	}

	// Wait for stop signal
	<-stop
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("docker", "run", "--rm", "alpine", "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func insertCommand(command, output string, err error) int {
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert command information into the database
	result, err := db.Exec(`
		INSERT INTO commands (command, status, log)
		VALUES (?, ?, ?)
	`, command, getStatus(err), output)
	if err != nil {
		log.Println(err)
		return 0
	}

	// Retrieve the ID of the inserted command
	id, _ := result.LastInsertId()
	return int(id)
}

func getStatus(err error) string {
	if err != nil {
		return "Failed"
	}
	return "Success"
}

func enqueueCommands(commandQueue chan<- string, commands []string) {
	for _, cmd := range commands {
		commandQueue <- cmd
		time.Sleep(time.Second) // Delay between enqueuing commands for demonstration
	}
	close(commandQueue)
}

func DoneWithItAll() {
	fmt.Println("Done with all commands!")
	// Add your logic here
}
