package exec

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/xlog"
)

// NewCmd creates a new command
func NewCmd(args ...string) (cmd *Cmd) {
	cmd = &Cmd{
		Args:              args,
		ChannelCancel:     make(chan bool),
		channelSignal:     make(chan os.Signal),
		channelSignalQuit: make(chan bool),
		Logger:            xlog.NopLogger,
	}
	return
}

// Cmd represents a command
type Cmd struct {
	Args              []string
	ChannelCancel     chan bool
	channelSignal     chan os.Signal
	channelSignalQuit chan bool
	handlingSignals   bool
	Logger            xlog.Logger
	Timeout           time.Duration
}

// String allows Cmd to implements the stringify interface
func (c *Cmd) String() string {
	return strings.Join(c.Args, " ")
}

// HandleSignals handles signals
func (c *Cmd) HandleSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func(c *Cmd) {
		for {
			select {
			case s := <-ch:
				c.channelSignal <- s
			case <-c.channelSignalQuit:
				return
			}
		}
	}(c)
	c.handlingSignals = true
}

// Close closes a command
func (c *Cmd) Close() {
	if c.handlingSignals {
		c.channelSignalQuit <- true
	}
}

// killCmdProcess kills a command process and logs the proper error
func killCmdProcess(cmd *exec.Cmd, m string) (err error) {
	if cmd.Process != nil {
		err = fmt.Errorf("%s, killing process %v", m, cmd.Process.Pid)
		cmd.Process.Kill()
	} else {
		err = fmt.Errorf("%s, no process to kill", m)
	}
	return
}

// Exec executes a command
var Exec = func(cmd *Cmd) (o []byte, d time.Duration, err error) {
	// Init
	defer func(t time.Time) {
		d = time.Since(t)
	}(time.Now())

	// Create exec command
	execCmd := exec.Command(cmd.Args[0], cmd.Args[1:]...)

	// Create channel that will be closed when execution is done
	done := make(chan bool)

	// Execute command in go routine
	go func(cmd *Cmd) {
		defer close(done)
		cmd.Logger.Debugf("Executing %s", cmd)
		o, err = execCmd.CombinedOutput()
	}(cmd)

	// Listen to either done, timeout, cancel or signal channels
	if cmd.Timeout > 0 {
		for {
			select {
			case <-done:
				return
			case <-time.After(cmd.Timeout):
				err = killCmdProcess(execCmd, fmt.Sprintf("Timeout of %s reached", cmd.Timeout))
				return
			case s := <-cmd.channelSignal:
				err = killCmdProcess(execCmd, fmt.Sprintf("Caught signal %s", s))
				return
			case <-cmd.ChannelCancel:
				err = killCmdProcess(execCmd, "Command was cancelled")
				return
			}
		}
	} else {

		for {
			select {
			case <-done:
				return
			case s := <-cmd.channelSignal:
				err = killCmdProcess(execCmd, fmt.Sprintf("Caught signal %s", s))
				return
			case <-cmd.ChannelCancel:
				err = killCmdProcess(execCmd, "Command was cancelled")
				return
			}
		}
	}
}
