package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type command struct {
	command             string
	usernameRequired    bool
	stockIDRequired     bool
	stockAmountRequired bool
	values              []string
}

var ip string
var port string
var url string

func init() {
	flag.StringVar(&ip, "ip", "localhost", "IP Address to send requests to webserver on. Default is localhost")
	flag.StringVar(&port, "port", "44440", "Port to send requests to the webserver on. Default is 44440")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No workload file provided")
		os.Exit(1)
	}

	filename := args[0]
	file, err := os.Open(filename)
	if err != nil {
		panic("Couldn't open file: " + filename)
	}

	defer file.Close()
	flag.Parse()
	url = fmt.Sprintf("http://%s:%s/result", ip, port)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line + ":")
		cmd, err := parseLine(line)
		if err != nil {
			fmt.Println("\tCoudln't parse row")
			continue
		}
		generateHTTPRequests(cmd)
	}

}

func parseLine(line string) (*command, error) {
	i := strings.Index(line, "]") + 1
	trimmedLine := strings.TrimSpace(string(line[i:]))
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
	switch cmd {
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
