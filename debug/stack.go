package debug

import (
	"bytes"
	"regexp"
	"runtime/debug"
	"strconv"
)

// Vars
var (
	byteLineDelimiter = []byte("\n")
	regexpFile        = regexp.MustCompile("(.+)\\:(\\d+) \\+(.+)$")
	regexpFunction    = regexp.MustCompile("(.+)\\(.*\\)$")
)

// DebugStack allows testing functions using it
var DebugStack = func() []byte {
	return debug.Stack()
}

// Stack represents a stack
type Stack []StackItem

// StackItem represents a stack item
type StackItem struct {
	Filename string
	Function string
	Line     int
}

// NewStack returns a new stack
func NewStack() (o Stack) {
	var i = &StackItem{}
	for _, line := range bytes.Split(DebugStack(), byteLineDelimiter) {
		// Trim line
		line = bytes.TrimSpace(line)

		// Check line type
		var r [][]string
		if r = regexpFunction.FindAllStringSubmatch(string(line), -1); len(r) > 0 && len(r[0]) > 1 {
			i.Function = r[0][1]
		} else if r = regexpFile.FindAllStringSubmatch(string(line), -1); len(r) > 0 && len(r[0]) > 2 {
			i.Filename = r[0][1]
			i.Line, _ = strconv.Atoi(r[0][2])
			o = append(o, *i)
			i = &StackItem{}
		}
	}
	return
}
