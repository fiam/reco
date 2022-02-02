package reco

import (
	"runtime"
	"runtime/debug"
	"strings"
)

type panicInfo struct {
	file      string
	line      int
	pc        uintptr
	pcSkip    int
	stackSkip int
}

func prettyStack(lines []string, toHTML bool) []string {
	return lines
}

func formatStack(skip int, toHTML bool) string {
	// Always skip the frames for formatStack
	skip += 1
	buf := debug.Stack()
	// Remove 2 * skip lines after first line, since they correspond
	// to the skipped frames
	lines := strings.Split(string(buf), "\n")
	end := 2*skip + 1
	if end > len(lines) {
		end = len(lines)
	}
	lines = append(lines[:1], lines[end:]...)
	lines = prettyStack(lines, toHTML)
	return strings.Join(lines, "\n")

}

func uppermostPanic() *panicInfo {
	skip := 0
	callers := make([]uintptr, 32)
	for {
		calls := callers[:runtime.Callers(skip, callers)]
		count := len(calls)
		if count == 0 {
			break
		}
		for ii := count - 1; ii >= 0; ii-- {
			if f := runtime.FuncForPC(calls[ii]); f != nil {
				name := f.Name()
				if strings.HasPrefix(name, "runtime.") && strings.Contains(name, "panic") {
					pcSkip := skip + ii - 1
					stackSkip := pcSkip
					switch name {
					case "runtime.panic":
					case "runtime.sigpanic":
						stackSkip -= 2
					default:
						stackSkip--
					}
					// Find the source location of the file that called panic, not
					// the call to panic, hence skip an extra frame
					_, file, line, _ := runtime.Caller(pcSkip + 1)
					return &panicInfo{
						file:      file,
						line:      line,
						pc:        calls[ii],
						pcSkip:    pcSkip,
						stackSkip: stackSkip,
					}
				}
			}
		}
		skip += count
	}
	return nil
}
