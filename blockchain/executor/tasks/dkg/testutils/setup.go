package testutils

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/MadBase/MadNet/blockchain/testutils"
	"math/big"

	ethereumInterfaces "github.com/MadBase/MadNet/blockchain/ethereum/interfaces"
	dkgObjects "github.com/MadBase/MadNet/blockchain/executor/tasks/dkg/objects"
	"github.com/MadBase/MadNet/blockchain/executor/tasks/dkg/utils"
	"github.com/MadBase/MadNet/blockchain/monitor/events"

	"github.com/MadBase/MadNet/consensus/objs"
	"github.com/MadBase/MadNet/constants"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/MadBase/MadNet/bridge/bindings"
	"github.com/MadBase/MadNet/crypto"
	"github.com/MadBase/MadNet/crypto/bn256"
	"github.com/MadBase/MadNet/crypto/bn256/cloudflare"
)

const SETUP_GROUP int = 13

type adminHandlerMock struct {
	snapshotCalled     bool
	privateKeyCalled   bool
	validatorSetCalled bool
	registerSnapshot   bool
	setSynchronized    bool
}

func (ah *adminHandlerMock) AddPrivateKey([]byte, constants.CurveSpec) error {
	ah.privateKeyCalled = true
	return nil
}

func (ah *adminHandlerMock) AddSnapshot(*objs.BlockHeader, bool) error {
	ah.snapshotCalled = true
	return nil
}

func (ah *adminHandlerMock) AddValidatorSet(*objs.ValidatorSet) error {
	ah.validatorSetCalled = true
	return nil
}

func (ah *adminHandlerMock) RegisterSnapshotCallback(func(*objs.BlockHeader) error) {
	ah.registerSnapshot = true
}

func (ah *adminHandlerMock) SetSynchronized(v bool) {
	ah.setSynchronized = true
}

func InitializeNewDetDkgStateInfo(n int) ([]*dkgObjects.DkgState, []*ecdsa.PrivateKey) {
	return InitializeNewDkgStateInfo(n, true)
}

func InitializeNewNonDetDkgStateInfo(n int) ([]*dkgObjects.DkgState, []*ecdsa.PrivateKey) {
	return InitializeNewDkgStateInfo(n, false)
}

func InitializeNewDkgStateInfo(n int, deterministicShares bool) ([]*dkgObjects.DkgState, []*ecdsa.PrivateKey) {
	// Get private keys for validators
	privKeys := testutils.SetupPrivateKeys(n)
	accountsArray := testutils.SetupAccounts(privKeys)
	dkgStates := []*dkgObjects.DkgState{}
	threshold := crypto.CalcThreshold(n)

	// Make base for secret key
	baseSecretBytes := make([]byte, 32)
	baseSecretBytes[0] = 101
	baseSecretBytes[31] = 101
	baseSecretValue := new(big.Int).SetBytes(baseSecretBytes)

	// Make base for transport key
	baseTransportBytes := make([]byte, 32)
	baseTransportBytes[0] = 1
	baseTransportBytes[1] = 1
	baseTransportValue := new(big.Int).SetBytes(baseTransportBytes)

	// Beginning dkgState initialization
	for k := 0; k < n; k++ {
		bigK := big.NewInt(int64(k))
		// Get base DkgState
		dkgState := dkgObjects.NewDkgState(accountsArray[k])
		// Set Index
		dkgState.Index = k + 1
		// Set Number of Validators
		dkgState.NumberOfValidators = n
		dkgState.ValidatorThreshold = threshold

		// Setup TransportKey
		transportPrivateKey := new(big.Int).Add(baseTransportValue, bigK)
		dkgState.TransportPrivateKey = transportPrivateKey
		transportPublicKeyG1 := new(cloudflare.G1).ScalarBaseMult(dkgState.TransportPrivateKey)
		transportPublicKey, err := bn256.G1ToBigIntArray(transportPublicKeyG1)
		if err != nil {
			panic(err)
		}
		dkgState.TransportPublicKey = transportPublicKey

		// Append to state array
		dkgStates = append(dkgStates, dkgState)
	}

	// Generate Participants
	for k := 0; k < n; k++ {
		participantList := GenerateParticipantList(dkgStates)
		for _, p := range participantList {
			dkgStates[k].Participants[p.Address] = p
		}
	}

	// Prepare secret shares
	for k := 0; k < n; k++ {
		bigK := big.NewInt(int64(k))
		// Set SecretValue and PrivateCoefficients
		dkgState := dkgStates[k]
		if deterministicShares {
			// Deterministic shares
			secretValue := new(big.Int).Add(baseSecretValue, bigK)
			privCoefs := GenerateDeterministicPrivateCoefficients(n)
			privCoefs[0].Set(secretValue) // Overwrite constant term
			dkgState.SecretValue = secretValue
			dkgState.PrivateCoefficients = privCoefs
		} else {
			// Random shares
			_, privCoefs, _, err := utils.GenerateShares(dkgState.TransportPrivateKey, dkgState.GetSortedParticipants())
			if err != nil {
				panic(err)
			}
			dkgState.SecretValue = new(big.Int)
			dkgState.SecretValue.Set(privCoefs[0])
			dkgState.PrivateCoefficients = privCoefs
		}
	}

	return dkgStates, privKeys
}

func GenerateParticipantList(dkgStates []*dkgObjects.DkgState) dkgObjects.ParticipantList {
	n := len(dkgStates)
	participants := make(dkgObjects.ParticipantList, int(n))
	for idx := 0; idx < n; idx++ {
		addr := dkgStates[idx].Account.Address
		publicKey := [2]*big.Int{}
		publicKey[0] = new(big.Int)
		publicKey[1] = new(big.Int)
		publicKey[0].Set(dkgStates[idx].TransportPublicKey[0])
		publicKey[1].Set(dkgStates[idx].TransportPublicKey[1])
		participant := &dkgObjects.Participant{}
		participant.Address = addr
		participant.PublicKey = publicKey
		participant.Index = dkgStates[idx].Index
		participants[idx] = participant
	}
	return participants
}

func GenerateEncryptedSharesAndCommitments(dkgStates []*dkgObjects.DkgState) {
	n := len(dkgStates)
	for k := 0; k < n; k++ {
		dkgState := dkgStates[k]
		publicCoefs := GeneratePublicCoefficients(dkgState.PrivateCoefficients)
		encryptedShares := GenerateEncryptedShares(dkgStates, k)
		// Loop through entire list and save in map
		for ell := 0; ell < n; ell++ {
			dkgStates[ell].Participants[dkgState.Account.Address].Commitments = publicCoefs
			dkgStates[ell].Participants[dkgState.Account.Address].EncryptedShares = encryptedShares
		}
	}
}

func GenerateDeterministicPrivateCoefficients(n int) []*big.Int {
	threshold := crypto.CalcThreshold(n)
	privCoefs := []*big.Int{}
	privCoefs = append(privCoefs, big.NewInt(0))
	for k := 1; k <= threshold; k++ {
		privCoef := big.NewInt(1)
		privCoefs = append(privCoefs, privCoef)
	}
	return privCoefs
}

func GeneratePublicCoefficients(privCoefs []*big.Int) [][2]*big.Int {
	publicCoefsG1 := cloudflare.GeneratePublicCoefs(privCoefs)
	publicCoefs := [][2]*big.Int{}
	for k := 0; k < len(publicCoefsG1); k++ {
		coefG1 := publicCoefsG1[k]
		coef, err := bn256.G1ToBigIntArray(coefG1)
		if err != nil {
			panic(err)
		}
		publicCoefs = append(publicCoefs, coef)
	}
	return publicCoefs
}

func GenerateEncryptedShares(dkgStates []*dkgObjects.DkgState, idx int) []*big.Int {
	dkgState := dkgStates[idx]
	// Get array of public keys and convert to cloudflare.G1
	publicKeysBig := [][2]*big.Int{}
	for k := 0; k < len(dkgStates); k++ {
		publicKeysBig = append(publicKeysBig, dkgStates[k].TransportPublicKey)
	}
	publicKeysG1, err := bn256.BigIntArraySliceToG1(publicKeysBig)
	if err != nil {
		panic(err)
	}

	// Get public key for caller
	publicKeyBig := dkgState.TransportPublicKey
	publicKey, err := bn256.BigIntArrayToG1(publicKeyBig)
	if err != nil {
		panic(err)
	}
	privCoefs := dkgState.PrivateCoefficients
	secretShares, err := cloudflare.GenerateSecretShares(publicKey, privCoefs, publicKeysG1)
	if err != nil {
		panic(err)
	}
	encryptedShares, err := cloudflare.GenerateEncryptedShares(secretShares, dkgState.TransportPrivateKey, publicKeysG1)
	if err != nil {
		panic(err)
	}
	return encryptedShares
}

func GenerateKeyShares(dkgStates []*dkgObjects.DkgState) {
	n := len(dkgStates)
	for k := 0; k < n; k++ {
		dkgState := dkgStates[k]
		g1KeyShare, g1Proof, g2KeyShare, err := utils.GenerateKeyShare(dkgState.SecretValue)
		if err != nil {
			panic(err)
		}
		addr := dkgState.Account.Address
		// Loop through entire list and save in map
		for ell := 0; ell < n; ell++ {
			dkgStates[ell].Participants[addr].KeyShareG1s = g1KeyShare
			dkgStates[ell].Participants[addr].KeyShareG1CorrectnessProofs = g1Proof
			dkgStates[ell].Participants[addr].KeyShareG2s = g2KeyShare
		}
	}
}

// GenerateMasterPublicKey computes the mpk for the protocol.
// This computes this by using all of the secret values from dkgStates.
func GenerateMasterPublicKey(dkgStates []*dkgObjects.DkgState) []*dkgObjects.DkgState {
	n := len(dkgStates)
	msk := new(big.Int)
	for k := 0; k < n; k++ {
		msk.Add(msk, dkgStates[k].SecretValue)
	}
	msk.Mod(msk, cloudflare.Order)
	for k := 0; k < n; k++ {
		mpkG2 := new(cloudflare.G2).ScalarBaseMult(msk)
		mpk, err := bn256.G2ToBigIntArray(mpkG2)
		if err != nil {
			panic(err)
		}
		dkgStates[k].MasterPublicKey = mpk
	}
	return dkgStates
}

func GenerateGPKJ(dkgStates []*dkgObjects.DkgState) {
	n := len(dkgStates)
	for k := 0; k < n; k++ {
		dkgState := dkgStates[k]

		encryptedShares := make([][]*big.Int, n)
		for idx, participant := range dkgState.GetSortedParticipants() {
			p, present := dkgState.Participants[participant.Address]
			if present && idx >= 0 && idx < n {
				encryptedShares[idx] = p.EncryptedShares
			} else {
				panic("Encrypted share state broken")
			}
		}

		groupPrivateKey, groupPublicKey, err := utils.GenerateGroupKeys(dkgState.TransportPrivateKey, dkgState.PrivateCoefficients,
			encryptedShares, dkgState.Index, dkgState.GetSortedParticipants())
		if err != nil {
			panic("Could not generate group keys")
		}

		dkgState.GroupPrivateKey = groupPrivateKey

		// Loop through entire list and save in map
		for ell := 0; ell < n; ell++ {
			dkgStates[ell].Participants[dkgState.Account.Address].GPKj = groupPublicKey
		}
	}
}

func GetETHDKGRegistrationOpened(logs []*types.Log, eth ethereumInterfaces.IEthereum) (*bindings.ETHDKGRegistrationOpened, error) {
	eventMap := events.GetETHDKGEvents()
	eventInfo, ok := eventMap["RegistrationOpened"]
	if !ok {
		return nil, fmt.Errorf("event not found: %v", eventInfo.Name)
	}

	var event *bindings.ETHDKGRegistrationOpened
	var err error
	for _, log := range logs {
		for _, topic := range log.Topics {
			if topic.String() == eventInfo.ID.String() {
				event, err = eth.Contracts().Ethdkg().ParseRegistrationOpened(*log)
				if err != nil {
					continue
				}

				return event, nil
			}
		}
	}
	return nil, fmt.Errorf("event not found")
}
