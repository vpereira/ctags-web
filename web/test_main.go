package main
import "testing"

func TestOpenDB(t *testing.T) {
    env := Env{}
    _ , err := env.OpenDB("db", "test", "test")
    if err != nil {
        t.Error("Cannot open connection")
    }
}

func TestSetDB(t *testing.T) {
  env := Env{}
  env.OpenDB("db", "test", "test")
  col := env.SetDB("test", "bar")

  if col.Name != "bar" {
    t.Error("Couldnt set collection")
  }
}
