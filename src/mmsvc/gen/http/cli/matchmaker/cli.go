// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// matchmaker HTTP client CLI support package
//
// Command:
// $ goa gen github.com/proepkes/speeddate/src/mmsvc/design

package cli

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	matchmakingc "github.com/proepkes/speeddate/src/mmsvc/gen/http/matchmaking/client"
	goa "goa.design/goa"
	goahttp "goa.design/goa/http"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `matchmaking insert
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` matchmaking insert` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, interface{}, error) {
	var (
		matchmakingFlags = flag.NewFlagSet("matchmaking", flag.ContinueOnError)

		matchmakingInsertFlags = flag.NewFlagSet("insert", flag.ExitOnError)
	)
	matchmakingFlags.Usage = matchmakingUsage
	matchmakingInsertFlags.Usage = matchmakingInsertUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if len(os.Args) < flag.NFlag()+3 {
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = os.Args[1+flag.NFlag()]
		switch svcn {
		case "matchmaking":
			svcf = matchmakingFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(os.Args[2+flag.NFlag():]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = os.Args[2+flag.NFlag()+svcf.NFlag()]
		switch svcn {
		case "matchmaking":
			switch epn {
			case "insert":
				epf = matchmakingInsertFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if len(os.Args) > 2+flag.NFlag()+svcf.NFlag() {
		if err := epf.Parse(os.Args[3+flag.NFlag()+svcf.NFlag():]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "matchmaking":
			c := matchmakingc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "insert":
				endpoint = c.Insert()
				data = nil
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// matchmakingUsage displays the usage of the matchmaking command and its
// subcommands.
func matchmakingUsage() {
	fmt.Fprintf(os.Stderr, `matchmaking.
Usage:
    %s [globalflags] matchmaking COMMAND [flags]

COMMAND:
    insert: .

Additional help:
    %s matchmaking COMMAND --help
`, os.Args[0], os.Args[0])
}
func matchmakingInsertUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] matchmaking insert

.

Example:
    `+os.Args[0]+` matchmaking insert
`, os.Args[0])
}