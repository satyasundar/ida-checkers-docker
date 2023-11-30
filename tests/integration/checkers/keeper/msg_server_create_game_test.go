package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/satya/checkers/x/checkers/types"
)

func (suite *IntegrationTestSuite) TestCreate1GameHasSaved() {
	suite.setupSuiteWithBalances()
	goCtx := sdk.WrapSDKContext(suite.ctx)

	suite.msgServer.CreateGame(goCtx, &types.MsgCreateGame{
		Creator: alice,
		Red:     bob,
		Black:   carol,
		Wager:   45,
	})
	keeper := suite.app.CheckersKeeper
	systemInfo, found := keeper.GetSystemInfo(suite.ctx)
	suite.Require().True(found)
	suite.Require().EqualValues(types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: "1",
		FifoTailIndex: "1",
	}, systemInfo)
}
