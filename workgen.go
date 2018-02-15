package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type command struct {
	command             string
	usernameRequired    bool
	stockIDRequired     bool
	stockAmountRequired bool
	values              []string
}

var (
	ip       = flag.String("ip", "localhost", "IP Address to send requests to webserver on. Default is localhost")
	port     = flag.Int("p", 44440, "Port to send requests to the webserver on. Default is 44440")
	filename = flag.String("f", "", "file to execute workload commands from")
	rate     = flag.Int("r", 50, "Delay (in ms) between successive commands")
	useTCP   = flag.Bool("tcp", false, "Sends the request as a TCP message instead of HTTP")
	cmd      = flag.String("c", "", "single user command to execute")
	addr     string
	url      string
)

func main() {

	flag.Parse()
	addr = fmt.Sprintf("%s:%d", *ip, *port)
	url = fmt.Sprintf("http://%s/result", addr)
	var file *os.File
	var scanner *bufio.Scanner
	var err error

	if *filename != "" {
		file, err = os.Open(*filename)
		if err != nil {
			panic("Couldn't open file: " + *filename)
		}
		scanner = bufio.NewScanner(file)
		defer file.Close()
	} else if *cmd != "" {
		scanner = bufio.NewScanner(strings.NewReader(*cmd))
	} else {
		panic("No commands to process")
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Sent: %s\n", line)

		if *useTCP {
			generateTCPRequest(line)
		} else {
			cmd, err := parseLine(line)
			if err != nil {
				fmt.Printf("\t%s: \"%s\"\n", err.Error(), line)
				continue
			}
			generateHTTPRequests(cmd)
		}
		time.Sleep(time.Millisecond * time.Duration(*rate))
	}

}

func removeBrackets(line string) string {
	i := strings.Index(line, "]") + 1
	return strings.TrimSpace(string(line[i:]))
}

func parseLine(line string) (*command, error) {
	trimmedLine := removeBrackets(line)
	args := strings.Split(trimmedLine, ",")
	cmd, err := checkForValidCommand(args[0])
	if err != nil {
		return nil, err
	}
	cmd.values = args
	return cmd, nil
}

func createCommandStruct(c string, uname bool, stock bool, amt bool) *command {
	return &command{c, uname, stock, amt, nil}
}

func checkForValidCommand(cmd string) (c *command, e error) {
	switch strings.ToUpper(cmd) {
	case "ADD":
		c, e = createCommandStruct(cmd, true, false, true), nil
	case "BUY":
		c, e = createCommandStruct(cmd, true, true, true), nil
	case "SELL":
		c, e = createCommandStruct(cmd, true, true, true), nil
	case "QUOTE":
		c, e = createCommandStruct(cmd, true, true, false), nil
	case "COMMIT_BUY":
		c, e = createCommandStruct(cmd, true, false, false), nil
	case "COMMIT_SELL":
		c, e = createCommandStruct(cmd, true, false, false), nil
	case "CANCEL_BUY":
		c, e = createCommandStruct(cmd, true, false, false), nil
	case "CANCEL_SELL":
		c, e = createCommandStruct(cmd, true, false, false), nil
	case "SET_BUY_AMOUNT":
		c, e = createCommandStruct(cmd, true, true, true), nil
	case "SET_BUY_TRIGGER":
		c, e = createCommandStruct(cmd, true, true, true), nil
	case "CANCEL_SET_BUY":
		c, e = createCommandStruct(cmd, true, true, false), nil
	case "SET_SELL_AMOUNT":
		c, e = createCommandStruct(cmd, true, true, true), nil
	case "SET_SELL_TRIGGER":
		c, e = createCommandStruct(cmd, true, true, true), nil
	case "CANCEL_SET_SELL":
		c, e = createCommandStruct(cmd, true, true, false), nil
	case "DUMPLOG":
		c, e = createCommandStruct(cmd, true, false, false), nil
	case "DISPLAY_SUMMARY":
		c, e = createCommandStruct(cmd, true, false, false), nil
	default:
		c, e = nil, errors.New("Invalid Command")
	}
	return
}

func pop(s []string) []string {
	return s[1:]
}

func generateMapFromCommand(c *command) (m map[string][]string) {
	m = make(map[string][]string)
	m["command"] = c.values[0:1]

	// Get Username
	if c.usernameRequired {
		c.values = pop(c.values)
		m["username"] = c.values[0:1]
	}
	// Get Stock ID
	if c.stockIDRequired {
		c.values = pop(c.values)
		m["stock"] = c.values[0:1]
	}

	// Get stock amount
	if c.stockAmountRequired {
		c.values = pop(c.values)
		m["amount"] = c.values[0:1]
	}
	return
}

func generateTCPRequest(line string) {
	trimmedLine := removeBrackets(line)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		// rip
	}
	defer conn.Close()
	fmt.Fprintf(conn, "%s\n", trimmedLine)
}

func generateHTTPRequests(c *command) {
	values := generateMapFromCommand(c)
	resp, err := http.PostForm(url, values)
	if err != nil {
		fmt.Println("\t" + err.Error())
		return
	}

	fmt.Println("\t" + resp.Status)
	resp.Body.Close()

}
