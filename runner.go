package simple

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type SimpleRunner struct {
	cmd    *exec.Cmd
	Scene  chan []byte
	quit   chan bool
	Events chan Event
	Errors chan error
	Exited chan bool
}

func NewSimpleRunner(errorChannel bool) *SimpleRunner {
	var errs chan error = nil
	if errorChannel {
		errs = make(chan error)
	}

	return &SimpleRunner{
		cmd:    nil,
		Scene:  make(chan []byte),
		quit:   make(chan bool),
		Events: make(chan Event),
		Errors: errs,
		Exited: make(chan bool),
	}
}

func (sr *SimpleRunner) Start() {
	go func() {
		for {
			select {
			case scene := <-sr.Scene:
				sr.kill()
				sr.run(scene)
			case <-sr.quit:
				sr.kill()
			}
		}
	}()
}

func (sr *SimpleRunner) Stop() {
	sr.quit <- true
}

func (sr *SimpleRunner) kill() {
	if sr.cmd != nil && sr.cmd.Process != nil {
		sr.cmd.Process.Kill()
	}
}

func (sr *SimpleRunner) error(msg string, err error) {
	wrapped := fmt.Errorf("%s: %w", msg, err)
	if sr.Errors == nil {
		log.Println(wrapped.Error())
	} else {
		sr.Errors <- wrapped
	}
}

func (sr *SimpleRunner) run(scene []byte) {
	sr.kill()
	sr.cmd = exec.Command("simple")
	sr.cmd.Stderr = os.Stderr
	log.Println("simple started")

	stdin, err := sr.cmd.StdinPipe()
	if err != nil {
		sr.error("couldn't connect to simple's stdin", err)
		return
	}

	stdout, err := sr.cmd.StdoutPipe()
	if err != nil {
		sr.error("couldn't connect to simple's stdout", err)
		return
	}

	err = sr.cmd.Start()
	if err != nil {
		sr.error("couldn't start simple", err)
	}

	stdin.Write([]byte(scene))
	stdin.Close()
	log.Println("stdin written and closed")

	go sr.outputParser(stdout)

	go func() {
		sr.Events <- RenderedEvent{}
		err := sr.cmd.Wait()
		log.Println("simple exited:", err)
		sr.Exited <- true
	}()
}

func (sr *SimpleRunner) outputParser(stdout io.Reader) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("outputParser: ", line)
		split := strings.SplitN(line, ":", 2)
		if len(split) < 2 {
			continue
		}

		switch strings.TrimSpace(split[0]) {
		case "selected":
			sr.Events <- SelectedEvent{strings.TrimSpace(split[1])}
		case "input":
			split := strings.SplitN(split[1], ":", 2)
			if len(split) != 2 {
				continue
			}
			// should I really trim space from the value here?
			sr.Events <- InputEvent{
				strings.TrimSpace(split[0]),
				strings.TrimSpace(split[1]),
			}
		default:
			sr.Events <- UnknownEvent{
				Raw: line,
			}
		}
	}
	log.Println("ending output parser")
}
