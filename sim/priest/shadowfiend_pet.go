package priest

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

type Shadowfiend struct {
	core.Pet

	Priest          *Priest
	Shadowcrawl     *core.Spell
	ShadowcrawlAura *core.Aura
}

var basePercentageStats = stats.Stats{
	// with 3% crit debuff, shadowfiend crits around 9-12% (TODO: verify and narrow down)
	stats.MeleeCrit: 8,
}

func (priest *Priest) NewShadowfiend() *Shadowfiend {
	baseStats := core.PetBaseStats[core.Pet_Shadowfiend][priest.Level].Stats

	shadowfiend := &Shadowfiend{
		Pet:    core.NewPet("Shadowfiend", &priest.Character, baseStats, basePercentageStats, priest.shadowfiendStatInheritance(), false, false),
		Priest: priest,
	}

	manaMetric := priest.NewManaMetrics(core.ActionID{SpellID: 34433})
	_ = core.MakePermanent(shadowfiend.GetOrRegisterAura(core.Aura{
		Label: "Autoattack mana regen",
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			restoreMana := priest.MaxMana() * 0.05
			priest.AddMana(sim, restoreMana, manaMetric)
		},
	}))

	actionID := core.ActionID{SpellID: 63619}
	shadowfiend.ShadowcrawlAura = shadowfiend.GetOrRegisterAura(core.Aura{
		Label:    "Shadowcrawl",
		ActionID: actionID,
		Duration: time.Second * 5,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			shadowfiend.PseudoStats.DamageDealtMultiplier *= 1.15
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			shadowfiend.PseudoStats.DamageDealtMultiplier /= 1.15
		},
	})

	shadowfiend.Shadowcrawl = shadowfiend.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolMagic,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagNoLogs,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: time.Second * 6,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			shadowfiend.ShadowcrawlAura.Activate(sim)
		},
	})

	shadowfiend.PseudoStats.DamageTakenMultiplier *= 0.1

	shadowfiend.EnableAutoAttacks(shadowfiend, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin:        float64(shadowfiend.Level)*2.5 - float64(shadowfiend.Level),
			BaseDamageMax:        float64(shadowfiend.Level)*2.5 + float64(shadowfiend.Level),
			SwingSpeed:           1.5,
			NormalizedSwingSpeed: 1.5,
			CritMultiplier:       2,
			SpellSchool:          core.SpellSchoolShadow,
		},
		AutoSwingMelee: true,
	})

	shadowfiend.AddStatDependency(stats.Strength, stats.AttackPower, 2.0)

	core.ApplyPetConsumeEffects(&shadowfiend.Character, priest.Consumes)

	priest.AddPet(shadowfiend)

	return shadowfiend
}

func (priest *Priest) shadowfiendStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats, ownerPseudoStats stats.PseudoStats) stats.Stats {

		return stats.Stats{ //still need to nail down shadow fiend crit scaling, but removing owner crit scaling after further investigation
			// 3 x sp to ap, lol
			stats.AttackPower: (ownerStats[stats.SpellPower] + ownerPseudoStats.ShadowSpellPower) * 3,
			stats.SpellPower:  (ownerStats[stats.SpellPower] + ownerPseudoStats.ShadowSpellPower) * 0.3,
			stats.MeleeHit:    priest.CalculateHitInheritance(stats.SpellHit, stats.MeleeHit),
			stats.Expertise:   priest.CalculateHitInheritance(stats.SpellHit, stats.Expertise),
			stats.Stamina:     ownerStats[stats.Stamina] * 0.65,
			stats.Intellect:   ownerStats[stats.Intellect] * 0.3,
		}
	}
}

func (shadowfiend *Shadowfiend) Initialize() {
}

func (shadowfiend *Shadowfiend) ExecuteCustomRotation(sim *core.Simulation) {
	shadowfiend.Shadowcrawl.Cast(sim, nil)
}

func (shadowfiend *Shadowfiend) Reset(sim *core.Simulation) {
	shadowfiend.ShadowcrawlAura.Deactivate(sim)
	shadowfiend.Disable(sim)
}

func (shadowfiend *Shadowfiend) OnPetDisable(sim *core.Simulation) {
	shadowfiend.ShadowcrawlAura.Deactivate(sim)
}

func (shadowfiend *Shadowfiend) GetPet() *core.Pet {
	return &shadowfiend.Pet
}
