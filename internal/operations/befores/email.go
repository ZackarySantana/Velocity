package befores

import (
	"bufio"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func Email(ctx *cli.Context) error {
	providedEmail := ctx.Args().First()

	if providedEmail == "" {
		scanner := bufio.NewScanner(os.Stdin)

		// Prompt the user for their email address
		fmt.Print("Enter your email address: ")

		// Read the user's input
		scanner.Scan()
		providedEmail = scanner.Text()
	}

	ctx.App.Metadata["email"] = providedEmail
	return nil
}

func GetEmail(c *cli.Context) (*string, error) {
	email, ok := c.App.Metadata["email"].(string)
	if !ok {
		return nil, cli.Exit("error getting email", 1)
	}
	return &email, nil
}
