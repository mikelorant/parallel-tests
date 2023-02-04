package shell

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os/exec"
)

func main() {}

func Stdlib(w io.Writer, command string, args []string) error {
	var buf bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		var execError *exec.Error

		if errors.As(err, &execError) {
			return fmt.Errorf("unable to exec command: %v: %w", execError.Name, execError.Err)
		}

		return fmt.Errorf("unable to exec command: %w", err)
	}

	if _, err = io.Copy(w, &buf); err != nil {
		var pathError *fs.PathError

		if !errors.As(err, &pathError) {
			return fmt.Errorf("unable to copy output: %w", err)
		}

		//nolint:errorlint
		pathError = err.(*fs.PathError)
		if pathError.Path != "/dev/ptmx" {
			return fmt.Errorf("unable to copy output: %v: %w", pathError.Path, pathError.Err)
		}
	}

	return nil
}
