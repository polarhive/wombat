package main

import (
	"log"
	"sync"
	"time"

	"github.com/polarhive/wombat/db"
	"github.com/polarhive/wombat/fetch"
)

const (
	seedURL       = "https://en.wikipedia.org/wiki/Palm_nut"
	// seedURL       = "https://en.wikipedia.org/wiki/Wombat"
	maxGoroutines = 5
	timeout       = 5
)

func main() {
	log.Println("Initializing database...")
	db.InitializeSQLiteDB("graph.db")
	defer db.CloseDB()

	log.Println("Starting to fetch and store links from seed URL:", seedURL)
	fetchAndStoreLinks(seedURL)
}

func fetchAndStoreLinks(url string) {
	pageName := fetch.ExtractPageName(url)
	if pageName == "" {
		log.Println("Invalid URL, cannot extract page name")
		return
	}

	log.Printf("Fetching links for page: %s\n", pageName)
	links, err := fetch.FetchLinksFromWikipedia(pageName)
	if err != nil {
		log.Printf("Error fetching links for page '%s': %v\n", pageName, err)
		return
	}

	if len(links) == 0 {
		log.Printf("No links found for page '%s'.\n", pageName)
	}

	log.Printf("Storing page '%s' into the database\n", pageName)
	db.SaveNode(pageName)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxGoroutines)

	for _, link := range links {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(link string) {
			defer wg.Done()
			log.Printf("Storing link: %s\n", link)
			db.SaveNode(link)
			db.SaveLink(pageName, link)
			time.Sleep(timeout * time.Millisecond)
			<-semaphore
		}(link)
	}

	wg.Wait()
	log.Printf("Finished processing links for page: %s\n", pageName)
}
