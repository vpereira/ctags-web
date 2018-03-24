package main
import "testing"

func TestEnv(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping test in short mode.")
    }
}
