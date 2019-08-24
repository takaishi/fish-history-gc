package main

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func Test_readEntries(t *testing.T) {
	input := `- cmd: foo
  when: 1565245208
- cmd: bar
  when: 1565245270
- cmd: echo "hello: hoge"
  when: 1566649936
`

	expected := Entries{
		{
			Cmd:  "foo",
			When: 1565245208,
		},
		{
			Cmd:  "bar",
			When: 1565245270,
		},
		{
			Cmd:  "echo \"hello: hoge\"",
			When: 1566649936,
		},
	}

	buf := bytes.NewBufferString(input)

	actual, err := readEntries(buf)
	if err != nil {
		t.Fatalf("failed to readEntries: %s", err)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("mismatch: (-want +got):\n %s", diff)
	}
}

func Test_removeDupEntries(t *testing.T) {
	entries := Entries{
		{
			Cmd:  "foo",
			When: 1566219481,
		},
		{
			Cmd:  "bar",
			When: 1566219608,
		},
		{
			Cmd:  "foo",
			When: 1566219614,
		},
	}

	expected := Entries{
		{
			Cmd:  "foo",
			When: 1566219614,
		},
		{
			Cmd:  "bar",
			When: 1566219608,
		},
	}

	actual := removeDupEntries(entries)

	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %+v, actual: %+v", expected, actual)
	}
}
