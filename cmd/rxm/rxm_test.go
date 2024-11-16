package main

import (
	"flag"
	"io"
	"os"
	"testing"
)

func captureOutput(f func()) string {
	originalStdout := os.Stdout
	read, write, _ := os.Pipe()
	os.Stdout = write
	f()
	os.Stdout = originalStdout
	write.Close()
	out, _ := io.ReadAll(read)
	return string(out)
}

func runMain(args []string) string {
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	os.Args = append([]string{"cmd"}, args...)
	return captureOutput(func() { main() })
}

func TestMainFunction(t *testing.T) {
	os.Remove("./test.db")
	runMain([]string{"-db", "./test.db", "add", "foo", "1", "2"})
	runMain([]string{"-db", "./test.db", "ls"})
	runMain([]string{"-db", "./test.db", "ls", "foo"})
	runMain([]string{"-db", "./test.db", "up", "foo", "quantity", "2"})
	runMain([]string{"-db", "./test.db", "up", "foo", "rate", "1"})
	runMain([]string{"-db", "./test.db", "rm", "foo"})
	os.Remove("./test.db")
}
