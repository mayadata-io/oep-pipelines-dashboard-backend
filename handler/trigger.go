package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/mayadata-io/oep-pipelines-dashboard-backend/database"
)

// openshiftCommit from gitlab api and store to database
func getPlatformData(token, project, branch, pipelineTable, jobTable string) {
	pipelineData, err := getPipelineData(token, project, branch)
	if err != nil {
		glog.Error(err)
		return
	}
	for i := range pipelineData {
		pipelineJobsData, err := getPipelineJobsData(pipelineData[i].ID, token, project)
		if err != nil {
			glog.Error(err)
			return
		}
		glog.Infoln("pipelieID :->  " + strconv.Itoa(pipelineData[i].ID) + " || JobSLegth :-> " + strconv.Itoa(len(pipelineJobsData)))
		sqlStatement := fmt.Sprintf("INSERT INTO %s (pipelineid, sha, ref, status, web_url) VALUES ($1, $2, $3, $4, $5)"+
			"ON CONFLICT (pipelineid) DO UPDATE SET status = $4 RETURNING pipelineid;", pipelineTable)
		id := 0
		err = database.Db.QueryRow(sqlStatement,
			pipelineData[i].ID,
			pipelineData[i].Sha,
			pipelineData[i].Ref,
			pipelineData[i].Status,
			pipelineData[i].WebURL,
		).Scan(&id)
		if err != nil {
			glog.Error(err)
		}
		glog.Infof("New record ID for %s Pipeline: %d", project, id)

		// Add pipeline jobs data to Database
		for j := range pipelineJobsData {
			sqlStatement := fmt.Sprintf("INSERT INTO %s (pipelineid, id, status, stage, name, ref, created_at, started_at, finished_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"+
				"ON CONFLICT (id) DO UPDATE SET status = $3, stage = $4, name = $5, ref = $6, created_at = $7, started_at = $8, finished_at = $9 RETURNING id;", jobTable)
			id := 0
			err = database.Db.QueryRow(sqlStatement,
				pipelineData[i].ID,
				pipelineJobsData[j].ID,
				pipelineJobsData[j].Status,
				pipelineJobsData[j].Stage,
				pipelineJobsData[j].Name,
				pipelineJobsData[j].Ref,
				pipelineJobsData[j].CreatedAt,
				pipelineJobsData[j].StartedAt,
				pipelineJobsData[j].FinishedAt,
			).Scan(&id)
			if err != nil {
				glog.Error(err)
			}
			glog.Infof("New record ID for %s pipeline Jobs: %d", project, id)
		}
	}
}
func getPipelineData(token, project, branch string) (Pipeline, error) {
	URL := BaseURL + "api/v4/projects/" + project + "/pipelines?ref=" + branch
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Connection", "close")
	req.Header.Add("PRIVATE-TOKEN", token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var obj Pipeline
	json.Unmarshal(data, &obj)
	return obj, nil
}

func getPipelineJobsData(pipelineID int, token string, project string) (Jobs, error) {
	// Generate pipeline jobs api url using BaseURL, pipelineID and OPENSHIFTID
	urlTmp := BaseURL + "api/v4/projects/" + project + "/pipelines/" + strconv.Itoa(pipelineID) + "/jobs?page="
	var obj Jobs
	for i := 1; ; i++ {
		url := urlTmp + strconv.Itoa(i)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Close = true
		// Set header for api request
		req.Header.Set("Connection", "close")
		req.Header.Add("PRIVATE-TOKEN", token)
		client := http.Client{
			Timeout: time.Minute * time.Duration(2),
		}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		if string(body) == "[]" {
			break
		}
		var tmpObj Jobs
		err = json.Unmarshal(body, &tmpObj)
		glog.Infoln("error ", err)
		obj = append(obj, tmpObj...)
	}
	return obj, nil
}
