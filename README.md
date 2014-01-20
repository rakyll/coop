# coop

coop contains some of the most common concurrent program flows I personally use in Go. I'm suggesting you to use this package as a snippets reference/cheat sheet instead of a library. The functionally provided in this package can be obtained in many different ways, and frankly in  more performant ways depending on the type of your problem.

coop contains implementations for the following flows:

* coop.At(time, fn)
* coop.Until(time, fn)
* coop.After(duration, fn)
* coop.Every(duration, fn)
* coop.Timeout(duration, fn)
* coop.Die(fn)
* coop.All(fns...)
* coop.AllWithThrottle(num, fns...)
* coop.Replicate(n, fn)

## License

Copyright 2014 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License. ![Analytics](https://ga-beacon.appspot.com/UA-46881978-1/coop?pixel)
