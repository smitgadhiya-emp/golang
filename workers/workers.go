package workers

import "log"

// StartAll launches every background worker. Call it once at startup, after
// RabbitMQ is connected and the queues are declared.
func StartAll() {
	if err := StartEmailWorker(); err != nil {
		log.Fatalf("Failed to start email worker: %v", err)
	}
}
