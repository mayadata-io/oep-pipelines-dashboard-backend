package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/mayadata-io/oep-pipelines-dashboard-backend/database"
)

// MetricsDailyPipelines testing data
func MetricsDailyPipelines(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
	w.Header().Set("Access-Control-Expose-Headers", "Authorization") // Double check it's a post request being made

	// if r.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	fmt.Fprintf(w, "invalid_http_method")
	// 	return
	// }

	platform := r.Form.Get("platform")
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// fmt.Println(r)
	fmt.Println("EndpointHit hit for " + platform + " Platform . from :" + start + "  to : " + end + " ;")
	// select pipelineid, created_at from aws_pipelines where created_at between '2020-06-01' and '2020-06-06' ORDER BY created_at DESC;
	outpp := MetricDashboard{}
	err := Testingxyz(&outpp, platform, start, end)
	if err != nil {
		http.Error(w, err.Error(), 500)
		glog.Error(err)
	}
	out, err := json.Marshal(outpp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		glog.Error(err)
	}
	w.Write(out)

}

// Testingxyz test function
func Testingxyz(outpp *MetricDashboard, platform, start, end string) error {
	pipelineQuery := fmt.Sprintf("select pipelineid, created_at from %s where created_at between %s and %s ORDER BY created_at DESC;", platform, start, end)
	pipelinerows, err := database.Db.Query(pipelineQuery)
	if err != nil {
		return err
	}
	defer pipelinerows.Close()
	for pipelinerows.Next() {
		pipelinedata := PostDailyMetrics{}
		err = pipelinerows.Scan(
			&pipelinedata.PipelineID,
			&pipelinedata.Date,
		)
		if err != nil {
			return err
		}
		outpp.Dashboard = append(outpp.Dashboard, pipelinedata)
	}
	err = pipelinerows.Err()
	if err != nil {
		return err
	}
	return nil

}
