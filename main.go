package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	i := 0
requestLabel:
	for i < 4 {
		i++
		response, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
		if err != nil {
			continue
		}

		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil || response.StatusCode != 200 {
			continue
		}

		values := strings.Split(string(body), ",")
		if len(values) != 7 {
			continue
		}
		numbers := make([]int, len(values))
		for n, val := range values {
			parsedNumber, err := strconv.Atoi(val)
			if err != nil {
				continue requestLabel
			}
			numbers[n] = parsedNumber
		}

		if numbers[0] > 30 {
			fmt.Printf("Load Average is too high: %d\n", numbers[0])
		}

		if usage := (numbers[2] / numbers[1]) * 100; usage > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", usage)
		}

		if (numbers[4]/numbers[3])*100 > 90 {
			fmt.Printf("Free disk space is too low: %d Mb left\n", (numbers[3]-numbers[4])/1000/1000)
		}

		if (numbers[6]/numbers[5])*100 > 90 {
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", (numbers[5]-numbers[6])/1000/1000)
		}
		break
	}
	if i == 4 {
		fmt.Println("Unable to fetch server statistic")
	}
}
