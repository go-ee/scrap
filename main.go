package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func download(url string, base string, pattern string, to string) {
	if url == "" || base == "" || pattern == "" || to == "" {
		log.Println("missing url/base/pattern/to arguments")
		return
	}
	log.Println("visiting", url)

	c := colly.NewCollector()

	visit := make(chan string)

	// count links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link != "" {
			visit <- link
		}
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	go func() {
		toVisit := regexp.MustCompile(base)
		toDownload := regexp.MustCompile(pattern)
		links := make(map[string]int)

		for {
			link := <-visit
			links[link]++
			if val, _ := links[link]; val == 1 && toVisit.MatchString(link) {
				log.Printf("visit %v\n", link)

				if toDownload.MatchString(link) {
					go downloadFromUrl(link, to)
				}

				go c.Visit(link)
			}
		}
	}()

	c.Visit(url)

	time.Sleep(3 * time.Hour)
}

func downloadFromUrl(url string, to string) {
	tokens := strings.Split(url, "/")
	fileName := fmt.Sprintf("%v/%v", to, tokens[len(tokens)-1])
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}

func main() {
	args := os.Args[1:]
	if len(args) == 4 {
		download(args[0], args[1], args[2], args[3])
	} else {
		log.Println("missing url/base/pattern/to arguments")
	}
}
