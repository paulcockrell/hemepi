package main

import (
	"regexp"
	"testing"
)

func TestAssetPath(t *testing.T) {
	got, err := assetPath()
	if err != nil {
		t.Fatal(err)
	}
	want := ".*/assets/$"
	matched, err := regexp.MatchString(want, got)
	if err != nil {
		t.Fatal(err)
	}

	if !matched {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestBuildLines(t *testing.T) {
	data := mockResponse
	lines := BuildLines(data)
	if len(lines) != 6 {
		t.Errorf("Expected %d lines got %d", len(lines), 6)
	}
}
