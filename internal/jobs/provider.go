package jobs

import (
	"context"

	"github.com/zackarysantana/velocity/internal/db"
)

type JobProvider interface {
	Next(num int) ([]Job, error)
	Update(result JobResult) error
	Finished() (bool, error)
	Cleanup() error
}

type MongoDBJobProvider struct {
	c db.Connection
}

func NewMongoDBJobProvider(client db.Connection) *MongoDBJobProvider {
	return &MongoDBJobProvider{client}
}

func (p *MongoDBJobProvider) Next(num int) ([]Job, error) {
	// TODO should this be another context?
	dbJobs, err := p.c.DequeueNJobs(context.TODO(), int64(num))
	if err != nil {
		return nil, err
	}

	jobs := []Job{}
	for i := 0; i < len(dbJobs); i++ {
		job := NewCommandJob(dbJobs[i].Id.String(), dbJobs[i].Image, dbJobs[i].Command, dbJobs[i].SetupCommands, dbJobs[i].Status, nil)
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (p *MongoDBJobProvider) Update(result JobResult) error {
	return nil
}

func (p *MongoDBJobProvider) Cleanup() error {
	return nil
}

func (p *MongoDBJobProvider) Finished() (bool, error) {
	return false, nil
}

type InMemoryJobProvider struct {
	jobs    []Job
	i       int
	results []JobResult
}

func NewInMemoryJobProvider(jobs []Job) *InMemoryJobProvider {
	return &InMemoryJobProvider{jobs, 0, []JobResult{}}
}

func (p *InMemoryJobProvider) Next(num int) ([]Job, error) {
	jobs := []Job{}
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
