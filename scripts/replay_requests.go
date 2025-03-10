package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Item struct {
	User      *string           `json:"user,omitempty"`
	Total     int               `json:"total"`
	Hits      int               `json:"hits"`
	Time      string            `json:"time"`
	TimeMs    int               `json:"time_ms"`
	Endpoint  string            `json:"endpoint"`
	Query     *string           `json:"query"`
	Filters   map[string]string `json:"filters"`
	Page      *string           `json:"page"`
	Pagesize  *string           `json:"pagesize"`
	URL       string            `json:"url"`
	NewTime   *string           `json:"new_time,omitempty"`
	NewTimeMs *int              `json:"new_time_ms,omitempty"`
	NewHits   *string           `json:"new_hits,omitempty"`
}

func parseRawURL(rawURL string) (*string, string, *string, *string, *string, map[string]string, error) {
	fmt.Println("Raw URL (as is):", rawURL)
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, "", nil, nil, nil, nil, err
	}

	endpoint := strings.ReplaceAll(u.Path, "/", "")
	queryParams := u.Query()
	dictArguments := make(map[string]string)
	var query, page, pagesize, user *string

	for key, values := range queryParams {
		val := strings.TrimSpace(values[0])
		if len(val) >= 1 {
			switch key {
			case "search":
				queryStr := val
				query = &queryStr
			case "page":
				pageStr := val
				page = &pageStr
			case "pagesize":
				pageSizeStr := val
				pagesize = &pageSizeStr
			default:
				dictArguments[strings.TrimSpace(key)] = val
			}
		}
	}
	return user, endpoint, query, page, pagesize, dictArguments, nil
}

func createItem(line string) (*Item, error) {
	parts := strings.Split(line, " url: ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid log line format: %s", line)
	}

	infoPart := parts[0]
	urlPart := parts[1]

	re := regexp.MustCompile(`total: (\d+) hits: (\d+) time: (\d+)(ms|µs|s)`)
	matches := re.FindStringSubmatch(infoPart)
	if len(matches) != 5 {
		return nil, fmt.Errorf("failed to parse info part: %s", infoPart)
	}

	total, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid total value: %s", matches[1])
	}
	hits, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid hits value: %s", matches[2])
	}
	timeStr := matches[3] + matches[4]
	timeMs, err := parseTimeMs(matches[3], matches[4])
	if err != nil {
		return nil, err
	}

	user, endpoint, query, page, pagesize, filters, err := parseRawURL(urlPart)
	if err != nil {
		return nil, fmt.Errorf("parseRawURL error: %w", err)
	}

	item := &Item{
		User:     user,
		Total:    total,
		Hits:     hits,
		Time:     timeStr,
		TimeMs:   timeMs,
		Endpoint: endpoint,
		Query:    query,
		Filters:  filters,
		Page:     page,
		Pagesize: pagesize,
		URL:      strings.TrimSpace(urlPart),
	}
	return item, nil
}

func parseTimeMs(valueStr, unit string) (int, error) {
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid time value: %s", valueStr)
	}

	switch unit {
	case "ms":
		return value, nil
	case "µs":
		return value / 1000, nil
	case "s":
		return value * 1000, nil
	default:
		return 0, fmt.Errorf("unknown time unit: %s", unit)
	}
}

func valid(item *Item) bool {
	return true
}

func runner(scanner *bufio.Scanner, matches []string) []*Item {
	var items []*Item
	for scanner.Scan() {
		line := scanner.Text()
		for _, match := range matches {
			if strings.Contains(line, match) {
				item, err := createItem(line)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error creating item: %v, line: %s\n", err, line)
					continue
				}
				items = append(items, item)
				break
			}
		}
	}
	return items
}

func checkValid(value *string, data map[string]interface{}) bool {
	if value == nil {
		return true
	}
	dataData, ok := data["data"].([]interface{})
	if !ok {
		return false
	}
	for _, itemData := range dataData {
		itemMap, ok := itemData.(map[string]interface{})
		if !ok {
			continue
		}
		for _, val := range itemMap {
			if strVal, ok := val.(string); ok && strings.Contains(strVal, *value) {
				return true
			}
		}
	}
	return false
}

func runItemAgainstStaging(item *Item, checkValidityResult bool) *Item {
	host := "http://127.0.0.1:8127" // staging
	urlStr := item.URL
	headers := map[string]string{
		"X-auth": "eyJ0",
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", host+urlStr, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		return nil
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	// Print headers for debugging
	// for k, v := range resp.Header {
	// 	fmt.Printf("%s: %s\n", k, v)
	// }

	if checkValidityResult {
		var responseData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding JSON response: %v\n", err)
			return nil // or continue without validation check if JSON decode fails
		}
		if !checkValid(item.Query, responseData) {
			fmt.Println("failed result")
			jsonItem, _ := json.MarshalIndent(item, "", "  ")
			fmt.Println(string(jsonItem))
			jsonResp, _ := json.MarshalIndent(responseData, "", "  ")
			fmt.Println(string(jsonResp))
		}

		resp.Body.Close()
		resp, err = client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error re-making request for header extraction: %v\n", err)
			return nil
		}
		defer resp.Body.Close()
	}

	queryDuration := resp.Header.Get("Query-Duration")
	totalItems := resp.Header.Get("Total-Items")

	if queryDuration == "" || totalItems == "" {
		return nil
	}

	timeMsStr := strings.TrimSuffix(queryDuration, "ms")
	timeMs, err := strconv.Atoi(timeMsStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting Query-Duration to int: %v\n", err)
		return nil
	}

	item.NewTime = &queryDuration
	item.NewTimeMs = &timeMs
	item.NewHits = &totalItems

	return item
}

func calculatePercentile(times []int, percentile float64) int {
	if len(times) == 0 {
		return 0 // or handle empty case as needed
	}
	sort.Ints(times)
	index := int(float64(len(times)-1) * percentile)
	return times[index]
}

func calculateMedian(times []int) int {
	if len(times) == 0 {
		return 0
	}
	sort.Ints(times)
	n := len(times)
	if n%2 == 0 {
		return (times[n/2-1] + times[n/2]) / 2
	}
	return times[n/2]
}

func calculateAverage(times []int) float64 {
	if len(times) == 0 {
		return 0
	}
	sum := 0
	for _, t := range times {
		sum += t
	}
	return float64(sum) / float64(len(times))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	matches := []string{"search", "typeahead", "list"}
	checkValidityResult := false // TODO handle case
	totalRuntime := 0
	newTotalRuntime := 0
	count := 0

	var oldRuntimes []int
	var newRuntimes []int

	items := runner(scanner, matches)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	for _, item := range items {
		if valid(item) {
			count++
			updatedItem := runItemAgainstStaging(item, checkValidityResult)
			if updatedItem == nil {
				continue
			}
			totalRuntime += item.TimeMs
			newTotalRuntime += *updatedItem.NewTimeMs

			oldRuntimes = append(oldRuntimes, item.TimeMs)
			newRuntimes = append(newRuntimes, *updatedItem.NewTimeMs)

			jsonOutput, err := json.MarshalIndent(updatedItem, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
				continue
			}
			fmt.Println(string(jsonOutput))

		}
	}

	oldMedian := calculateMedian(oldRuntimes)
	newMedian := calculateMedian(newRuntimes)
	oldAverage := calculateAverage(oldRuntimes)
	newAverage := calculateAverage(newRuntimes)
	old95th := calculatePercentile(oldRuntimes, 0.95)
	new95th := calculatePercentile(newRuntimes, 0.95)

	old99th := calculatePercentile(oldRuntimes, 0.99)
	new99th := calculatePercentile(newRuntimes, 0.99)

	fmt.Println("\n--- Report ---")
	fmt.Printf("%-35s %-15d %-15s\n", "Total calls", count, "")

	fmt.Printf("%-35s %-15s %-15s\n", "Metric", "Old", "New")
	fmt.Printf("%-35s %-15d %-15d\n", "Total runtime (ms)", totalRuntime, newTotalRuntime)
	fmt.Printf("%-35s %-15.2f %-15.2f\n", "Average runtime (ms)", oldAverage, newAverage)
	fmt.Printf("%-35s %-15d %-15d\n", "Median runtime (ms)", oldMedian, newMedian)
	fmt.Printf("%-35s %-15d %-15d\n", "95th percentile runtime (ms)", old95th, new95th)
	fmt.Printf("%-35s %-15d %-15d\n", "99th percentile runtime (ms)", old99th, new99th)
}
