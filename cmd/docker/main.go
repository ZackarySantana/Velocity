package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {

	// docker build -t my_hello_image - <<EOF
	// FROM alpine/git

	// git clone https://github.com/zackarysantana/velocity.git app

	// cd app

	// git checkout 1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z

	// RUN echo hello
	// EOF

	buildCmd := exec.Command("docker", "build", "-t", "velocity_repository_clone", "-")
	dockerfileContent := `
		FROM alpine/git
		RUN git clone https://github.com/zackarysantana/velocity.git app
		WORKDIR app
		RUN git checkout c8dc99dfc0b62842b0a524fe34112c3df27f7e86
	`
	buildCmd.Stdin = strings.NewReader(dockerfileContent)

	buildOutput, err := buildCmd.CombinedOutput()
	log.Printf("Build output: %s", buildOutput)
	if err != nil {
		fmt.Println("error")
		log.Fatal(err)
	}

	fmt.Println("Done with velocity_repository_clone")

	// Build the test image
	testBuildCmd := exec.Command("docker", "build", "-t", "velocity_test_name", "-")
	dockerfileContent = `
		FROM golang:latest
		COPY --from=velocity_repository_clone /git/app /app
		WORKDIR /app
		RUN go mod vendor
		CMD ["go", "mod", "vendor"]
	`
	testBuildCmd.Stdin = strings.NewReader(dockerfileContent)

	testBuildOutput, err := testBuildCmd.CombinedOutput()
	log.Printf("Test build output: %s", testBuildOutput)
	if err != nil {
		fmt.Println("error")
		log.Fatal(err)
	}
	// Build the test image
	testCmd := exec.Command("docker", "run", "--rm", "velocity_test_name", "sh", "-c", "go test -v ./...")

	testOutput, err := testCmd.CombinedOutput()
	log.Printf("Test output: %s", testOutput)
	if err != nil {
		fmt.Println("error")
		log.Fatal(err)
	}
	// Run `docker rmi velocity_test_name` to remove the test image
}
