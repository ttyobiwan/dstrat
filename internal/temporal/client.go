package temporal

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
)

type Client struct {
	client client.Client
}

func NewClient(ctx context.Context, options client.Options) (*Client, error) {
	c, err := client.DialContext(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("creating client: %v", err)
	}
	return &Client{client: c}, nil
}

func (c *Client) Execute(ctx context.Context, task any, options client.StartWorkflowOptions, args ...any) (client.WorkflowRun, error) {
	we, err := c.client.ExecuteWorkflow(ctx, options, task, args...)
	if err != nil {
		return nil, fmt.Errorf("executing workflow: %v", err)
	}
	return we, err
}

func (c *Client) AwaitResult(ctx context.Context, task client.WorkflowRun, dest any) error {
	err := task.Get(ctx, &dest)
	if err != nil {
		return fmt.Errorf("getting workflow result: %v", err)
	}
	return nil
}

func (c *Client) Close() error {
	c.client.Close()
	return nil
}
