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

func parseFilter(arg string) map[string]string {
	filters := make(map[string]string)
	if len(arg) > 0 {
		for _, filter_string := range strings.Split(arg, ",") {
			filter := strings.SplitN(filter_string, ":", 2)
			filters[filter[0]] = filter[1]
		}
	}

	return filters
}

func filterAndDisplay(file *os.File, c *Condition) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line, passing := parseLineOfLtsv(scanner.Text(), c)
		if passing {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func parseLineOfLtsv(line string, c *Condition) (edited string, passing bool) {
	if len(c.filters) > 0 {
		filters := c.copiedFilters()
		for _, factor_string := range strings.Split(line, "\t") {
			factor := strings.SplitN(factor_string, ":", 2)
			value, exist := filters[factor[0]]
			if exist {
				if factor[1] == value {
					delete(filters, factor[0])
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

		if len(filters) > 0 {
			edited = line
			passing = false
		} else {
			edited = line
			passing = true
		}
	} else {
		edited = line
		passing = true
	}
	return
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

