package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"text/tabwriter"
)

type TimeZone struct {
	name       string
	addr       string
	latestTime string
}

type ClockUpdate struct {
	index int
	time  string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: clockwall {TIMEZONE}={ipAddr:port} ...")
		os.Exit(1)
	}

	cmd_data := os.Args[1:]
	timezones := parseTimezone(cmd_data)

	updates := make(chan ClockUpdate)

	// Dial each server in a goroutine
	for i, tz := range timezones {
		go connectToClockServer(i, tz.name, tz.addr, updates)
	}

	// Main loop: receive updates and refresh the table
	for update := range updates {
		timezones[update.index].latestTime = update.time
		printTable(timezones)
	}
}

func parseTimezone(data []string) (tzs []*TimeZone) {
	tzs = make([]*TimeZone, 0)
	for _, tz := range data {
		parts := strings.Split(strings.TrimSpace(tz), "=")
		if len(parts) == 2 {
			tzs = append(tzs, &TimeZone{name: parts[0], addr: parts[1]})
		}
	}
	return tzs
}

func connectToClockServer(index int, name, addr string, ch chan<- ClockUpdate) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		ch <- ClockUpdate{index: index, time: "Connection Failed"}
		return
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ch <- ClockUpdate{index: index, time: scanner.Text()}
	}

	if err := scanner.Err(); err != nil {
		ch <- ClockUpdate{index: index, time: "Error Reading"}
	} else {
		ch <- ClockUpdate{index: index, time: "Disconnected"}
	}
}

func printTable(tzs []*TimeZone) {
	// Clear screen and move cursor to top-left
	fmt.Print("\033[H\033[2J")

	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)

	// Print Header
	fmt.Fprintf(tw, "Location\tClock\t\n")
	fmt.Fprintf(tw, "--------\t-----\t\n")

	// Print Rows
	for _, tz := range tzs {
		t := tz.latestTime
		if t == "" {
			t = "Waiting..."
		}
		fmt.Fprintf(tw, "%s\t%s\t\n", tz.name, t)
	}
	tw.Flush()
}
