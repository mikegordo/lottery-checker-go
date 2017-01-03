package main

import (
	"log"
	"os"
)

func main() {
	os.Setenv("TZ", "America/New_York")
	//l := LotteryData{}

	/* fetch data from remote */
	//l.Fetch()

	b := Builder{}
	b.Initialize()
	b.Populate()
	log.Print(b.Numbers)

}
