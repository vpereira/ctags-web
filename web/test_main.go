package main

import (
	"context"
	"testing"
)

func TestEnv(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
}

func TestEnvOpenDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ctx := context.Background()
	env := &Env{}
	// When no MongoDB is available, expect an error
	err := env.OpenDB(ctx, "mongodb://localhost:27017", "test", "test")
	if err == nil {
		t.Log("OpenDB succeeded (MongoDB may be running)")
		defer env.Client.Disconnect(ctx)
	}
}
