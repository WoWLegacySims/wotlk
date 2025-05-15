package core

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/coreinfo"
)

func (character *Character) GetMaxProfessionRank() int32 {
	var rank int32
	switch {
	case character.Level >= 55:
		rank = 5
	case character.Level >= 40:
		rank = 4
	case character.Level >= 25:
		rank = 3
	case character.Level >= 10:
		rank = 2
	default:
		rank = 1
	}

	switch {
	case character.Expansion == proto.Expansion_ExpansionTbc:
		rank = min(rank, 4)
	case character.Expansion == proto.Expansion_ExpansionVanilla:
		rank = min(rank, 3)
	}

	return rank
}

// This is just the static bonuses. Most professions are handled elsewhere.
func (character *Character) applyProfessionEffects() {
	if character.HasProfession(proto.Profession_Mining) {
		bonus := []float64{3, 5, 7, 10, 30, 60}[character.GetMaxProfessionRank()]
		character.AddStat(stats.Stamina, bonus)
	}

	if character.HasProfession(proto.Profession_Skinning) {
		bonus := []float64{3, 6, 9, 12, 20, 40}[character.GetMaxProfessionRank()]
		character.AddStats(stats.Stats{stats.MeleeCrit: bonus, stats.SpellCrit: bonus})
	}

	if character.HasProfession(proto.Profession_Herbalism) {
		dbc := coreinfo.LifebloodInfos.SpellInfos[character.GetMaxProfessionRank()]
		heal := dbc.Effects[0].BasePoints + 1
		actionID := ActionID{SpellID: dbc.SpellID}
		healthMetrics := character.NewHealthMetrics(actionID)

		spell := character.RegisterSpell(SpellConfig{
			ActionID:    actionID,
			SpellSchool: SpellSchoolNature,
			Cast: CastConfig{
				CD: Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 3,
				},
			},
			ApplyEffects: func(sim *Simulation, _ *Unit, _ *Spell) {
				amount := heal + character.MaxHealth()*0.0032
				StartPeriodicAction(sim, PeriodicActionOptions{
					Period:   time.Second,
					NumTicks: 5,
					OnAction: func(sim *Simulation) {
						character.GainHealth(sim, amount*character.PseudoStats.HealingTakenMultiplier, healthMetrics)
					},
				})
			},
		})
		character.AddMajorCooldown(MajorCooldown{
			Type:  CooldownTypeSurvival,
			Spell: spell,
		})
	}
}
