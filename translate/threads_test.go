//+build threads

package main

import "testing"

// TestThread test
func TestThread(t *testing.T) {
	threadChan := make(chan string)
	go getMainThread("1576725066.000300", "DHS7SGTDX", threadChan)

	threadTs := <-threadChan
	expected := "1576725059.000200"

	if threadTs != expected {
		t.Errorf("Expecting %s, got %s", expected, threadTs)
	}
}

// TestThreadReply test
func TestThreadReply(t *testing.T) {
	threadChan := make(chan string)
	go getMainThread("1576725059.000200", "DHS7SGTDX", threadChan)

	threadTs := <-threadChan
	expected := "1576725059.000200"

	if threadTs != expected {
		t.Errorf("Expecting %s, got %s", expected, threadTs)
	}
}
