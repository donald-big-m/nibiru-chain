package keeper

import (
	"time"

	"github.com/NibiruChain/nibiru/collections/keys"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/NibiruChain/nibiru/x/common"
	"github.com/NibiruChain/nibiru/x/vpool/types"
)

/*
GetSpotPrice retrieves the price of the base asset denominated in quote asset.

The convention is the amount of quote assets required to buy one base asset.

e.g. If the tokenPair is BTC:NUSD, the method would return sdk.Dec(40,000.00)
because the instantaneous tangent slope on the vpool curve is 40,000.00,
so it would cost ~40,000.00 to buy one BTC:NUSD perp.

args:
  - ctx: cosmos-sdk context
  - pair: the token pair to get price for

ret:
  - price: the price of the token pair as sdk.Dec
  - err: error
*/
func (k Keeper) GetSpotPrice(ctx sdk.Context, pair common.AssetPair) (sdk.Dec, error) {
	pool, err := k.Pools.Get(ctx, pair)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	if pool.BaseAssetReserve.IsNil() || pool.BaseAssetReserve.IsZero() ||
		pool.QuoteAssetReserve.IsNil() || pool.QuoteAssetReserve.IsZero() {
		return sdk.ZeroDec(), nil
	}

	return pool.QuoteAssetReserve.Quo(pool.BaseAssetReserve), nil
}

/*
So how much stablecoin you would get if you sold baseAssetAmount amount of perpetual contracts.

Returns the amount of quote assets required to achieve a move of baseAssetAmount in a direction.
e.g. if removing <baseAssetAmount> base assets from the pool, returns the amount of quote assets do so.

args:
  - ctx: cosmos-sdk context
  - pair: the trading token pair
  - dir: add or remove
  - baseAssetAmount: the amount of base asset

ret:
  - quoteAmount: the amount of quote assets required to make the desired swap
  - err: error
*/
func (k Keeper) GetBaseAssetPrice(
	ctx sdk.Context,
	pair common.AssetPair,
	dir types.Direction,
	baseAssetAmount sdk.Dec,
) (quoteAmount sdk.Dec, err error) {
	pool, err := k.Pools.Get(ctx, pair)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return pool.GetQuoteAmountByBaseAmount(dir, baseAssetAmount)
}

/*
Returns the amount of base assets required to achieve a move of quoteAmount in a direction.
e.g. if removing <quoteAmount> quote assets from the pool, returns the amount of base assets do so.

args:
  - ctx: cosmos-sdk context
  - pair: the trading token pair
  - dir: add or remove
  - quoteAmount: the amount of quote asset

ret:
  - baseAssetAmount: the amount of base assets required to make the desired swap
  - err: error
*/
func (k Keeper) GetQuoteAssetPrice(
	ctx sdk.Context,
	pair common.AssetPair,
	dir types.Direction,
	quoteAmount sdk.Dec,
) (baseAssetAmount sdk.Dec, err error) {
	pool, err := k.Pools.Get(ctx, pair)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return pool.GetBaseAmountByQuoteAmount(dir, quoteAmount)
}

/*
Returns the twap of the spot price (y/x).

args:
  - ctx: cosmos-sdk context
  - pair: the token pair
  - direction: add or remove
  - baseAssetAmount: amount of base asset to add or remove
  - lookbackInterval: how far back to calculate TWAP

ret:
  - quoteAssetAmount: the amount of quote asset to make the desired move, as sdk.Dec
  - err: error
*/
func (k Keeper) GetSpotTWAP(
	ctx sdk.Context,
	pair common.AssetPair,
	lookbackInterval time.Duration,
) (quoteAssetAmount sdk.Dec, err error) {
	return k.calcTwap(
		ctx,
		pair,
		types.TwapCalcOption_SPOT,
		types.Direction_DIRECTION_UNSPECIFIED, // unused
		sdk.ZeroDec(),                         // unused
		lookbackInterval,
	)
}

/*
Returns the amount of quote assets required to achieve a move of baseAssetAmount in a direction,
based on historical snapshots.
e.g. if removing <baseAssetAmount> base assets from the pool, returns the amount of quote assets do so.

args:
  - ctx: cosmos-sdk context
  - pair: the token pair
  - direction: add or remove
  - baseAssetAmount: amount of base asset to add or remove
  - lookbackInterval: how far back to calculate TWAP

ret:
  - quoteAssetAmount: the amount of quote asset to make the desired move, as sdk.Dec
  - err: error
*/
func (k Keeper) GetBaseAssetTWAP(
	ctx sdk.Context,
	pair common.AssetPair,
	direction types.Direction,
	baseAssetAmount sdk.Dec,
	lookbackInterval time.Duration,
) (quoteAssetAmount sdk.Dec, err error) {
	return k.calcTwap(
		ctx,
		pair,
		types.TwapCalcOption_BASE_ASSET_SWAP,
		direction,
		baseAssetAmount,
		lookbackInterval,
	)
}

/*
Returns the amount of base assets required to achieve a move of quoteAssetAmount in a direction,
based on historical snapshots.
e.g. if removing <quoteAssetAmount> quote assets from the pool, returns the amount of base assets do so.

args:
  - ctx: cosmos-sdk context
  - pair: the token pair
  - direction: add or remove
  - quoteAssetAmount: amount of base asset to add or remove
  - lookbackInterval: how far back to calculate TWAP

ret:
  - baseAssetAmount: the amount of quote asset to make the desired move, as sdk.Dec
  - err: error
*/
func (k Keeper) GetQuoteAssetTWAP(
	ctx sdk.Context,
	pair common.AssetPair,
	direction types.Direction,
	quoteAssetAmount sdk.Dec,
	lookbackInterval time.Duration,
) (baseAssetAmount sdk.Dec, err error) {
	return k.calcTwap(
		ctx,
		pair,
		types.TwapCalcOption_QUOTE_ASSET_SWAP,
		direction,
		quoteAssetAmount,
		lookbackInterval,
	)
}

/*
Gets the time-weighted average price from [ ctx.BlockTime() - interval, ctx.BlockTime() )
Note the open-ended right bracket.

args:
  - ctx: cosmos-sdk context
  - pair: the token pair
  - twapCalcOption: one of SPOT, QUOTE_ASSET_SWAP, or BASE_ASSET_SWAP
  - direction: add or remove, only required for QUOTE_ASSET_SWAP or BASE_ASSET_SWAP
  - assetAmount: amount of asset to add or remove, only required for QUOTE_ASSET_SWAP or BASE_ASSET_SWAP
  - lookbackInterval: how far back to calculate TWAP

ret:
  - price: TWAP as sdk.Dec
  - err: error
*/
func (k Keeper) calcTwap(
	ctx sdk.Context,
	pair common.AssetPair,
	twapCalcOption types.TwapCalcOption,
	direction types.Direction,
	assetAmount sdk.Dec,
	lookbackInterval time.Duration,
) (price sdk.Dec, err error) {
	// earliest timestamp we'll look back until
	lowerLimitTimestampMs := ctx.BlockTime().Add(-1 * lookbackInterval).UnixMilli()

	prefix := keys.PairPrefix[common.AssetPair, keys.Uint64Key](pair)
	end := keys.PairSuffix[common.AssetPair, keys.Uint64Key](keys.Uint64(uint64(ctx.BlockTime().UnixMilli())))
	rng := keys.NewRange[keys.Pair[common.AssetPair, keys.Uint64Key]]().
		Prefix(prefix).
		End(keys.Inclusive(end)).
		Descending()
	iter := k.ReserveSnapshots.Iterate(ctx, rng)
	defer iter.Close()

	var snapshots []types.ReserveSnapshot
	for ; iter.Valid(); iter.Next() {
		s := iter.Value()
		snapshots = append(snapshots, s)
		if s.TimestampMs <= lowerLimitTimestampMs {
			break
		}
	}

	if len(snapshots) == 0 {
		return sdk.OneDec().Neg(), types.ErrNoValidTWAP
	}

	return calcTwap(ctx, snapshots, lowerLimitTimestampMs, twapCalcOption, direction, assetAmount)
}

// calcTwap walks through a slice of PriceSnapshots and tallies up the prices weighted by the amount of time they were active for.
// Callers of this function should already check if the snapshot slice is empty. Passing an empty snapshot slice will result in a panic.
func calcTwap(ctx sdk.Context, snapshots []types.ReserveSnapshot, lowerLimitTimestampMs int64, twapCalcOption types.TwapCalcOption, direction types.Direction, assetAmt sdk.Dec) (sdk.Dec, error) {
	// circuit-breaker when there's only one snapshot to process
	if len(snapshots) == 1 {
		return getPriceWithSnapshot(
			snapshots[0],
			snapshotPriceOptions{
				pair:           snapshots[0].Pair,
				twapCalcOption: twapCalcOption,
				direction:      direction,
				assetAmount:    assetAmt,
			},
		)
	}

	prevTimestampMs := ctx.BlockTime().UnixMilli()
	cumulativePrice := sdk.ZeroDec()

	for i, s := range snapshots {
		sPrice, err := getPriceWithSnapshot(
			s,
			snapshotPriceOptions{
				pair:           s.Pair,
				twapCalcOption: twapCalcOption,
				direction:      direction,
				assetAmount:    assetAmt,
			},
		)
		if err != nil {
			return sdk.Dec{}, err
		}
		var startTimestampMs int64
		if i == len(snapshots)-1 {
			// if we're at the oldest snapshot, then consider that price as starting from the lower limit timestamp
			startTimestampMs = lowerLimitTimestampMs
		} else {
			startTimestampMs = s.TimestampMs
		}
		cumulativePrice = cumulativePrice.Add(sPrice.MulInt64(prevTimestampMs - startTimestampMs))
		prevTimestampMs = s.TimestampMs
	}
	return cumulativePrice.QuoInt64(ctx.BlockTime().UnixMilli() - lowerLimitTimestampMs), nil
}
