package schedular

import (
	"log"
	"time"
)

//this schedular job by using gorotines
func StartStudentFetchJob() {
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		 for range ticker.C {
			log.Println("Running schedular student fetch job...")
		 }
	}()
}