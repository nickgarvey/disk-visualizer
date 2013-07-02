package main

import "bufio"
import "fmt"
import "os"
import "os/exec"
import "strconv"
import "strings"

type blkTrace struct {
	Time   float32
	Action string
	IoType string
	Sector uint64
	Blocks uint64
}

var masks = []string{"read", "write", "complete", "queue"}

var blkparseFormat = []struct {
	name   string
	format string
}{
	{"time", "%T.%4t"},
	{"action", "%a"},
	{"iotype", "%d"},
	{"sector", "%S"},
	{"blocks", "%n"},
}

func traceBlocks(traceCh chan blkTrace, errCh chan error) {
	traceArgs := []string{"/dev/sda", "-o", "-"}
	for _, mask := range masks {
		traceArgs = append(traceArgs, "-a")
		traceArgs = append(traceArgs, mask)
	}

	parseArgs := []string{"-i", "-", "-q", "-f"}
	formatArg := ""
	for i, s := range blkparseFormat {
		formatArg += s.format
		if i != len(blkparseFormat)-1 {
			formatArg += " "
		}
	}
	formatArg += "\n"

	parseArgs = append(parseArgs, formatArg)

	blktraceCmd := exec.Command("blktrace", traceArgs...)
	blkparseCmd := exec.Command("blkparse", parseArgs...)

	var err error

	blkparseCmd.Stdin, err = blktraceCmd.StdoutPipe()
	if err != nil {
		errCh <- err
		panic(err)
	}

	parsePipe, err := blkparseCmd.StdoutPipe()
	if err != nil {
		errCh <- err
		panic(err)
	}

	blktraceCmd.Stderr = os.Stderr
	blkparseCmd.Stderr = os.Stderr

	for _, cmd := range [](*exec.Cmd){blktraceCmd, blkparseCmd} {
		err = cmd.Start()
		if err != nil {
			errCh <- err
		}
	}

	blkparseReader := bufio.NewReader(parsePipe)

	var line string
	for {
		line, err = blkparseReader.ReadString('\n')
		if err != nil {
			break
		}

		go func(l string) {
			trace, traceErr := buildTrace(l)
			if traceErr != nil {
				errCh <- traceErr
				return
			}
			traceCh <- trace
		}(line)
	}

	// If we got here, then there was an error above
	errCh <- err

	panic(err)
}

func buildTrace(line string) (trace blkTrace, err error) {
	fields := strings.Fields(line)
	if len(fields) != 5 {
		err = fmt.Errorf("Unable to parse line: %s", line)
		return
	}

	time, err := strconv.ParseFloat(fields[0], 32)
	if err != nil {
		return
	}
	trace.Time = float32(time)

	trace.Action = fields[1]

	trace.IoType = fields[2]

	trace.Sector, err = strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return
	}

	trace.Blocks, err = strconv.ParseUint(fields[4], 10, 64)
	if err != nil {
		return
	}

	return trace, nil
}
