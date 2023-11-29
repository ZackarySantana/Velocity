package jobs

type Job interface {
	SetupCommand() []string
	GetImage() string
	GetCommand() string
	GetName() string
}

type CommandJob struct {
	Image     string
	Directory *string

	Command string
	Name    string
}

func (j *CommandJob) SetupCommand() []string {
	cmds := []string{}
	if j.Directory != nil {
		cmds = append(cmds, getDirectoryCommands(*j.Directory)...)
	}
	return cmds
}

func (j *CommandJob) GetImage() string {
	return j.Image
}

func (j *CommandJob) GetCommand() string {
	return j.Command
}

func (j *CommandJob) GetName() string {
	return j.Name
}

type FrameworkJob struct {
	Image     *string
	Directory *string

	Language  string
	Framework string
	Name      string
}

func (j *FrameworkJob) SetupCommand() []string {
	cmds := []string{}
	if j.Directory != nil {
		cmds = append(cmds, getDirectoryCommands(*j.Directory)...)
	}
	cmds = append(cmds, getLanguageAndFrameworkDefaults(j.Language, j.Framework).SetupCommands...)
	return cmds
}

func (j *FrameworkJob) GetImage() string {
	if j.Image != nil {
		return *j.Image
	}
	return getLanguageAndFrameworkDefaults(j.Language, j.Framework).Image
}

func (j *FrameworkJob) GetCommand() string {
	return getLanguageAndFrameworkDefaults(j.Language, j.Framework).Command
}

func (j *FrameworkJob) GetName() string {
	return j.Name
}
