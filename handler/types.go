package handler

// PlatformID assignment from gitlab repository
const (
	RANCHERID = "36"
	KONVOYID  = "37"
	AWSID     = "39"
)

// BranchName assignment from gitlab repository
const (
	GROUPNAME   = "openebs"
	MAYABRANCH  = "master"
	JIVABRANCH  = "master"
	ISTGTBRANCH = "replication"
	ZFSBRANCH   = "develop"
)

var project string

const (
	token = "TOKEN"
)

// BaseURL for gitlab
var BaseURL = "https://gitlab.mayadata.io/"

// PlatformPipeline struct
type PlatformPipeline struct {
	ID     int    `json:"id"`
	Sha    string `json:"sha"`
	Ref    string `json:"ref"`
	Status string `json:"status"`
	WebURL string `json:"web_url"`
}

// Pipeline struct
type Pipeline []struct {
	ID        int    `json:"id"`
	Sha       string `json:"sha"`
	Ref       string `json:"ref"`
	Status    string `json:"status"`
	WebURL    string `json:"web_url"`
	CreatedAt string `json:"created_at"`
	Jobs      Jobs   `json:"jobs"`
}

// Commit  dsdsdss
type Commit []struct {
	ID     string `json:"id"`
	Sha    string `json:"short_id"`
	Ref    string `json:"title"`
	Status string `json:"message"`
	WebURL string `json:"author_name"`
}

//CommitData wdsd
type CommitData struct {
	ProjectID      int    `json:"project_id"`
	Sha            string `json:"id"`
	CommittedDate  string `json:"committed_date"`
	AuthorName     string `json:"author_name"`
	AuthorEmail    string `json:"author_email"`
	CommitterName  string `json:"committer_name"`
	CommitterEmail string `json:"committer_email"`
	CommitMessage  string `json:"title"`
	CommitPipeline struct {
		ID     int    `json:"id"`
		Sha    string `json:"sha"`
		Status string `json:"status"`
		WebURL string `json:"web_url"`
		Ref    string `json:"ref"`
	} `json:"last_pipeline"`
}

// Jobs struct
type Jobs []struct {
	ID         int    `json:"id"`
	Status     string `json:"status"`
	Stage      string `json:"stage"`
	Name       string `json:"name"`
	Ref        string `json:"ref"`
	CreatedAt  string `json:"created_at"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
	WebURL     string `json:"web_url"`
}

// Jobssummary struct
type Jobssummary struct {
	PipelineID  int    `json:"pipelineid"`
	ID          int    `json:"id"`
	Status      string `json:"status"`
	Stage       string `json:"stage"`
	Name        string `json:"name"`
	Ref         string `json:"ref"`
	CreatedAt   string `json:"created_at"`
	StartedAt   string `json:"started_at"`
	FinishedAt  string `json:"finished_at"`
	TestCaseURL string `json:"test_case_URL"`
	// JobLogURL  string `json:"job_log_url"`
}

// TriggredID contains the details of a gitlab pipelines
type TriggredID struct {
	ID          int    `json:"id"`
	ProjectID   int    `json:"project_id"`
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	Message     string `json:"message"`
}

// pipelineSummary contains the details of a gitlab pipelines
type pipelineSummary struct {
	PipelineID int `json:"pipelineid"`
	// ProjectID   int           `json:"projectid"`
	Sha               string        `json:"sha"`
	Ref               string        `json:"ref"`
	Status            string        `json:"status"`
	WebURL            string        `json:"web_url"`
	ReleaseTag        string        `json:"release_tag"`
	Percentage        string        `json:"coverage"`
	Total             string        `json:"total_coverage_count"`
	ValidTestCount    string        `json:"valid_test_count"`
	KubernetesVersion string        `json:"kubernetes_version"`
	CreatedAt         string        `json:"created_at"`
	Jobs              []Jobssummary `json:"jobs"`
}
type dashboard struct {
	Dashboard     []pipelineSummary `json:"dashboard"`
	PipelineCount int               `json:"pipelineCount"`
}

// BuildJobs struct
type BuildJobs []struct {
	ID         int    `json:"id"`
	Status     string `json:"status"`
	Stage      string `json:"stage"`
	Name       string `json:"name"`
	Ref        string `json:"ref"`
	CreatedAt  string `json:"created_at"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
	Message    string `json:"message"`
	AuthorName string `json:"author_name"`
	WebURL     string `json:"web_url"`
	Commit     struct {
		Message    string `json:"message"`
		AuthorName string `json:"author_name"`
	} `json:"commit"`
}

// BuildJobssummary contains the details of builds job for database
type BuildJobssummary struct {
	PipelineID int    `json:"pipelineid"`
	ID         int    `json:"id"`
	Status     string `json:"status"`
	Stage      string `json:"stage"`
	Name       string `json:"name"`
	Ref        string `json:"ref"`
	CreatedAt  string `json:"created_at"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
	Message    string `json:"message"`
	AuthorName string `json:"author_name"`
}

// BuildpipelineSummary contains the details of a builds pipelines
type BuildpipelineSummary struct {
	Project       int    `json:"project"`
	ID            int    `json:"id"`
	Sha           string `json:"sha"`
	Ref           string `json:"ref"`
	Status        string `json:"status"`
	WebURL        string `json:"web_url"`
	Committeddate string `json:"committeddate"`
	AuthorName    string `json:"author_name"`
	AuthorEmail   string `json:"author_email"`
	ComitterName  string `json:"comitter_name"`
	CommitTitle   string `json:"commit_title"`
	CommitMessage string `json:"commit_message"`
}

// Builddashboard contains the details related to builds
type Builddashboard struct {
	Dashboard []BuildpipelineSummary `json:"dashboard"`
}

// PostDailyMetrics contains the details of postData
type PostDailyMetrics struct {
	PipelineID string `json:"pipelineid"`
	Date       string `json:"created_at"`
}

// MetricDashboard contains the details of postData
type MetricDashboard struct {
	Dashboard []PostDailyMetrics `json:"dashboard"`
}
