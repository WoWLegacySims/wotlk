package core

import (
	"github.com/WoWLegacySims/wotlk/sim/core/proto"
	"github.com/WoWLegacySims/wotlk/sim/core/stats"
)

type BaseStatsKey struct {
	Race  proto.Race
	Class proto.Class
}

var BaseStats = map[BaseStatsKey]stats.Stats{}

// To calculate base stats, get a naked level 70 of the race/class you want, ideally without any talents to mess up base stats.
//  Basic stats are as-shown (str/agi/stm/int/spirit)

// Base Spell Crit is calculated by
//   1. Take as-shown value (troll shaman have 3.5%)
//   2. Calculate the bonus from int (for troll shaman that would be 104/78.1=1.331% crit)
//   3. Subtract as-shown from int bouns (3.5-1.331=2.169)
//   4. 2.169*22.08 (rating per crit percent) = 47.89 crit rating.

// Base mana can be looked up here: https://wowwiki-archive.fandom.com/wiki/Base_mana

// These are also scattered in various dbc/casc files,
// `octbasempbyclass.txt`, `combatratings.txt`, `chancetospellcritbase.txt`, etc.

var RaceOffsets = map[proto.Race]stats.Stats{
	proto.Race_RaceUnknown: stats.Stats{},
	proto.Race_RaceHuman:   stats.Stats{},
	proto.Race_RaceOrc: {
		stats.Agility:   -3,
		stats.Strength:  3,
		stats.Intellect: -3,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
	proto.Race_RaceDwarf: {
		stats.Agility:   -4,
		stats.Strength:  5,
		stats.Intellect: -1,
		stats.Spirit:    -1,
		stats.Stamina:   1,
	},
	proto.Race_RaceNightElf: {
		stats.Agility:   4,
		stats.Strength:  -4,
		stats.Intellect: 0,
		stats.Spirit:    0,
		stats.Stamina:   0,
	},
	proto.Race_RaceUndead: {
		stats.Agility:   -2,
		stats.Strength:  -1,
		stats.Intellect: -2,
		stats.Spirit:    5,
		stats.Stamina:   0,
	},
	proto.Race_RaceTauren: {
		stats.Agility:   -4,
		stats.Strength:  5,
		stats.Intellect: -4,
		stats.Spirit:    2,
		stats.Stamina:   1,
	},
	proto.Race_RaceGnome: {
		stats.Agility:   2,
		stats.Strength:  -5,
		stats.Intellect: 3,
		stats.Spirit:    0,
		stats.Stamina:   0,
	},
	proto.Race_RaceTroll: {
		stats.Agility:   2,
		stats.Strength:  1,
		stats.Intellect: -4,
		stats.Spirit:    1,
		stats.Stamina:   0,
	},
	proto.Race_RaceBloodElf: {
		stats.Agility:   2,
		stats.Strength:  -3,
		stats.Intellect: 3,
		stats.Spirit:    -2,
		stats.Stamina:   0,
	},
	proto.Race_RaceDraenei: {
		stats.Agility:   -3,
		stats.Strength:  1,
		stats.Intellect: 0,
		stats.Spirit:    2,
		stats.Stamina:   0,
	},
}

var BaseAp = map[proto.Class]stats.Stats{
	proto.Class_ClassUnknown:     {},
	proto.Class_ClassWarrior:     {stats.AttackPower: 220.0},
	proto.Class_ClassPaladin:     {stats.AttackPower: 220.0},
	proto.Class_ClassHunter:      {stats.AttackPower: 140.0, stats.RangedAttackPower: 150.0},
	proto.Class_ClassRogue:       {stats.AttackPower: 140.0},
	proto.Class_ClassPriest:      {},
	proto.Class_ClassDeathknight: {stats.AttackPower: 220.0},
	proto.Class_ClassShaman:      {stats.AttackPower: 140.0},
	proto.Class_ClassMage:        {},
	proto.Class_ClassWarlock:     {stats.AttackPower: -10},
	proto.Class_ClassDruid:       {stats.AttackPower: -20},
}

func AddBaseStatsCombo(r proto.Race, c proto.Class) {
	BaseStats[BaseStatsKey{Race: r, Class: c}] = ClassBaseStats[c].Add(RaceOffsets[r]).Add(BaseCrit[c]).Add(BaseAp[c])
}

func init() {
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassDruid)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassDruid)

	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassDeathknight)
	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassDeathknight)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassHunter)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassHunter)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassMage)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassMage)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassPaladin)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassPaladin)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassPaladin)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassPaladin)

	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassPriest)
	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassPriest)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassRogue)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassRogue)

	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassShaman)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassShaman)

	AddBaseStatsCombo(proto.Race_RaceBloodElf, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassWarlock)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassWarlock)

	AddBaseStatsCombo(proto.Race_RaceDraenei, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceDwarf, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceGnome, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceHuman, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceNightElf, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceOrc, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceTauren, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceTroll, proto.Class_ClassWarrior)
	AddBaseStatsCombo(proto.Race_RaceUndead, proto.Class_ClassWarrior)
}
