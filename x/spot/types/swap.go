package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/x/spot/math"
)

/*
CalcOutAmtGivenIn Calculates the amount of tokenOut given tokenIn, deducting the swap fee.
Solved using the SolveConstantProductInvariant AMM curve.
Only supports single asset swaps.

args:
  - tokenIn: the amount of tokens to swap
  - tokenOutDenom: the target token denom
  - noFee: whether we want to bypass swap fee (for single asset join)

ret:
  - tokenOut: the tokens received from the swap
  - fee: the fee deducted from the swap
  - err: error if any
*/
func (pool Pool) CalcOutAmtGivenIn(tokenIn sdk.Coin, tokenOutDenom string, noFee bool) (
	tokenOut sdk.Coin, fee sdk.Coin, err error,
) {
	_, poolAssetIn, err := pool.getPoolAssetAndIndex(tokenIn.Denom)
	if err != nil {
		return tokenOut, fee, err
	}

	_, poolAssetOut, err := pool.getPoolAssetAndIndex(tokenOutDenom)
	if err != nil {
		return tokenOut, fee, err
	}

	var tokenAmountInAfterFee sdkmath.LegacyDec
	if noFee {
		tokenAmountInAfterFee = sdkmath.LegacyNewDecFromInt(tokenIn.Amount)
	} else {
		tokenAmountInAfterFee = sdkmath.LegacyNewDecFromInt(tokenIn.Amount).Mul(
			sdkmath.LegacyOneDec().Sub(pool.PoolParams.SwapFee),
		)
	}
	feeAmount := sdkmath.LegacyNewDecFromInt(tokenIn.Amount).Sub(tokenAmountInAfterFee)
	fee = sdk.NewCoin(tokenIn.Denom, feeAmount.TruncateInt())

	poolTokenInBalance := sdkmath.LegacyNewDecFromInt(poolAssetIn.Token.Amount)
	poolTokenInBalancePostSwap := poolTokenInBalance.Add(tokenAmountInAfterFee)

	// deduct swapfee on the in asset
	// delta balanceOut is positive(tokens inside the pool decreases)
	var tokenAmountOut sdkmath.Int
	if pool.PoolParams.PoolType == PoolType_STABLESWAP {
		tokenAmountOut, err = pool.Exchange(sdk.NewCoin(tokenIn.Denom, tokenAmountInAfterFee.TruncateInt()), tokenOutDenom)

		if err != nil {
			return
		}
	} else if pool.PoolParams.PoolType == PoolType_BALANCER {
		tokenAmountOut = math.SolveConstantProductInvariant(
			/*xPrior=*/ poolTokenInBalance,
			/*xAfter=*/ poolTokenInBalancePostSwap,
			/*xWeight=*/ sdkmath.LegacyNewDecFromInt(poolAssetIn.Weight),
			/*yPrior=*/ sdkmath.LegacyNewDecFromInt(poolAssetOut.Token.Amount),
			/*yWeight=*/ sdkmath.LegacyNewDecFromInt(poolAssetOut.Weight),
		).TruncateInt()
	}

	if tokenAmountOut.IsZero() {
		return tokenOut, fee, fmt.Errorf("tokenIn (%s) must be higher to perform a swap", tokenIn.Denom)
	}

	return sdk.NewCoin(tokenOutDenom, tokenAmountOut), fee, nil
}

/*
Calculates the amount of tokenIn required to obtain tokenOut coins from a swap,
accounting for additional fees.
Only supports single asset swaps.
This function is the inverse of CalcOutAmtGivenIn.

args:
  - tokenOut: the amount of tokens to swap
  - tokenInDenom: the target token denom

ret:
  - tokenIn: the tokens received from the swap
  - err: error if any
*/
func (pool Pool) CalcInAmtGivenOut(tokenOut sdk.Coin, tokenInDenom string) (
	tokenIn sdk.Coin, err error,
) {
	if pool.PoolParams.PoolType == PoolType_BALANCER {
		return pool.CalcInAmtGivenOutBalancer(tokenOut, tokenInDenom)
	} else if pool.PoolParams.PoolType == PoolType_STABLESWAP {
		return pool.CalcInAmtGivenOutStableswap(tokenOut, tokenInDenom)
	}
	return sdk.Coin{}, ErrInvalidPoolType
}

/*
Calculates the amount of tokenIn required to obtain tokenOut coins from a swap,
accounting for additional fees. This is not implemented yet in curve and in Nibiru.
*/
func (pool Pool) CalcInAmtGivenOutStableswap(tokenOut sdk.Coin, tokenInDenom string) (
	tokenIn sdk.Coin, err error,
) {
	return sdk.Coin{}, ErrNotImplemented
}

/*
Calculates the amount of tokenIn required to obtain tokenOut coins from a swap,
accounting for additional fees.
Only supports single asset swaps.
This function is the inverse of CalcOutAmtGivenIn.

args:
  - tokenOut: the amount of tokens to swap
  - tokenInDenom: the target token denom

ret:
  - tokenIn: the tokens received from the swap
  - err: error if any
*/
func (pool Pool) CalcInAmtGivenOutBalancer(tokenOut sdk.Coin, tokenInDenom string) (
	tokenIn sdk.Coin, err error,
) {
	_, poolAssetOut, err := pool.getPoolAssetAndIndex(tokenOut.Denom)
	if err != nil {
		return tokenIn, err
	}

	_, poolAssetIn, err := pool.getPoolAssetAndIndex(tokenInDenom)
	if err != nil {
		return tokenIn, err
	}

	// assuming the user wishes to withdraw 'tokenOut', the balance of 'tokenOut' post swap will be lower
	poolTokenOutBalance := sdkmath.LegacyNewDecFromInt(poolAssetOut.Token.Amount)
	poolTokenOutBalancePostSwap := poolTokenOutBalance.Sub(sdkmath.LegacyNewDecFromInt(tokenOut.Amount))
	// (x_0)(y_0) = (x_0 + in)(y_0 - out)
	tokenAmountIn := math.SolveConstantProductInvariant(
		/*xPrior=*/ poolTokenOutBalance,
		/*xAfter=*/ poolTokenOutBalancePostSwap,
		/*xWeight=*/ sdkmath.LegacyNewDecFromInt(poolAssetOut.Weight),
		/*yPrior=*/ sdkmath.LegacyNewDecFromInt(poolAssetIn.Token.Amount),
		/*yWeight=*/ sdkmath.LegacyNewDecFromInt(poolAssetIn.Weight),
	).Neg()

	// We deduct a swap fee on the input asset. The swap happens by following the invariant curve on the input * (1 - swap fee)
	// and then the swap fee is added to the pool.
	// Thus in order to give X amount out, we solve the invariant for the invariant input. However invariant input = (1 - swapfee) * trade input.
	// Therefore we divide by (1 - swapfee) here
	tokenAmountInBeforeFee := tokenAmountIn.Quo(sdkmath.LegacyOneDec().Sub(pool.PoolParams.SwapFee)).Ceil().TruncateInt()
	return sdk.NewCoin(tokenInDenom, tokenAmountInBeforeFee), nil
}

/*
Applies a swap to the pool by adding tokenIn and removing tokenOut from pool asset balances.

args:
  - tokenIn: the amount of token to deposit
  - tokenOut: the amount of token to withdraw

ret:
  - err: error if any
*/
func (pool *Pool) ApplySwap(tokenIn sdk.Coin, tokenOut sdk.Coin) (err error) {
	if tokenIn.Amount.LTE(sdkmath.ZeroInt()) {
		return fmt.Errorf("tokenIn (%s) cannot be zero", tokenIn.Denom)
	}
	if tokenOut.Amount.LTE(sdkmath.ZeroInt()) {
		return fmt.Errorf("tokenOut (%s) cannot be zero", tokenOut.Denom)
	}

	_, poolAssetIn, err := pool.getPoolAssetAndIndex(tokenIn.Denom)
	if err != nil {
		return err
	}

	_, poolAssetOut, err := pool.getPoolAssetAndIndex(tokenOut.Denom)
	if err != nil {
		return err
	}

	poolAssetIn.Token.Amount = poolAssetIn.Token.Amount.Add(tokenIn.Amount)
	poolAssetOut.Token.Amount = poolAssetOut.Token.Amount.Sub(tokenOut.Amount)

	return pool.updatePoolAssetBalances(poolAssetIn.Token, poolAssetOut.Token)
}
