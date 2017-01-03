package main

import (
	"log"
	"os"
)

func main() {
	os.Setenv("TZ", "America/New_York")
	l := LotteryData{}

	/* fetch data from remote */
	l.Fetch()

	b := Builder{}
	b.Initialize()
	b.Populate()

	r := Frequency{}
	r.Analyse(l)
	z := r.CheckSet(b.Numbers)

	log.Println(b.Numbers)

	log.Println(z)
}
