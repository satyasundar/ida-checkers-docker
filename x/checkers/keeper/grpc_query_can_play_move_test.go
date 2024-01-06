package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/satya/checkers/testutil/keeper"
	"github.com/satya/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
)

type canPlayGameCase struct {
	desc     string
	game     types.StoredGame
	request  *types.QueryCanPlayMoveRequest
	response *types.QueryCanPlayMoveResponse
	err      string
}

var (
	canPlayOkResponse = &types.QueryCanPlayMoveResponse{
		Possible: true,
		Reason:   "ok",
	}

	canPlayTestRange = []canPlayGameCase{
		{
			desc: "first move by black",
			game: types.StoredGame{
				Index:  "1",
				Board:  "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
				Turn:   "b",
				Winner: "*",
			},
			request: &types.QueryCanPlayMoveRequest{
				GameIndex: "1",
				Player:    "b",
				FromX:     1,
				FromY:     2,
				ToX:       2,
				ToY:       3,
			},
			response: canPlayOkResponse,
			err:      "nil",
		},
		{
			desc: "Nil request, wrong",
			game: types.StoredGame{
				Index:  "1",
				Board:  "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
				Turn:   "b",
				Winner: "*",
			},
			request:  nil,
			response: nil,
			err:      "rpc error: code = InvalidArgument desc = invalid request",
		},
		{
			desc: "First move by red, wrong",
			game: types.StoredGame{
				Index:  "1",
				Board:  "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
				Turn:   "b",
				Winner: "*",
			},
			request: &types.QueryCanPlayMoveRequest{
				GameIndex: "1",
				Player:    "r",
				FromX:     1,
				FromY:     2,
				ToX:       2,
				ToY:       3,
			},
			response: &types.QueryCanPlayMoveResponse{
				Possible: false,
				Reason:   "player tried to play out of turn: red",
			},
			err: "nil",
		},
	}
)

func TestCanPlayCasesAsExpected(t *testing.T) {
	for _, testCase := range canPlayTestRange {
		keeper, ctx := keepertest.CheckersKeeper(t)
		goCtx := sdk.WrapSDKContext(ctx)
		t.Run(testCase.desc, func(t *testing.T) {
			keeper.SetStoredGame(ctx, testCase.game)
			response, err := keeper.CanPlayMove(goCtx, testCase.request)
			if testCase.response == nil {
				require.Nil(t, response)
			} else {
				require.EqualValues(t, testCase.response, response)
			}
			if testCase.err == "nil" {
				require.Nil(t, err)
			} else {
				require.EqualError(t, err, testCase.err)
			}
		})
	}

}
