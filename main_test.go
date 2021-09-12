package main

import (
	"testing"

	"github.com/softashell/lewdbot-discord/brain"
	"github.com/softashell/lewdbot-discord/config"
)

func Benchmark_fillBrain(b *testing.B) {
	for n := 0; n < b.N; n++ {
		config.Init()
		brain.Init()

		fillBrain()
	}
}