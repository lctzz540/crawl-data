package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func crawlDataByID(id int) ([]float64, error) {
	url := "https://vietnamnet.vn/giao-duc/diem-thi/tra-cuu-diem-thi-vao-lop-10-2023/" + strconv.Itoa(id) + ".html"

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("URL returned 404")
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	math, _ := strconv.ParseFloat(doc.Find("table tbody tr:first-child td:nth-child(2)").Text(), 64)
	literature, _ := strconv.ParseFloat(doc.Find("table tbody tr:nth-child(2) td:nth-child(2)").Text(), 64)
	english, _ := strconv.ParseFloat(doc.Find("table tbody tr:nth-child(3) td:nth-child(2)").Text(), 64)
	specialized, _ := strconv.ParseFloat(doc.Find("table tbody tr:nth-child(4) td:nth-child(2)").Text(), 64)

	return []float64{float64(id), math, literature, english, specialized}, nil
}

func insertData(data []float64, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	existingData, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read existing data: %w", err)
	}

	id := int(data[0])
	for _, row := range existingData {
		existingID, err := strconv.Atoi(row[0])
		if err != nil {
			return fmt.Errorf("failed to convert existing ID: %w", err)
		}
		if id == existingID {
			return nil // ID already exists, no need to insert
		}
	}

	writer := csv.NewWriter(file)
	stringData := make([]string, len(data))
	for i, val := range data {
		stringData[i] = strconv.FormatFloat(val, 'f', -1, 64)
	}

	err = writer.Write(stringData)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func updateData() {
	var wg sync.WaitGroup
	idCh := make(chan int)

	// Define the number of threads (goroutines)
	numThreads := 10

	// Launch goroutines
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range idCh {
				data, err := crawlDataByID(id)
				if err != nil {
					continue
				}
				err = insertData(data, "data.csv")
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("added", id)
				}
			}
		}()
	}

	// Send IDs to the channel
	for id := 90000; id < 1000000; id++ {
		idCh <- id
	}

	// Close the channel to signal that no more IDs will be sent
	close(idCh)

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("Data inserted successfully!")
}

func removeDuplicateData(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	allData, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	// Map to track existing IDs
	existingIDs := make(map[int]bool)
	var uniqueData [][]string

	for _, row := range allData {
		id, err := strconv.Atoi(row[0])
		if err != nil {
			return fmt.Errorf("failed to convert ID: %w", err)
		}

		// Check if ID already exists
		if !existingIDs[id] {
			existingIDs[id] = true
			uniqueData = append(uniqueData, row)
		}
	}

	// Print the number of data entries before and after removing duplicates
	fmt.Println("Number of data entries before removing duplicates:", len(allData))
	fmt.Println("Number of data entries after removing duplicates:", len(uniqueData))

	// Truncate the file
	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	// Move file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	// Write the unique data back to the file
	writer := csv.NewWriter(file)
	err = writer.WriteAll(uniqueData)
	if err != nil {
		return fmt.Errorf("failed to write unique data: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func main() {
	update := flag.Bool("update", false, "Update Data")
	plot := flag.Bool("plot", false, "Plot Graph")
	removeDuplicate := flag.Bool("remove_duplicate", false, "Remove duplicate data from data.csv")

	flag.Parse()

	if *update {
		fmt.Println("Updating data...")
		updateData()
	} else if *plot {
		fmt.Println("Plotting graph...")
		plotGraph()
	} else if *removeDuplicate {
		fmt.Println("Removing duplicate data...")
		_ = removeDuplicateData("data.csv")
	} else {
		flag.PrintDefaults()
	}
}
