package server

import (
	"fmt"
	"os/exec"

	"github.com/Please-Change/backend/pkg/types"
)

type Examiner struct {
	DockerProcess *exec.Cmd
	Abort         func()
}

func NewExaminer() Examiner {
	var e Examiner
	return e
}

func (e Examiner) RunExam(program string, language types.Language, challenge string) string {
	var expectedOutput = challenge

	switch language {
	case types.C:
		{
			compilation := exec.Command("docker", "run", "-it", "runner_c:latest",
				"bash", "-c", fmt.Sprintf("echo %s > program.c; gcc program.c; ./a.out", program))

			compilationPipe, err := compilation.StdoutPipe()

			if err != nil {
				return err.Error()
			}

			var compilationOut []byte
			_, err = compilationPipe.Read(compilationOut)
			if err != nil {
				return err.Error()
			}

			out := string(compilationOut)
			if out == expectedOutput {
				return ""
			} else {
				if out == "" {
					return "Empty file"
				}
				return out
			}
		}
	case types.Cpp:
		{
			compilation := exec.Command("docker", "run", "-it", "runner_cpp:latest",
				"bash", "-c", fmt.Sprintf("echo %s > program.cpp; gcc program.cpp; ./a.out", program))

			compilationPipe, err := compilation.StdoutPipe()

			if err != nil {
				return err.Error()
			}

			var compilationOut []byte
			_, err = compilationPipe.Read(compilationOut)
			if err != nil {
				return err.Error()
			}

			out := string(compilationOut)
			if out == expectedOutput {
				return ""
			} else {
				if out == "" {
					return "Empty file"
				}
				return out
			}
		}
	case types.Rust:
		{
			compilation := exec.Command("docker", "run", "-it", "runner_rust:latest",
				"bash", "-c", fmt.Sprintf("echo %s > program.rs; rustc program.rs; ./program", program))

			compilationPipe, err := compilation.StdoutPipe()

			if err != nil {
				return err.Error()
			}

			var compilationOut []byte
			_, err = compilationPipe.Read(compilationOut)
			if err != nil {
				return err.Error()
			}

			out := string(compilationOut)
			if out == expectedOutput {
				return ""
			} else {
				if out == "" {
					return "Empty file"
				}
				return out
			}
		}
	case types.Go:
		{
			compilation := exec.Command("docker", "run", "-it", "runner_go:latest",
				"bash", "-c", fmt.Sprintf("echo %s > main.go; go run .", program))

			compilationPipe, err := compilation.StdoutPipe()

			if err != nil {
				return err.Error()
			}

			var compilationOut []byte
			_, err = compilationPipe.Read(compilationOut)
			if err != nil {
				return err.Error()
			}

			out := string(compilationOut)
			if out == expectedOutput {
				return ""
			} else {
				if out == "" {
					return "Empty file"
				}
				return out
			}
		}
	case types.JavaScript:
		{
			compilation := exec.Command("docker", "run", "-it", "runner_js:latest",
				"bash", "-c", fmt.Sprintf("echo %s > program.js; node program.js", program))

			compilationPipe, err := compilation.StdoutPipe()

			if err != nil {
				return err.Error()
			}

			var compilationOut []byte
			_, err = compilationPipe.Read(compilationOut)
			if err != nil {
				return err.Error()
			}

			out := string(compilationOut)
			if out == expectedOutput {
				return ""
			} else {
				if out == "" {
					return "Empty file"
				}
				return out
			}
		}
	case types.Python:
		{
			compilation := exec.Command("docker", "run", "-it", "runner_py:latest",
				"bash", "-c", fmt.Sprintf("echo %s > program.py; python program.py", program))

			compilationPipe, err := compilation.StdoutPipe()

			if err != nil {
				return err.Error()
			}

			var compilationOut []byte
			_, err = compilationPipe.Read(compilationOut)
			if err != nil {
				return err.Error()
			}

			out := string(compilationOut)
			if out == expectedOutput {
				return ""
			} else {
				if out == "" {
					return "Empty file"
				}
				return out
			}
		}
	default:
		{
			return "Unknown filetype"
		}
	}

}
