package client

import (
	"flag"
	"fmt"
	"ngrok/version"
	"os"
)

const usage1 string = `Usage: %s [OPTIONS] <local port or address>
Options:
`

const usage2 string = `
Examples:
	ngrok 80
	ngrok -subdomain=example 8080
	ngrok -proto=tcp 22
	ngrok -hostname="example.com" -httpauth="user:password" 10.0.0.1


Advanced usage: ngrok [OPTIONS] <command> [command args] [...]
Commands:
	ngrok start [tunnel] [...]    Start tunnels by name from config file
	ngork start-all               Start all tunnels defined in config file
	ngrok list                    List tunnel names from config file
	ngrok help                    Print help
	ngrok version                 Print ngrok version

Examples:
	ngrok start www api blog pubsub
	ngrok -log=stdout -config=ngrok.yml start ssh
	ngrok start-all
	ngrok version

`

type Options struct {
	config     string
	logto      string
	loglevel   string
	authtoken  string
	httpauth   string
	hostname   string
	protocol   string
	subdomain  string
	server     string
	remotePort string
	command    string
	args       []string
}

func ParseArgs(args []string) (opts *Options, err error) {
	var flags = flag.NewFlagSet(args[0], flag.PanicOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, usage1, args[0])
		flags.PrintDefaults()
		fmt.Fprintf(os.Stderr, usage2)
	}

	config := flags.String(
		"config",
		"",
		"Path to ngrok configuration file. (default: $HOME/.ngrok)")

	logto := flags.String(
		"log",
		"none",
		"Write log messages to this file. 'stdout' and 'none' have special meanings")

	loglevel := flags.String(
		"log-level",
		"DEBUG",
		"The level of messages to log. One of: DEBUG, INFO, WARNING, ERROR")

	authtoken := flags.String(
		"authtoken",
		"",
		"Authentication token for identifying an ngrok.com account")

	httpauth := flags.String(
		"httpauth",
		"",
		"username:password HTTP basic auth creds protecting the public tunnel endpoint")

	subdomain := flags.String(
		"subdomain",
		"",
		"Request a custom subdomain from the ngrok server. (HTTP only)")

	server := flags.String(
		"server",
		"",
		"Connect to the specified host:port")

	remotePort := flags.String(
		"remotePort",
		"",
		"Use the specified remotePort")

	hostname := flags.String(
		"hostname",
		"",
		"Request a custom hostname from the ngrok server. (HTTP only) (requires CNAME of your DNS)")

	protocol := flags.String(
		"proto",
		"http",
		"The protocol of the traffic over the tunnel {'http', 'https', 'tcp'} (default: 'http')")

	flags.Parse(args[1:])

	opts = &Options{
		config:     *config,
		logto:      *logto,
		loglevel:   *loglevel,
		httpauth:   *httpauth,
		subdomain:  *subdomain,
		server:     *server,
		remotePort: *remotePort,
		protocol:   *protocol,
		authtoken:  *authtoken,
		hostname:   *hostname,
		command:    flags.Arg(0),
	}

	switch opts.command {
	case "list":
		opts.args = flags.Args()[1:]
	case "start":
		opts.args = flags.Args()[1:]
	case "start-all":
		opts.args = flags.Args()[1:]
	case "version":
		fmt.Println(version.MajorMinor())
		panic("version")
		//os.Exit(0)
	case "help":
		flags.Usage()
		panic("help")
		//os.Exit(0)
	case "":
		err = fmt.Errorf("Error: Specify a local port to tunnel to, or " +
			"an ngrok command.\n\nExample: To expose port 80, run " +
			"'ngrok 80'")
		return

	default:
		if len(flags.Args()) > 1 {
			err = fmt.Errorf("You may only specify one port to tunnel to on the command line, got %d: %v",
				len(flags.Args()),
				flags.Args())
			return
		}

		opts.command = "default"
		opts.args = flags.Args()
	}

	return
}
