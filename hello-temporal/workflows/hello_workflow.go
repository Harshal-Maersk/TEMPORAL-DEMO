package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// TemporalDemoWorkflow is a simple Temporal workflow that returns a greeting message.
func TemporalDemoWorkflow(ctx workflow.Context, name string) (string, error) {
	// Workflow logic: Just returns "Hello, {name}!" after 1 second
	workflow.Sleep(ctx, 1*time.Second)
	return "Hello, " + name + "!", nil
}
