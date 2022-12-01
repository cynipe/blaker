package blaker

import (
	"fmt"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/go-cmd/cmd"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/cynipe/blaker/pkg/clock"
)

type Blaker struct {
	db        *dynamodb.DynamoDB
	clock     clock.Clock
	configKey string
}

func New(db *dynamodb.DynamoDB, clock clock.Clock, configKey string) *Blaker {
	return &Blaker{
		db:        db,
		clock:     clock,
		configKey: configKey,
	}
}

type RunCmdInput struct {
	Command      string
	Args         []string
	Stdout       io.Writer
	Stderr       io.Writer
	WaitDuration time.Duration
	Verbose      bool
	NoDelay      bool
}

func (b *Blaker) RunCmd(input *RunCmdInput) (cmd.Status, error) {
	breakTime, err := b.GetBreakTime()
	if err != nil {
		return cmd.Status{}, err
	}
	if breakTime != nil && b.clock.Now().After(*breakTime) {
		return cmd.Status{}, NewBreakError(*breakTime, input)
	}

	options := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	c := cmd.NewCmdOptions(options, input.Command, input.Args...)
	go func(stdout, stderr io.Writer) {
		for {
			select {
			case line := <-c.Stdout:
				if _, err := fmt.Fprintln(stdout, line); err != nil {
					fmt.Println(err)
				}
			case line := <-c.Stderr:
				if _, err := fmt.Fprintln(stderr, line); err != nil {
					fmt.Println(err)
				}
			}
		}
	}(input.Stdout, input.Stderr)
	status := <-c.Start()
	// wait until goroutine prints all the message
	for len(c.Stdout) > 0 || len(c.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}
	return status, nil
}

func (b *Blaker) GetBreakTime() (*time.Time, error) {
	req, res := b.db.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String("blaker_config"),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {S: aws.String(b.configKey)},
		},
	})
	if err := req.Send(); err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve value for key `%s` from blaker_config, check your ddb.", b.configKey)
	}

	v := res.Item["value"]
	if v == nil {
		return nil, errors.Errorf("failed to retrieve value for key `%s` from blaker_config, check your ddb.", b.configKey)
	}

	val := aws.StringValue(v.S)
	if val == "" {
		return nil, nil
	}

	bt, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return nil, err
	}
	return &bt, nil
}
