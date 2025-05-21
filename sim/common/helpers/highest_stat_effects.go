package helpers

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

type HighestStatAura struct {
	statOptions []stats.Stat
	auras       []*core.Aura
	factory     func(stat stats.Stat) *core.Aura
}

func (hsa HighestStatAura) Init(character *core.Character) {
	for i, stat := range hsa.statOptions {
		hsa.auras[i] = hsa.factory(stat)
	}
}

func (hsa HighestStatAura) Get(character *core.Character) *core.Aura {
	bestValue := 0.0
	bestIdx := 0

	for i, stat := range hsa.statOptions {
		value := character.GetStat(stat)
		if value > bestValue {
			bestValue = value
			bestIdx = i
		}
	}

	a := hsa.auras[bestIdx]
	if a == nil {
		a = hsa.factory(hsa.statOptions[bestIdx])
		hsa.auras[bestIdx] = a
	}
	return a
}

func NewHighestStatAura(statOptions []stats.Stat, auraFactory func(stat stats.Stat) *core.Aura) HighestStatAura {
	return HighestStatAura{
		statOptions: statOptions,
		factory:     auraFactory,
		auras:       make([]*core.Aura, len(statOptions)),
	}
}
