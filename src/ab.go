package main

import (
	"ab"
)

func main() {
	done := make(chan int)
	ab.ParseConfig()
	ab.Do(done)
	<-done
}
