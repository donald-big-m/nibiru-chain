package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/x/common"
)

func (k Keeper) GetSupplyUSDM(
	ctx sdk.Context,
) sdk.Coin {
	return k.BankKeeper.GetSupply(ctx, common.StableDenom)
}

func (k Keeper) GetSupplyMTRX(
	ctx sdk.Context,
) sdk.Coin {
	return k.BankKeeper.GetSupply(ctx, common.GovDenom)
}
