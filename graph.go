package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateLabels(start, end, step float64) []string {
	labels := make([]string, 0)
	for i := start; i <= end; i += step {
		labels = append(labels, strconv.FormatFloat(i, 'f', 2, 64))
	}
	return labels
}

func generateBarItems(counts map[float64]int) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0.0; i <= 30.0; i += 0.25 {
		count := int(counts[i])
		items = append(items, opts.BarData{Value: count})
	}
	return items
}
func plotTotalGraph() {
	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var totalScores []float64
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")
		score := 0.0
		for i := 1; i <= 3; i++ {
			subjectScore, err := strconv.ParseFloat(strings.TrimSpace(fields[i]), 64)
			if err != nil {
				panic(err)
			}
			score += subjectScore
		}
		totalScores = append(totalScores, score)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	counts := make(map[float64]int)
	for _, score := range totalScores {
		for i := 0.0; i <= 30.0; i += 0.25 {
			if score >= i && score < i+0.25 {
				counts[i]++
				break
			}
		}
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Total Score Distribution",
			Subtitle: "Count of scores within specific ranges. Total: " + strconv.Itoa(len(totalScores)),
		}),
	)

	bar.SetXAxis(generateLabels(0.0, 30.0, 0.25)).AddSeries("Count", generateBarItems(counts))

	f, err := os.Create("bar_chart_total.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = bar.Render(f)
	if err != nil {
		panic(err)
	}
}

func plotGraph() {
	for subject := 1; subject <= 3; subject++ {
		file, err := os.Open("data.csv")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var subjectTitle string
		if subject == 1 {
			subjectTitle = "Math"
		} else if subject == 2 {
			subjectTitle = "Literature"
		} else if subject == 3 {
			subjectTitle = "English"
		} else {
			panic("Invalid subject number")
		}

		var subjectScores []float64
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Split(line, ",")
			score, err := strconv.ParseFloat(strings.TrimSpace(fields[subject]), 64)
			if err != nil {
				panic(err)
			}
			subjectScores = append(subjectScores, score)
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		counts := make(map[float64]int)
		for _, score := range subjectScores {
			for i := 0.0; i <= 10.0; i += 0.25 {
				if score >= i && score < i+0.25 {
					counts[i]++
					break
				}
			}
		}

		bar := charts.NewBar()
		bar.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title:    subjectTitle + " Score Distribution",
				Subtitle: "Count of scores within specific ranges. Total: " + strconv.Itoa(len(subjectScores)),
			}),
		)

		bar.SetXAxis(generateLabels(0.0, 10.0, 0.25)).AddSeries("Count", generateBarItems(counts))

		f, err := os.Create("bar_chart_" + strings.ToLower(subjectTitle) + ".html")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = bar.Render(f)
		if err != nil {
			panic(err)
		}
	}
}
