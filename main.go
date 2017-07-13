package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"sync"

	"github.com/fatih/color"
)

var wg sync.WaitGroup

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

	total := 1000000
	threads := 20
	maxcount = 5
	res := make([]Builts, threads*maxcount)
	j := 0

	for t := 0; t < threads; t++ {
		wg.Add(1)
		go func() {
			b := brs(l.Data[0].numbers, l.Data[0].mega, dist, total/threads)
			for i := 0; i < maxcount; i++ {
				res[j] = b[i]
				j++
			}
			wg.Done()
		}()
	}

	wg.Wait()

	s := make(ASorter, len(res))
	i := 0
	for _, v := range res {
		s[i] = Builts{v.Numbers, v.Total}
		i++
	}

	sort.Sort(s)
	res = s

	for _, v := range res {
		maxcount--
		if maxcount < 0 {
			break
		}
		errRng := rng.CheckSet(v.Numbers)
		if len(errRng) > 0 {
			fmt.Println("\nrange: fail")
		} else {
			fmt.Println("\nrange: pass")
		}

		err := freq.CheckSet(v.Numbers)
		if len(err) > 0 {
			fmt.Println("freg: fail")
		} else {
			fmt.Println("freg: pass")
		}

		normal, normMb, total := dist.CheckSet(v.Numbers)

		for pos, value := range v.Numbers.Set {
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

		fmt.Printf("mega %2d (%.2f)\n", v.Numbers.Mb, freqRes.Mb[v.Numbers.Mb])
		fmt.Printf("dist: %.6f\n", total)
		fmt.Printf("mega dist %2d (%.2f)\n", v.Numbers.Mb, normMb)

	}

}

func brs(lastSeq []int, lastMb int, dist Distance, size int) []Builts {
	builder := Builder{}
	builder.Initialize()
	builts := make([]Builts, size)
	predefined := make(map[string]Predefined)

	for i := 0; i < size; i++ {
		builder.Populate()
		bset := builder.Numbers
		if 0 != arrayDiff(bset.Set, lastSeq) || bset.Mb == lastMb {
			continue
		}

		res, resMb, resTotal := dist.CheckSet(bset)
		if resTotal > 0 {
			builts[i].Numbers = builder.Numbers
			builts[i].Total = resTotal

			tmp := Predefined{}
			tmp.Mb = resMb
			tmp.Numbers = res
			tmp.Total = resTotal

			predefined[builder.GetNumbersString()] = tmp
		}
	}

	/* sort jobs - and this is not trivial */
	if len(builts) > 1 {
		s := make(ASorter, len(builts))
		i := 0
		for _, v := range builts {
			s[i] = Builts{v.Numbers, v.Total}
			i++
		}

		sort.Sort(s)
		builts = s
	}

	return builts
}

type Predefined struct {
	Numbers [6]float32
	Mb      float32
	Total   float32
}

type Builts struct {
	Numbers Numbers
	Total   float32
}

type ASorter []Builts

func (a ASorter) Len() int           { return len(a) }
func (a ASorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ASorter) Less(i, j int) bool { return a[i].Total > a[j].Total }

func arrayDiff(a []int, b []int) int {
	result := 0
	for _, i := range a {
		if inArray(i, b) {
			result++
		}
	}

	return result
}
