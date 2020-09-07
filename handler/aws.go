package handler

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
)

// AWSHandler return packet pipeline data to /packet path
func AWSHandler(w http.ResponseWriter, r *http.Request) {
	// Allow cross origin request
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	pageNo := getparams(r, "page")
	datas := dashboard{}
	err := QueryData(&datas, "aws_pipelines", "aws_pipelines_jobs", pageNo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		glog.Error(err)
	}

	out, err := json.Marshal(datas)
	if err != nil {
		http.Error(w, err.Error(), 500)
		glog.Error(err)
	}
	w.Write(out)

}

func getparams(r *http.Request, param string) (response string) {
	pages, ok := r.URL.Query()[param]
	if !ok {
		return "0"
	}
	return pages[0]
}
