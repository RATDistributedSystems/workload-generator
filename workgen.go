package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type command struct {
	command             string
	rawCommandString    string
	usernameRequired    bool
	stockIDRequired     bool
	stockAmountRequired bool
	values              []string
}

type userChannels struct {
	input chan command
	start chan bool
}

var (
	ip       = flag.String("ip", "localhost", "IP Address to send requests to webserver on. Default is localhost")
	port     = flag.Int("p", 44440, "Port to send requests to the webserver on. Default is 44440")
	filename = flag.String("f", "", "file to execute workload commands from")
	rate     = flag.Int("r", 50, "Delay (in ms) between successive commands")
	useTCP   = flag.Bool("tcp", false, "Sends the request as a TCP message instead of HTTP")
	cmd      = flag.String("c", "", "single user command to execute")
	parallel = flag.Bool("para", false, "Whether to parallize the workload by user")
	useNum   = flag.Bool("num", true, "Use the Transaction Number provided by the workload file")
	addr     string
	url      string
)

func main() {

	flag.Parse()
	addr = fmt.Sprintf("%s:%d", *ip, *port)
	url = fmt.Sprintf("http://%s/result", addr)
	var file *os.File
	var scanner *bufio.Scanner
	var wg sync.WaitGroup
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

	// Parallization
	userCommandInput := make(map[string]userChannels)
	var parallelDumplog command
	userCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)

		if *parallel && *cmd == "" {
			comm, _ := parseLine(line)

			if comm.command == "DUMPLOG" {
				parallelDumplog = *comm
				wg.Done()
				continue
			}

			userChanns, exists := userCommandInput[comm.values[1]]
			if !exists {
				userChanns.input = make(chan command)
				userChanns.start = make(chan bool, 1)
				userCommandInput[comm.values[1]] = userChanns
				go parallelUserExecution(userChanns.input, userChanns.start, &wg)
				userCount++
				fmt.Printf("\rUser Count: %-4d", userCount)
			}

			userChanns.input <- *comm
			continue
		}

		if *useTCP {
			generateTCPRequest(line, &wg)
		} else {
			comm, _ := parseLine(line)
			generateHTTPRequests(comm, &wg)
		}
		time.Sleep(time.Millisecond * time.Duration(*rate))
	}

	if *parallel && *cmd == "" {
		fmt.Println("")
		countdown := 3
		for i := 0; i < countdown; i++ {
			fmt.Printf("\rParallel Execution Starting in %d", countdown-i)
			time.Sleep(time.Second)
		}
		fmt.Println("\nStarting...")
		for _, user := range userCommandInput {
			user.start <- true
		}
	}

	wg.Wait()

	if parallelDumplog.command != "" {
		wg.Add(1)
		if *useTCP {
			generateTCPRequest(parallelDumplog.rawCommandString, &wg)
		} else {
			generateHTTPRequests(&parallelDumplog, &wg)
		}
	}
}

func getTransactionNumber(line string) string {
	re := regexp.MustCompile(`(?s)\[(.*)\]`)
	m := re.FindAllStringSubmatch(line, -1)
	return m[0][1]
}

func removeBrackets(line string) string {
	i := strings.Index(line, "]") + 1
	return strings.TrimSpace(string(line[i:]))
}

func parseLine(line string) (*command, error) {
	transaction := getTransactionNumber(line)
	trimmedLine := removeBrackets(line)
	args := strings.Split(trimmedLine, ",")
	for i, items := range args {
		args[i] = strings.TrimSpace(items)
	}
	cmd, err := checkForValidCommand(args[0])

	if err != nil {
		return nil, err
	}

	cmd.rawCommandString = strings.Join(args, ",")
	args = append(args, transaction)
	cmd.values = args
	return cmd, nil
}

func createCommandStruct(c string, uname bool, stock bool, amt bool) *command {
	return &command{c, "", uname, stock, amt, nil}
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

	c.values = pop(c.values)
	if !*useNum {
		c.values[0] = "0"
	}
	m["transaction"] = c.values[0:1]
	return
}

func generateTCPRequest(line string, wg *sync.WaitGroup) {
	trimmedLine := removeBrackets(line)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err.Error())
		wg.Done()
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "%s\n", trimmedLine)
	log.Printf("TCP: %s\n", line)
	wg.Done()
}

func generateHTTPRequests(c *command, wg *sync.WaitGroup) {
	if c == nil {
		wg.Done()
		return
	}
	values := generateMapFromCommand(c)
	resp, err := http.PostForm(url, values)
	if err != nil {
		fmt.Println("\t" + err.Error())
		wg.Done()
		return
	}
	resp.Request.Close = true
	resp.Body.Close()
	log.Printf("HTTP: %s", c.rawCommandString)
	wg.Done()
}

func parallelUserExecution(line <-chan command, start <-chan bool, wg *sync.WaitGroup) {
	var lines []command

	for {
		select {
		case msg := <-line:
			lines = append(lines, msg)
		case <-start:
			for _, item := range lines {
				if *useTCP {
					generateTCPRequest(item.rawCommandString, wg)
					time.Sleep(time.Millisecond * time.Duration(*rate))

				} else {
					go generateHTTPRequests(&item, wg)
					time.Sleep(time.Millisecond * time.Duration(*rate))
				}
			}
			return
		}
	}
}
