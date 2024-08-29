package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <file_to_watch> <script_to_run>")
	}

	fileToWatch := os.Args[1]
	scriptToRun := os.Args[2]

	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Get the absolute path to the file
	fileToWatch, err = filepath.Abs(fileToWatch)
	if err != nil {
		log.Fatal(err)
	}

	// Add the file to the watcher
	err = watcher.Add(fileToWatch)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Watching file: %s", fileToWatch)

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// Check if the file was modified
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("File modified: %s", event.Name)
					// Execute the shell script
					runScript(scriptToRun)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	<-done
}

// runScript executes the specified shell script
func runScript(scriptPath string) {
	cmd := exec.Command("/bin/bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("Error running script: %v", err)
	}
}

