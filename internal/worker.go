package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/ratelimit"
)

func executeRequest(tasks *[]int, rl ratelimit.Limiter, prev time.Time, amount int) int {
	var wg sync.WaitGroup
	var errCount int

	errChan := make(chan error, 100)
	errors := 0

	for i := amount - 1; i >= 0; i-- {
		now := rl.Take()
		//fmt.Println(i, now.Sub(prev))

		wg.Add(1)
		go func() {
			defer wg.Done()

			resp, err := http.Get("http://localhost:8080/hc")
			if err != nil {
				errChan <- fmt.Errorf("error request: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errChan <- fmt.Errorf("error response status code: %v", resp.StatusCode)
				return
			}
		}()
		prev = now
		*tasks = append((*tasks)[:i], (*tasks)[i+1:]...)
	}

	wg.Wait()

	close(errChan)
	for err := range errChan {
		fmt.Println(err)
		errCount++
	}

	fmt.Printf("Counted Errors: %d\n", errCount)
	return errors
}

func getInitialRateLimit(tasks []int, hours int) int {
	total := len(tasks)
	byHour := total / hours
	byMinute := byHour / 60
	bySecond := byMinute / 60
	if bySecond == 0 {
		return 5
	}
	return bySecond

}

func shouldIncreaseAmount(prev, current int) bool {
	return current <= prev
}

func main() {
	var amountByExecution int
	var previousErrorsCount int

	hours := 2
	errorsCount := -1

	tasks := make([]int, 1000)
	for i := 0; i < len(tasks); i++ {
		tasks[i] = i + 1
	}

	rl := getInitialRateLimit(tasks, hours)
	amountByExecution = rl * 60
	rateLimit := ratelimit.New(rl)
	sizeTask := len(tasks)
	for i := 0; i < sizeTask; i += amountByExecution {
		if errorsCount >= 0 {
			if shouldIncreaseAmount(previousErrorsCount, errorsCount) {
				amountByExecution = int(float64(amountByExecution) * 1.2)
				if amountByExecution > len(tasks) {
					amountByExecution = len(tasks)
				}

				rl = amountByExecution / 60
				if rl == 0 {
					rl = 1
				}

				rateLimit = ratelimit.New(rl)
			}
			previousErrorsCount = errorsCount
		}
		prev := time.Now()
		fmt.Println("Throughput (req/s): ", rl)
		errorsCount = executeRequest(&tasks, rateLimit, prev, amountByExecution)

	}
}
