package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Attumm/settingo/settingo"
)

type Resp struct {
	Duration int
	//Query string
	Found int
	Body  string
	URL   string
}

func getDuration(duration_raw string) int {
	duration_str := duration_raw[:len(duration_raw)-2]
	duration_int, _ := strconv.Atoi(duration_str)
	return duration_int

}

func toInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n

}
func createQuery(host string, query string, page, pagesize int) string {
	return fmt.Sprintf(
		"%s/search/?page=%d&pagesize=%d&search=%s",
		host, page, pagesize, query,
	)
}

func worker(id int, urls <-chan string, responses chan<- *Resp) {
	for url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if DEBUG {
			fmt.Println(body)
			fmt.Println(resp.Status)

			for key, val := range resp.Header {
				fmt.Println(key, "=", val)
			}
		}

		responses <- &Resp{
			Duration: getDuration(resp.Header["Query-Duration"][0]),
			Found:    toInt(resp.Header["Total-Items"][0]),
			Body:     string(body),
			URL:      url,
		}
	}
}

// damn, no stats lib
func sum(ii []int) int {
	total := 0
	for _, i := range ii {
		total += i
	}
	return total
}

func mean(ii []int) int {
	return sum(ii) / len(ii)
}

func meanf(ii []int) float64 {
	return float64(sum(ii)) / float64(len(ii))
}

func max(ii []int) int {
	max := ii[0]
	for _, i := range ii {
		if i > max {
			max = i
		}
	}
	return max
}

func min(ii []int) int {
	min := ii[0]
	for _, i := range ii {
		if i < min {
			min = i
		}
	}
	return min
}

func stdev(ii []int) float64 {
	m := meanf(ii)
	var t float64
	for j := 0; j < len(ii); j++ {
		t += math.Pow(float64(ii[j])-m, 2)
	}
	return math.Sqrt(t / float64(len(ii)))

}

const DEBUG = false

func main() {
	settingo.SetInt("c", 5, "concurrency for the requests e.g workers")
	settingo.SetInt("n", 1000, "how many requests")
	settingo.Set("host", "http://127.0.0.1:8128", "host to run on")

	settingo.Parse()

	CONN := settingo.GetInt("c")
	REGS := settingo.GetInt("n")
	urls := make(chan string, 1000)
	responses := make(chan *Resp, 10)

	running := time.Now()
	for i := 1; i < CONN+1; i++ {
		go worker(i, urls, responses)
	}
	go func() {
		for i := 0; i < REGS; i++ {
			urls <- createQuery(settingo.Get("host"), "ams", 1, 10)
			urls <- createQuery(settingo.Get("host"), "drama", 1, 10)
			urls <- createQuery(settingo.Get("host"), "tvepis", 1, 10)
			urls <- createQuery(settingo.Get("host"), "amsterdam", 1, 10)
			urls <- createQuery(settingo.Get("host"), "hollywood", 1, 10)
		}
	}()

	durations := make(map[string][]int)

	for i := 0; i < (REGS * 5); i++ {
		resp := <-responses
		durations[resp.URL] = append(durations[resp.URL], resp.Duration)
	}

	fmt.Println("total,\tmean,\tmin,\tmax,\tstdev,\ttime,\tregs,\turl")
	for url, duration := range durations {
		fmt.Printf("%d,\t%d,\t%d,\t%d,\t%.2f,\t%v,\t%d,\t%s\n",
			sum(duration),
			mean(duration),
			min(duration),
			max(duration),
			stdev(duration),
			time.Now().Sub(running),
			len(durations),
			url,
		)
		//fmt.Println(
		//	"\nhost: ", createQuery(settingo.Get("host"), "ams", 1, 10),
		//	"\nurl: ", url,
		//	"\ntime:", time.Now().Sub(running),
		//	"\nreqs: ", len(durations),
		//	"\ntotal: ", sum(duration),
		//	"\nmean: ", mean(duration),
		//	"\nstdev: ", stdev(duration),
		//	"\nmin: ", min(duration),
		//	"\nmax: ", max(duration),
		//)
	}

	if DEBUG {
		fmt.Println("")
		for _, i := range durations {
			fmt.Print(",", i)
		}
		fmt.Println("")
	}

	//Query-Duration = [135ms]
	//Total-Items = [7628625]
	//Content-Type = [application/json]
	//Date = [Sun, 24 Oct 2021 18:32:38 GMT]
	//Total-Pages = [381432]

}
