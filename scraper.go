package main

import (
	"bufio"
	"io"
	"net/http"
)

type Report struct {
	LatestTopic string
	TotalTopics int
}

func GetReportStats() Report {
	resp, err := http.Get("https://rickvanfessem.nl/assets/travelReports/config.txt")

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	count := 0
	var topic string
	for _, line := range lines(resp.Body) {
		count++
		topic = line
	}

	return Report{
		LatestTopic: topic[:len(topic)-4],
		TotalTopics: count,
	}
}

func lines(reader io.Reader) []string {
	// Create new Scanner.
	scanner := bufio.NewScanner(reader)
	result := []string{}
	// Use Scan.
	for scanner.Scan() {
		line := scanner.Text()
		// Append line to result.
		result = append(result, line)
	}
	return result
}
