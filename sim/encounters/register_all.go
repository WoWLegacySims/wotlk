package encounters

import (
	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/encounters/icc"
	"github.com/WoWLegacySims/wotlk/sim/encounters/naxxramas"
	"github.com/WoWLegacySims/wotlk/sim/encounters/toc"
	"github.com/WoWLegacySims/wotlk/sim/encounters/ulduar"
)

func init() {
	naxxramas.Register()
	ulduar.Register()
	toc.Register()
	icc.Register()
}

func AddSingleTargetBossEncounter(presetTarget *core.PresetTarget) {
	core.AddPresetTarget(presetTarget)
	core.AddPresetEncounter(presetTarget.Config.Name, []string{
		presetTarget.Path(),
	})
}
