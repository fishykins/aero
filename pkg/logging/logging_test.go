package logging

import "testing"

func TestLogging(t *testing.T) {
	log := New().WithLevel(5)
	log.Trace("Hello, world!")
	log.Info("Hello, world!")
	log.WarnWith("Hello, world!", map[string]interface{}{"name": "fishykins"})
	log.Error("Error was logged")
}
