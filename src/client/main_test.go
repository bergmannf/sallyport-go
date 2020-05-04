package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestReadEvent(t *testing.T) {
	type args struct {
		buffer *bytes.Buffer
		reader *bufio.Reader
	}
	tests := []struct {
		name string
		args args
	}{
		{"Read one line events", args{bytes.NewBuffer([]byte(`data: one line
`)), bufio.NewReader(strings.NewReader(`data: one line

`))}},
		{"Read two line events", args{bytes.NewBuffer([]byte(`data: one line
data: two line
`)), bufio.NewReader(strings.NewReader(`data: one line
data: two line

`))}},
		{"Won't read without two newlines", args{bytes.NewBuffer([]byte("")), bufio.NewReader(strings.NewReader("data: one line\n"))}},
		{"Will read multiple lines", args{bytes.NewBuffer([]byte(`data: one line
data: two line
data: three line
`)), bufio.NewReader(strings.NewReader(`data: one line
data: two line
data: three line

`))}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytes.Buffer{}
			ReadEvent(&result, tt.args.reader)
			actualText := string(result.Bytes())
			expectedText := string(tt.args.buffer.Bytes())
			expected := expectedText == actualText
			if expected != true {
				t.Errorf("Expected: %s Got: %s", expectedText, actualText)
			}
		})
	}
}
