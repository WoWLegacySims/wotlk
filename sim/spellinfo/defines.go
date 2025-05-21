package spellinfo

import (
	"math"
	"slices"
	"time"
)

type SpellEffect struct {
	BasePoints     float64
	Die            float64
	LevelScaling   float64
	Coefficient    float64
	PointsPerCombo float64
	ChainAmplitude float64
	AuraPeriod     int32
	TriggerSpell   int32
	MiscValueA     float64
	MiscValueB     float64
}

type SpellInfo struct {
	SpellID  int32
	MinLevel int32
	MaxLevel int32
	Duration int32
	CastTime time.Duration
	BaseCost float64
	Effects  [3]SpellEffect
	Rank     int32
	ManaCost float64
}

type SpellDBC struct {
	SpellInfos []*SpellInfo
}

func (info *SpellInfo) GetCoefficient(effect int32) float64 {
	if effect > 2 || effect < 0 {
		panic("Effect has to be between 0 and 2")
	}

	return info.Effects[effect].Coefficient
}

func (info *SpellInfo) GetLevelPenalty(level int32) float64 {
	if info.MinLevel <= 0 || info.MinLevel >= info.MaxLevel {
		return 1.0
	}
	lvlPenalty := 0.0
	if info.MinLevel < 20 {
		lvlPenalty = float64(20-info.MinLevel) * 3.75
	}

	lvlFactor := min(float64(info.MinLevel+6)/float64(level), 1.0)

	lvlFactor -= lvlFactor * lvlPenalty / 100
	return lvlFactor
}

func (info *SpellInfo) GetBPDie(effect int32, level int32) (float64, float64) {
	if effect > 2 || effect < 0 {
		panic("Effect has to be between 0 and 2")
	}

	bp := info.Effects[effect].BasePoints
	if math.Abs(info.Effects[effect].Die) == 1 {
		bp += info.Effects[effect].Die
	}
	maxLevel := min(level, info.MaxLevel)
	if maxLevel == 0 {
		maxLevel = level
	}

	return bp + float64(level-info.MinLevel)*info.Effects[effect].LevelScaling, info.Effects[effect].Die
}

func (dbc *SpellDBC) GetByID(id int32) *SpellInfo {
	index := slices.IndexFunc(dbc.SpellInfos, func(info *SpellInfo) bool { return info.SpellID == id })
	if index == -1 {
		return nil
	}
	return dbc.SpellInfos[index]
}

func (spell *SpellDBC) GetMaxRank(level int32) *SpellInfo {
	var info *SpellInfo
	for _, s := range spell.SpellInfos {
		if level >= s.MinLevel {
			if info == nil || info.MinLevel < s.MinLevel {
				info = s
			}
		}
	}
	return info
}

func (spell *SpellDBC) GetDownRank(level int32) *SpellInfo {
	var info *SpellInfo
	var downRankIndex int
	for i, s := range spell.SpellInfos {
		if level >= s.MinLevel {
			if info == nil || info.MinLevel < s.MinLevel {
				info = s
				downRankIndex = i
			}
		}
	}
	if downRankIndex > 0 {
		return spell.SpellInfos[downRankIndex-1]
	}
	return nil
}
