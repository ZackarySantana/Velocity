package jobs

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
			"go mod vendor",
		},
		Command:   "go test./...",
		Image:     "golang",
		Language:  "golang",
		Framework: "std",
	}

	landuageAndFrameworkDefaults = []LanguageAndFrameworkDefaults{
		golang_std,
	}
)

func getLanguageAndFrameworkDefaults(language, framework string) LanguageAndFrameworkDefaults {
	for _, l := range landuageAndFrameworkDefaults {
		if l.Language == language && l.Framework == framework {
			return l
		}
	}
	return LanguageAndFrameworkDefaults{}
}
