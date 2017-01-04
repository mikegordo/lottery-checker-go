package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

func main() {
	os.Setenv("TZ", "America/New_York")
	l := LotteryData{}

	/* fetch data from remote */
	l.Fetch()

	fmt.Println("Fetched", len(l.Data), "sets")

	/* display fresh data */
	fmt.Println("== Frequency analysis ==")
	freq := Frequency{}
	ares := freq.Analyse(l)

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	maxcount := 5
	for _, v := range l.Data {
		if maxcount < 0 {
			break
		}
		maxcount--
		set := Numbers{Set: v.numbers, Mb: v.mega}

		fmt.Print(v.date, "\t")
		err := freq.CheckSet(set)
		if len(err) > 0 {
			fmt.Print("freg: fail")
		} else {
			fmt.Print("freg: pass")
		}
		fmt.Println("")

		for pos, value := range v.numbers {
			if inArray(pos, err) {
				fmt.Printf("%s\t", red(strconv.Itoa(value)))
			} else {
				fmt.Printf("%s\t", green(strconv.Itoa(value)))
			}

			fmt.Printf("%.3f\t", ares.Numbers[value][1])
			fmt.Printf("%.3f\t", ares.Numbers[value][2])
			fmt.Printf("%.3f\t", ares.Numbers[value][3])
			fmt.Printf("%.3f\t", ares.Numbers[value][4])
			fmt.Printf("%.3f\n", ares.Numbers[value][5])
		}

		fmt.Printf("mega %2d (%.2f)\n", v.mega, ares.Mb[v.mega])
	}
}
