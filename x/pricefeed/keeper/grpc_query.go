package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NibiruChain/nibiru/x/common"
	"github.com/NibiruChain/nibiru/x/pricefeed/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Price(goCtx context.Context, req *types.QueryPriceRequest) (*types.QueryPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pair, err := common.NewAssetPairFromStr(req.PairId)
	if err != nil {
		return nil, err
	}
	if !k.GetPairs(ctx).Contains(pair) {
		return nil, status.Error(codes.NotFound, "pair not in module params")
	}
	if !k.ActivePairsStore().getKV(ctx).Has([]byte(pair.AsString())) {
		return nil, status.Error(codes.NotFound, "invalid market ID")
	}

	tokens := common.DenomsFromPoolName(req.PairId)
	token0, token1 := tokens[0], tokens[1]
	currentPrice, sdkErr := k.GetCurrentPrice(ctx, token0, token1)
	if sdkErr != nil {
		return nil, sdkErr
	}

	return &types.QueryPriceResponse{
		Price: types.CurrentPriceResponse{PairID: req.PairId, Price: currentPrice.Price},
	}, nil
}

func (k Keeper) RawPrices(
	goCtx context.Context, req *types.QueryRawPricesRequest,
) (*types.QueryRawPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.IsActivePair(ctx, req.PairId) {
		return nil, status.Error(codes.NotFound, "invalid market ID")
	}

	var prices types.PostedPriceResponses
	for _, rp := range k.GetRawPrices(ctx, req.PairId) {
		prices = append(prices, types.PostedPriceResponse{
			PairID:        rp.PairID,
			OracleAddress: rp.OracleAddress.String(),
			Price:         rp.Price,
			Expiry:        rp.Expiry,
		})
	}

	return &types.QueryRawPricesResponse{
		RawPrices: prices,
	}, nil
}

func (k Keeper) Prices(goCtx context.Context, req *types.QueryPricesRequest) (*types.QueryPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var currentPrices types.CurrentPriceResponses
	for _, currentPrice := range k.GetCurrentPrices(ctx) {
		if currentPrice.PairID != "" {
			currentPrices = append(currentPrices, types.CurrentPriceResponse(currentPrice))
		}
	}

	return &types.QueryPricesResponse{
		Prices: currentPrices,
	}, nil
}

func (k Keeper) Oracles(goCtx context.Context, req *types.QueryOraclesRequest) (*types.QueryOraclesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := common.NewAssetPairFromStr(req.PairId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invalid market ID")
	}

	oracles := k.GetOraclesForPair(ctx, req.PairId)
	if len(oracles) == 0 {
		return &types.QueryOraclesResponse{}, nil
	}

	var strOracles []string
	for _, oracle := range oracles {
		strOracles = append(strOracles, oracle.String())
	}

	return &types.QueryOraclesResponse{
		Oracles: strOracles,
	}, nil
}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) Pairs(goCtx context.Context, req *types.QueryPairsRequest,
) (*types.QueryPairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var responses types.PairResponses
	for _, pairStr := range k.GetParams(ctx).Pairs {
		pair := common.MustNewAssetPairFromStr(pairStr)

		var oracleStrings []string
		for _, oracle := range k.OraclesStore().Get(ctx, pair) {
			oracleStrings = append(oracleStrings, oracle.String())
		}

		responses = append(responses, types.PairResponse{
			PairID:  pairStr,
			Token0:  pair.Token0,
			Token1:  pair.Token1,
			Oracles: oracleStrings,
			Active:  k.IsActivePair(ctx, pairStr),
		})
	}

	// TODO improve these variable names. PairResponse is confusing on field "Pairs"
	return &types.QueryPairsResponse{
		Pairs: responses,
	}, nil
}
