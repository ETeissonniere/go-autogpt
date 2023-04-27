package commands

import (
	"errors"
	"os"
)

var (
	ErrNoFileOrDirectory = errors.New("please provide a path")
	ErrNoFileOrContent   = errors.New("please provide a file and content")
)

type LsFilesCommand struct{}

func (c *LsFilesCommand) Name() string {
	return "fs-ls"
}

func (c *LsFilesCommand) Usage() string {
	return "list files in the specified directory. Example: fs-ls /tmp or fs-ls ./"
}

func (c *LsFilesCommand) Execute(args []string) (string, error) {
	if len(args) == 0 {
		args = []string{"."}
	}

	dir := args[0]
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", NewAgentError(err)
	}

	output := ""
	for i, entry := range entries {
		if i < len(entries)-2 {
			output += entry.Name() + ","
		} else {
			output += entry.Name()
		}
	}

	return output, nil
}

type ReadFileCommand struct{}

func (c *ReadFileCommand) Name() string {
	return "fs-read"
}

func (c *ReadFileCommand) Usage() string {
	return "read the contents of a file. Example: fs-read /tmp/file.txt"
}

func (c *ReadFileCommand) Execute(args []string) (string, error) {
	if len(args) == 0 {
		return "", NewAgentError(ErrNoFileOrDirectory)
	}

	file := args[0]
	content, err := os.ReadFile(file)
	if err != nil {
		return "", NewAgentError(err)
	}

	if len(content) == 0 {
		content = []byte("file is empty")
	}

	return string(content), nil
}

type WriteFileCommand struct{}

func (c *WriteFileCommand) Name() string {
	return "fs-write"
}

func (c *WriteFileCommand) Usage() string {
	return "write the provided content to a file, may overwrite if the file already exists. Example: fs-write /tmp/file.txt 'hello world'"
}

func (c *WriteFileCommand) Execute(args []string) (string, error) {
	if len(args) < 2 {
		return "", NewAgentError(ErrNoFileOrContent)
	}

	file := args[0]
	content := args[1]
	err := os.WriteFile(file, []byte(content), 0644)
	if err != nil {
		return "", NewAgentError(err)
	}

	return "file written", nil
}

type DeleteFileOrDirectoryCommand struct{}

func (c *DeleteFileOrDirectoryCommand) Name() string {
	return "fs-rm"
}

func (c *DeleteFileOrDirectoryCommand) Usage() string {
	return "delete the specified file or directory. Example: fs-rm /tmp/file.txt or fs-rm /tmp/dir"
}

func (c *DeleteFileOrDirectoryCommand) Execute(args []string) (string, error) {
	if len(args) == 0 {
		return "", NewAgentError(ErrNoFileOrDirectory)
	}

	file := args[0]
	err := os.RemoveAll(file)
	if err != nil {
		return "", NewAgentError(err)
	}

	return "file or directory deleted", nil
}

type CreateDirectoryCommand struct{}

func (c *CreateDirectoryCommand) Name() string {
	return "fs-mkdir"
}

func (c *CreateDirectoryCommand) Usage() string {
	return "create the specified directory. Example: fs-mkdir /tmp/dir"
}

func (c *CreateDirectoryCommand) Execute(args []string) (string, error) {
	if len(args) == 0 {
		return "", NewAgentError(ErrNoFileOrDirectory)
	}

	dir := args[0]
	err := os.Mkdir(dir, 0755)
	if err != nil {
		return "", NewAgentError(err)
	}

	return "directory created", nil
}

func init() {
	DefaultCommands = append(DefaultCommands, &LsFilesCommand{})
	DefaultCommands = append(DefaultCommands, &ReadFileCommand{})
	DefaultCommands = append(DefaultCommands, &WriteFileCommand{})
	DefaultCommands = append(DefaultCommands, &DeleteFileOrDirectoryCommand{})
	DefaultCommands = append(DefaultCommands, &CreateDirectoryCommand{})
}
