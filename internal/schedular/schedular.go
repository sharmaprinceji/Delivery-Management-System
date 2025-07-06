package schedular

import (
	"log"
	"time"

	"github.com/sharmaprinceji/delivery-management-system/internal/jobs"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)

func SchedularJob(storage storage.Storage) {
	go func() {
		for {
			now := time.Now()
			if now.Hour() == 7 {
				log.Println("Running allocation job...")
				err := jobs.AllocateOrders(storage)
				if err != nil {
					log.Printf("Job failed: %v", err)
				}
			}
			time.Sleep(time.Minute * 10)
		}
	}()
}
