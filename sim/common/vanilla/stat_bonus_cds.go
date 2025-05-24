package vanilla

import (
	"time"

	"github.com/WoWLegacySims/wotlk/sim/common/helpers"
)

func init() {

	helpers.NewHasteActive(9449, 500, time.Second*30, time.Hour)
}
