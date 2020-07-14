package handler

import (
	"fmt"

	"github.com/mayadata-io/oep-pipelines-dashboard-backend/database"
)

// QueryData fetches the builddashboard data from the db
func QueryData(datas *dashboard, pipelineTable string, jobsTable string) error {
	pipelineQuery := fmt.Sprintf("SELECT * FROM %s ORDER BY pipelineid DESC LIMIT 20;", pipelineTable)
	pipelinerows, err := database.Db.Query(pipelineQuery)
	if err != nil {
		return err
	}
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
	}
	err = pipelinerows.Err()
	if err != nil {
		return err
	}
	return nil
}
