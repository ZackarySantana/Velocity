package jobs

import "go.mongodb.org/mongo-driver/mongo"

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
		job := Job{
			Command: "echo hello world",
			Image:   "alpine",
		}
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
