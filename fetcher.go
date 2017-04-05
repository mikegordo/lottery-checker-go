package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
LotteryDay - represents a single lottery result
*/
type LotteryDay struct {
	date    string
	numbers []int
	mega    int
}

/*
LotteryData - all lottery data for all time
*/
type LotteryData struct {
	Data []LotteryDay
}

/*
Fetch - returns data from remote
*/
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

/*
convert - conterts data to our format
*/
func (l *LotteryData) convert(dat []map[string]string) {
	for _, v := range dat {
		if strings.Compare(v["draw_date"], "2013-10-15T00:00:00") < 0 {
			/* ignore too old data */
			continue
		}
		mega, _ := strconv.Atoi(v["mega_ball"])
		n_ := strings.Split(v["winning_numbers"], " ")
		var numbers []int
		numbers = make([]int, 5)
		for k, n := range n_ {
			numbers[k], _ = strconv.Atoi(n)
		}
		l.Data = append(l.Data, LotteryDay{date: v["draw_date"][:10], numbers: numbers, mega: mega})
	}
}
