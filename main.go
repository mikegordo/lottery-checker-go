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
	fmt.Println("== Frequency Analysis ==")
	freq := Frequency{}
	freqRes := freq.Analyse(l)

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	maxcount := 5
	for _, v := range l.Data {
		maxcount--
		if maxcount < 0 {
			break
		}
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

			fmt.Printf("%.3f\t", freqRes.Numbers[value][1])
			fmt.Printf("%.3f\t", freqRes.Numbers[value][2])
			fmt.Printf("%.3f\t", freqRes.Numbers[value][3])
			fmt.Printf("%.3f\t", freqRes.Numbers[value][4])
			fmt.Printf("%.3f\n", freqRes.Numbers[value][5])
		}

		fmt.Printf("mega %2d (%.2f)\n", v.mega, freqRes.Mb[v.mega])
	}

	fmt.Print("\n== Random Sets ==")

	rng := Range{}
	rng.Analyse(l)
	dist := Distance{}
	dist.Analyse(l)

	builder := Builder{}
	builder.Initialize()
	built := make([]Numbers, 5)

	for i := 0; i < 5; i++ {
		builder.Populate()
		built[i] = builder.Numbers
	}

	for _, v := range built {
		errRng := rng.CheckSet(v)
		if len(errRng) > 0 {
			fmt.Println("\nrange: fail")
		} else {
			fmt.Println("\nrange: pass")
		}

		err := freq.CheckSet(v)
		if len(err) > 0 {
			fmt.Println("freg: fail")
		} else {
			fmt.Println("freg: pass")
		}

		normal, normMb, total := dist.CheckSet(v)

		for pos, value := range v.Set {
			if inArray(pos, err) {
				fmt.Printf("%s\t", red(strconv.Itoa(value)))
			} else {
				fmt.Printf("%s\t", green(strconv.Itoa(value)))
			}

			fmt.Printf("%.3f\t", freqRes.Numbers[value][1])
			fmt.Printf("%.3f\t", freqRes.Numbers[value][2])
			fmt.Printf("%.3f\t", freqRes.Numbers[value][3])
			fmt.Printf("%.3f\t", freqRes.Numbers[value][4])
			fmt.Printf("%.3f\t", freqRes.Numbers[value][5])

			fmt.Printf("D: %.6f\n", normal[pos+1])
		}

		fmt.Printf("mega %2d (%.2f)\n", v.Mb, freqRes.Mb[v.Mb])
		fmt.Printf("dist: %.6f\n", total)
		fmt.Printf("mega dist %2d (%.2f)\n", v.Mb, normMb)

	}

	fmt.Println("\n== Best Random Sets ==")

	//lastSeq := l.Data[0].numbers
	built = make([]Numbers, 9999)
}
