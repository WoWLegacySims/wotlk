package hunter

import (
	"math"
	"time"

	"github.com/WoWLegacySims/wotlk/sim/core"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
	"github.com/WoWLegacySims/wotlk/sim/spellinfo/hunterinfo"
)

type PetAbilityType int

// Pet AI doesn't use abilities immediately, so model this with a 1.6s GCD.
const PetGCD = time.Millisecond * 1600

const (
	Unknown PetAbilityType = iota
	AcidSpit
	Bite
	Claw
	DemoralizingScreech
	FireBreath
	FuriousHowl
	FroststormBreath
	Gore
	LavaBreath
	LightningBreath
	MonstrousBite
	NetherShock
	Pin
	PoisonSpit
	Rake
	Ravage
	SavageRend
	ScorpidPoison
	Smack
	Snatch
	SonicBlast
	SpiritStrike
	SporeCloud
	Stampede
	Sting
	Swipe
	TendonRip
	VenomWebSpray
)

func (hp *HunterPet) NewPetAbility(abilityType PetAbilityType, isPrimary bool) *core.Spell {
	switch abilityType {
	case AcidSpit:
		return hp.newAcidSpit()
	case Bite:
		return hp.newBite()
	case Claw:
		return hp.newClaw()
	case DemoralizingScreech:
		return hp.newDemoralizingScreech()
	case FireBreath:
		return hp.newFireBreath()
	case FroststormBreath:
		return hp.newFroststormBreath()
	case FuriousHowl:
		return hp.newFuriousHowl()
	case Gore:
		return hp.newGore()
	case LavaBreath:
		return hp.newLavaBreath()
	case LightningBreath:
		return hp.newLightningBreath()
	case MonstrousBite:
		return hp.newMonstrousBite()
	case NetherShock:
		return hp.newNetherShock()
	case Pin:
		return hp.newPin()
	case PoisonSpit:
		return hp.newPoisonSpit()
	case Rake:
		return hp.newRake()
	case Ravage:
		return hp.newRavage()
	case SavageRend:
		return hp.newSavageRend()
	case ScorpidPoison:
		return hp.newScorpidPoison()
	case Smack:
		return hp.newSmack()
	case Snatch:
		return hp.newSnatch()
	case SonicBlast:
		return hp.newSonicBlast()
	case SpiritStrike:
		return hp.newSpiritStrike()
	case SporeCloud:
		return hp.newSporeCloud()
	case Stampede:
		return hp.newStampede()
	case Sting:
		return hp.newSting()
	case Swipe:
		return hp.newSwipe()
	case TendonRip:
		return hp.newTendonRip()
	case VenomWebSpray:
		return hp.newVenomWebSpray()
	case Unknown:
		return nil
	default:
		panic("Invalid pet ability type")
	}
}

func (hp *HunterPet) newFocusDump(spellID int32, bp float64, die float64) *core.Spell {
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskMeleeMHSpecial,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage,

		FocusCost: core.FocusCostOptions{
			Cost: 25,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
		},

		DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + 0.07*spell.MeleeAttackPower()
			baseDamage *= hp.killCommandMult()
			spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
		},
	})
}

func (hp *HunterPet) newBite() *core.Spell {
	dbc := hunterinfo.Bite.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)
	return hp.newFocusDump(dbc.SpellID, bp, die)
}
func (hp *HunterPet) newClaw() *core.Spell {
	dbc := hunterinfo.Claw.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)
	return hp.newFocusDump(dbc.SpellID, bp, die)
}
func (hp *HunterPet) newSmack() *core.Spell {
	dbc := hunterinfo.Smack.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)
	return hp.newFocusDump(dbc.SpellID, bp, die)
}

type PetSpecialAbilityConfig struct {
	Type    PetAbilityType
	Cost    float64
	SpellID int32
	School  core.SpellSchool
	GCD     time.Duration
	CD      time.Duration
	Bp      float64
	Die     float64
	APRatio float64

	Dot core.DotConfig

	OnSpellHitDealt func(*core.Simulation, *core.Spell, *core.SpellResult)
}

func (hp *HunterPet) newSpecialAbility(config PetSpecialAbilityConfig) *core.Spell {
	var flags core.SpellFlag
	var applyEffects core.ApplySpellResults
	var procMask core.ProcMask
	onSpellHitDealt := config.OnSpellHitDealt
	if config.School == core.SpellSchoolPhysical {
		flags = core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage
		procMask = core.ProcMaskSpellDamage
		applyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(config.Bp, config.Die) + config.APRatio*spell.MeleeAttackPower()
			baseDamage *= hp.killCommandMult()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)
			if onSpellHitDealt != nil {
				onSpellHitDealt(sim, spell, result)
			}
		}
	} else {
		procMask = core.ProcMaskMeleeMHSpecial
		applyEffects = func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(config.Bp, config.Die) + config.APRatio*spell.MeleeAttackPower()
			baseDamage *= 1 + 0.2*float64(hp.KillCommandAura.GetStacks())
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)
			if onSpellHitDealt != nil {
				onSpellHitDealt(sim, spell, result)
			}
		}
	}

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: config.SpellID},
		SpellSchool: config.School,
		ProcMask:    procMask,
		Flags:       flags,

		DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		FocusCost: core.FocusCostOptions{
			Cost: config.Cost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: config.GCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(config.CD),
			},
		},
		Dot:          config.Dot,
		ApplyEffects: applyEffects,
	})
}

func (hp *HunterPet) newAcidSpit() *core.Spell {
	dbc := hunterinfo.AcidSpit.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)
	acidSpitAuras := hp.NewEnemyAuraArray(core.AcidSpitAura)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    AcidSpit,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolNature,
		Bp:      bp,
		Die:     die,
		APRatio: 0.049,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				aura := acidSpitAuras.Get(result.Target)
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			}
		},
	})
}

func (hp *HunterPet) newDemoralizingScreech() *core.Spell {
	dbc := hunterinfo.DemoralizingScreech.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)
	debuffs := hp.NewEnemyAuraArray(core.DemoralizingScreechAura)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    DemoralizingScreech,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolPhysical,
		Bp:      bp,
		Die:     die,
		APRatio: 0.07,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					debuffs.Get(aoeTarget).Activate(sim)
				}
			}
		},
	})
}

func (hp *HunterPet) newFireBreath() *core.Spell {
	dbc := hunterinfo.FireBreath.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)
	dotBp, dotDie := dbc.GetBPDie(1, hp.Level)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    FireBreath,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolFire,
		Bp:      bp,
		Die:     die,
		APRatio: 0.049,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Fire Breath",
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(dotBp, dotDie) * hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				spell.Dot(result.Target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newFroststormBreath() *core.Spell {
	dbc := hunterinfo.FroststormBreath.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    FroststormBreath,
		Cost:    20,
		GCD:     0,
		CD:      time.Second * 10,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolFrost,
		Bp:      bp,
		Die:     die,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newFuriousHowl() *core.Spell {
	dbc := hunterinfo.FuriousHowl.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, _ := dbc.GetBPDie(0, hp.Level)
	actionID := core.ActionID{SpellID: dbc.SpellID}

	petAura := hp.NewTemporaryStatsAura("FuriousHowl", actionID, stats.Stats{stats.AttackPower: bp, stats.RangedAttackPower: bp}, time.Second*20)
	ownerAura := hp.hunterOwner.NewTemporaryStatsAura("FuriousHowl", actionID, stats.Stats{stats.AttackPower: bp, stats.RangedAttackPower: bp}, time.Second*20)

	howlSpell := hp.RegisterSpell(core.SpellConfig{
		ActionID: actionID,

		FocusCost: core.FocusCostOptions{
			Cost: 20,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 40),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			petAura.Activate(sim)
			ownerAura.Activate(sim)
		},
	})

	hp.hunterOwner.RegisterSpell(core.SpellConfig{
		ActionID:   actionID,
		SpellRanks: hunterinfo.FuriousHowl.GetAllIDs(),
		Flags:      core.SpellFlagAPL | core.SpellFlagMCD,
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return howlSpell.CanCast(sim, target)
		},
		ApplyEffects: func(sim *core.Simulation, target *core.Unit, _ *core.Spell) {
			howlSpell.Cast(sim, target)
		},
	})

	hp.hunterOwner.AddMajorCooldown(core.MajorCooldown{
		Spell: howlSpell,
		Type:  core.CooldownTypeDPS,
	})

	return nil
}

func (hp *HunterPet) newGore() *core.Spell {
	dbc := hunterinfo.Gore.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Gore,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolPhysical,
		Bp:      bp,
		Die:     die,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newLavaBreath() *core.Spell {
	dbc := hunterinfo.LavaBreath.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    LavaBreath,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolFire,
		Bp:      bp,
		Die:     die,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newLightningBreath() *core.Spell {
	dbc := hunterinfo.LightningBreath.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    LightningBreath,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolNature,
		Bp:      bp,
		Die:     die,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newMonstrousBite() *core.Spell {
	dbc := hunterinfo.MonstrousBite.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	procAura := hp.RegisterAura(core.Aura{
		Label:     "Monstrous Strength",
		ActionID:  core.ActionID{SpellID: dbc.SpellID},
		Duration:  time.Second * 12,
		MaxStacks: 3,
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= math.Pow(1.03, float64(oldStacks))
			aura.Unit.PseudoStats.DamageDealtMultiplier *= math.Pow(1.03, float64(newStacks))
		},
	})

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    MonstrousBite,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: 55499,
		School:  core.SpellSchoolPhysical,
		Bp:      bp,
		Die:     die,
		APRatio: 0.07,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				procAura.Activate(sim)
				procAura.AddStack(sim)
			}
		},
	})
}

func (hp *HunterPet) newNetherShock() *core.Spell {
	dbc := hunterinfo.NetherShock.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    NetherShock,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolShadow,
		Bp:      bp,
		Die:     die,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newPin() *core.Spell {
	dbc := hunterinfo.Pin.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(1, hp.Level)

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskEmpty,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 40),
			},
		},

		DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Pin",
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(bp, die) + 0.07*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(result.Target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newPoisonSpit() *core.Spell {
	dbc := hunterinfo.PoisonSpit.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,

		FocusCost: core.FocusCostOptions{
			Cost: 20,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 10),
			},
		},

		DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "PoisonSpit",
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(bp, die) + (0.049/4)*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(result.Target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newRake() *core.Spell {
	dbc := hunterinfo.Rake.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)
	dotBp, dotDie := dbc.GetBPDie(1, hp.Level)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Rake,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 10,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolPhysical,
		Bp:      bp,
		Die:     die,
		APRatio: 0.0175,
		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Rake",
			},
			NumberOfTicks: 3,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(dotBp, dotDie) + 0.0175*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				spell.Dot(result.Target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newRavage() *core.Spell {
	dbc := hunterinfo.Ravage.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Ravage,
		Cost:    0,
		CD:      time.Second * 40,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolPhysical,
		Bp:      bp,
		Die:     die,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newSavageRend() *core.Spell {
	dbc := hunterinfo.SavageRend.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)
	dotBp, dotDie := dbc.GetBPDie(1, hp.Level)

	actionID := core.ActionID{SpellID: dbc.SpellID}

	procAura := hp.RegisterAura(core.Aura{
		Label:    "Savage Rend",
		ActionID: actionID,
		Duration: time.Second * 30,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier *= 1.1
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.DamageDealtMultiplier /= 1.1
		},
	})

	srSpell := hp.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolPhysical,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       core.SpellFlagMeleeMetrics | core.SpellFlagIncludeTargetBonusDamage | core.SpellFlagApplyArmorReduction,

		FocusCost: core.FocusCostOptions{
			Cost: 20,
		},
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 60),
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return hp.IsEnabled()
		},

		DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		CritMultiplier:   2,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SavageRend",
			},
			NumberOfTicks: 3,
			TickLength:    time.Second * 5,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(dotBp, dotDie) + 0.07*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(bp, die) + 0.07*spell.MeleeAttackPower()
			baseDamage *= hp.killCommandMult()
			result := spell.CalcAndDealDamage(sim, target, baseDamage, spell.OutcomeMeleeSpecialHitAndCrit)

			if result.Landed() {
				spell.Dot(target).Apply(sim)
				if result.DidCrit() {
					procAura.Activate(sim)
				}
			}
		},
	})

	hp.hunterOwner.AddMajorCooldown(core.MajorCooldown{
		Spell: srSpell,
		Type:  core.CooldownTypeDPS,
	})

	return nil
}

func (hp *HunterPet) newScorpidPoison() *core.Spell {
	dbc := hunterinfo.ScorpidPoison.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(1, hp.Level)

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,

		FocusCost: core.FocusCostOptions{
			Cost: 20,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 10),
			},
		},

		DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "ScorpidPoison",
			},
			NumberOfTicks: 5,
			TickLength:    time.Second * 2,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(bp, die) + (0.07/5)*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newSnatch() *core.Spell {
	dbc := hunterinfo.Snatch.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Snatch,
		Cost:    20,
		CD:      time.Second * 60,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolPhysical,
		Bp:      bp,
		Die:     die,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newSonicBlast() *core.Spell {
	dbc := hunterinfo.SonicBlast.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    SonicBlast,
		Cost:    80,
		CD:      time.Second * 60,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolNature,
		Bp:      bp,
		Die:     die,
		APRatio: 0.049,
	})
}

func (hp *HunterPet) newSpiritStrike() *core.Spell {
	dbc := hunterinfo.SpiritStrike.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    SpiritStrike,
		Cost:    20,
		GCD:     0,
		CD:      time.Second * 10,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolArcane,
		Bp:      bp,
		Die:     die,
		APRatio: 0.04,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "SpiritStrike",
			},
			NumberOfTicks: 1,
			TickLength:    time.Second * 6,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(bp, die) + 0.04*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				spell.Dot(result.Target).Apply(sim)
			}
		},
	})
}

func (hp *HunterPet) newSporeCloud() *core.Spell {
	dbc := hunterinfo.SporeCloud.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	debuffs := hp.NewEnemyAuraArray(core.SporeCloudAura)
	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskSpellDamage,

		FocusCost: core.FocusCostOptions{
			Cost: 20,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: PetGCD,
			},
			IgnoreHaste: true,
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 10),
			},
		},

		DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			IsAOE: true,
			Aura: core.Aura{
				Label: "SporeCloud",
			},
			NumberOfTicks: 3,
			TickLength:    time.Second * 3,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(bp, die) + (0.049/3)*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, aoeTarget, dot.OutcomeTick)
				}
			},
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
			spell.AOEDot().Apply(sim)
			for _, target := range spell.Unit.Env.Encounter.TargetUnits {
				debuffs.Get(target).Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) newStampede() *core.Spell {
	dbc := hunterinfo.Stampede.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	debuffs := hp.NewEnemyAuraArray(core.StampedeAura)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Stampede,
		Cost:    0,
		CD:      time.Second * 60,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolPhysical,
		Bp:      bp,
		Die:     die,
		APRatio: 0.07,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				debuffs.Get(result.Target).Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) newSting() *core.Spell {
	dbc := hunterinfo.Sting.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	debuffs := hp.NewEnemyAuraArray(core.StingAura)
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Sting,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 6,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolNature,
		Bp:      bp,
		Die:     die,
		APRatio: 0.049,
		OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.Landed() {
				debuffs.Get(result.Target).Activate(sim)
			}
		},
	})
}

func (hp *HunterPet) newSwipe() *core.Spell {
	dbc := hunterinfo.Swipe.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	// TODO: This is frontal cone, but might be more realistic as single-target
	// since pets are hard to control.
	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    Swipe,
		Cost:    20,
		GCD:     PetGCD,
		CD:      time.Second * 5,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolPhysical,
		Bp:      bp,
		Die:     die,
		APRatio: 0.07,
	})
}

func (hp *HunterPet) newTendonRip() *core.Spell {
	dbc := hunterinfo.TendonRip.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(0, hp.Level)

	return hp.newSpecialAbility(PetSpecialAbilityConfig{
		Type:    TendonRip,
		Cost:    20,
		CD:      time.Second * 20,
		SpellID: dbc.SpellID,
		School:  core.SpellSchoolPhysical,
		Bp:      bp,
		Die:     die,
		APRatio: 0,
	})
}

func (hp *HunterPet) newVenomWebSpray() *core.Spell {
	dbc := hunterinfo.VenomWebSpray.GetMaxRank(hp.Level)
	if dbc == nil {
		return nil
	}
	bp, die := dbc.GetBPDie(1, hp.Level)

	return hp.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: dbc.SpellID},
		SpellSchool: core.SpellSchoolNature,
		ProcMask:    core.ProcMaskEmpty,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    hp.NewTimer(),
				Duration: hp.hunterOwner.applyLongevity(time.Second * 40),
			},
		},

		DamageMultiplier: 1 * hp.hunterOwner.markedForDeathMultiplier(),
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "VenomWebSpray",
			},
			NumberOfTicks: 4,
			TickLength:    time.Second * 1,
			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.SnapshotBaseDamage = sim.Roll(bp, die) + 0.07*dot.Spell.MeleeAttackPower()
				dot.SnapshotBaseDamage *= hp.killCommandMult()
				dot.SnapshotAttackerMultiplier = dot.Spell.AttackerDamageMultiplier(dot.Spell.Unit.AttackTables[target.UnitIndex])
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealOutcome(sim, target, spell.OutcomeMeleeSpecialHit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}
