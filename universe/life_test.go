package universe

import (
	"sync"
	"testing"
)

func BenchmarkApplyRulesSequential(b *testing.B) {
	SpawnUniverse()
	b.ResetTimer()
	for b.Loop() {
		ApplyRules()
	}
}

func BenchmarkApplyRulesConcurrent(b *testing.B) {
	SpawnUniverse()
	wg := sync.WaitGroup{}
	b.ResetTimer()
	for b.Loop() {
		ApplyRulesInParallel(&wg)
		wg.Wait()
	}
}
