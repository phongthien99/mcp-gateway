package main

import (
	"fmt"
	"log"
	"os"

	"mcp-gateway/src/workflow"
)

func main() {
	runConfig := "runs/export-task-csv.yaml"
	if len(os.Args) > 1 {
		runConfig = os.Args[1]
	}

	runner, err := workflow.NewRunner(runConfig)
	if err != nil {
		log.Fatalf("init: %v", err)
	}

	if err := runner.Run(); err != nil {
		log.Fatalf("run: %v", err)
	}

	fmt.Println("\nPOC complete — check artifacts/ for generated files.")
}
