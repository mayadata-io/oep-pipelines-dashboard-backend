[![Go Report Card](https://goreportcard.com/badge/github.com/mayadata-io/oep-pipelines-dashboard-backend)](https://goreportcard.com/report/github.com/mayadata-io/oep-pipelines-dashboard-backend)
#[![Build status](https://img.shields.io/gitlab/pipeline/mayadata-io/oep-pipelines-dashboard-backend.svg?color=green&gitlab_url=https%3A%2F%2Fgitlab.openebs.ci&style=plastic)](https://gitlab.openebs.ci/openebs/oep-pipelines-dashboard-backend/pipelines)
#[![BCH compliance](https://bettercodehub.com/edge/badge/openebs/oep-pipelines-dashboard-backend?branch=master)](https://bettercodehub.com/)

OpenEBS CI-E2E Status

## Pre-requisites for k8s cluster

```bash
kubectl apply -f https://openebs.github.io/charts/openebs-operator-0.8.0.yaml
kubectl apply -f https://raw.githubusercontent.com/openebs/openebs/0.8/k8s/openebs-storageclasses.yaml
https://github.com/openebs/openebs.git
cd openebs/k8s/demo/crunchy-postgres/
ls -ltr
./run.sh
```

## Pre-requisite for localhost

1 - Postgress running in local
2 - Export variable like following ...

```bash
export DBHOST=<db_host>
export DBPORT=<db_port>
export DBUSER=<db_user>
export DBPASS=<db_password>
export DBNAME=<db_name>
export TOKEN=<gitlab_token>
```

3 - run the main file

example:

```bash
go run main.go
```
