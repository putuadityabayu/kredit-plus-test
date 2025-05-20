/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package response

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// A StackFrame contains all necessary information about to generate a line
// in a callstack.
type stackFrame struct {
	// The path to the file containing this ProgramCounter
	File string
	// The LineNumber in that file
	LineNumber int
	// The Name of the function that contains this ProgramCounter
	Name string
	// The Package that contains this function
	Package string
	// The underlying ProgramCounter
	ProgramCounter uintptr
}

// NewStackFrame popoulates a stack frame object from the program counter.
func newStackFrame(pc uintptr) (frame stackFrame) {

	frame = stackFrame{ProgramCounter: pc}
	if frame.function() == nil {
		return
	}
	frame.Package, frame.Name = packageAndName(frame.function())

	// pc -1 because the program counters we use are usually return addresses,
	// and we want to show the line that corresponds to the function call
	frame.File, frame.LineNumber = frame.function().FileLine(pc - 1)
	return

}

// Func returns the function that contained this frame.
func (frame *stackFrame) function() *runtime.Func {
	if frame.ProgramCounter == 0 {
		return nil
	}
	return runtime.FuncForPC(frame.ProgramCounter)
}

// String returns the stackframe formatted in the same way as go does
// in runtime/debug.Stack()
func (frame *stackFrame) String() string {
	str := fmt.Sprintf("%s:%d (0x%x)\n", frame.File, frame.LineNumber, frame.ProgramCounter)

	source, err := frame.sourceLine()
	if err != nil {
		return str
	}

	return str + fmt.Sprintf("\t%s: %s\n", frame.Name, source)
}

// sourceLine gets the line of code (from File and Line) of the original source if possible.
func (frame *stackFrame) sourceLine() (string, error) {
	if frame.LineNumber <= 0 {
		return "???", nil
	}

	file, err := os.Open(frame.File)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 1
	for scanner.Scan() {
		if currentLine == frame.LineNumber {
			return string(bytes.Trim(scanner.Bytes(), " \t")), nil
		}
		currentLine++
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}

	return "???", nil
}

func packageAndName(fn *runtime.Func) (string, string) {
	name := fn.Name()
	pkg := ""

	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		pkg += name[:lastslash] + "/"
		name = name[lastslash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.Replace(name, "·", ".", -1)
	return pkg, name
}

// Stack returns the callstack formatted the same way that go does
// in runtime/debug.Stack()
func (e *ErrorResponse) Stack() []byte {
	buf := bytes.Buffer{}

	for _, frame := range e.StackFrames() {
		buf.WriteString(frame.String())
	}

	return buf.Bytes()
}

// Callers satisfies the bugsnag ErrorWithCallerS() interface
// so that the stack can be read out.
func (e *ErrorResponse) Callers() []uintptr {
	return e.stack
}

// ErrorStack returns a string that contains both the
// error message and the callstack.
func (e *ErrorResponse) ErrorStack() string {
	return string(e.Stack())
}

// StackFrames returns an array of frames containing information about the
// stack.
func (e *ErrorResponse) StackFrames() []stackFrame {
	if e.frames == nil {
		e.frames = make([]stackFrame, len(e.stack))

		for i, pc := range e.stack {
			sf := newStackFrame(pc)
			e.frames[i] = sf
		}
	}

	return e.frames
}
