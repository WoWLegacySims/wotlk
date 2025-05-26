package hunter

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

func (hunter *Hunter) registerKillCommandCD() {
	if hunter.pet == nil || hunter.Level < 66 {
		return
	}

	var beastLordProcAura *core.Aura
	if hunter.HasSetBonus(ItemSetBeastLord, 4) {
		beastLordProcAura = hunter.NewTemporaryStatsAura("Beast Lord Proc", core.ActionID{SpellID: 37483}, stats.Stats{stats.ArmorPenetration: 85}, time.Second*15)
	}

	actionID := core.ActionID{SpellID: 34026}
	bonusPetSpecialCrit := 10 * float64(hunter.Talents.FocusedFire)

	hunter.pet.KillCommandAura = hunter.pet.RegisterAura(core.Aura{
		Label:     "Kill Command",
		ActionID:  actionID,
		Duration:  time.Second * 30,
		MaxStacks: 3,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			hunter.pet.focusDump.BonusCrit += bonusPetSpecialCrit
			if hunter.pet.specialAbility != nil {
				hunter.pet.specialAbility.BonusCrit += bonusPetSpecialCrit
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			hunter.pet.focusDump.BonusCrit -= bonusPetSpecialCrit
			if hunter.pet.specialAbility != nil {
				hunter.pet.specialAbility.BonusCrit -= bonusPetSpecialCrit
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeSpecial | core.ProcMaskSpellDamage) {
				aura.RemoveStack(sim)
			}
		},
	})

	hunter.KillCommand = hunter.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		Flags:       core.SpellFlagNoOnCastComplete,

		ManaCost: core.ManaCostOptions{
			BaseCost: 0.03,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hunter.NewTimer(),
				Duration: time.Minute - time.Second*10*time.Duration(hunter.Talents.CatlikeReflexes),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hunter.pet.IsEnabled()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			hunter.pet.KillCommandAura.Activate(sim)
			hunter.pet.KillCommandAura.SetStacks(sim, 3)
			if beastLordProcAura != nil {
				beastLordProcAura.Activate(sim)
			}
		},
	})

	hunter.AddMajorCooldown(core.MajorCooldown{
		Spell: hunter.KillCommand,
		Type:  core.CooldownTypeDPS,
	})
}
