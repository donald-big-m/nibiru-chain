package common

import (
	"fmt"
	"sort"
	"strings"
)

var (
	GovDenom    = "unibi"
	CollDenom   = "uust"
	StableDenom = "unusd"

	TreasuryPoolModuleAccount = "treasury_pool"

	PairSeparator = ":"

	WhitelistedColl = []string{CollDenom}

	GovStablePool  = AssetPair{Token0: GovDenom, Token1: StableDenom}
	CollStablePool = AssetPair{Token0: CollDenom, Token1: StableDenom}

	ErrInvalidTokenPair = fmt.Errorf("invalid token pair")
)

type AssetPair struct {
	Token0 string
	Token1 string
}

// name is the name of the pool that corresponds to the two assets on this pair.
func (pair AssetPair) Name() string {
	return PoolNameFromDenoms([]string{pair.Token0, pair.Token1})
}

func (pair AssetPair) PairID() string {
	return pair.Name()
}

func (pair AssetPair) String() string {
	return fmt.Sprintf("%s:%s", pair.Token0, pair.Token1)
}

func (pair AssetPair) IsProperOrder() bool {
	return pair.Name() == pair.String()
}

func (pair AssetPair) Inverse() AssetPair {
	return AssetPair{pair.Token1, pair.Token0}
}

func (pair AssetPair) Proper() AssetPair {
	if pair.IsProperOrder() {
		return pair
	} else {
		return pair.Inverse()
	}
}

func DenomsFromPoolName(pool string) (denoms []string) {
	return strings.Split(pool, ":")
}

// PoolNameFromDenoms returns a sorted string representing a pool of assets
func PoolNameFromDenoms(denoms []string) string {
	sort.Strings(denoms) // alphabetically sort in-place
	return RawPoolNameFromDenoms(denoms)
}

// RawPoolNameFromDenoms returns a string representing a pool of assets in the
// exact order the denoms were given as args
func RawPoolNameFromDenoms(denoms []string) string {
	poolName := denoms[0]
	for idx, denom := range denoms {
		if idx != 0 {
			poolName += fmt.Sprintf(":%s", denom)
		}
	}
	return poolName
}

type Pair string

func NewPairFromStr(pair string) (Pair, error) {
	split := strings.Split(pair, PairSeparator)
	if len(split) != 2 {
		return "", ErrInvalidTokenPair
	}

	return Pair(pair), nil
}

func (p Pair) GetBaseToken() string {
	return strings.Split(string(p), ":")[0]
}

func (p Pair) GetQuoteToken() interface{} {
	return strings.Split(string(p), ":")[1]
}
