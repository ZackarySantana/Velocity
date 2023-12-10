package jobs

import (
	"context"
	"strings"

	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/internal/jobs/jobtypes"
	"github.com/zackarysantana/velocity/src/clients"
	"github.com/zackarysantana/velocity/src/clients/v1types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobProvider interface {
	Next(num int) ([]*Job, error)
	Update(result JobResult) error
	Finished() (bool, error)
	Cleanup() error
}

type VelocityJobProvider struct {
	v *clients.VelocityClientV1
}

func NewVelocityJobProvider(v *clients.VelocityClientV1) *VelocityJobProvider {
	return &VelocityJobProvider{v}
}

func (p *VelocityJobProvider) Next(num int) ([]*Job, error) {
	resp, err := p.v.PostJobsDequeue(v1types.PostJobsDequeueRequest{}, v1types.PostJobsDequeueQueryParams{})
	if err != nil {
		return nil, err
	}
	dbJobs := resp.Jobs

	jobs := []*Job{}
	for i := 0; i < len(dbJobs); i++ {
		var job Job = NewCommandJob(dbJobs[i].Id.Hex(), dbJobs[i].Image, dbJobs[i].Command, dbJobs[i].SetupCommands, dbJobs[i].Status, nil)
		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (p *VelocityJobProvider) Update(result JobResult) error {
	j := v1types.PostJobResultRequest{
		Id: result.Job.GetName(),
	}
	if result.Failed != nil {
		err := result.Failed.Error.Error()
		j.Error = &err
	}
	if result.Success != nil {
		logs := strings.TrimSuffix(result.Success.Logs, "\n")
		j.Logs = &logs
	}
	_, err := p.v.PostJobsResults(j)

	return err
}

func (p *VelocityJobProvider) Cleanup() error {
	return nil
}

func (p *VelocityJobProvider) Finished() (bool, error) {
	return false, nil
}

type MongoDBJobProvider struct {
	c db.Connection
}

func NewMongoDBJobProvider(client db.Connection) *MongoDBJobProvider {
	return &MongoDBJobProvider{client}
}

func (p *MongoDBJobProvider) Next(num int) ([]*Job, error) {
	// TODO should this be another context?
	dbJobs, err := p.c.DequeueNJobs(context.TODO(), int64(num))
	if err != nil {
		return nil, err
	}

	jobs := []*Job{}
	for i := 0; i < len(dbJobs); i++ {
		var job Job = NewCommandJob(dbJobs[i].Id.Hex(), dbJobs[i].Image, dbJobs[i].Command, dbJobs[i].SetupCommands, dbJobs[i].Status, nil)
		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (p *MongoDBJobProvider) Update(result JobResult) error {
	id, err := primitive.ObjectIDFromHex(result.Job.GetName())
	if err != nil {
		return err
	}
	// TODO: Replace this context?
	j := &db.Job{
		Id:     id,
		Status: jobtypes.JobStatusCompleted,
	}
	if result.Failed != nil {
		j.Error = result.Failed.Error.Error()
	}
	if result.Success != nil {
		j.Logs = strings.TrimSuffix(result.Success.Logs, "\n")
	}
	_, err = p.c.UpdateJob(context.TODO(), j)

	return err
}

func (p *MongoDBJobProvider) Cleanup() error {
	return nil
}

func (p *MongoDBJobProvider) Finished() (bool, error) {
	return false, nil
}

type InMemoryJobProvider struct {
	jobs    []*Job
	i       int
	results []JobResult
}

func NewInMemoryJobProvider(jobs []*Job) *InMemoryJobProvider {
	return &InMemoryJobProvider{jobs, 0, []JobResult{}}
}

func (p *InMemoryJobProvider) Next(num int) ([]*Job, error) {
	jobs := []*Job{}
	for limit := p.i + num; p.i < limit && p.i < len(p.jobs); p.i++ {
		jobs = append(jobs, p.jobs[p.i])
	}
	return jobs, nil
}

func (p *InMemoryJobProvider) Update(result JobResult) error {
	p.results = append(p.results, result)
	return nil
}

func (p *InMemoryJobProvider) Cleanup() error {
	return nil
}

func (p *InMemoryJobProvider) Finished() (bool, error) {
	return p.i == len(p.jobs), nil
}

func (p *InMemoryJobProvider) Results() []JobResult {
	return p.results
}
