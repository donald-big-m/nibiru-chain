package binding_test

import (
	"testing"
	"time"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/NibiruChain/nibiru/app"
	"github.com/NibiruChain/nibiru/x/common/denoms"
	"github.com/NibiruChain/nibiru/x/common/testutil"
	"github.com/NibiruChain/nibiru/x/common/testutil/testapp"
	"github.com/NibiruChain/nibiru/x/wasm/binding/wasmbin"
)

// StoreContract submits Wasm bytecode for storage on the chain.
func StoreContract(
	t *testing.T,
	contractWasm wasmbin.WasmKey,
	ctx sdk.Context,
	nibiru *app.NibiruApp,
	sender sdk.AccAddress,
) (codeId uint64) {
	pathToWasmBin := wasmbin.GetPackageDir(t) + "/wasmbin"
	wasmBytecode, err := contractWasm.ToByteCode(pathToWasmBin)
	require.NoError(t, err)

	// The "Create" fn is private on the nibiru.WasmKeeper. By placing it as the
	// decorated keeper in PermissionedKeeper type, we can access "Create" as a
	// public fn.
	wasmPermissionedKeeper := wasmkeeper.NewDefaultPermissionKeeper(nibiru.WasmKeeper)
	instantiateAccess := &wasmtypes.AccessConfig{
		Permission: wasmtypes.AccessTypeEverybody,
	}
	codeId, _, err = wasmPermissionedKeeper.Create(
		ctx, sender, wasmBytecode, instantiateAccess,
	)
	require.NoError(t, err)
	return codeId
}

func InstantiateContract(
	t *testing.T, ctx sdk.Context, nibiru *app.NibiruApp, codeId uint64,
	initMsg []byte, sender sdk.AccAddress, label string, deposit sdk.Coins,
) (contractAddr sdk.AccAddress) {
	wasmPermissionedKeeper := wasmkeeper.NewDefaultPermissionKeeper(nibiru.WasmKeeper)
	contractAddr, _, err := wasmPermissionedKeeper.Instantiate(
		ctx, codeId, sender, sender, initMsg, label, deposit,
	)
	require.NoError(t, err)
	return contractAddr
}

func InstantiatePerpBindingContract(
	t *testing.T, ctx sdk.Context, nibiru *app.NibiruApp, codeId uint64,
	sender sdk.AccAddress, deposit sdk.Coins,
) (contractAddr sdk.AccAddress) {
	initMsg := []byte("{}")
	label := "x/perp module bindings"
	return InstantiateContract(
		t, ctx, nibiru, codeId, initMsg, sender, label, deposit,
	)
}

// ContractMap is a map from WasmKey to contract address
type ContractMapType = map[wasmbin.WasmKey]sdk.AccAddress

var ContractMap = make(map[wasmbin.WasmKey]sdk.AccAddress)

// SetupAllContracts stores and instantiates all of wasm binding contracts.
func SetupAllContracts(
	t *testing.T, sender sdk.AccAddress, nibiru *app.NibiruApp, ctx sdk.Context,
) (*app.NibiruApp, sdk.Context) {
	codeId := StoreContract(t, wasmbin.WasmKeyPerpBinding, ctx, nibiru, sender)

	deposit := sdk.NewCoins(sdk.NewCoin(denoms.NIBI, sdk.NewInt(1)))
	contract := InstantiatePerpBindingContract(t, ctx, nibiru, codeId, sender, deposit)

	ContractMap[wasmbin.WasmKeyPerpBinding] = contract

	return nibiru, ctx
}

func TestSetupContracts(t *testing.T) {
	sender := testutil.AccAddress()
	nibiru, _ := testapp.NewNibiruTestAppAndContext(true)
	ctx := nibiru.NewContext(false, tmproto.Header{
		Height:  1,
		ChainID: "nibiru-wasmnet-1",
		Time:    time.Now().UTC(),
	})
	coins := sdk.NewCoins(sdk.NewCoin(denoms.NIBI, sdk.NewInt(10)))
	require.NoError(t, testapp.FundAccount(nibiru.BankKeeper, ctx, sender, coins))
	_, _ = SetupAllContracts(t, sender, nibiru, ctx)
}
