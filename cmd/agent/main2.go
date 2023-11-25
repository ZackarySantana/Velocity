package main

import (
	"database/sql"
	"log"
	"os/exec"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseFile   = "commands.db"
	dockerImage    = "alpine"
	maxConcurrency = 5
	pollInterval   = 5 * time.Second // Adjust the polling interval as needed
)

// CommandInfo represents information about a command.
type CommandInfo struct {
	ID      int
	Image   string
	Command string
	Status  string
	Log     string
}

func main() {
	// Continuously poll for new commands
	for {
		processNewCommands()
		time.Sleep(pollInterval)
	}
}

func processNewCommands() {
	// Open the database
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query for new commands
	rows, err := db.Query(`
		SELECT id, image, command
		FROM commands
		WHERE status IS NULL
		LIMIT ?
	`, maxConcurrency)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	// Use a semaphore to limit concurrency
	semaphore := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	// Process each new command asynchronously
	for rows.Next() {
		var cmdInfo CommandInfo
		if err := rows.Scan(&cmdInfo.ID, &cmdInfo.Image, &cmdInfo.Command); err != nil {
			log.Println(err)
			continue
		}

		// Acquire a semaphore slot
		semaphore <- struct{}{}
		wg.Add(1)

		go func(cmdInfo CommandInfo) {
			defer func() {
				// Release the semaphore slot
				<-semaphore
				wg.Done()
			}()

			// Execute the command
			output, err := executeCommand(cmdInfo.Command)

			// Update the database with the result
			_, err = db.Exec(`
				UPDATE commands
				SET status = ?, log = ?
				WHERE id = ?
			`, getStatus(err), output, cmdInfo.ID)
			if err != nil {
				log.Println(err)
			}
		}(cmdInfo)
	}

	// Wait for all asynchronous commands to complete
	wg.Wait()
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("docker", "run", "--rm", dockerImage, "sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func getStatus(err error) string {
	if err != nil {
		return "Failed"
	}
	return "Success"
}
