package jobs

import (
	"fmt"
	"os/exec"
	"strings"
)

type JobExecutor interface {
	Execute(ctx Context, job Job) (string, error)
}

type DockerJobExecutor struct{}

func (e *DockerJobExecutor) Execute(ctx Context, job Job) (string, error) {
	// Build base image that most jobs will use
	buildCmd := exec.Command("docker", "build", "-t", "velocity_repository_clone", "-")
	dockerfileContent := fmt.Sprintf(`
		FROM alpine/git
		RUN git clone %s app
		WORKDIR app
		RUN git fetch --all
		RUN git checkout %s
	`, ctx.RepositoryURL, ctx.CommitHash)
	buildCmd.Stdin = strings.NewReader(dockerfileContent)
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error building repository clone image: error '%s' output '%s'", err, output)
	}

	// Build the test image
	testBuildCmd := exec.Command("docker", "build", "-t", "velocity_test_name", "-")
	dockerfileContent = fmt.Sprintf(`
		FROM %s
		COPY --from=velocity_repository_clone /git/app /app
		WORKDIR /app
		%s
		CMD %s
	`, job.GetImage(), parseCommands(job.SetupCommand()), job.GetCommand())
	testBuildCmd.Stdin = strings.NewReader(dockerfileContent)
	output, err = testBuildCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error building test image: error '%s' output '%s'", err, output)
	}

	// Run the test image
	testCmd := exec.Command("docker", "run", "--rm", "velocity_test_name")
	testOutput, err := testCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running test image: error '%s' output '%s'", err, testOutput)
	}
	return string(testOutput), err
}

func parseCommands(commands []string) string {
	return strings.Join(commands, "\n")
}
