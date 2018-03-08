package httpServer

import "testing"

func TestSetWebRoot(t *testing.T) {
	expectedResult := "test"
	SetWebRoot(expectedResult)
	actualResult := webRoot

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func TestGetWebRoot(t *testing.T) {
	expectedResult := "test2"
	webRoot = expectedResult
	actualResult := GetWebRoot()

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}
