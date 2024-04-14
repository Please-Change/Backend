package server

import (
	"os/exec"
)

type Examiner struct {
	DockerProcess *exec.Cmd
	Abort         func()
}

func NewExaminer() Examiner {
	// Docker container started
	var e Examiner
	e.DockerProcess = exec.Command("ls", ".")
	e.Abort = func() { e.DockerProcess.Cancel() }
	return e
}

func (e Examiner) RunExam(program string, output chan string, isDone chan bool,
	isSuccess chan bool) {

}
