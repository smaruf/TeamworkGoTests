package customerimporter

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sync"
)

func ProcessWithConcurrentStreaming(file *os.File) (map[string]int, error) {
	reader := csv.NewReader(file)
	ch := make(chan map[string]int)
	var domainCounts sync.Map
	var wg sync.WaitGroup

	processChunks(reader, ch, &wg)
	collectResults(ch, &domainCounts)

	return convertSyncMapToRegularMap(&domainCounts), nil
}

func processChunks(reader *csv.Reader, ch chan map[string]int, wg *sync.WaitGroup) {
	var chunk [][]string

	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV: %v", err)
			continue
		}
		chunk = append(chunk, fields)
		if len(chunk) >= ChunkSize {
			wg.Add(1)
			go processChunk(chunk, ch, wg)
			chunk = nil
		}
	}

	if len(chunk) > 0 {
		wg.Add(1)
		go processChunk(chunk, ch, wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
}

func processChunk(chunk [][]string, ch chan map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()
	localCounts := make(map[string]int)

	for _, fields := range chunk {
		if len(fields) < 3 {
			continue
		}
		domain := extractDomain(fields[2])
		if domain != "" {
			localCounts[domain]++
		}
	}
	ch <- localCounts
}

func collectResults(ch chan map[string]int, domainCounts *sync.Map) {
	for localCounts := range ch {
		for domain, count := range localCounts {
			domainCounts.Store(domain, count)
		}
	}
}

func convertSyncMapToRegularMap(domainCounts *sync.Map) map[string]int {
	finalCounts := make(map[string]int)
	domainCounts.Range(func(key, value interface{}) bool {
		finalCounts[key.(string)] = value.(int)
		return true
	})
	return finalCounts
}

// this code written with help of co-pilot
// tested with unit tests
