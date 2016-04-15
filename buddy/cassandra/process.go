package cassandra

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Process interface {
	Start() error
	Stop() error
	Running() bool
}

func New(cfg *Config) Process {
	return &cassandraProcess{
		cfg: cfg,
	}
}

type cassandraProcess struct {
	cfg    *Config
	cmd    *exec.Cmd
	stdout *bytes.Buffer
}

func (c *cassandraProcess) Start() error {
	cmd := exec.Command(c.cfg.Executable, "-f")
	cmd.Env = c.buildEnv()

	c.cmd = cmd
	c.stdout = new(bytes.Buffer)
	c.cmd.Stdout = c.stdout
	err := c.cmd.Start()
	log.Println(c.cmd, err)
	log.Println(c.cmd.Process)
	log.Println(c.cmd.ProcessState)
	return err
}

func (c *cassandraProcess) Stop() error {
	if !c.Running() {
		return nil
	}
	log.Println(c.cmd.Process)
	log.Println(c.cmd.ProcessState)
	log.Println("Sending SIGTERM to", c.cmd.Process.Pid)

	c.cmd.Process.Signal(syscall.SIGTERM)

	if err := c.cmd.Wait(); err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			if status, ok := execErr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 143 {
					return nil
				}
			}
		}
		return err
	}
	return nil
}

func (c *cassandraProcess) Running() bool {
	process, err := os.FindProcess(c.cmd.Process.Pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}
	return true
}

func (c *cassandraProcess) buildEnv() []string {
	e := os.Environ()
	cfg := c.cfg
	e = append(e, cfg.Env()...)
	return e
}
