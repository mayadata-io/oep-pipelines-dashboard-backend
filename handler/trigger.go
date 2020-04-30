package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/mayadata-io/oep-pipelines-dashboard-backend/database"
)

// openshiftCommit from gitlab api and store to database
func getPlatformData(token, project, branch, pipelineTable, jobTable string) {
	var releaseImageTag, percentageCoverage, totalTestCoverage, validTestCount string
	// var percentageCoverage string

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
		releaseImageTag, err = getReleaseImageTag(pipelineJobsData, token)
		if err != nil {
			glog.Error(err)
		}
		percentageCoverage, totalTestCoverage, validTestCount, err = percentageCoverageFunc(pipelineJobsData, token)
		if err != nil {
			glog.Error(err)
			return
		}
		glog.Infoln("pipeline :- "+strconv.Itoa(pipelineData[i].ID)+" \n Total Coverage :- ", totalTestCoverage+" : Percentage :- "+percentageCoverage+" validTestCount:- "+validTestCount)
		glog.Infoln("pipelieID :->  " + strconv.Itoa(pipelineData[i].ID) + " || JobSLegth :-> " + strconv.Itoa(len(pipelineJobsData)))
		sqlStatement := fmt.Sprintf("INSERT INTO %s (pipelineid, sha, ref, status, web_url, release_tag, coverage, total_coverage_count, valid_test_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"+
			"ON CONFLICT (pipelineid) DO UPDATE SET status = $4, release_tag = $6, coverage = $7, total_coverage_count = $8, valid_test_count = $9 RETURNING pipelineid;", pipelineTable)
		id := 0
		err = database.Db.QueryRow(sqlStatement,
			pipelineData[i].ID,
			pipelineData[i].Sha,
			pipelineData[i].Ref,
			pipelineData[i].Status,
			pipelineData[i].WebURL,
			releaseImageTag,
			percentageCoverage,
			totalTestCoverage,
			validTestCount,
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
func getReleaseImageTag(jobsData Jobs, token string) (string, error) {
	var jobURL string
	for _, value := range jobsData {
		if strings.Contains(value.Name, "TCID-DIR-INSTALL") {
			jobURL = value.WebURL + "/raw"
		}
	}
	if jobURL == "" {
		return "NA", nil
	}
	req, err := http.NewRequest("GET", jobURL, nil)
	if err != nil {
		return "NA", err
	}
	req.Close = true
	req.Header.Set("Connection", "close")
	client := http.Client{
		Timeout: time.Minute * time.Duration(1),
	}
	res, err := client.Do(req)
	if err != nil {
		return "NA", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	data := string(body)
	if data == "" {
		return "NA", err
	}
	re := regexp.MustCompile("release: [^ ]*")
	value := re.FindString(data)
	if value == "" {
		return "NA", nil
	}
	result := strings.Split(string(value), ":")
	if result != nil && len(result) != 0 {
		if result[1] == "" {
			return "NA", nil
		}
		s := result[1]
		t := strings.Replace(s, "\n", "", -1)
		return t, nil
	}
	return "NA", nil
}

func percentageCoverageFunc(jobsData Jobs, token string) (string, string, string, error) {
	// var jobURL = "https://gitlab.mayadata.io/oep/oep-e2e-gcp/-/jobs/38871/raw"
	var jobURL string
	for _, value := range jobsData {
		if value.Name == "e2e-metrics" {
			jobURL = value.WebURL + "/raw"
		}
	}
	if jobURL != "" {
		req, err := http.NewRequest("GET", jobURL, nil)
		if err != nil {
			return "", "NA", "NA", err
		}
		req.Close = true
		req.Header.Set("Connection", "close")
		client := http.Client{
			Timeout: time.Minute * time.Duration(1),
		}
		res, err := client.Do(req)
		if err != nil {
			return "", "NA", "NA", err
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		data := string(body)
		if data == "" {
			return "", "NA", "NA", err
		}
		re := regexp.MustCompile("coverage: [^ ]*")
		value := re.FindString(data)
		totalCount := regexp.MustCompile("count: [^ ]*")
		totalValue := totalCount.FindString(data)
		validTestsRegex := regexp.MustCompile("\\WvalidTestCount: [^ ]*")
		validTestCountValue := validTestsRegex.FindString(data)
		var totalAutomatedTests, coveragePercentage, validTestCount string
		if validTestCountValue != "" {
			validTestCountString := strings.Split(string(validTestCountValue), ":")
			if len(validTestCountString) != 0 {
				splitValidCount := strings.Split(validTestCountString[1], "\n")
				validTestCount = splitValidCount[0]
			} else {
				validTestCount = "NA"
			}
		}
		if totalValue != "" {
			splitCountString := strings.Split(string(totalValue), ":")
			if len(splitCountString) != 0 {
				splitLine := strings.Split(splitCountString[1], "\n")
				totalAutomatedTests = splitLine[0]
			} else {
				totalAutomatedTests = "NA"
			}
		}
		if value != "" {
			splitTotalString := strings.Split(string(value), ":")
			if len(splitTotalString) != 0 {
				coveragePercentage = splitTotalString[1]
			} else {
				coveragePercentage = "NA"
			}
		}
		return coveragePercentage, totalAutomatedTests, validTestCount, nil
	}
	return "", "NA", "NA", nil
}
