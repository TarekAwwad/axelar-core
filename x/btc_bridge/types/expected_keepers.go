package types

import (
	"crypto/ecdsa"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tss "github.com/axelarnetwork/axelar-core/x/tss/exported"
	voting "github.com/axelarnetwork/axelar-core/x/vote/exported"
)

type Voter interface {
	voting.Voter
}

type Signer interface {
	GetSig(ctx sdk.Context, sigID string) (tss.Signature, bool)
	GetKey(ctx sdk.Context, keyID string) (ecdsa.PublicKey, bool)
	GetLatestMasterKey(ctx sdk.Context, chain string) (ecdsa.PublicKey, bool)
}