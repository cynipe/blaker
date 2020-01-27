package main

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-cmd/cmd"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/cynipe/blaker/pkg/blaker"
	"github.com/cynipe/blaker/pkg/clock"
)

var (
	writer    io.Writer = os.Stdout
	errWriter io.Writer = os.Stderr

	usageError = 1
	breakError = 250
	// tool error not the wrapped command error
	blakerError = 255
)

func main() {
	app := blakerApp()
	err := app.Run(os.Args)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

func blakerApp() *cli.App {
	app := cli.NewApp()
	app.Version = "0.0.3"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "profile, p",
			Usage: "aws profile to use",
		},
		&cli.StringFlag{
			Name:  "region, r",
			Usage: "aws region to use",
		},
		&cli.StringFlag{
			Name:  "config-key, c",
			Usage: "config key for ddb table",
			Value: "default",
		},
		&cli.BoolFlag{
			Name:  "error-on-break, E",
			Usage: "return non-zero (250) if on break time",
		},
	}
	app.Writer = writer
	app.ErrWriter = errWriter
	app.Usage = "<command> [command args]"
	app.Action = run
	return app
}

func run(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		_ = cli.ShowAppHelp(ctx)
		return cli.NewExitError("required args not specified", usageError)
	}

	sess, err := awsSession(ctx)
	if err != nil {
		return err
	}

	args := ctx.Args()
	var cmdArgs []string
	if ctx.NArg() > 1 {
		cmdArgs = args[1:]
	}

	b := blaker.New(dynamodb.New(sess), clock.New(), ctx.String("config-key"))

	status, err := b.RunCmd(&blaker.RunCmdInput{
		Command: args.First(),
		Args:    cmdArgs,
		Stdout:  ctx.App.Writer,
		Stderr:  ctx.App.ErrWriter,
	})

	return handleError(err, ctx, status)
}

func handleError(err error, ctx *cli.Context, status cmd.Status) error {
	if err != nil {
		switch err.(type) {
		case *blaker.BreakError:
			// on break-time
			if ctx.Bool("error-on-break") {
				return cli.NewExitError(err, breakError)
			}
			if _, werr := fmt.Fprintln(ctx.App.Writer, err); werr != nil {
				return cli.NewExitError(errors.Wrapf(werr, "failed to write skipped log: %s", err), blakerError)
			}
			return nil
		default:
			return cli.NewExitError(err, blakerError)
		}
	}

	merr := make([]error, 0)
	// go-cmd or os/exec internal error
	if status.Error != nil {
		merr = append(merr, cli.NewExitError(status.Error, blakerError))
	}

	// command failure
	if status.Exit != 0 {
		merr = append(merr, cli.NewExitError(status.Error, status.Exit))
	}

	// aggregate all the errors
	if len(merr) > 0 {
		return cli.NewMultiError(merr...)
	}
	return nil
}

func awsSession(ctx *cli.Context) (*session.Session, error) {
	options := session.Options{}
	if profile := ctx.String("profile"); profile != "" {
		fmt.Printf("using profile: %s\n", profile)
		options.SharedConfigState = session.SharedConfigEnable
		options.Profile = profile
	}
	if region := ctx.String("region"); region != "" {
		fmt.Printf("using region: %s\n", region)
		options.Config.Region = aws.String(region)
	}
	return session.NewSessionWithOptions(options)
}
