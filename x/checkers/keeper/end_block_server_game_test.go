package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/satya/checkers/x/checkers/types"
	"github.com/stretchr/testify/require"
)

func TestForfeitUnplayed(t *testing.T) {
	_, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)

	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, systemInfo)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "*"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeitPlayedOnce(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)

	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, systemInfo)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "*"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeitPlayedTwice(t *testing.T) {
	msgServer, keeper, context, ctrl, escrow, board := setupMsgServerWithOneGameForPlayMoveWithMock(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	payBob := escrow.ExpectPay(context, bob, 45).Times(1)
	payCarol := escrow.ExpectPay(context, carol, 45).Times(1).After(payBob)
	escrow.ExpectRefund(context, carol, 90).Times(1).After(payCarol)
	board.Expectwin(context, carol).Times(1)
	board.ExpectForfeit(context, bob).Times(1)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	game1, found = keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		MoveCount:   uint64(2),
		BeforeIndex: "-1",
		AfterIndex:  "-1",
		Deadline:    oldDeadline,
		Winner:      "r",
		Wager:       45,
		Denom:       "stake",
	}, game1)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, systemInfo)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "r"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeit2OldestPlayedOnceIn1Call(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)

	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
		Wager:   46,
		Denom:   "coin",
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "2",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Black:   alice,
		Red:     bob,
		Wager:   47,
		Denom:   "gold",
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	game2, found := keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	game2.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game2)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)
	_, found = keeper.GetStoredGame(ctx, "2")
	require.False(t, found)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        4,
		FifoHeadIndex: "3",
		FifoTailIndex: "3",
	}, systemInfo)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	event := events[0]
	require.EqualValues(t,
		sdk.StringEvent{
			Type: "game-forfeited",
			Attributes: []sdk.Attribute{
				{Key: "game-index", Value: "1"},
				{Key: "winner", Value: "*"},
				{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
				{Key: "game-index", Value: "2"},
				{Key: "winner", Value: "*"},
				{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
			},
		}, event)
}
