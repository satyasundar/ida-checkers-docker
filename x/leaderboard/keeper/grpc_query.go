package keeper

import (
	"github.com/satya/checkers/x/leaderboard/types"
)

var _ types.QueryServer = Keeper{}
