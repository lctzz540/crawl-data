package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

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
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	existingData, err := reader.ReadAll()
	if err != nil {
		return err
	}

	id := int(data[0])
	for _, row := range existingData {
		existingID, err := strconv.Atoi(row[0])
		if err != nil {
			return err
		}
		if id == existingID {
			return nil
		}
	}

	writer := csv.NewWriter(file)
	stringData := make([]string, len(data))
	for i, val := range data {
		stringData[i] = strconv.FormatFloat(val, 'f', -1, 64)
	}
	err = writer.Write(stringData)
	if err != nil {
		return err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}

func updateData() {
	for id := 90000; id < 100000; id++ {
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

	fmt.Println("Data inserted successfully!")
}

func main() {
	updateData := flag.Bool("update", false, "Update Data")
	plotGraph := flag.Bool("plot", false, "Plot Graph")

	flag.Parse()

	if *updateData {
		// Call the updateData() function here
		fmt.Println("Updating data...")
		// Add your code logic for updating the data
	} else if *plotGraph {
		// Call the plotgraph() function here
		fmt.Println("Plotting graph...")
		// Add your code logic for plotting the graph
	} else {
		flag.PrintDefaults()
	}
}
