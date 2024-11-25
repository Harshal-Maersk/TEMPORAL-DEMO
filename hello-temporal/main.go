package main

import (
	"context"
	"hello-temporal/workflows"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create a Temporal client
	c, err := client.NewClient(client.Options{
		HostPort: "localhost:7233", // Temporal service address
	})
	if err != nil {
		log.Fatalf("Unable to create Temporal client: %v", err)
	}
	defer c.Close()

	// Start a worker
	w := worker.New(c, "hello-world-task-queue", worker.Options{})
	// Register the workflow
	w.RegisterWorkflow(workflows.TemporalDemoWorkflow)

	// Start the worker in a goroutine
	go func() {
		if err := w.Start(); err != nil {
			log.Fatalf("Failed to start worker: %v", err)
		}
	}()
	log.Println("Worker started and polling for tasks...")

	// Run the workflow
	go func() {
		workflowOptions := client.StartWorkflowOptions{
			TaskQueue: "hello-world-task-queue",
		}
		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.TemporalDemoWorkflow, "World")
		if err != nil {
			log.Fatalf("Failed to start workflow: %v", err)
		}

		// Get the workflow result
		var result string
		if err := we.Get(context.Background(), &result); err != nil {
			log.Fatalf("Failed to get workflow result: %v", err)
		}

		log.Printf("Workflow result: %s", result)
	}()

	// Keep the application running and listen for OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	// Graceful shutdown of the worker
	log.Println("Shutting down worker...")
	w.Stop()
	log.Println("Worker stopped. Exiting application.")
}
