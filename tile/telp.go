package tile

import (
	"math"
	"testing"
)

//AssetAreEqual check if two float64 numbers are roughly equal, and fails the test if not.
func AssertAreEqual(t *testing.T, expected float64, actual float64, message string) {
	if math.Abs(expected-actual) > 0.001 {
		t.Log(message)
		t.Log("Should be ", expected, " was ", actual)
		t.Fail()
	}
}
