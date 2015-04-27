package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

type Condition struct {
	filters map[string]string
}

func newCondition(filters map[string]string) *Condition {
	return &Condition{
		filters: filters,
	}
}

func (condition *Condition) copiedFilters() map[string]string {
	filters := make(map[string]string)
	for key, value := range condition.filters {
		filters[key] = value
	}
	return filters
}

func mainAction(c *cli.Context) {
	filters := parseFilter(c.String("filter"))
	condition := newCondition(filters)

	if len(c.Args()) > 0 {
		for _, filename := range c.Args() {
			file, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			filterAndDisplay(file, condition)
		}
	} else {
		filterAndDisplay(os.Stdin, condition)
	}
}

func filterAndDisplay(file *os.File, c *Condition) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line, filtering := parseLineOfLtsv(scanner.Text(), c)
		if filtering {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func parseFilter(filters string) map[string]string {
	parsed_filters := make(map[string]string)
	if len(filters) == 0 {
		return parsed_filters
	}

	for _, filter := range strings.Split(filters, ",") {
		splited_filter := strings.SplitN(filter, ":", 2)
		parsed_filters[splited_filter[0]] = splited_filter[1]
	}
	return parsed_filters
}

func parseLineOfLtsv(line string, c *Condition) (string, bool) {
	if len(c.filters) == 0 {
		return line, true
	}

	filters := c.copiedFilters()
	for _, factor := range strings.Split(line, "\t") {
		splited_factor := strings.SplitN(factor, ":", 2)
		condition, exist := filters[splited_factor[0]]
		if exist {
			if splited_factor[1] == condition {
				delete(filters, splited_factor[0])
			} else {
				break
			}
		} else {
			continue
		}

		if len(filters) == 0 {
			break
		}
	}

	if len(filters) != 0 {
		return line, false
	} else {
		return line, true
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "ltsvf"
	app.Version = "0.0.1"
	app.Usage = "LTSV filter"
	app.Author = "i2bskn"
	app.Email = "i2bskn@gmail.com"
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "filter, f",
			Usage: "Filtering the value of specific key.",
		},
	}
	app.Action = mainAction
	app.Run(os.Args)
}

