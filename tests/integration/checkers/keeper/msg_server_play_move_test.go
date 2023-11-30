package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/satya/checkers/x/checkers/testutil"
	"github.com/satya/checkers/x/checkers/types"
)

func (suite *IntegrationTestSuite) setupSuiteWithOneGameForPlayMove() {
	suite.setupSuiteWithBalances()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.CreateGame(goCtx, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
		Wager:   45,
	})
}

func (suite *IntegrationTestSuite) TestPlayMovePlayerPaid() {
	suite.setupSuiteWithOneGameForPlayMove()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalance(balBob, bob)
	suite.RequireBankBalance(balCarol, carol)
	suite.RequireBankBalance(0, checkersModuleAddress)

	suite.msgServer.PlayMove(goCtx, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})

	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalance(balBob-45, bob)
	suite.RequireBankBalance(balCarol, carol)
	suite.RequireBankBalance(45, checkersModuleAddress)
}


func (suite *IntegrationTestSuite) TestPlayMoveToWinnerBankPaidDifferentTokens() {
	suite.setupSuiteWithOneGameForPlayMove()
	goCtx := sdk.WrapSDKContext(suite.ctx)
	suite.msgServer.CreateGame(goCtx, &types.MsgCreateGame{
		Creator: alice,
		Black:   bob,
		Red:     carol,
		Wager:   46,
		Denom:   "coin",
	})
	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalanceWithDenom(0, "coin", alice)
	suite.RequireBankBalance(balBob, bob)
	suite.RequireBankBalanceWithDenom(balBob, "coin", bob)
	suite.RequireBankBalance(balCarol, carol)
	suite.RequireBankBalanceWithDenom(balCarol, "coin", carol)
	suite.RequireBankBalance(0, checkersModuleAddress)
	testutil.PlayAllMoves(suite.T(), suite.msgServer, sdk.WrapSDKContext(suite.ctx), "1", bob, carol, testutil.Game1Moves)
	testutil.PlayAllMoves(suite.T(), suite.msgServer, sdk.WrapSDKContext(suite.ctx), "2", bob, carol, testutil.Game1Moves)
	suite.RequireBankBalance(balAlice, alice)
	suite.RequireBankBalanceWithDenom(0, "coin", alice)
	suite.RequireBankBalance(balBob+45, bob)
	suite.RequireBankBalanceWithDenom(balBob+46, "coin", bob)
	suite.RequireBankBalance(balCarol-45, carol)
	suite.RequireBankBalanceWithDenom(balCarol-46, "coin", carol)
	suite.RequireBankBalance(0, checkersModuleAddress)
	suite.RequireBankBalanceWithDenom(0, "coin", checkersModuleAddress)
}