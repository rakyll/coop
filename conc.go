package conc

import (
	"time"
)

// Runs fn forever.
func Forever(fn func()) {
	go func() {
		for {
			fn()
		}
	}()
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
func Timeout(duration time.Duration, fn func()) (done <-chan bool) {
	timeout := make(chan bool)
	go func() {
		<-time.After(duration)
		timeout <- true
	}()
	go fn() // TODO: make sure fn has run
}

// Starts a job and returns a channel for cancellation signal.
// Once a message is sent to the channel, stops the fn.
func Cancel(fn func()) (cancel chan<- bool) {

}

// Runs fn at the specified time.
func At(time time.Time, fn func()) {

}

// Starts to run the given list of fns concurrently.
func Batch(fn ...func()) (done <-chan bool) {

}

// Starts to run the given list of fns concurrently,
// at most n fns at a time.
func BatchWithLimit(n int, fn ...func()) (done <-chan bool) {

}

// Run the same function with n copies.
func Duplicate(n int, fn func()) (done <-chan bool) {

}
