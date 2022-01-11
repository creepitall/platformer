package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	t.Run("", func(t *testing.T) {
		exp := &Config{
			Title:       "little story: the knight",
			VSync:       true,
			EnableFPS:   false,
			FPS:         time.Duration(60),
			WindowsSize: WindowSize{X: 1920, Y: 960},
		}

		cfg := &Config{}
		cfg.ReadFile("./../../../cfg/values.yaml")

		assert.Equal(t, exp, cfg)
	})
}
