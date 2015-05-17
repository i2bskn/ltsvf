package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ltsvf"
	app.Version = "0.0.1"
	app.Usage = "LTSV filter"
	app.Author = "i2bskn"
	app.Email = "i2bskn@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "filter, f",
			Usage: "Filtering by value of specific key.",
		},
		cli.StringFlag{
			Name:  "keys, k",
			Usage: "Display only specified keys.",
		},
	}
	app.Action = func(c *cli.Context) {
		filters := parseFilter(c.String("filter"))
		keys := parseKeys(c.String("keys"))
		condition := newCondition(filters, keys)

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
	app.Run(os.Args)
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

func parseKeys(arg string) []string {
	keys := make([]string, 0, 0)
	if len(arg) > 0 {
		keys = strings.Split(arg, ",")
	}

	return keys
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
		restricted_factor := make([]string, 0, 0)

		for _, factor_string := range strings.Split(line, "\t") {
			factor := strings.SplitN(factor_string, ":", 2)

			if c.displayKey(factor[0]) {
				restricted_factor = append(restricted_factor, factor_string)
			}

			if value, exist := filters[factor[0]]; exist {
				if factor[1] == value {
					delete(filters, factor[0])
				} else {
					break
				}
			} else {
				continue
			}

			if len(filters) == 0 && len(c.keys) == len(restricted_factor) {
				break
			}
		}

		if len(c.keys) > 0 {
			edited = strings.Join(restricted_factor, "\t")
		} else {
			edited = line
		}

		if len(filters) > 0 {
			passing = false
		} else {
			passing = true
		}
	} else {
		if len(c.keys) > 0 {
			restricted_factor := make([]string, 0, 0)
			for _, factor_string := range strings.Split(line, "\t") {
				factor := strings.SplitN(factor_string, ":", 2)
				if c.displayKey(factor[0]) {
					restricted_factor = append(restricted_factor, factor_string)
				}
			}

			edited = strings.Join(restricted_factor, "\t")
		} else {
			edited = line
		}
		passing = true
	}
	return
}
