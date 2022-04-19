package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/x/common"
)

func (k Keeper) GetSupplyNUSD(
	ctx sdk.Context,
) sdk.Coin {
	return k.BankKeeper.GetSupply(ctx, common.StableDenom)
}

func (k Keeper) GetSupplyNIBI(
	ctx sdk.Context,
) sdk.Coin {
	return k.BankKeeper.GetSupply(ctx, common.GovDenom)
}

func (k Keeper) GetStableMarketCap(ctx sdk.Context) sdk.Int {
	return k.GetSupplyNUSD(ctx).Amount
}

func (k Keeper) GetGovMarketCap(ctx sdk.Context) (sdk.Int, error) {
	pairID, err := k.DexKeeper.GetFromPair(ctx, common.GovDenom, common.StableDenom)
	if err != nil {
		return sdk.Int{}, err
	}

	pool := k.DexKeeper.FetchPool(ctx, pairID)

	price, err := pool.CalcSpotPrice(common.GovDenom, common.StableDenom)
	if err != nil {
		return sdk.Int{}, err
	}

	nibiSupply := k.GetSupplyNIBI(ctx)

	return nibiSupply.Amount.ToDec().Mul(price).RoundInt(), nil
}

// GetLiquidityRatio returns the liquidity ratio defined as govMarketCap / stableMarketCap
func (k Keeper) GetLiquidityRatio(ctx sdk.Context) (sdk.Dec, error) {
	govMarketCap, err := k.GetGovMarketCap(ctx)
	if err != nil {
		return sdk.Dec{}, err
	}

	stableMarketCap := k.GetStableMarketCap(ctx)
	if stableMarketCap.Equal(sdk.ZeroInt()) {
		return sdk.Dec{}, fmt.Errorf("stable maket cap is equal to zero")
	}

	return govMarketCap.ToDec().Quo(stableMarketCap.ToDec()), nil
}
