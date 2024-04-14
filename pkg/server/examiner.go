package server

import (
	"fmt"
	"log"
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
			outBin, err := exec.Command("docker", "run", "-it", "runner_c:latest",
				"bash", "-c", fmt.Sprintf("echo %q > program.c; gcc program.c; ./a.out", program)).Output()

			if err != nil {
				log.Printf("Err: %s", err)
				return err.Error()
			}

			out := string(outBin)
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
			outBin, err := exec.Command("docker", "run", "-it", "runner_cpp:latest",
				"bash", "-c", fmt.Sprintf("echo %q > program.cpp; gcc program.cpp; ./a.out", program)).Output()

			if err != nil {
				log.Printf("Err: %s", err)
				return err.Error()
			}

			out := string(outBin)
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
			outBin, err :=
				exec.Command("docker", "run", "-i", "runner_rust:latest",
					"bash", "-c", fmt.Sprintf("echo %q > program.rs; rustc program.rs; ./program", program)).Output()

			if err != nil {
				log.Printf("Err: %s", err)
				return err.Error()
			}

			out := string(outBin)
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
			outBin, err :=
				exec.Command("docker", "run", "-i", "runner_go:latest",
					"bash", "-c", fmt.Sprintf("echo %q > main.go; go run .", program)).Output()

			if err != nil {
				log.Printf("Err: %s", err)
				return err.Error()
			}

			out := string(outBin)
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
			outBin, err :=
				exec.Command("docker", "run", "-i", "runner_js:latest",
					"bash", "-c", fmt.Sprintf("echo %q > program.js; node program.js", program)).Output()

			if err != nil {
				log.Printf("Err: %s", err)
				return err.Error()
			}

			out := string(outBin)
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
			outBin, err :=
				exec.Command("docker", "run", "-i", "runner_py:latest",
					"bash", "-c", fmt.Sprintf("echo %q > program.py; python program.py", program)).Output()

			if err != nil {
				log.Printf("Err: %s", err)
				return err.Error()
			}

			out := string(outBin)
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
