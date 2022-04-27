package keeper

import (
	"errors"

	"github.com/NibiruChain/nibiru/x/dex/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) updatePoolForSwap(
	ctx sdk.Context,
	pool types.Pool,
	sender sdk.AccAddress,
	tokenIn sdk.Coin,
	tokenOut sdk.Coin,
) (err error) {
	if err = pool.ApplySwap(tokenIn, tokenOut); err != nil {
		return err
	}
	k.SetPool(ctx, pool)

	if err = k.bankKeeper.SendCoins(
		ctx,
		/*from=*/ sender,
		/*to=*/ pool.GetAddress(),
		/*coins=*/ sdk.Coins{tokenIn},
	); err != nil {
		return err
	}

	if err = k.bankKeeper.SendCoins(
		ctx,
		/*from=*/ pool.GetAddress(),
		/*to=*/ sender,
		/*coins=*/ sdk.Coins{tokenOut},
	); err != nil {
		return err
	}

	k.RecordTotalLiquidityIncrease(ctx, sdk.Coins{tokenIn})
	k.RecordTotalLiquidityDecrease(ctx, sdk.Coins{tokenOut})

	return err
}

/*
Given a poolId and the amount of tokens to swap in, returns the number of tokens out
received, specified by the tokenOutDenom.

For example, if pool 1 has 100foo and 100bar, this function can be called with
tokenIn=10foo and tokenOutDenom=bar.

args:
  - ctx: the cosmos-sdk context
  - sender: the address wishing to perform the swap
  - poolId: the pool id number
  - tokenIn: the amount of tokens to given to the pool
  - tokenOutDenom: the denom of the token taken out of the pool

ret:
  - tokenOut: the amount of tokens taken out of the pool
  - err: error if any
*/
func (k Keeper) SwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	tokenIn sdk.Coin,
	tokenOutDenom string,
) (tokenOut sdk.Coin, err error) {
	if tokenIn.Denom == tokenOutDenom {
		return sdk.Coin{}, types.ErrSameTokenDenom
	}

	if err = k.checkEnoughBalance(ctx, tokenIn, sender); err != nil {
		return sdk.Coin{}, err
	}

	pool, err := k.FetchPool(ctx, poolId)
	if err != nil {
		return sdk.Coin{}, err
	}

	tokenOut, err = pool.CalcOutAmtGivenIn(tokenIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, err
	}

	if tokenOut.Amount.LTE(sdk.ZeroInt()) {
		return sdk.Coin{}, errors.New("tokenOut amount must be greater than zero")
	}

	err = k.updatePoolForSwap(ctx, pool, sender, tokenIn, tokenOut)
	if err != nil {
		return sdk.Coin{}, err
	}

	return tokenOut, nil
}
