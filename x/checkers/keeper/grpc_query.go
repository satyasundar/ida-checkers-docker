package keeper

import (
	"github.com/satya/checkers/x/checkers/types"
)

var _ types.QueryServer = Keeper{}
