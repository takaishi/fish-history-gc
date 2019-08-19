package main

import (
	"reflect"
	"testing"
)

func Test_removeDupEntries(t *testing.T) {
	entries := Entries{
		{
			Cmd: "foo",
			When: 1566219481,
		},
		{
			Cmd: "bar",
			When: 1566219608,
		},
		{
			Cmd: "foo",
			When: 1566219614,
		},
	}

	expected := Entries{
		{
			Cmd: "foo",
			When: 1566219614,
		},
		{
			Cmd: "bar",
			When: 1566219608,
		},
	}

	actual := removeDupEntries(entries)

	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %+v, actual: %+v", expected, actual)
	}
}