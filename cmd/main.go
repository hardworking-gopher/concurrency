package main

import (
	"fmt"
	"net/http"
)

func main() {
	links := []string{
		"https://www.google.com",
		"https://www.bing.com",
		"https://www.x.com",
		"https://www.facebook.com",
		"https://www.linkedin.com",
	}

	c := make(chan string)

	for _, v := range links {
		go ping(v, c)
	}

	for i := 0; i < len(links); i++ {
		fmt.Println(i+1, <-c)
	}

	fmt.Println("Done")
}

func ping(s string, c chan string) {
	_, err := http.Get(s)
	if err != nil {
		panic(err)
	} else {
		c <- s
	}
}
