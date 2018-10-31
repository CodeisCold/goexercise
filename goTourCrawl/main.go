package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func crawl(url string, depth int, fetcher Fetcher, vl *sync.WaitGroup) {
	if depth <= 0 {
		return
	}
	vl.Add(1)
	//body, urls, err := fetcher.Fetch(url)
	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("found: %s %q\n", url, body)
	fmt.Printf("found: %s\n", url)
	for _, u := range urls {
		go crawl(u, depth-1, fetcher, vl)
	}
	vl.Done()
	return
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	var groupLock sync.WaitGroup
	//fmt.Println(groupLock)
	crawl(url, depth, fetcher, &groupLock)
	//fmt.Println(groupLock)
	//groupLock.Add(1)
	groupLock.Wait()
	return
}

func getUrls(body string) ([] string) {
	v := regexp.MustCompile(`href="(https?.*?)"`)
	s := v.FindAllStringSubmatch(body, -1)
	urls := make([]string, len(s), len(s))
	for i, a := range s {
		urls[i] = a[1]
	}
	return urls
}

func main() {
	Crawl("http://www.baidu.com", 2, fetcher)
	for i, v := range fetcher.res {
		fmt.Println(i)
		fmt.Println(v.urls)
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher struct {
	mux sync.Mutex
	res map[string]*fakeResult
}

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	f.mux.Lock()
	_, ok := f.res[url]
	f.mux.Unlock()
	// if exist
	if ok {
		f.mux.Lock()
		defer f.mux.Unlock()
		return f.res[url].body, f.res[url].urls, nil
	} else {
		var fRes *fakeResult
		response, err := http.Get(url)
		if err != nil {
			fRes = &fakeResult{"http error: " + err.Error(), nil}
		}
		res, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			fRes = &fakeResult{"read error: " + err.Error(), nil}
		}
		body := string(res)
		urls := getUrls(body)
		fRes = &fakeResult{body, urls}

		f.mux.Lock()
		f.res[url] = fRes
		f.mux.Unlock()
		return fRes.body, fRes.urls, nil
	}
}

var fetcher = fakeFetcher{
	res: make(map[string]*fakeResult),
}

//// fetcher is a populated fakeFetcher.
//var fetcher = fakeFetcher{
//	"https://golang.org/": &fakeResult{
//		"The Go Programming Language",
//		[]string{
//			"https://golang.org/pkg/",
//			"https://golang.org/cmd/",
//		},
//	},
//	"https://golang.org/pkg/": &fakeResult{
//		"Packages",
//		[]string{
//			"https://golang.org/",
//			"https://golang.org/cmd/",
//			"https://golang.org/pkg/fmt/",
//			"https://golang.org/pkg/os/",
//		},
//	},
//	"https://golang.org/pkg/fmt/": &fakeResult{
//		"Package fmt",
//		[]string{
//			"https://golang.org/",
//			"https://golang.org/pkg/",
//		},
//	},
//	"https://golang.org/pkg/os/": &fakeResult{
//		"Package os",
//		[]string{
//			"https://golang.org/",
//			"https://golang.org/pkg/",
//		},
//	},
//}
