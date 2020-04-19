FROM golang:latest 

RUN go get github.com/lib/pq

RUN go get github.com/golang/glog

RUN mkdir /app

RUN mkdir $GOPATH/src/github.com/mayadata-io

RUN mkdir $GOPATH/src/github.com/mayadata-io/oep-pipelines-dashboard-backend

ADD . $GOPATH/src/github.com/mayadata-io/oep-pipelines-dashboard-backend

WORKDIR $GOPATH/src/github.com/mayadata-io/oep-pipelines-dashboard-backend

RUN go build -o /app/main .

CMD ["/app/main"]

EXPOSE 3000