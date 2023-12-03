package jobs

import (
	"errors"
	"fmt"
)

type JobStatus string

var (
	JobStatusQueued    JobStatus = "queued"
	JobStatusInactive  JobStatus = "inactive"
	JobStatusActive    JobStatus = "active"
	JobStatusCompleted JobStatus = "completed"

	JobStatuses = []JobStatus{
		JobStatusQueued,
		JobStatusInactive,
		JobStatusActive,
		JobStatusCompleted,
	}
)

func JobStatusFromString(s string) (JobStatus, error) {
	for _, status := range JobStatuses {
		if s == string(status) {
			return status, nil
		}
	}
	return "", fmt.Errorf("invalid status %s", s)
}

type Job interface {
	SetupCommand() []string
	GetImage() string
	GetCommand() string
	GetName() string
	GetStatus() JobStatus
	SetStatus(JobStatus)
	Validate() error
}

type BaseJob struct {
	SetupCommands []string
	Image         string
	Command       string
	Name          string
	Status        JobStatus
}

func (j *BaseJob) SetupCommand() []string {
	return j.SetupCommands
}

func (j *BaseJob) GetImage() string {
	return j.Image
}

func (j *BaseJob) GetCommand() string {
	return j.Command
}

func (j *BaseJob) GetName() string {
	return j.Name
}

func (j *BaseJob) GetStatus() JobStatus {
	return j.Status
}

func (j *BaseJob) SetStatus(status JobStatus) {
	j.Status = status
}

func (j *BaseJob) Validate() error {
	if j.Name == "" {
		return errors.New("missing name")
	}
	if j.Image == "" {
		return fmt.Errorf("missing image on %s", j.Name)
	}
	if j.Command == "" {
		return fmt.Errorf("missing command %s", j.Name)
	}
	if j.Status == "" {
		return fmt.Errorf("missing status %s", j.Name)
	}
	if j.SetupCommands == nil {
		j.SetupCommands = []string{}
	}
	return nil
}

func (j *BaseJob) Populate() *BaseJob {
	if j.SetupCommands == nil {
		j.SetupCommands = []string{}
	}
	return j
}

type CommandJob struct {
	BaseJob
}

type CommandJobOptions struct {
	Directory *string
}

func NewCommandJob(name string, image string, command string, setupCommands []string, status JobStatus, opts *CommandJobOptions) *CommandJob {
	j := &CommandJob{
		BaseJob: BaseJob{
			SetupCommands: setupCommands,
			Image:         image,
			Command:       command,
			Name:          name,
			Status:        status,
		},
	}
	j.Populate()
	if opts == nil {
		return j
	}

	if opts.Directory != nil {
		j.SetupCommands = append(j.SetupCommands, getDirectoryCommands(*opts.Directory)...)
	}

	return j
}

type FrameworkJob struct {
	BaseJob
}

type FrameworkJobOptions struct {
	Directory *string
	Image     *string
}

func NewFrameworkJob(name, language, framework string, status JobStatus, opts *FrameworkJobOptions) *FrameworkJob {
	i := getLanguageAndFrameworkDefaults(language, framework)
	j := &FrameworkJob{
		BaseJob: BaseJob{
			SetupCommands: i.SetupCommands,
			Image:         i.Image,
			Command:       i.Command,
			Name:          name,
			Status:        status,
		},
	}
	j.Populate()
	if opts == nil {
		return j
	}

	if opts.Directory != nil {
		// Reverse the order of the setup commands so that the directory is cd'd into first
		j.SetupCommands = append(i.SetupCommands, j.SetupCommands...)
	}

	if opts.Image != nil {
		j.Image = *opts.Image
	}

	return j
}
