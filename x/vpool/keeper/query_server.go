package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NibiruChain/nibiru/x/common"
	"github.com/NibiruChain/nibiru/x/vpool/types"
)

type queryServer struct {
	Keeper
}

func NewQuerier(k Keeper) queryServer {
	return queryServer{Keeper: k}
}

var _ types.QueryServer = queryServer{}

func (q queryServer) ReserveAssets(
	goCtx context.Context,
	req *types.QueryReserveAssetsRequests,
) (resp *types.QueryReserveAssetsResponse, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	tokenPair, err := common.NewAssetPairFromStr(req.Pair)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pool, err := q.getPool(ctx, tokenPair)
	if err != nil {
		return nil, err
	}

	return &types.QueryReserveAssetsResponse{
		BaseAssetReserve:  pool.BaseAssetReserve,
		QuoteAssetReserve: pool.QuoteAssetReserve,
	}, nil
}

func (q queryServer) AllPools(
	goCtx context.Context,
	req *types.QueryAllPoolsRequests,
) (resp *types.QueryAllPoolsResponse, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pools := q.GetAllPools(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryAllPoolsResponse{
		Pools: pools,
	}, nil
}

func (q queryServer) BaseAssetPrice(
	goCtx context.Context,
	req *types.QueryBaseAssetRequest,
) (resp *types.QueryBaseAssetResponse, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pair, err := common.NewAssetPairFromStr(req.Pair)
	if err != nil {
		return nil, err
	}

	price, err := q.GetBaseAssetPrice(
		ctx,
		pair,
		types.Direction(req.Direction),
		req.BaseAssetAmount,
	)
	if err != nil {
		return nil, err
	}

	// TODO: add underylying price
	// currently it gives error in cli_test
	// rpc error: code = InvalidArgument desc = All input prices are expired: invalid request
	// underlyingPrice, err := q.GetUnderlyingPrice(
	// 	ctx,
	// 	pair,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	return &types.QueryBaseAssetResponse{
		Price: price,
		// UnderlyingPrice: underlyingPrice,
	}, nil
}
