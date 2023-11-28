package jobs

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type JobProvider interface {
	Next(num int) ([]Job, error)
	Update(result JobResult) error
	Cleanup() error
}

type MongoDBJobProvider struct {
	client mongo.Client
}

func NewMongoDBJobProvider(client mongo.Client) *MongoDBJobProvider {
	return &MongoDBJobProvider{client}
}

func (p *MongoDBJobProvider) Next(num int) ([]Job, error) {
	jobs := []Job{}
	for i := 0; i < num; i++ {
		job := CommandJob{
			Command: "echo hello world",
			Image:   "alpine",
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

func (p *MongoDBJobProvider) Update(result JobResult) error {
	return nil
}

func (p *MongoDBJobProvider) Cleanup() error {
	return nil
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

func (p *InMemoryJobProvider) Results() []JobResult {
	return p.results
}
