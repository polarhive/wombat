package export

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
)

// CSV with a hierarchical color scheme.
func ExportToCSV(dbPath, outputFile string) {
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer database.Close()

	csvFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error creating CSV file: %v\n", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	header := []string{"source", "target", "color"}
	writer.Write(header)

	// get the links between nodes (source and target)
	rows, err := database.Query(`
		SELECT n1.id AS source, n2.id AS target 
		FROM links l
		JOIN nodes n1 ON l.source = n1.id
		JOIN nodes n2 ON l.target = n2.id
	`)
	defer rows.Close()

	// (each node will have a unique color)
	nodeColors := make(map[string]string)

	for rows.Next() {
		var source, target string
		if err := rows.Scan(&source, &target); err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}

		sourceColor, exists := nodeColors[source]
		if !exists {
			sourceColor = generateRandomColor()
			nodeColors[source] = sourceColor
		}
		targetColor, exists := nodeColors[target]
		if !exists {
			targetColor = generateRandomColor()
			nodeColors[target] = targetColor
		}

		record := []string{source, target, targetColor}
		if err := writer.Write(record); err != nil {
			log.Fatalf("Error writing row to CSV file: %v\n", err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating over rows: %v\n", err)
	}

	log.Println("Export completed successfully!")
}

// a random hex color code.
func generateRandomColor() string {
	r := fmt.Sprintf("%02X", randInt(0, 255))
	g := fmt.Sprintf("%02X", randInt(0, 255))
	b := fmt.Sprintf("%02X", randInt(0, 255))
	return "#" + r + g + b
}

// min and max (inclusive).
func randInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}
