package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestStan(t *testing.T) {
	tests := make([]*exec.Cmd, 0)
	test := exec.Command("./stan", "-s", "3", "-c")
	tests = append(tests, test)
	for i := 1; i <= 9; i++ {
		t := "./test" + strconv.Itoa(i) + ".sh"
		test = exec.Command("bash", t)
		tests = append(tests, test)
	}
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err)
		}
		file := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := os.ReadFile(file)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}
