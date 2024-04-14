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

func Lang2Compiler(language types.Language) string {
	switch language {
	case types.Go:
		{
			return "go"
		}
	case types.C:
		{
			return "gcc"
		}
	case types.Cpp:
		{
			return "gcc"
		}
	case types.Rust:
		{
			return "rustc"
		}
	case types.Python:
		{
			return "python3"
		}
	case types.JavaScript:
		{
			return "node"
		}
	default:
		{
			return "exit"
		}
	}
}

func NewExaminer() Examiner {
	var e Examiner
	return e
}

func (e Examiner) RunExam(program string, output chan string, isDone chan bool,
	isSuccess chan bool, language types.Language, challenge string) {
	switch language {
	case types.C:
		{
			exec.Command("docker", "run", "-it", "runner_c:latest",
				"bash", "-c", fmt.Sprintf("\"cd $PWD;echo %s > program.c;"+
					"echo $(gcc program.c) > compiler.txt;\"", program))

			compilation := exec.Command("docker", "run", "-it",
				"runner_c:latest", "bash", "-c", "\"cd $PWD; echo $?;\"")

			compilationPipe, err := compilation.StdoutPipe()

			if err != nil {

			}

			var compilationOut []byte
			_, err = compilationPipe.Read(compilationOut)
			if err != nil {
			}

			if string(compilationOut) == "1" {
				exec.Command("docker", "run", "-it", "runner_c:latest", "bash",
					"-c", "\"cd $PWD; cat compiler.txt\"")
				compilationPipe, err := compilation.StdoutPipe()

				if err != nil {

				}

				var compilationOut []byte
				_, err = compilationPipe.Read(compilationOut)
				if err != nil {
				}

				output <- string(compilationOut)
				isSuccess <- false
				isDone <- true
			}

			runnerCmd := exec.Command("docker", "run", "-it", "runner_c:latest",
				"bash", "-c", "\"cd $PWD; echo $(./a.out) > result.txt;")

			resultCheck :=
				exec.Command("docker", "run", "-it", "runner_c:latest",
					"diff -q "+"result.txt /usr/share/coderunner/solutions/"+
						fmt.Sprintf("solution_%s.txt;", challenge)+" echo "+
						"$?;\"")

			runnerReader, err := runnerCmd.StdoutPipe()
			if err != nil {

			}
			var runnerOut []byte
			_, err = runnerReader.Read(runnerOut)
			if err != nil {

			}
		}
	case types.Cpp:
		{

		}
	case types.Rust:
		{

		}
	case types.Go:
		{

		}
	case types.JavaScript:
		{

		}
	case types.Python:
		{

		}
	default:
		{

		}
	}

}
