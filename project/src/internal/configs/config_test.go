package configs

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	description := "If the load config scenarios are working correctely"
	defer func() {
		log.Printf("Test: %s\n", description)
		log.Println("Deferred tearing down.")
	}()

	t.Run("test success load", func(t *testing.T) {
		t.Setenv("ENV", "test")
		t.Setenv("HTTP_PORT", "8080")
		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, "test", cfg.Env)
		assert.Equal(t, "8080", cfg.Port)
	})

	t.Run("test fail missing env but got default", func(t *testing.T) {
		t.Setenv("ENV", "")
		t.Setenv("HTTP_PORT", "8080")
		cfg, err := LoadConfig()
		assert.Nil(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, cfg.Env, "development")
	})

	t.Run("test fail missing http port", func(t *testing.T) {
		t.Setenv("ENV", "test")
		t.Setenv("HTTP_PORT", "")
		cfg, err := LoadConfig()
		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.EqualError(t, err, "HTTP_PORT is required")
	})
}
