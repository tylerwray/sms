package ingest

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sync"

	"github.com/tylerwray/sms/internal/messenger"
)

func FromCSV(ms *messenger.Service, fileName string) {
	f, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)

	var wg sync.WaitGroup
	maxGoroutines := 50
	guard := make(chan struct{}, maxGoroutines)

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)

		guard <- struct{}{}

		go func(smsXID, status string) {
			defer wg.Done()

			err = ms.UpdateMessageStatus(smsXID, status)

			if err != nil {
				log.Printf("ERROR: Could not update message status for smsXID:%s. Error: %v\n", smsXID, err)
			}

			<-guard
		}(record[0], record[1])
	}

	wg.Wait()
}
