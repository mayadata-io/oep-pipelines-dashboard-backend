package handler

import (
	"os"
	"time"

	"github.com/golang/glog"
)

// UpdateDatabase will update the latest pipelines detail and status
// TODO
func UpdateDatabase() {
	// Read token environment variable
	token, ok := os.LookupEnv(token)
	if !ok {
		glog.Fatalf("TOKEN environment variable required")
	}
	go getPlatformData(token, RANCHERID, "oep-release", "rancher_pipelines", "rancher_pipelines_jobs") //e2e-rancher
	// Update the database, This wil run only first time
	// loop will iterate at every 2nd minute and update the database
	tick := time.Tick(2 * time.Minute)
	for range tick {
		// commitData(token)
		// pipelineData(token)

	}
}
