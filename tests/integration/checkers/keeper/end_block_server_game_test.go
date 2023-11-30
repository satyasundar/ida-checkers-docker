package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/satya/checkers/x/checkers/types"
)

func (suite *IntegrationTestSuite) TestForfeitPlayedTwicePaidEmitted() {
	suite.setupSuiteWithOneGameForPlayMove()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.PlayMove(goCtx, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	suite.msgServer.PlayMove(goCtx, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	keeper := suite.app.CheckersKeeper
	game1, found := keeper.GetStoredGame(suite.ctx, "1")
	suite.Require().True(found)
	oldDeadline := types.FormatDeadline(suite.ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(suite.ctx, game1)
	keeper.ForfeitExpiredGames(goCtx)

	events := sdk.StringifyEvents(suite.ctx.EventManager().ABCIEvents())
	suite.Require().Len(events, 7)

	forfeitEvent := events[2]
	suite.Require().EqualValues(sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "r"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, forfeitEvent)

	transferEvent := events[6]
	suite.Require().Equal(transferEvent.Type, "transfer")
	suite.Require().EqualValues([]sdk.Attribute{
		{Key: "recipient", Value: carol},
		{Key: "sender", Value: checkersModuleAddress},
		{Key: "amount", Value: "90stake"},
	}, transferEvent.Attributes[6:])
}
