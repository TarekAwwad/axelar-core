package types

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/axelarnetwork/axelar-core/testutils/rand"
	nexus "github.com/axelarnetwork/axelar-core/x/nexus/exported"
	tss "github.com/axelarnetwork/axelar-core/x/tss/exported"
	tssTestUtils "github.com/axelarnetwork/axelar-core/x/tss/exported/testutils"
)

func TestCreateApproveContractCallWithMintCommand(t *testing.T) {
	chainID := sdk.NewInt(1)
	keyID := tssTestUtils.RandKeyID()
	sourceChain := nexus.ChainName("polygon")
	txID := Hash(common.HexToHash("0x5bb45dc24ddd6b90fa37f26eecfcf203328427c3226db29d1c01051b965ca93b"))
	index := uint64(99)
	sourceAddress := "0x68B93045fe7D8794a7cAF327e7f855CD6Cd03BB8"
	contractAddress := common.HexToAddress("0x956dA338C1518a7FB213042b70c60c021aeBd554")
	payloadHash := common.HexToHash("0x7c6498469c4e2d466b6fc9af3c910587f6c0bdade714a16ab279a08a759a5c14")
	symbol := "testA"
	amount := sdk.NewUint(20000)
	event := EventContractCallWithToken{
		Sender:          Address(common.HexToAddress(sourceAddress)),
		ContractAddress: contractAddress.Hex(),
		PayloadHash:     Hash(payloadHash),
		Symbol:          rand.NormalizedStrBetween(1, 5),
	}

	expectedParams := "00000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000140000000000000000000000000956da338c1518a7fb213042b70c60c021aebd5547c6498469c4e2d466b6fc9af3c910587f6c0bdade714a16ab279a08a759a5c1400000000000000000000000000000000000000000000000000000000000001a00000000000000000000000000000000000000000000000000000000000004e205bb45dc24ddd6b90fa37f26eecfcf203328427c3226db29d1c01051b965ca93b00000000000000000000000000000000000000000000000000000000000000630000000000000000000000000000000000000000000000000000000000000007706f6c79676f6e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002a3078363842393330343566653744383739346137634146333237653766383535434436436430334242380000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000057465737441000000000000000000000000000000000000000000000000000000"
	actual, err := CreateApproveContractCallWithMintCommand(
		chainID,
		keyID,
		sourceChain,
		txID,
		index,
		event,
		amount,
		symbol,
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedParams, hex.EncodeToString(actual.Params))

	actualSourceChain, actualSourceAddress, actualContractAddress, actualPayloadHash, actualSymbol, actualAmount, actualSourceTxID, actualSourceEventIndex, err := decodeApproveContractCallWithMintParams(actual.Params)
	assert.NoError(t, err)
	assert.Equal(t, sourceChain.String(), actualSourceChain)
	assert.Equal(t, sourceAddress, actualSourceAddress)
	assert.Equal(t, contractAddress, actualContractAddress)
	assert.Equal(t, payloadHash, actualPayloadHash)
	assert.Equal(t, symbol, actualSymbol)
	assert.Equal(t, amount.BigInt(), actualAmount)
	assert.Equal(t, txID, Hash(actualSourceTxID))
	assert.Equal(t, index, actualSourceEventIndex.Uint64())
}

func TestNewCommandBatchMetadata(t *testing.T) {
	chainID := sdk.NewInt(1)
	commands := []Command{
		{
			ID:      CommandID(common.HexToHash("0xc5baf525fe191e3e9e35c2012ff5f86954c04677a1e4df56079714fc4949409f")),
			Command: AxelarGatewayCommandDeployToken,
			Params:  common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000271000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010416e20417765736f6d6520546f6b656e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034141540000000000000000000000000000000000000000000000000000000000"),
		},
	}

	expectedData := common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000001c5baf525fe191e3e9e35c2012ff5f86954c04677a1e4df56079714fc4949409f00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000b6465706c6f79546f6b656e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000271000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010416e20417765736f6d6520546f6b656e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000034141540000000000000000000000000000000000000000000000000000000000")
	actual, err := NewCommandBatchMetadata(
		rand.PosI64(),
		chainID,
		tssTestUtils.RandKeyID(),
		tss.MasterKey,
		commands,
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, actual.Data)
}

func TestDeployToken(t *testing.T) {
	chainID := sdk.NewInt(1)
	keyID := tssTestUtils.RandKeyID()

	details := TokenDetails{
		TokenName: rand.Str(10),
		Symbol:    rand.Str(3),
		Decimals:  uint8(rand.I64Between(3, 10)),
		Capacity:  sdk.NewIntFromBigInt(big.NewInt(rand.I64Between(100, 100000))),
	}
	address := Address(common.BytesToAddress(rand.Bytes(common.AddressLength)))
	asset := rand.Str(5)

	capBz := make([]byte, 8)
	binary.BigEndian.PutUint64(capBz, details.Capacity.Uint64())
	capHex := hex.EncodeToString(capBz)

	expectedParams := fmt.Sprintf("00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000%s%s000000000000000000000000%s000000000000000000000000000000000000000000000000000000000000000a%s000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003%s0000000000000000000000000000000000000000000000000000000000",
		hex.EncodeToString([]byte{byte(details.Decimals)}),
		strings.Repeat("0", 64-len(capHex))+capHex,
		hex.EncodeToString(address.Bytes()),
		hex.EncodeToString([]byte(details.TokenName)),
		hex.EncodeToString([]byte(details.Symbol)),
	)
	expectedCommandID := NewCommandID([]byte(asset+"_"+details.Symbol), chainID)
	actual, err := CreateDeployTokenCommand(chainID, keyID, asset, details, address)

	assert.NoError(t, err)
	assert.Equal(t, expectedParams, hex.EncodeToString(actual.Params))
	assert.Equal(t, expectedCommandID, actual.ID)

	decodedName, decodedSymbol, decodedDecs, decodedCap, tokenAddress, err := decodeDeployTokenParams(actual.Params)
	assert.NoError(t, err)
	assert.Equal(t, details.TokenName, decodedName)
	assert.Equal(t, details.Symbol, decodedSymbol)
	assert.Equal(t, details.Decimals, decodedDecs)
	assert.Equal(t, details.Capacity.BigInt(), decodedCap)
	assert.Equal(t, address, Address(tokenAddress))
}

func TestCreateMintTokenCommand(t *testing.T) {
	chainID := sdk.NewInt(1)
	keyID := tssTestUtils.RandKeyID()
	commandID := NewCommandID(rand.Bytes(32), chainID)
	symbol := rand.Str(3)
	address := common.BytesToAddress(rand.Bytes(common.AddressLength))
	amount := big.NewInt(rand.I64Between(100, 100000))

	amountBz := make([]byte, 8)
	binary.BigEndian.PutUint64(amountBz, amount.Uint64())
	amountHex := hex.EncodeToString(amountBz)

	expectedParams := fmt.Sprintf("0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000%s%s0000000000000000000000000000000000000000000000000000000000000003%s0000000000000000000000000000000000000000000000000000000000",
		hex.EncodeToString(address.Bytes()),
		strings.Repeat("0", 64-len(amountHex))+amountHex,
		hex.EncodeToString([]byte(symbol)),
	)
	actual, err := CreateMintTokenCommand(keyID, commandID, symbol, address, amount)

	assert.NoError(t, err)
	assert.Equal(t, expectedParams, hex.EncodeToString(actual.Params))

	decodedSymbol, decodedAddr, decodedAmount, err := decodeMintTokenParams(actual.Params)
	assert.NoError(t, err)
	assert.Equal(t, symbol, decodedSymbol)
	assert.Equal(t, address, decodedAddr)
	assert.Equal(t, amount, decodedAmount)

}

func TestCreateBurnTokenCommand(t *testing.T) {
	chainID := sdk.NewInt(1)
	keyID := tssTestUtils.RandKeyID()
	symbol := rand.Str(3)
	salt := common.BytesToHash(rand.Bytes(common.HashLength))
	height := rand.I64Between(100, 10000)

	expectedParams := fmt.Sprintf("0000000000000000000000000000000000000000000000000000000000000040%s0000000000000000000000000000000000000000000000000000000000000003%s0000000000000000000000000000000000000000000000000000000000",
		hex.EncodeToString(salt.Bytes()),
		hex.EncodeToString([]byte(symbol)),
	)
	actual, err := CreateBurnTokenCommand(
		chainID,
		keyID,
		height,
		BurnerInfo{Symbol: symbol, Salt: Hash(salt)},
		false,
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedParams, hex.EncodeToString(actual.Params))

	decodedSymbol, decodedSalt, err := decodeBurnTokenParams(actual.Params)
	assert.NoError(t, err)
	assert.Equal(t, symbol, decodedSymbol)
	assert.Equal(t, salt, decodedSalt)
}

func TestCreateSinglesigTransferCommand_Ownership(t *testing.T) {
	chainID := sdk.NewInt(1)
	keyID := tssTestUtils.RandKeyID()
	newOwnerAddr := common.BytesToAddress(rand.Bytes(common.AddressLength))

	expectedParams := fmt.Sprintf("000000000000000000000000%s", hex.EncodeToString(newOwnerAddr.Bytes()))
	actual, err := CreateSinglesigTransferCommand(
		Ownership,
		chainID,
		keyID,
		newOwnerAddr,
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedParams, hex.EncodeToString(actual.Params))

	decodedAddr, err := decodeTransferSinglesigParams(actual.Params)
	assert.NoError(t, err)
	assert.Equal(t, newOwnerAddr, decodedAddr)

	_, _, err = decodeTransferMultisigParams(actual.Params)
	assert.Error(t, err)
}

func TestCreateSinglesigTransferCommand_Operatorship(t *testing.T) {
	chainID := sdk.NewInt(1)
	keyID := tssTestUtils.RandKeyID()
	newOperatorAddr := common.BytesToAddress(rand.Bytes(common.AddressLength))

	expectedParams := fmt.Sprintf("000000000000000000000000%s", hex.EncodeToString(newOperatorAddr.Bytes()))
	actual, err := CreateSinglesigTransferCommand(
		Operatorship,
		chainID,
		keyID,
		newOperatorAddr,
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedParams, hex.EncodeToString(actual.Params))

	decodedAddr, err := decodeTransferSinglesigParams(actual.Params)
	assert.NoError(t, err)
	assert.Equal(t, newOperatorAddr, decodedAddr)

	_, _, err = decodeTransferMultisigParams(actual.Params)
	assert.Error(t, err)
}

func TestCreateMultisigTransferCommand_Ownership(t *testing.T) {
	chainID := sdk.NewInt(1)
	keyID := tssTestUtils.RandKeyID()

	addresses := []common.Address{
		common.HexToAddress("0xd59ca627Af68D29C547B91066297a7c469a7bF72"),
		common.HexToAddress("0xc2FCc7Bcf743153C58Efd44E6E723E9819E9A10A"),
		common.HexToAddress("0x2ad611e02E4F7063F515C8f190E5728719937205"),
	}
	threshold := uint8(2)

	expectedParams := "000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003000000000000000000000000d59ca627af68d29c547b91066297a7c469a7bf72000000000000000000000000c2fcc7bcf743153c58efd44e6e723e9819e9a10a0000000000000000000000002ad611e02e4f7063f515c8f190e5728719937205"
	actual, err := CreateMultisigTransferCommand(
		Ownership,
		chainID,
		keyID,
		threshold,
		addresses...,
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedParams, hex.EncodeToString(actual.Params))

	decodedAddrs, decodedThreshold, err := decodeTransferMultisigParams(actual.Params)
	assert.NoError(t, err)
	assert.ElementsMatch(t, addresses, decodedAddrs)
	assert.Equal(t, threshold, decodedThreshold)

	_, err = decodeTransferSinglesigParams(actual.Params)
	assert.Error(t, err)
}
func TestCreateMultisigTransferCommand_Operatorship(t *testing.T) {
	chainID := sdk.NewInt(1)
	keyID := tssTestUtils.RandKeyID()

	addresses := []common.Address{
		common.HexToAddress("0xd59ca627Af68D29C547B91066297a7c469a7bF72"),
		common.HexToAddress("0xc2FCc7Bcf743153C58Efd44E6E723E9819E9A10A"),
		common.HexToAddress("0x2ad611e02E4F7063F515C8f190E5728719937205"),
	}
	threshold := uint8(2)

	expectedParams := "000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003000000000000000000000000d59ca627af68d29c547b91066297a7c469a7bf72000000000000000000000000c2fcc7bcf743153c58efd44e6e723e9819e9a10a0000000000000000000000002ad611e02e4f7063f515c8f190e5728719937205"
	actual, err := CreateMultisigTransferCommand(
		Ownership,
		chainID,
		keyID,
		threshold,
		addresses...,
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedParams, hex.EncodeToString(actual.Params))

	decodedAddrs, decodedThreshold, err := decodeTransferMultisigParams(actual.Params)
	assert.NoError(t, err)
	assert.ElementsMatch(t, addresses, decodedAddrs)
	assert.Equal(t, threshold, decodedThreshold)

	_, err = decodeTransferSinglesigParams(actual.Params)
	assert.Error(t, err)
}

func TestGetSignHash(t *testing.T) {
	data := common.FromHex("0000000000000000000000000000000000000000000000000000000000000001ec78d9c22c08bb9f0ecd5d95571ae83e3f22219c5a9278c3270691d50abfd91b000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000096d696e74546f6b656e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000014141540000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000063fc2ad3d021a4d7e64323529a55a9442c444da00000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000270f")

	expected := "0xe7bce8f57491e71212d930096bacf9288c711e5f27200946edd570e3a93546bf"
	actual := GetSignHash(data)

	assert.Equal(t, expected, actual.Hex())
}

func TestCreateExecuteDataSinglesig(t *testing.T) {
	commandData := common.FromHex("0000000000000000000000000000000000000000000000000000000000000001ec78d9c22c08bb9f0ecd5d95571ae83e3f22219c5a9278c3270691d50abfd91b000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000096d696e74546f6b656e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000014141540000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000063fc2ad3d021a4d7e64323529a55a9442c444da00000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000270f")
	commandSig := Signature{}
	copy(commandSig[:], common.FromHex("42b936b3c37fb7deed86f52154798d0c9abfe5ba838b2488f4a7e5193a9bb60b5d8c521f5c8c64f9442fc745ecd3bc496b04dc03a81b4e89c72342ab5903284d1c"))

	expected := "09c5eabe000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002e00000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000026000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000001ec78d9c22c08bb9f0ecd5d95571ae83e3f22219c5a9278c3270691d50abfd91b000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000096d696e74546f6b656e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000014141540000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000063fc2ad3d021a4d7e64323529a55a9442c444da00000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000270f000000000000000000000000000000000000000000000000000000000000004142b936b3c37fb7deed86f52154798d0c9abfe5ba838b2488f4a7e5193a9bb60b5d8c521f5c8c64f9442fc745ecd3bc496b04dc03a81b4e89c72342ab5903284d1c00000000000000000000000000000000000000000000000000000000000000"
	actual, err := CreateExecuteDataSinglesig(commandData, commandSig)

	assert.NoError(t, err)
	assert.Equal(t, expected, common.Bytes2Hex(actual))
}

func TestCreateExecuteDataMultisig(t *testing.T) {
	commandData := common.FromHex("0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000000186c71b9698cc55f8238266b026414ed9880bcd3dafd254cfc1079f1d4c2098800000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000096d696e74546f6b656e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000006000000000000000000000000017ec8597ff92c3f44523bdc65bf0f1be632917ff000000000000000000000000000000000000000000000000000000000152a1c000000000000000000000000000000000000000000000000000000000000000034141540000000000000000000000000000000000000000000000000000000000")
	commandSigs := make([]Signature, 2)
	copy(commandSigs[0][:], common.FromHex("226f548e306ba150c2895f192c71de4e455655508bb0762d6808756ac5cae9dd41145781fa6f7bcd52c3a71d492b3bf15d8792c431568e1b379b8d52a479b0971c"))
	copy(commandSigs[1][:], common.FromHex("44e9e6a66df68d798802914c41f57c0ef488e0ca5f244afa60e3438a5078356803213e2de2f4d41a4002fb3115722b5804ff8cd0a5101d7b37ba97fadd223fc51b"))

	expected := "09c5eabe00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000002a000000000000000000000000000000000000000000000000000000000000002400000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000000186c71b9698cc55f8238266b026414ed9880bcd3dafd254cfc1079f1d4c2098800000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000096d696e74546f6b656e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000006000000000000000000000000017ec8597ff92c3f44523bdc65bf0f1be632917ff000000000000000000000000000000000000000000000000000000000152a1c0000000000000000000000000000000000000000000000000000000000000000341415400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000041226f548e306ba150c2895f192c71de4e455655508bb0762d6808756ac5cae9dd41145781fa6f7bcd52c3a71d492b3bf15d8792c431568e1b379b8d52a479b0971c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004144e9e6a66df68d798802914c41f57c0ef488e0ca5f244afa60e3438a5078356803213e2de2f4d41a4002fb3115722b5804ff8cd0a5101d7b37ba97fadd223fc51b00000000000000000000000000000000000000000000000000000000000000"
	actual, err := CreateExecuteDataMultisig(commandData, commandSigs...)

	assert.NoError(t, err)
	assert.Equal(t, expected, common.Bytes2Hex(actual))
}