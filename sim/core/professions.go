package core

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func (character *Character) GetMaxProfessionRank() int32 {
	switch {
	case character.Level >= 55:
		return 5
	case character.Level >= 40:
		return 4
	case character.Level >= 25:
		return 3
	case character.Level >= 10:
		return 2
	default:
		return 1
	}
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
		dbc := LifebloodInfos[character.GetMaxProfessionRank()]
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
