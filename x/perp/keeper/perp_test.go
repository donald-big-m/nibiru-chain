package keeper_test

import (
	"testing"
	"time"

	"github.com/NibiruChain/nibiru/collections"
	"github.com/NibiruChain/nibiru/collections/keys"

	simapp2 "github.com/NibiruChain/nibiru/simapp"

	"github.com/cosmos/cosmos-sdk/simapp"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/x/common"

	"github.com/stretchr/testify/require"

	"github.com/NibiruChain/nibiru/x/perp/types"
	"github.com/NibiruChain/nibiru/x/testutil/sample"
)

func TestKeeperClosePosition(t *testing.T) {
	// TODO(mercilex): simulate funding payments
	t.Run("success", func(t *testing.T) {
		t.Log("Setup Nibiru app, pair, and trader")
		nibiruApp, ctx := simapp2.NewTestNibiruAppAndContext(true)
		pair := common.MustNewAssetPair("xxx:yyy")

		t.Log("Set vpool defined by pair on VpoolKeeper")
		vpoolKeeper := &nibiruApp.VpoolKeeper
		vpoolKeeper.CreatePool(
			ctx,
			pair,
			/*tradeLimitRatio*/ sdk.MustNewDecFromStr("0.9"),
			/*quoteAssetReserve*/ sdk.NewDec(10_000_000),
			/*baseAssetReserve*/ sdk.NewDec(5_000_000),
			/*fluctuationLimitRatio*/ sdk.MustNewDecFromStr("0.1"),
			/*maxOracleSpreadRatio*/ sdk.MustNewDecFromStr("0.1"),
			/* maintenanceMarginRatio */ sdk.MustNewDecFromStr("0.0625"),
			/* maxLeverage */ sdk.MustNewDecFromStr("15"),
		)
		require.True(t, vpoolKeeper.ExistsPool(ctx, pair))
		nibiruApp.PricefeedKeeper.ActivePairsStore().Set(ctx, pair, true)

		t.Log("Set vpool defined by pair on PerpKeeper")
		setPairMetadata(nibiruApp.PerpKeeper, ctx, types.PairMetadata{
			Pair: pair,
			CumulativeFundingRates: []sdk.Dec{
				sdk.MustNewDecFromStr("0.2")},
		},
		)

		t.Log("open position for alice - long")
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(time.Now().Add(time.Minute))

		alice := sample.AccAddress()
		err := simapp.FundAccount(nibiruApp.BankKeeper, ctx, alice,
			sdk.NewCoins(sdk.NewInt64Coin("yyy", 300)))
		require.NoError(t, err)

		aliceSide := types.Side_BUY
		aliceQuote := sdk.NewInt(60)
		aliceLeverage := sdk.NewDec(10)
		aliceBaseLimit := sdk.NewDec(150)
		_, err = nibiruApp.PerpKeeper.OpenPosition(
			ctx, pair, aliceSide, alice, aliceQuote, aliceLeverage, aliceBaseLimit)
		require.NoError(t, err)

		t.Log("open position for bob - long")
		// force funding payments
		setPairMetadata(nibiruApp.PerpKeeper, ctx, types.PairMetadata{
			Pair: pair,
			CumulativeFundingRates: []sdk.Dec{
				sdk.MustNewDecFromStr("0.3")},
		})
		bob := sample.AccAddress()
		err = simapp.FundAccount(nibiruApp.BankKeeper, ctx, bob,
			sdk.NewCoins(sdk.NewInt64Coin("yyy", 62)))
		require.NoError(t, err)

		bobSide := types.Side_BUY
		bobQuote := sdk.NewInt(60)
		bobLeverage := sdk.NewDec(10)
		bobBaseLimit := sdk.NewDec(150)

		_, err = nibiruApp.PerpKeeper.OpenPosition(
			ctx, pair, bobSide, bob, bobQuote, bobLeverage, bobBaseLimit)
		require.NoError(t, err)

		t.Log("testing close position")
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).
			WithBlockTime(ctx.BlockTime().Add(1 * time.Minute))

		posResp, err := nibiruApp.PerpKeeper.ClosePosition(ctx, pair, alice)
		require.NoError(t, err)
		require.True(t, posResp.BadDebt.IsZero())
		require.True(t, !posResp.FundingPayment.IsZero() && posResp.FundingPayment.IsPositive())

		position, err := nibiruApp.PerpKeeper.Positions.Get(ctx, keys.Join(pair, keys.String(alice.String())))
		require.ErrorIs(t, err, collections.ErrNotFound)
		require.Empty(t, position)

		// this tests the following issue https://github.com/NibiruChain/nibiru/issues/645
		// in which opening a position from the same address on the same pair
		// was not possible after calling close position, due to bad data clearance.
		_, err = nibiruApp.PerpKeeper.OpenPosition(ctx, pair, aliceSide, alice, aliceQuote, aliceLeverage, aliceBaseLimit)
		require.NoError(t, err)
	})
}
