package ingest

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/tylerwray/sms/internal/messenger"
)

func FromCSV(ms *messenger.Service, fileName string) {
	f, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)

	// var wg sync.WaitGroup

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		var smsXID = record[0]
		var status = record[1]
		// wg.Add(1)

		// go func(smsXID, status string) {
		// defer wg.Done()

		err = ms.UpdateMessageStatus(smsXID, status)

		if err != nil {
			log.Printf("ERROR: Could not update message status for smsXID:%s. Error: %v\n", smsXID, err)
		}
		// }(record[0], record[1])
	}

	// wg.Wait()
}
