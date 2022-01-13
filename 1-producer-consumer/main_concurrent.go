//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"

	ms "github.com/go-concurrency-exercises/1-producer-consumer/mockstream"
)

func producerConc(stream ms.Stream, tweets chan<- *ms.Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ms.ErrEOF {
			close(tweets)
		} else {
			tweets <- tweet
		}
	}
}

func consumerConc(tweets <-chan *ms.Tweet) {
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func concurrent() {
	start := time.Now()
	stream := ms.GetMockStream()

	// tweets channel
	tweets := make(chan *ms.Tweet)

	// Producer
	go func() {
		producerConc(stream, tweets)
	}()

	// Consumer
	consumerConc(tweets)

	fmt.Printf("Process took %s\n", time.Since(start))
}
