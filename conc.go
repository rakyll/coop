package conc

import (
	"sync"
	"time"
)

// Runs fn at the specified time.
func At(t time.Time, fn func()) {
	After(t.Sub(time.Now()), fn)
}

// Runs until time in every dur.
func Util(t time.Time, dur time.Duration, fn func()) {
	if time.Now().Sub(t) > 0 {
		After(dur, func() {
			fn()
			Util(t, dur, fn)
		})
	}
}

// Runs fn after duration. Similar to time.AfterFunc
func After(duration time.Duration, fn func()) {
	time.AfterFunc(duration, fn)
}

// Runs fn in every specified duration.
func Every(dur time.Duration, fn func()) {
	time.AfterFunc(dur, func() {
		fn()
		Every(dur, fn)
	})
}

// Runs fn and times out if it runs longer than the provided
// duration. It will send false to the returning
// channel if timeout occurs.
// TODO: cancel if timeout occurs
func Timeout(duration time.Duration, fn func()) (done <-chan bool) {
	timeout := make(chan bool)
	ch := make(chan bool)
	go func() {
		<-time.After(duration)
		timeout <- true
	}()
	go func() {
		fn()
		ch <- true
	}()
	return ch
}

// Starts a job and returns a channel for cancellation signal.
// Once a message is sent to the channel, stops the fn.
func Cancel(fn func()) (cancel chan<- bool, done <-chan bool) {
	ch := make(chan bool)
	cch := make(chan bool)
	go func() {
		select {
		case <-cch:
			return
		default:
			fn()
			ch <- true
		}
	}()
	return cch, ch
}

// Starts to run the given list of fns concurrently.
func All(fns ...func()) (done <-chan bool) {
	var wg sync.WaitGroup
	wg.Add(len(fns))

	ch := make(chan bool)
	for _, fn := range fns {
		go func() {
			fn()
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}

// Starts to run the given list of fns concurrently,
// at most n fns at a time.
func AllWithLimit(n int, fn ...func()) (done <-chan bool) {
	ch := make(chan bool)
	panic("not implemented")
	return ch
}

// Run the same function with n copies.
func Duplicate(n int, fn func()) (done <-chan bool) {
	var wg sync.WaitGroup
	wg.Add(n)

	ch := make(chan bool)
	for i := 0; i < n; i++ {
		go func() {
			fn()
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}
