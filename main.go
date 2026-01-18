package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	errors := 0
requestLabel:
	for {
		if errors >= 3 {
			fmt.Println("Unable to fetch server statistic")
		}

		response, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
		if err != nil {
			errors++
			continue
		}

		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil || response.StatusCode != 200 {
			errors++
			continue
		}

		values := strings.Split(string(body), ",")
		if len(values) != 7 {
			errors++
			continue
		}
		numbers := make([]int, len(values))
		for n, val := range values {
			parsedNumber, err := strconv.Atoi(val)
			if err != nil {
				errors++
				continue requestLabel
			}
			numbers[n] = parsedNumber
		}

		if numbers[0] > 30 {
			fmt.Printf("Load Average is too high: %d\n", numbers[0])
		}

		if usage := float64(numbers[2]) / float64(numbers[1]) * 100; usage > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", int32(usage))
		}

		if (float64(numbers[4])/float64(numbers[3]))*100 > 90 {
			fmt.Printf("Free disk space is too low: %d Mb left\n", (numbers[3]-numbers[4])/1024/1024)
		}

		if (float64(numbers[6])/float64(numbers[5]))*100 > 90 {
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", (numbers[5]-numbers[6])/1024/1024)
		}
	}
}
