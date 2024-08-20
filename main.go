package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	workers, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	RPS, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		panic(err)
	}
	var interval = float64(time.Second) * float64(workers) / RPS

	beginning := time.Now()

	db := NewDatabase()
	closer := db.Run()
	defer func() { close(closer) }()

	var wg sync.WaitGroup
	chans := make([]chan struct{}, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		done := make(chan struct{})
		chans[i] = done
		go func() {
			rng := rand.New(rand.NewSource(int64(i)))
			next := time.Now().Add(time.Duration(rng.Float64() * interval))
			for {
				sleep := time.Until(next)
				timer := time.NewTimer(sleep)
				select {
				case <-done:
					wg.Done()
					return
				case <-timer.C:
					start := time.Now()
					db.Request()
					fmt.Printf("%d\t%f\t%f\t%f\n", workers, RPS, time.Since(beginning).Seconds(), time.Since(start).Seconds())
					next = start.Add(time.Duration(rng.ExpFloat64() * interval))
				}
			}
		}()
	}

	time.Sleep(10 * time.Second)
	for i := 0; i < workers; i++ {
		chans[i] <- struct{}{}
		close(chans[i])
	}

	wg.Wait()
}
