package schedular

import (
	"log"
	"time"

	"github.com/sharmaprinceji/delivery-management-system/internal/jobs"
	"github.com/sharmaprinceji/delivery-management-system/internal/storage"
)

func SchedularJob(s storage.Storage) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			duration := next.Sub(now)

			log.Printf("⏰ Allocation job scheduled at: %v", next)

			time.Sleep(duration)

			log.Println("🚚 Running auto allocation job...")
			if err := jobs.AllocateOrders(s); err != nil {
				log.Printf("❌ Auto allocation error: %v", err)
			} else {
				log.Println("✅ Auto allocation completed.")
			}
		}
	}()
}
