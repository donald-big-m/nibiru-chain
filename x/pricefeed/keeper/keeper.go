package keeper

import (
	"fmt"
	"sort"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/NibiruChain/nibiru/x/common"
	"github.com/NibiruChain/nibiru/x/pricefeed/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetPrice updates the posted price for a specific oracle
func (k Keeper) SetPrice(
	ctx sdk.Context,
	oracle sdk.AccAddress,
	pairStr string,
	price sdk.Dec,
	expiry time.Time,
) (postedPrice types.PostedPrice, err error) {
	// If the posted price expires before the current block, it is invalid.
	if expiry.Before(ctx.BlockTime()) {
		return postedPrice, types.ErrExpired
	}

	if !price.IsPositive() {
		return postedPrice, fmt.Errorf("price must be positive, not: %s", price)
	}

	pair, err := common.NewAssetPair(pairStr)
	if err != nil {
		return postedPrice, err
	}

	// Set inverse price if the oracle gives the wrong string
	if k.IsActivePair(ctx, pair.Inverse().AsString()) {
		pair = pair.Inverse()
		price = sdk.OneDec().Quo(price)
	}

	if !k.IsWhitelistedOracle(ctx, pair.AsString(), oracle) {
		return types.PostedPrice{}, fmt.Errorf(
			"oracle %s cannot post on pair %v", oracle, pair.AsString())
	}

	newPostedPrice := types.NewPostedPrice(pair.AsString(), oracle, price, expiry)

	// Emit an event containing the oracle's new price

	err = ctx.EventManager().EmitTypedEvent(&types.EventOracleUpdatePrice{
		PairId:    pair.AsString(),
		Oracle:    oracle.String(),
		PairPrice: price,
		Expiry:    expiry,
	})
	if err != nil {
		panic(err)
	}

	// Sets the raw price for a single oracle instead of an array of all oracle's raw prices
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RawPriceKey(pair.AsString(), oracle), k.cdc.MustMarshal(&newPostedPrice))
	return newPostedPrice, nil
}

// SetCurrentPrices updates the price of an asset to the median of all valid oracle inputs
func (k Keeper) SetCurrentPrices(ctx sdk.Context, token0 string, token1 string) error {
	assetPair := common.AssetPair{Token0: token0, Token1: token1}
	pairID := assetPair.AsString()

	if !k.IsActivePair(ctx, pairID) {
		return sdkerrors.Wrap(types.ErrInvalidPair, pairID)
	}
	// store current price
	validPrevPrice := true
	prevPrice, err := k.GetCurrentPrice(ctx, token0, token1)
	if err != nil {
		validPrevPrice = false
	}

	postedPrices := k.GetRawPrices(ctx, pairID)

	var notExpiredPrices []types.CurrentPrice
	// filter out expired prices
	for _, post := range postedPrices {
		if post.Expiry.After(ctx.BlockTime()) {
			notExpiredPrices = append(
				notExpiredPrices,
				types.NewCurrentPrice(token0, token1, post.Price))
		}
	}

	if len(notExpiredPrices) == 0 {
		// NOTE: The current price stored will continue storing the most recent (expired)
		// price if this is not set.
		// This zero's out the current price stored value for that market and ensures
		// that CDP methods that GetCurrentPrice will return error.
		k.setCurrentPrice(ctx, pairID, types.CurrentPrice{})
		return types.ErrNoValidPrice
	}

	medianPrice := k.CalculateMedianPrice(notExpiredPrices)

	// check case that market price was not set in genesis
	if validPrevPrice && !medianPrice.Equal(prevPrice.Price) {
		// only emit event if price has changed
		err = ctx.EventManager().EmitTypedEvent(&types.EventPairPriceUpdated{
			PairId:    pairID,
			PairPrice: medianPrice,
		})
		if err != nil {
			panic(err)
		}
	}

	currentPrice := types.NewCurrentPrice(token0, token1, medianPrice)
	k.setCurrentPrice(ctx, pairID, currentPrice)

	// Update the TWA prices
	err = k.updateTWAPPrice(ctx, pairID)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) setCurrentPrice(ctx sdk.Context, pairID string, currentPrice types.CurrentPrice) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CurrentPriceKey(pairID), k.cdc.MustMarshal(&currentPrice))
}

/* updateTWAPPrice updates the twap price for a token0, token1 pair
We use the blocktime to update the twap price.

Calculation is done as follow:
	$$P_{TWAP} = \frac { \sum P_j \times t_B }{ \sum{t_B} } $$
With
	P_j: current posted price for the pair of tokens
	t_B: current block timestamp

*/

func (k Keeper) updateTWAPPrice(ctx sdk.Context, pairID string) error {
	tokens := common.DenomsFromPoolName(pairID)
	token0, token1 := tokens[0], tokens[1]

	currentPrice, err := k.GetCurrentPrice(ctx, token0, token1)
	if err != nil {
		return err
	}

	currentTWAP, err := k.GetCurrentTWAPPrice(ctx, token0, token1)
	// Err there means no twap price have been set yet for this pair
	if err != nil {
		currentTWAP = types.CurrentTWAP{
			PairID:      pairID,
			Numerator:   sdk.MustNewDecFromStr("0"),
			Denominator: sdk.NewInt(0),
			Price:       sdk.MustNewDecFromStr("0"),
		}
	}

	blockUnixTime := sdk.NewInt(ctx.BlockTime().Unix())

	newDenominator := currentTWAP.Denominator.Add(blockUnixTime)
	newNumerator := currentTWAP.Numerator.Add(currentPrice.Price.Mul(sdk.NewDecFromInt(blockUnixTime)))

	newTWAP := types.CurrentTWAP{
		PairID:      pairID,
		Numerator:   newNumerator,
		Denominator: newDenominator,
		Price:       newNumerator.Quo(sdk.NewDecFromInt(newDenominator)),
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CurrentTWAPPriceKey("twap-"+pairID), k.cdc.MustMarshal(&newTWAP))

	return nil
}

// CalculateMedianPrice calculates the median prices for the input prices.
func (k Keeper) CalculateMedianPrice(prices []types.CurrentPrice) sdk.Dec {
	l := len(prices)

	if l == 1 {
		// Return immediately if there's only one price
		return prices[0].Price
	}
	// sort the prices
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Price.LT(prices[j].Price)
	})
	// for even numbers of prices, the median is calculated as the mean of the two middle prices
	if l%2 == 0 {
		median := k.calculateMeanPrice(prices[l/2-1], prices[l/2])
		return median
	}
	// for odd numbers of prices, return the middle element
	return prices[l/2].Price
}

func (k Keeper) calculateMeanPrice(priceA, priceB types.CurrentPrice) sdk.Dec {
	sum := priceA.Price.Add(priceB.Price)
	mean := sum.Quo(sdk.NewDec(2))
	return mean
}

// GetCurrentPrice fetches the current median price of all oracles for a specific market
func (k Keeper) GetCurrentPrice(ctx sdk.Context, token0 string, token1 string,
) (currPrice types.CurrentPrice, err error) {
	pair := common.AssetPair{Token0: token0, Token1: token1}
	givenIsActive := k.IsActivePair(ctx, pair.AsString())
	inverseIsActive := k.IsActivePair(ctx, pair.Inverse().AsString())
	if !givenIsActive && inverseIsActive {
		pair = pair.Inverse()
	}

	// Retrieve current price from the KV store
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CurrentPriceKey(pair.AsString()))
	if bz == nil {
		return types.CurrentPrice{}, types.ErrNoValidPrice
	}
	var price types.CurrentPrice
	k.cdc.MustUnmarshal(bz, &price)
	if price.Price.Equal(sdk.ZeroDec()) {
		return types.CurrentPrice{}, types.ErrNoValidPrice
	}

	if inverseIsActive {
		// Return the inverse price if the tokens are not in params order.
		inversePrice := sdk.OneDec().Quo(price.Price)
		return types.NewCurrentPrice(
			/* token0 */ token1,
			/* token1 */ token0,
			/* price */ inversePrice), nil
	}

	return price, nil
}

// GetCurrentTWAPPrice fetches the current median price of all oracles for a specific market
func (k Keeper) GetCurrentTWAPPrice(ctx sdk.Context, token0 string, token1 string) (currPrice types.CurrentTWAP, err error) {
	pair := common.AssetPair{Token0: token0, Token1: token1}
	givenIsActive := k.IsActivePair(ctx, pair.AsString())
	inverseIsActive := k.IsActivePair(ctx, pair.Inverse().AsString())
	if !givenIsActive && inverseIsActive {
		pair = pair.Inverse()
	}

	// Ensure we still have valid prices
	_, err = k.GetCurrentPrice(ctx, token0, token1)
	if err != nil {
		return types.CurrentTWAP{}, types.ErrNoValidPrice
	}

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CurrentTWAPPriceKey("twap-" + pair.AsString()))

	if bz == nil {
		return types.CurrentTWAP{}, types.ErrNoValidTWAP
	}

	var price types.CurrentTWAP
	k.cdc.MustUnmarshal(bz, &price)
	if price.Price.IsZero() {
		return types.CurrentTWAP{}, types.ErrNoValidPrice
	}

	if inverseIsActive {
		// Return the inverse price if the tokens are not in "proper" order.
		inversePrice := sdk.OneDec().Quo(price.Price)
		return types.NewCurrentTWAP(
			/* token0 */ token1,
			/* token1 */ token0,
			/* numerator */ price.Numerator,
			/* denominator */ price.Denominator,
			/* price */ inversePrice), nil
	}

	return price, nil
}

// IterateCurrentPrices iterates over all current price objects in the store and performs a callback function
func (k Keeper) IterateCurrentPrices(ctx sdk.Context, cb func(cp types.CurrentPrice) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.CurrentPricePrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var cp types.CurrentPrice
		k.cdc.MustUnmarshal(iterator.Value(), &cp)
		if cb(cp) {
			break
		}
	}
}

// GetCurrentPrices returns all current price objects from the store
func (k Keeper) GetCurrentPrices(ctx sdk.Context) types.CurrentPrices {
	var cps types.CurrentPrices
	k.IterateCurrentPrices(ctx, func(cp types.CurrentPrice) (stop bool) {
		cps = append(cps, cp)
		return false
	})
	return cps
}

// GetRawPrices fetches the set of all prices posted by oracles for an asset
func (k Keeper) GetRawPrices(ctx sdk.Context, pairStr string) types.PostedPrices {
	inversePair := common.MustNewAssetPair(pairStr).Inverse()
	if k.IsActivePair(ctx, inversePair.AsString()) {
		panic(fmt.Errorf(
			`cannot fetch posted prices using inverse pair, %v ;
			Use pair, %v, instead`, inversePair.AsString(), pairStr))
	}

	var pps types.PostedPrices
	k.IterateRawPricesByPair(ctx, pairStr, func(pp types.PostedPrice) (stop bool) {
		pps = append(pps, pp)
		return false
	})
	return pps
}

// IterateRawPrices iterates over all raw prices in the store and performs a callback function
func (k Keeper) IterateRawPricesByPair(ctx sdk.Context, marketId string, cb func(record types.PostedPrice) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.RawPriceIteratorKey((marketId)))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var record types.PostedPrice
		k.cdc.MustUnmarshal(iterator.Value(), &record)
		if cb(record) {
			break
		}
	}
}
