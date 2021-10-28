package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/iovisor/gobpf/bcc"
	"os"
	"os/signal"
)

var binaryProg string
var functionName string

type event struct {
	Arg1 string
	Arg2 string
	Arg3 string
	Arg4 string
	Arg5 string
	Ret  string
}

func init() {
	flag.StringVar(&binaryProg, "binary", "", "The binary to probe")
	flag.StringVar(&functionName, "fn", "", "The function to probe")
}

func main() {
	flag.Parse()
	if len(binaryProg) == 0 {
		panic("Argument --binary needs to be specified")
	}
	if len(functionName) == 0 {
		panic("Argument --fn needs to be specified")
	}

	bccMod := bcc.NewModule(bpfProgram, []string{})
	uprobeFDStart, err := bccMod.LoadUprobe("start")
	if err != nil {
		panic(err)
	}
	uprobeFDEnd, err := bccMod.LoadUprobe("end")
	if err != nil {
		panic(err)
	}

	// Attach the uprobe to be called everytime main.add is called.
	// We need to specify the path to the binary so it can be patched.
	err = bccMod.AttachUprobe(binaryProg, functionName, uprobeFDStart, -1)
	if err != nil {
		panic(err)
	}
	err = bccMod.AttachUretprobe(binaryProg, functionName, uprobeFDEnd, -1)
	if err != nil {
		panic(err)
	}

	// Create the output table named "trace" that the BPF program writes to.
	table := bcc.NewTable(bccMod.TableId("trace"), bccMod)
	ch := make(chan []byte)

	pm, err := bcc.InitPerfMap(table, ch, nil)
	if err != nil {
		panic(err)
	}

	// Watch Ctrl-C so we can quit this program.
	intCh := make(chan os.Signal, 1)
	signal.Notify(intCh, os.Interrupt)

	pm.Start()
	defer pm.Stop()

	for {
		select {
		case <-intCh:
			fmt.Println("Terminating")
			os.Exit(0)
		case data := <-ch:
			var ev event
			err := binary.Read(bytes.NewBuffer(data), bcc.GetHostByteOrder(), &ev)
			if err != nil {
				fmt.Printf("failed to decode received data: %s\n", err)
				continue
			}
			fmt.Printf("arg1: %v arg2: %v arg3: %v arg4: %v arg5: %v ret: %v\n", ev.Arg1, ev.Arg2, ev.Arg3, ev.Arg4, ev.Arg5, ev.Ret)
		}
	}
}
