package main

//Studying interfaces and methods

import (
	"fmt"
)

type Executor interface {
	Exec() bool
}

type Runner struct {
	id  int
	cmd string
}

type SSHRunner struct {
	Runner
	AllowedHosts []string
	RSAKeyPath   string
	Audit        bool
}

type HTTPRunner struct {
	Runner
	MethodsAllowed []string
	URL            string
}

// todo
func (r *Runner) Exec() bool {
	fmt.Println("I fulfill the requirement")
	fmt.Println(r.cmd)
	return true
}

func AuditBeforeRun(e Executor) bool {
	if ssh, ok := e.(*SSHRunner); ok {
		fmt.Println("Auditing SSH CMD:")
		fmt.Println(ssh.cmd)
		ssh.Audit = true
		return ssh.Exec()
	}
	return false
}

func RunExecutor(e Executor) bool {
	if ssh, ok := e.(*SSHRunner); ok {
		fmt.Printf("Command ID: %d must be audited first!\n", ssh.id)
		return false
	}

	if http, ok := e.(*HTTPRunner); ok {
		return http.Exec()
	}

	return false
}

func multiply(a1 int, a2 int) int {
	return a1 * a2
}

func main() {
	//b := [5]int{1, 2, 3, 4, 5}
	//fmt.Printf("%d", multiply(b[(len(b)-1)], (1000/2)))
	fmt.Println("::TESTING::")
	sshRunner := new(SSHRunner)
	sshRunner.id = 1
	sshRunner.cmd = "echo"

	RunExecutor(sshRunner)
	AuditBeforeRun(sshRunner)

	// TODO
}
