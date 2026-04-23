package main

import (
	"fmt"
	"log"
	"os"

	"mcp-gateway/src/workflow"
)

func main() {
	workflowPath := "workflows/export-task-csv.yaml"
	projectID := "demo-project"
	featureID := "export-task-csv"

	if len(os.Args) > 1 {
		workflowPath = os.Args[1]
	}
	if len(os.Args) > 2 {
		projectID = os.Args[2]
	}
	if len(os.Args) > 3 {
		featureID = os.Args[3]
	}

	runner, err := workflow.NewRunner(workflowPath, projectID, featureID)
	if err != nil {
		log.Fatalf("init: %v", err)
	}

	extraArgs := map[string]string{
		"project_id": projectID,
		"feature_id": featureID,
		"request":    "Tôi muốn thêm chức năng export danh sách task ra file CSV. Người dùng bấm nút Export để tải file về. CSV gồm: title, status, assignee, due_date.",
	}

	if err := runner.Run(extraArgs); err != nil {
		log.Fatalf("run: %v", err)
	}

	fmt.Println("\nPOC complete — check artifacts/ for generated files.")
}
