package keeper_test

import (
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/NibiruChain/nibiru/x/common"
	types2 "github.com/NibiruChain/nibiru/x/dex/types"
	"github.com/NibiruChain/nibiru/x/stablecoin/mock"
	"github.com/NibiruChain/nibiru/x/stablecoin/types"
	"github.com/NibiruChain/nibiru/x/testutil"
	"github.com/NibiruChain/nibiru/x/testutil/sample"
)

func TestKeeper_GetStableMarketCap(t *testing.T) {
	matrixApp, ctx := testutil.NewNibiruApp(false)
	k := matrixApp.StablecoinKeeper

	// We set some supply
	err := k.BankKeeper.MintCoins(ctx, types.ModuleName, sdktypes.NewCoins(
		sdktypes.NewInt64Coin(common.StableDenom, 1_000_000),
	))
	require.NoError(t, err)

	// We set some supply
	marketCap := k.GetStableMarketCap(ctx)

	require.Equal(t, sdktypes.NewInt(1_000_000), marketCap)
}

func TestKeeper_GetGovMarketCap(t *testing.T) {
	matrixApp, ctx := testutil.NewNibiruApp(false)
	keeper := matrixApp.StablecoinKeeper

	poolAccountAddr := sample.AccAddress()
	poolParams := types2.PoolParams{
		SwapFee: sdktypes.NewDecWithPrec(3, 2),
		ExitFee: sdktypes.NewDecWithPrec(3, 2),
	}
	poolAssets := []types2.PoolAsset{
		{
			Token:  sdktypes.NewInt64Coin(common.GovDenom, 2_000_000),
			Weight: sdktypes.NewInt(100),
		},
		{
			Token:  sdktypes.NewInt64Coin(common.StableDenom, 1_000_000),
			Weight: sdktypes.NewInt(100),
		},
	}

	pool, err := types2.NewPool(1, poolAccountAddr, poolParams, poolAssets)
	require.NoError(t, err)
	keeper.DexKeeper = mock.NewKeeper(pool)

	// We set some supply
	err = keeper.BankKeeper.MintCoins(ctx, types.ModuleName, sdktypes.NewCoins(
		sdktypes.NewInt64Coin(common.GovDenom, 1_000_000),
	))
	require.NoError(t, err)

	marketCap, err := keeper.GetGovMarketCap(ctx)
	require.NoError(t, err)

	require.Equal(t, sdktypes.NewInt(2_000_000), marketCap) // 1 * 10^6 * 2 (price of gov token)
}
