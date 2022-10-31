package types

import (
	"errors"
	math "math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
For 2 asset pools, swap first to maximize the amount of tokens deposited in the pool.
A user can deposit either one or 2 tokens, and we will swap first the biggest individual share and then join the pool.

args:
  - tokensIn: the tokens to add to the pool

ret:
  - numShares: the number of LP shares given to the user for the deposit
  - remCoins: the number of coins remaining after the deposit
  - err: error if any
*/
func (pool *Pool) SwapForSwapAndJoin(tokensIn sdk.Coins) (
	out sdk.Coin, err error,
) {
	if len(pool.PoolAssets) != 2 {
		err = errors.New("swap and add tokens to pool only available for 2 assets pool")
		return
	}

	var x sdk.Int
	var xDenom string
	var y sdk.Int

	// check who's x and y (x/)
	if len(tokensIn) == 1 {
		x = tokensIn[0].Amount
		xDenom = tokensIn[0].Denom

		y = sdk.ZeroInt()
	} else {
		// 2 assets
		poolLiquidity := pool.PoolBalances()

		s0 := tokensIn[0].Amount.ToDec().Quo(poolLiquidity.AmountOfNoDenomValidation(tokensIn[0].Denom).ToDec())
		s1 := tokensIn[1].Amount.ToDec().Quo(poolLiquidity.AmountOfNoDenomValidation(tokensIn[1].Denom).ToDec())

		if s0.GTE(s1) {
			x = tokensIn[0].Amount
			y = tokensIn[1].Amount

			xDenom = tokensIn[0].Denom
		} else {
			x = tokensIn[1].Amount
			y = tokensIn[0].Amount

			xDenom = tokensIn[1].Denom
		}
	}

	xIndex, xPoolAsset, err := pool.getPoolAssetAndIndex(xDenom)
	lx := xPoolAsset.Token.Amount
	ly := pool.PoolAssets[1-xIndex].Token.Amount

	// x'=\sqrt{\frac{xk+kl_x}{y+l_y}}-l_x;\:x'=-\sqrt{\frac{xk+kl_x}{y+l_y}}-l_x
	k := lx.Mul(ly)

	xSwap := sdk.NewInt(
		int64(math.Sqrt(
			(x.Mul(k).Add(k.Mul(lx))).ToDec().Quo(
				y.Add(ly).ToDec()).MustFloat64()))).Sub(lx)

	return sdk.NewCoin(pool.PoolAssets[xIndex].Token.Denom, xSwap), err
}

/*
Takes a pool and the amount of tokens desired to add to the pool,
and calculates the number of pool shares and remaining coins after theoretically
adding the tokensIn to the pool.

Note that this function is pure/read-only. It only calculates the theoretical amoount
and doesn't modify the actual state.

args:
  - tokensIn: a slice of coins to add to the pool

ret:
  - numShares: the number of LP shares representing the maximal number of tokens added to the pool
  - remCoins: the remaining number of coins after adding the tokens
  - err: error if any
*/
func (pool Pool) numSharesOutFromTokensIn(tokensIn sdk.Coins) (
	numShares sdk.Int, remCoins sdk.Coins, err error,
) {
	coinShareRatios := make([]sdk.Dec, len(tokensIn))
	minShareRatio := sdk.MaxSortableDec
	maxShareRatio := sdk.ZeroDec()

	poolLiquidity := pool.PoolBalances()

	for i, coin := range tokensIn {
		shareRatio := coin.Amount.ToDec().QuoInt(
			poolLiquidity.AmountOfNoDenomValidation(coin.Denom),
		)
		if shareRatio.LT(minShareRatio) {
			minShareRatio = shareRatio
		}
		if shareRatio.GT(maxShareRatio) {
			maxShareRatio = shareRatio
		}
		coinShareRatios[i] = shareRatio
	}

	if minShareRatio.Equal(sdk.MaxSortableDec) {
		return sdk.ZeroInt(), sdk.NewCoins(), errors.New("unexpected error in balancer maximalExactRatioJoin")
	}

	if minShareRatio.IsZero() {
		return sdk.ZeroInt(), tokensIn, nil
	}

	numShares = minShareRatio.MulInt(pool.TotalShares.Amount).TruncateInt()
	remCoins = sdk.Coins{}

	// if we have multiple shares, calculate remCoins
	if !minShareRatio.Equal(maxShareRatio) {
		// we have to calculate remCoins
		for i, coin := range tokensIn {
			if !coinShareRatios[i].Equal(minShareRatio) {
				usedAmount := minShareRatio.MulInt(
					poolLiquidity.AmountOfNoDenomValidation(coin.Denom)).Ceil().TruncateInt()
				remainingAmount := coin.Amount.Sub(usedAmount)
				// add to RemCoins
				if !remainingAmount.IsZero() {
					remCoins = remCoins.Add(sdk.Coin{Denom: coin.Denom, Amount: remainingAmount})
				}
			}
		}
	}

	return numShares, remCoins, nil
}

/*
Calculates the number of tokens to remove from liquidity given LP shares returned to the pool.

Note that this function is pure/read-only. It only calculates the theoretical amoount
and doesn't modify the actual state.

args:
  - numSharesIn: number of LP shares to return to the pool

ret:
  - tokensOut: the tokens withdrawn from the pool
  - err: error if any
*/
func (pool Pool) TokensOutFromPoolSharesIn(numSharesIn sdk.Int) (
	tokensOut sdk.Coins, err error,
) {
	if numSharesIn.IsZero() {
		return nil, errors.New("num shares in must be greater than zero")
	}

	shareRatio := numSharesIn.ToDec().QuoInt(pool.TotalShares.Amount)
	if shareRatio.IsZero() {
		return nil, errors.New("share ratio must be greater than zero")
	}
	if shareRatio.GT(sdk.OneDec()) {
		return nil, errors.New("share ratio cannot be greater than one")
	}

	poolLiquidity := pool.PoolBalances()
	tokensOut = make(sdk.Coins, len(poolLiquidity))
	for i, coin := range poolLiquidity {
		// tokenOut = shareRatio * poolTokenAmt * (1 - exitFee)
		tokenOutAmt := shareRatio.MulInt(coin.Amount).Mul(
			sdk.OneDec().Sub(pool.PoolParams.ExitFee),
		).TruncateInt()
		tokensOut[i] = sdk.NewCoin(coin.Denom, tokenOutAmt)
	}

	return tokensOut, nil
}

/*
Adds new liquidity to the pool and increments the total number of shares.

args:
  - numShares: the number of LP shares to increment
  - newLiquidity: the new tokens to deposit into the pool
*/
func (pool *Pool) incrementBalances(numShares sdk.Int, newLiquidity sdk.Coins) (
	err error,
) {
	for _, coin := range newLiquidity {
		i, poolAsset, err := pool.getPoolAssetAndIndex(coin.Denom)
		if err != nil {
			return err
		}
		poolAsset.Token.Amount = poolAsset.Token.Amount.Add(coin.Amount)
		pool.PoolAssets[i] = poolAsset
	}
	pool.TotalShares.Amount = pool.TotalShares.Amount.Add(numShares)
	return nil
}
