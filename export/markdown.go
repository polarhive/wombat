package export

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

func ExportMarkdownFromDB(dbPath string) error {
	database, _ := sql.Open("sqlite3", dbPath)
	defer database.Close()

	outputDir := "out/vault"
	os.MkdirAll(outputDir, os.ModePerm)

	nodes, _ := database.Query("SELECT DISTINCT id FROM nodes")
	defer nodes.Close()

	var wg sync.WaitGroup
	workerPool := make(chan struct{}, 50)

	for nodes.Next() {
		var nodeID string
		nodes.Scan(&nodeID)

		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			workerPool <- struct{}{}
			defer func() { <-workerPool }()

			fmt.Printf("Processing node: %s\n", id) // the current processing article

			nodeFilePath := filepath.Join(outputDir, id+".md")

			relatedLinks, _ := database.Query(`
				SELECT n2.id AS link
				FROM links l
				JOIN nodes n1 ON l.source = n1.id
				JOIN nodes n2 ON l.target = n2.id
				WHERE n1.id = ?`, id)
			defer relatedLinks.Close()

			linksToAdd := make(map[string]bool)
			for relatedLinks.Next() {
				var link string
				relatedLinks.Scan(&link)
				linksToAdd[link] = true
			}

			var contentBuilder strings.Builder
			contentBuilder.WriteString(fmt.Sprintf("# %s\n\nRelated Links:\n\n", id))

			var existingContent []byte
			if _, err := os.Stat(nodeFilePath); err == nil {
				existingContent, _ = os.ReadFile(nodeFilePath)
			}

			for link := range linksToAdd {
				if !strings.Contains(string(existingContent), fmt.Sprintf("[[%s]]", link)) {
					contentBuilder.WriteString(fmt.Sprintf("- [[%s]]\n", link))
				}
			}

			if len(contentBuilder.String()) > 0 {
				os.WriteFile(nodeFilePath, append(existingContent, []byte(contentBuilder.String())...), 0644)
			}
		}(nodeID)
	}

	wg.Wait()
	return nil
}
