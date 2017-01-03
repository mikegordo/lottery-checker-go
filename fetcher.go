package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type LotteryDay struct {
	date    string
	numbers [5]int
	mega    int
}

type LotteryData struct {
	Data []LotteryDay
}

func (l *LotteryData) Fetch() {
	url := "http://data.ny.gov/resource/5xaw-6ayf.json"

	reader := strings.NewReader(``)
	request, err := http.NewRequest("GET", url, reader)

	if err != nil {
		log.Fatalf("Unable to fetch %v", err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("Unable to fetch %v", err.Error())
	}

	byteArray, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Fatalf("Unable to fetch %v", err.Error())
	}

	var dat []map[string]string
	if err := json.Unmarshal(byteArray, &dat); err != nil {
		log.Fatalf("Unable to convert %v", err.Error())
	}

	l.convert(dat)
}

func (l *LotteryData) convert(dat []map[string]string) {
	for _, v := range dat {
		mega, _ := strconv.Atoi(v["mega_ball"])
		n_ := strings.Split(v["winning_numbers"], " ")
		var numbers [5]int
		for k, n := range n_ {
			numbers[k], _ = strconv.Atoi(n)
		}
		l.Data = append(l.Data, LotteryDay{date: v["draw_date"], numbers: numbers, mega: mega})
	}
}
