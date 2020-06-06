package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/mayadata-io/oep-pipelines-dashboard-backend/database"
	"github.com/mayadata-io/oep-pipelines-dashboard-backend/handler"
)

func main() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	// Initailze Db connection
	database.InitDb()
	http.HandleFunc("/api/rancher", handler.RancherHandler)
	http.HandleFunc("/api/aws", handler.AWSHandler)
	http.HandleFunc("/api/konvoy", handler.KonvoyHandler)
	// postReq endpoints
	http.HandleFunc("/api/metrics/", handler.MetricsDailyPipelines)

	// OepPipelineHandler
	glog.Infof("Listening on http://0.0.0.0:3000")

	// Trigger db update function
	go handler.UpdateDatabase()
	glog.Info(http.ListenAndServe(":"+"3000", nil))
}
