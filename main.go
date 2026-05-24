// github.com/wpxq
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
)

var (
	bluePrefix   = color.New(color.FgBlue, color.Bold).SprintFunc()   // Running
	yellowPrefix = color.New(color.FgYellow, color.Bold).SprintFunc() // Exited
	redPrefix    = color.New(color.FgRed, color.Bold).SprintFunc()    // Dead
	cyanPrefix   = color.New(color.FgCyan, color.Bold).SprintFunc()   // Paused
	labelStyle   = color.New(color.FgWhite, color.Bold).SprintFunc()
	valueStyle   = color.New(color.FgHiBlack, color.Bold).SprintFunc()
	errorStyle   = color.New(color.FgRed).SprintFunc()
)

func listRunningContainers(cli *client.Client) {
	f := filters.NewArgs()
	f.Add("status", "running")
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{Filters: f})
	if err != nil {
		fmt.Printf("%s %s\n", errorStyle("[!]"), labelStyle(err.Error()))
		return
	}
	fmt.Println(labelStyle(">> RUNNING CONTAINERS"))
	for _, c := range containers {
		name := "N/A"
		if len(c.Names) > 0 {
			name = c.Names[0][1:]
		}
		fmt.Printf("%s %s %-15s %s\n", bluePrefix("[*]"), valueStyle(c.ID[:12]), labelStyle(name), valueStyle(c.Image))
		stats, errStats := cli.ContainerStatsOneShot(context.Background(), c.ID)
		info, errInspect := cli.ContainerInspect(context.Background(), c.ID)
		if errStats == nil && errInspect == nil {
			status := info.State.Status
			var ports []string
			for p := range info.NetworkSettings.Ports {
				ports = append(ports, string(p))
			}
			startedAt, _ := time.Parse(time.RFC3339, info.State.StartedAt)
			uptime := time.Since(startedAt).Round(time.Second)
			var v struct {
				MemoryStats struct {
					Usage uint64 `json:"usage"`
				} `json:"memory_stats"`
			}
			if err := json.NewDecoder(stats.Body).Decode(&v); err == nil {
				memMB := v.MemoryStats.Usage / 1024 / 1024
				fmt.Printf("     %s %s\n", labelStyle(">> STATUS:"), status)
				if len(ports) > 0 {
					fmt.Printf("     %s %s\n", labelStyle(">> PORTS:"), ports)
				}
				fmt.Printf("     %s %dMB\n", labelStyle(">> RAM:"), memMB)
				fmt.Printf("     %s %s\n", labelStyle(">> UPTIME:"), uptime.String())
			}
			stats.Body.Close()
		}
	}
}

func listExitedContainers(cli *client.Client) {
	f := filters.NewArgs()
	f.Add("status", "exited")
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{Filters: f})
	if err != nil {
		fmt.Printf("%s %s\n", errorStyle("[!]"), valueStyle(err.Error()))
		return
	}
	fmt.Println(labelStyle(">> EXITED CONTAINERS"))
	for _, c := range containers {
		name := "N/A"
		if len(c.Names) > 0 {
			name = c.Names[0][1:]
		}
		fmt.Printf("%s %s %s %s\n", yellowPrefix("[!]"), valueStyle(c.ID[:12]), labelStyle(name), valueStyle(c.Image))
	}
}

func listDeadContainers(cli *client.Client) {
	f := filters.NewArgs()
	f.Add("status", "dead")
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{Filters: f})
	if err != nil {
		fmt.Printf("%s %s\n", errorStyle("[!]"), valueStyle(err.Error()))
		return
	}
	fmt.Println(labelStyle(">> DEAD CONTAINERS"))
	for _, c := range containers {
		name := "N/A"
		if len(c.Names) > 0 {
			name = c.Names[0][1:]
		}
		fmt.Printf("%s %s %s %s\n", redPrefix("[x]"), valueStyle(c.ID[:12]), labelStyle(name), valueStyle(c.Image))
	}
}

func listPausedContainers(cli *client.Client) {
	f := filters.NewArgs()
	f.Add("status", "paused")
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{Filters: f})
	if err != nil {
		fmt.Printf("%s %s\n", errorStyle("[!]"), valueStyle(err.Error()))
		return
	}
	fmt.Println(labelStyle(">> PAUSED CONTAINERS"))
	for _, c := range containers {
		name := "N/A"
		if len(c.Names) > 0 {
			name = c.Names[0][1:]
		}
		fmt.Printf("%s %s %s %s\n", cyanPrefix("[=]"), valueStyle(c.ID[:12]), labelStyle(name), valueStyle(c.Image))
	}
}

func main() {
	listRunning := flag.Bool("list-running", false, "Lists running containers")
	listExited := flag.Bool("list-exited", false, "Lists stopped containers")
	listDead := flag.Bool("list-dead", false, "Lists dead containers")
	listPaused := flag.Bool("list-paused", false, "Lists paused containers")
	flag.Usage = func() {
		fmt.Printf("\n%s\n", labelStyle("[ dockstat ]"))
		fmt.Printf("Usage: dockstat %s", labelStyle("<option>\n"))

		fmt.Println(labelStyle("\nOptions:"))
		fmt.Printf("  %-15s %s\n", "--list-running", "Lists all running containers")
		fmt.Printf("  %-15s %s\n", "--list-exited", "Lists all stopped containers")
		fmt.Printf("  %-15s %s\n", "--list-dead", "Lists all dead containers")
		fmt.Printf("  %-15s %s\n", "--list-paused", "Lists all paused containers")

		fmt.Printf("\n%s %s\n", labelStyle("Example:"), "dockstat --list-running")
		fmt.Println()
	}
	flag.Parse()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("%s %s", errorStyle("[!]"), labelStyle(err.Error()))
		return
	}
	defer cli.Close()

	if *listRunning {
		listRunningContainers(cli)
		return
	}
	if *listExited {
		listExitedContainers(cli)
		return
	}
	if *listDead {
		listDeadContainers(cli)
		return
	}
	if *listPaused {
		listPausedContainers(cli)
		return
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		fmt.Printf("%s %s", errorStyle("[!]"), labelStyle(err.Error()))
	}

	var running, exited, dead, paused int
	for _, c := range containers {
		switch c.State {
		case "running":
			running++
		case "exited":
			exited++
		case "dead":
			dead++
		case "paused":
			paused++
		}
	}
	fmt.Println(labelStyle(":: DOCKER INFRASTRUCTURE"))
	fmt.Printf("%s %s %s\n", bluePrefix("[+]"), labelStyle("Running:"), valueStyle(fmt.Sprintf("%d", running)))
	fmt.Printf("%s %s %s\n", yellowPrefix("[!]"), labelStyle("Exited:"), valueStyle(fmt.Sprintf("%d", exited)))
	fmt.Printf("%s %s %s\n", redPrefix("[x]"), labelStyle("Dead:"), valueStyle(fmt.Sprintf("%d", dead)))
	fmt.Printf("%s %s %s\n", cyanPrefix("[=]"), labelStyle("Paused"), valueStyle(fmt.Sprintf("%d", paused)))
}
