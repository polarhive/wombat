package crawler

import (
	"log"
	"sync"
	"time"

	"github.com/polarhive/wombat/db"
	"github.com/polarhive/wombat/fetch"
)

const (
	maxGoroutines = 100
	timeout       = 0
)

var dbMutex sync.Mutex

// FetchAndStoreLinks crawls a Wikipedia page and stores links recursively.
func FetchAndStoreLinks(url string, depth int, output string) {
	if depth <= 0 {
		return
	}

	pageName := fetch.ExtractPageName(url)
	if pageName == "" {
		log.Println("Invalid URL, cannot extract page name")
		return
	}

	log.Printf("Fetching links for page: %s\n", pageName)
	links, _ := fetch.FetchLinksFromWikipedia(pageName)

	if len(links) == 0 {
		log.Printf("No links found for page '%s'.\n", pageName)
	}

	log.Printf("Storing page '%s' into the database\n", pageName)

	// Lock the mutex before interacting with the database
	dbMutex.Lock()
	db.SaveNode(pageName)
	dbMutex.Unlock()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxGoroutines)

	for _, link := range links {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(link string) {
			defer wg.Done()
			log.Printf("Storing link: %s\n", link)

			// ensure only one goroutine writes to the DB at a time
			dbMutex.Lock()
			db.SaveNode(link)
			db.SaveLink(pageName, link)
			dbMutex.Unlock()

			time.Sleep(timeout * time.Millisecond)

			// Recursively process the link as a new seed URL
			FetchAndStoreLinks("https://en.wikipedia.org/wiki/"+link, depth-1, output)

			<-semaphore
		}(link)
	}

	wg.Wait()
	log.Printf("Finished processing links for page: %s\n", pageName)

	if output == "json" {
	} else {
		// Default DB output
		log.Println("Data stored in the database.")
	}
}
