package jobs

import "strings"

type LanguageAndFrameworkDefaults struct {
	SetupCommands []string
	Command       string
	Image         string

	Language  string
	Framework string
}

var (
	none = LanguageAndFrameworkDefaults{
		SetupCommands: []string{},
		Command:       "echo 'That language/framework combination is not supported.'",
		Image:         "alpine",
		Language:      "",
		Framework:     "",
	}

	golang_std = LanguageAndFrameworkDefaults{
		SetupCommands: []string{
			"RUN go mod vendor",
		},
		Command:   "go test ./...",
		Image:     "golang",
		Language:  "golang",
		Framework: "std",
	}

	javascript_jest = LanguageAndFrameworkDefaults{
		SetupCommands: []string{
			"RUN npm install",
		},
		Command:   "npm test",
		Image:     "node",
		Language:  "javascript",
		Framework: "jest",
	}

	python_pytest = LanguageAndFrameworkDefaults{
		SetupCommands: []string{
			"RUN pip install -r requirements.txt",
		},
		Command:   "pytest",
		Image:     "python",
		Language:  "python",
		Framework: "pytest",
	}

	landuageAndFrameworkDefaults = []LanguageAndFrameworkDefaults{
		golang_std,
		javascript_jest,
		python_pytest,
	}
)

func getLanguageAndFrameworkDefaults(language, framework string) LanguageAndFrameworkDefaults {
	l := strings.ToLower(language)
	f := strings.ToLower(framework)
	for _, d := range landuageAndFrameworkDefaults {
		if d.Language == l && d.Framework == f {
			return d
		}
	}
	return LanguageAndFrameworkDefaults{}
}

func getDirectoryCommands(directory string) []string {
	// TODO: Should we remove a leading slash?
	if directory == "" {
		return []string{}
	}
	return []string{
		"WORKDIR " + directory,
	}
}
