# coop

[![Build Status](https://travis-ci.org/rakyll/coop.png?branch=master)](https://travis-ci.org/rakyll/coop)

coop contains some of the most common concurrent program flows I personally use in Go. I'm suggesting you to use this package as a snippets reference/cheat sheet instead of a library. The functionally provided in this package can be obtained in many different ways, and frankly with more performant implementations depending on the type of your problem.

coop contains implementations for the following flows:

### coop.At(time, fn)

Runs fn at t, returns a boolean channel that will receive a message after fn returns. The following example, prints "Hello World" in a minute and blocks the goroutine its running in until fn is completed.

~~~ go
done := coop.At(time.Now().Add(time.Minute), func() {
    fmt.Println("Hello world")
})
<-done // wait for fn to be done
~~~

### coop.Until(time, duration, fn)

Runs fn once in every provided duration until t, returns a boolean channel that will receive a message after fn returns. The following example prints "Hello world" every minute until tomorrow, and blocks the goroutine its running in until the job is completed.

~~~ go
done := coop.Until(time.Now().Add(24*time.Hour), time.Minute, func() {
    fmt.Println("Hello world")
})
<-done
~~~

### coop.After(duration, fn)

Runs fn after duration, returns a boolean channel that will receive a message after fn returns. The following example prints "Hello world" after a second and blocks until fn is completed.

~~~ go
done := coop.After(time.Second, func() {
    fmt.Println("Hello world")
})
<-done
~~~

### coop.Every(duration, fn)

Runs fn once in every duration, and never stops. The following example will print "Hello World" once in every second.

~~~ go
coop.Every(time.Second, func() {
    fmt.Println("Hello world")
})
~~~

### coop.Timeout(duration, fn)
Runs fn, and cancels the running job if timeout is exceeded. The following example will timeout and fn will return immediately ("Hello world will not printed"), the value read from the done channel will be false if timeout occurs, true if fn is completed.

~~~ go
done := coop.Timeout(time.Second, func() {
    time.Sleep(time.Hour)
    fmt.Println("Hello world")
})
<-done // will return false, because timeout occurred
~~~

### coop.All(fns...)
Runs the list of fns concurrently, returns a boolean channel that will receive a message after all of the fns are completed. The following example will start 4 printing jobs concurrently and wait until all of them are completed.

~~~ go
printFn := func() {
    fmt.Println("Hello world")
}
<-coop.All(printFn, printFn, printFn, printFn)
~~~

### coop.AllWithThrottle(num, fns...)
Similar to coop.All, but with limiting. Runs the list of fns concurrently, but at most num fns at a time. Returns a boolean channel that will receive a message after all of the fns are completed. The following example will start 3 printing jobs immediately, and run the left out one once the first 3 is completed. It will block the goroutine until all 4 are finished.

~~~ go
printFn := func() {
    fmt.Println("Hello world")
}
<-coop.AllWithThrottle(3, printFn, printFn, printFn, printFn)
~~~

### coop.Replicate(n, fn)

Runs fn n time concurrently, returns a boolean channel that indicates all runs are completed. The following example prints "Hello world" 5 times, and waits for all printing jobs are finished.

~~~ go
<-coop.Replicate(5, func() {
    fmt.Println("Hello world")
})
~~~

## License

Copyright 2014 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License. ![Analytics](https://ga-beacon.appspot.com/UA-46881978-1/coop?pixel)
