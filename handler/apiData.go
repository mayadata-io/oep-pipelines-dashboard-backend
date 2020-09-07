package handler

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"

	"github.com/mayadata-io/oep-pipelines-dashboard-backend/database"
)

// QueryData fetches the builddashboard data from the db
func QueryData(datas *dashboard, pipelineTable string, jobsTable string, page string) error {
	pageno, err := strconv.Atoi(page)
	if err != nil {
		glog.Infof("given page number is not a interger , %s replace this string with a number", page)
	}
	pipelineQuery := fmt.Sprintf("SELECT * FROM %s ORDER BY pipelineid DESC LIMIT 20 OFFSET %d;", pipelineTable, pageno*20)
	countQuery := fmt.Sprintf("SELECT count(*) FROM %s ;", pipelineTable)
	pipelinerows, err := database.Db.Query(pipelineQuery)

	if err != nil {
		return err
	}
	var counter int
	database.Db.QueryRow(countQuery).Scan(&counter)

	defer pipelinerows.Close()
	for pipelinerows.Next() {
		pipelinedata := pipelineSummary{}
		err = pipelinerows.Scan(
			&pipelinedata.PipelineID,
			&pipelinedata.Sha,
			&pipelinedata.Ref,
			&pipelinedata.Status,
			&pipelinedata.WebURL,
			&pipelinedata.ReleaseTag,
			&pipelinedata.Percentage,
			&pipelinedata.Total,
			&pipelinedata.ValidTestCount,
			&pipelinedata.KubernetesVersion,
			&pipelinedata.CreatedAt,
		)
		if err != nil {
			return err
		}
		jobsquery := fmt.Sprintf("SELECT pipelineid, id, status , stage , name , ref , created_at , started_at , finished_at , test_case_URL  FROM %s WHERE pipelineid = $1 ORDER BY id;", jobsTable)
		jobsrows, err := database.Db.Query(jobsquery, pipelinedata.PipelineID)
		if err != nil {
			return err
		}
		defer jobsrows.Close()
		jobsdataarray := []Jobssummary{}
		for jobsrows.Next() {
			jobsdata := Jobssummary{}
			err = jobsrows.Scan(
				&jobsdata.PipelineID,
				&jobsdata.ID,
				&jobsdata.Status,
				&jobsdata.Stage,
				&jobsdata.Name,
				&jobsdata.Ref,
				&jobsdata.CreatedAt,
				&jobsdata.StartedAt,
				&jobsdata.FinishedAt,
				&jobsdata.TestCaseURL,
			)
			if err != nil {
				return err
			}
			jobsdataarray = append(jobsdataarray, jobsdata)
			pipelinedata.Jobs = jobsdataarray
		}
		datas.Dashboard = append(datas.Dashboard, pipelinedata)
		datas.PipelineCount = counter
	}
	err = pipelinerows.Err()
	if err != nil {
		return err
	}
	return nil
}
