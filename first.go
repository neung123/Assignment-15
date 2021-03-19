package main

import (
	"errors"
	"time"
	"math/rand"
	"fmt"
	"golang.org/x/talks/content/2016/applicative/google"
)

var (
    replicatedWeb   = First(Web1, Web2)
    replicatedImage = First(Image1, Image2)
    replicatedVideo = First(Video1, Video2)
)

var (
	Web1   = FakeSearch("web1", "The Go Programming Language", "http://golang.org")
	Web2   = FakeSearch("web2", "The Go Programming Language", "http://golang.org")
	Image1 = FakeSearch("image1", "The Go gopher", "https://blog.golang.org/gopher/gopher.png")
	Image2 = FakeSearch("image2", "The Go gopher", "https://blog.golang.org/gopher/gopher.png")
	Video1 = FakeSearch("video1", "Concurrency is not Parallelism",
		"https://www.youtube.com/watch?v=cN_DpYBzKso")
	Video2 = FakeSearch("video2", "Concurrency is not Parallelism",
		"https://www.youtube.com/watch?v=cN_DpYBzKso")
)

type Result struct { Title, URL string }

var (
	Web   = FakeSearch("web", "The Go Programming Language", "http://golang.org")
	Image = FakeSearch("image", "The Go gopher", "https://blog.golang.org/gopher/gopher.png")
	Video = FakeSearch("video", "Concurrency is not Parallelism", "https://www.youtube.com/watch?v=cN_DpYBzKso")
)

type SearchFunc func(query string) Result

func FakeSearch(kind, title, url string) SearchFunc {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond) // HL
		return Result{
			Title: fmt.Sprintf("%s(%q): %s", kind, query, title),
			URL:   url,
		}
	}
}

func First(replicas ...SearchFunc) SearchFunc { // HL
	return func(query string) Result {
		c := make(chan Result, len(replicas))
		searchReplica := func(i int) {
			c <- replicas[i](query)
		}
		for i := range replicas {
			go searchReplica(i) // HL
		}
		return <-c
	}
}

func SearchReplicated(query string, timeout time.Duration) ([]Result, error) {
	timer := time.After(timeout)
	c := make(chan Result, 3)
    go func() { c <- replicatedWeb(query) }()
    go func() { c <- replicatedImage(query) }()
    go func() { c <- replicatedVideo(query) }()
	// STOP OMIT

	var results []Result
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timer:
			return results, errors.New("timed out")
		}
	}
	return results, nil
}

func main() {
    start := time.Now()
	results, err := google.SearchReplicated("golang", 80*time.Millisecond)
    elapsed := time.Since(start)
    fmt.Println(results)
    fmt.Println(elapsed, err)
}