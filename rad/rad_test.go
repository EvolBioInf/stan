package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestRad(t *testing.T) {
	var tests []*exec.Cmd
	f := "test.fasta"
	test := exec.Command("./rad", "-s", "3", f)
	tests = append(tests, test)
	test = exec.Command("./rad", "-s", "3", "-l", "100", f)
	tests = append(tests, test)
	test = exec.Command("./rad", "-s", "3", "-d", "100", f)
	tests = append(tests, test)
	test = exec.Command("./rad", "-s", "3", "-n", "1", f)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err)
		}
		fn := "r" + strconv.Itoa(i+1) + ".fasta"
		want, err := os.ReadFile(fn)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
