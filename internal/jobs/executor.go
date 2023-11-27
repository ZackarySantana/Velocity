package jobs

import "os/exec"

type JobExecutor interface {
	Execute(job Job) (string, error)
}

type DockerJobExecutor struct{}

func (e *DockerJobExecutor) Execute(job Job) (string, error) {
	cmd := exec.Command("docker", "run", "--rm", job.Image, "sh", "-c", job.Command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
