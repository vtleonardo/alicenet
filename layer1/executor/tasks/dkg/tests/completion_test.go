//go:build integration

package tests

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/MadBase/MadNet/crypto"
	bn256 "github.com/MadBase/MadNet/crypto/bn256/cloudflare"
	"github.com/alicenet/alicenet/consensus/objs"
	"github.com/alicenet/alicenet/constants"
	"github.com/alicenet/alicenet/layer1/executor/tasks/dkg"
	"github.com/alicenet/alicenet/layer1/executor/tasks/dkg/state"
	"github.com/alicenet/alicenet/layer1/monitor/objects"
	"github.com/alicenet/alicenet/layer1/tests"
	"github.com/alicenet/alicenet/layer1/transaction"
	"github.com/alicenet/alicenet/logging"
	"github.com/alicenet/alicenet/test/mocks"
	"github.com/alicenet/alicenet/utils"
	"github.com/stretchr/testify/assert"
)

func makeSigners(t *testing.T) ([]byte, []*crypto.BNGroupSigner, [][]byte, []*crypto.Secp256k1Signer, [][]byte) {
	s := new(crypto.BNGroupSigner)
	msg := []byte("A message to sign")

	secret1 := big.NewInt(100)
	secret2 := big.NewInt(101)
	secret3 := big.NewInt(102)
	secret4 := big.NewInt(103)

	msk := big.NewInt(0)
	msk.Add(msk, secret1)
	msk.Add(msk, secret2)
	msk.Add(msk, secret3)
	msk.Add(msk, secret4)
	msk.Mod(msk, bn256.Order)
	mpk := new(bn256.G2).ScalarBaseMult(msk)

	big1 := big.NewInt(1)
	big2 := big.NewInt(2)

	privCoefs1 := []*big.Int{secret1, big1, big2}
	privCoefs2 := []*big.Int{secret2, big1, big2}
	privCoefs3 := []*big.Int{secret3, big1, big2}
	privCoefs4 := []*big.Int{secret4, big1, big2}

	share1to1 := bn256.PrivatePolyEval(privCoefs1, 1)
	share1to2 := bn256.PrivatePolyEval(privCoefs1, 2)
	share1to3 := bn256.PrivatePolyEval(privCoefs1, 3)
	share1to4 := bn256.PrivatePolyEval(privCoefs1, 4)
	share2to1 := bn256.PrivatePolyEval(privCoefs2, 1)
	share2to2 := bn256.PrivatePolyEval(privCoefs2, 2)
	share2to3 := bn256.PrivatePolyEval(privCoefs2, 3)
	share2to4 := bn256.PrivatePolyEval(privCoefs2, 4)
	share3to1 := bn256.PrivatePolyEval(privCoefs3, 1)
	share3to2 := bn256.PrivatePolyEval(privCoefs3, 2)
	share3to3 := bn256.PrivatePolyEval(privCoefs3, 3)
	share3to4 := bn256.PrivatePolyEval(privCoefs3, 4)
	share4to1 := bn256.PrivatePolyEval(privCoefs4, 1)
	share4to2 := bn256.PrivatePolyEval(privCoefs4, 2)
	share4to3 := bn256.PrivatePolyEval(privCoefs4, 3)
	share4to4 := bn256.PrivatePolyEval(privCoefs4, 4)

	groupShares := make([][]byte, 4)
	for k := 0; k < len(groupShares); k++ {
		groupShares[k] = make([]byte, len(mpk.Marshal()))
	}

	listOfSS1 := []*big.Int{share1to1, share2to1, share3to1, share4to1}
	gsk1 := bn256.GenerateGroupSecretKeyPortion(listOfSS1)
	gpk1 := new(bn256.G2).ScalarBaseMult(gsk1)
	groupShares[0] = gpk1.Marshal()
	s1 := new(crypto.BNGroupSigner)
	err := s1.SetPrivk(gsk1.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	sig1, err := s1.Sign(msg)
	if err != nil {
		t.Fatal(err)
	}

	listOfSS2 := []*big.Int{share1to2, share2to2, share3to2, share4to2}
	gsk2 := bn256.GenerateGroupSecretKeyPortion(listOfSS2)
	gpk2 := new(bn256.G2).ScalarBaseMult(gsk2)
	groupShares[1] = gpk2.Marshal()
	s2 := new(crypto.BNGroupSigner)
	err = s2.SetPrivk(gsk2.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	sig2, err := s2.Sign(msg)
	if err != nil {
		t.Fatal(err)
	}

	listOfSS3 := []*big.Int{share1to3, share2to3, share3to3, share4to3}
	gsk3 := bn256.GenerateGroupSecretKeyPortion(listOfSS3)
	gpk3 := new(bn256.G2).ScalarBaseMult(gsk3)
	groupShares[2] = gpk3.Marshal()
	s3 := new(crypto.BNGroupSigner)
	err = s3.SetPrivk(gsk3.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	sig3, err := s3.Sign(msg)
	if err != nil {
		t.Fatal(err)
	}

	listOfSS4 := []*big.Int{share1to4, share2to4, share3to4, share4to4}
	gsk4 := bn256.GenerateGroupSecretKeyPortion(listOfSS4)
	gpk4 := new(bn256.G2).ScalarBaseMult(gsk4)
	groupShares[3] = gpk4.Marshal()
	s4 := new(crypto.BNGroupSigner)
	err = s4.SetPrivk(gsk4.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	sig4, err := s4.Sign(msg)
	if err != nil {
		t.Fatal(err)
	}

	sigs := make([][]byte, 4)
	for k := 0; k < len(sigs); k++ {
		sigs[k] = make([]byte, 192)
	}
	sigs[0] = sig1
	sigs[1] = sig2
	sigs[2] = sig3
	sigs[3] = sig4

	err = s.SetGroupPubk(mpk.Marshal())
	if err != nil {
		t.Fatal(err)
	}

	// Make bad sigs array
	sigsBad := make([][]byte, 2)
	for k := 0; k < len(sigsBad); k++ {
		sigsBad[k] = make([]byte, 192)
	}
	sigsBad[0] = sig1
	sigsBad[1] = sig2
	_, err = s.Aggregate(sigsBad, groupShares)
	if err == nil {
		t.Fatal("Should have raised an error for too few signatures!")
	}

	// Finally submit signature
	grpsig, err := s.Aggregate(sigs, groupShares)
	if err != nil {
		t.Fatal(err)
	}

	bnVal := &crypto.BNGroupValidator{}
	groupk, err := bnVal.PubkeyFromSig(grpsig)
	if err != nil {
		t.Fatal(err)
	}

	bnSigners := []*crypto.BNGroupSigner{}
	bnSigners = append(bnSigners, s1)
	bnSigners = append(bnSigners, s2)
	bnSigners = append(bnSigners, s3)
	bnSigners = append(bnSigners, s4)

	secpSigners := []*crypto.Secp256k1Signer{}
	secpPubks := [][]byte{}
	for _, share := range groupShares {
		signer, pubk := makeSecpSigner(share)
		secpPubks = append(secpPubks, pubk)
		secpSigners = append(secpSigners, signer)
	}

	for _, signer := range bnSigners {
		err := signer.SetGroupPubk(groupk)
		if err != nil {
			t.Fatal(err)
		}
	}
	mpkBin := mpk.Marshal()
	fmt.Printf("MPK: %x\n", mpkBin)
	fmt.Printf("GroupK: %x\n", groupk)
	if bytes.Equal(mpkBin, groupk) {
		t.Fatal("Error")
	}

	return groupk, bnSigners, groupShares, secpSigners, secpPubks
}

func makeSecpSigner(seed []byte) (*crypto.Secp256k1Signer, []byte) {
	secpSigner := &crypto.Secp256k1Signer{}
	err := secpSigner.SetPrivk(crypto.Hasher(seed))
	if err != nil {
		panic(err)
	}
	secpKey, _ := secpSigner.Pubkey()
	return secpSigner, secpKey
}

// We complete everything correctly, happy path
func TestCompletion_Group_1_AllGood(t *testing.T) {
	numValidators := 4
	ecdsaPrivateKeys := tests.InitializePrivateKeys(numValidators)
	fixture := setupEthereum(t, numValidators)
	suite := StartFromGPKjPhase(t, fixture, []int{}, []int{}, 100)
	ctx := context.Background()

	monState := objects.NewMonitorState()
	accounts := suite.Eth.GetKnownAccounts()
	for idx := 0; idx < numValidators; idx++ {
		monState.PotentialValidators[accounts[idx].Address] = objects.PotentialValidator{
			Account: accounts[idx].Address,
		}
	}

	for idx := 0; idx < numValidators; idx++ {
		err := monState.PersistState(suite.DKGStatesDbs[idx])
		assert.Nil(t, err)
	}

	for idx := 0; idx < numValidators; idx++ {
		for j := 0; j < numValidators; j++ {
			disputeGPKjTask := suite.DisputeGPKjTasks[idx][j]

			err := disputeGPKjTask.Initialize(ctx, nil, suite.DKGStatesDbs[idx], fixture.Logger, suite.Eth, fixture.Contracts, "disputeGPKjTask", "task-id", nil)
			assert.Nil(t, err)
			err = disputeGPKjTask.Prepare(ctx)
			assert.Nil(t, err)

			shouldExecute, err := disputeGPKjTask.ShouldExecute(ctx)
			assert.Nil(t, err)
			assert.True(t, shouldExecute)

			txn, taskError := disputeGPKjTask.Execute(ctx)
			assert.Nil(t, taskError)
			assert.Nil(t, txn)
		}
	}

	dkgState, err := state.GetDkgState(suite.DKGStatesDbs[0])
	assert.Nil(t, err)
	tests.AdvanceTo(suite.Eth, dkgState.PhaseStart+dkgState.PhaseLength)
	myStr := `
		import { ValidatorRawData } from "../../ethdkg/setup";
		export const validatorsSnapshots: ValidatorRawData[] = [
	`
	for idx := 0; idx < numValidators; idx++ {
		completionTask := suite.CompletionTasks[idx]

		err := completionTask.Initialize(ctx, nil, suite.DKGStatesDbs[idx], fixture.Logger, suite.Eth, fixture.Contracts, "CompletionTask", "task-id", nil)
		assert.Nil(t, err)
		err = completionTask.Prepare(ctx)
		assert.Nil(t, err)

		dkgState, err := state.GetDkgState(suite.DKGStatesDbs[idx])
		assert.Nil(t, err)

		shouldExecute, err := completionTask.ShouldExecute(ctx)
		assert.Nil(t, err)
		if shouldExecute {
			txn, taskError := completionTask.Execute(ctx)
			amILeading, err := utils.AmILeading(
				suite.Eth,
				ctx,
				fixture.Logger,
				int(completionTask.GetStart()),
				completionTask.StartBlockHash[:],
				numValidators,
				// we need -1 since ethdkg indexes start at 1 while leader election expect index starting at 0.
				dkgState.Index-1,
				constants.ETHDKGDesperationFactor,
				constants.ETHDKGDesperationDelay,
			)
			assert.Nil(t, err)
			if amILeading {
				assert.Nil(t, taskError)
				rcptResponse, err := fixture.Watcher.Subscribe(ctx, txn, nil)
				assert.Nil(t, err)
				tests.WaitGroupReceipts(t, suite.Eth, []transaction.ReceiptResponse{rcptResponse})
			} else {
				assert.Nil(t, txn)
				assert.NotNil(t, taskError)
				assert.True(t, taskError.IsRecoverable())
			}
		}
		encryptedShares := `[`
		for _, share := range dkgState.Participants[accounts[idx].Address].EncryptedShares {
			encryptedShares += fmt.Sprintf(`"0x%x",`, share)
		}
		encryptedShares += "]"
		commitments := `[`
		for _, commit := range dkgState.Participants[accounts[idx].Address].Commitments {
			commitments += fmt.Sprintf(`["0x%x","0x%x"],`, commit[0], commit[1])
		}
		commitments += "]"
		myStr += fmt.Sprintf(`
			{
			privateKey: "0x%x",
			address: "0x%x",
			alicenetPublicKey:["0x%x", "0x%x"],
			encryptedShares:%s,
			commitments:%s,
			keyShareG1: ["0x%x", "0x%x"],
			keyShareG1CorrectnessProof: ["0x%x", "0x%x"],
			keyShareG2: ["0x%x", "0x%x", "0x%x", "0x%x"],
			mpk: ["0x%x", "0x%x", "0x%x", "0x%x"],
			gpkj: ["0x%x", "0x%x", "0x%x", "0x%x"],
			},
		`,
			ecdsaPrivateKeys[idx].D,
			dkgState.Account.Address,
			dkgState.TransportPublicKey[0], dkgState.TransportPublicKey[1],
			encryptedShares,
			commitments,
			dkgState.Participants[accounts[idx].Address].KeyShareG1s[0], dkgState.Participants[accounts[idx].Address].KeyShareG1s[1],
			dkgState.Participants[accounts[idx].Address].KeyShareG1CorrectnessProofs[0], dkgState.Participants[accounts[idx].Address].KeyShareG1CorrectnessProofs[1],
			dkgState.Participants[accounts[idx].Address].KeyShareG2s[0], dkgState.Participants[accounts[idx].Address].KeyShareG2s[1], dkgState.Participants[accounts[idx].Address].KeyShareG2s[2], dkgState.Participants[accounts[idx].Address].KeyShareG2s[3],
			dkgState.MasterPublicKey[0], dkgState.MasterPublicKey[1], dkgState.MasterPublicKey[2], dkgState.MasterPublicKey[3],
			dkgState.Participants[accounts[idx].Address].GPKj[0], dkgState.Participants[accounts[idx].Address].GPKj[1], dkgState.Participants[accounts[idx].Address].GPKj[2], dkgState.Participants[accounts[idx].Address].GPKj[3],
		)
	}
	myStr += "];\n"
	fmt.Print(myStr)
	bnSigners := []*crypto.BNGroupSigner{}
	for idx := 0; idx < n; idx++ {
		state := dkgStates[idx]
		signer := &crypto.BNGroupSigner{}
		signer.SetPrivk(state.GroupPrivateKey.Bytes())
		bnSigners = append(bnSigners, signer)
		groupKey, err := bn256A.MarshalBigIntSlice(state.MasterPublicKey[:])
		if err != nil {
			t.Fatal(err)
		}
		err = signer.SetGroupPubk(groupKey)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Valid at 1024
	grpSig1024, bClaimsBin1024, err := GenerateSnapshotData(1, 1024, bnSigners, n, state.MasterPublicKey[:], false)
	if err != nil {
		t.Fatal(err)
	}

	//Valid block at 2048
	grpSig2048, bClaimsBin2048, err := GenerateSnapshotData(1, 2048, bnSigners, n, state.MasterPublicKey[:], false)
	if err != nil {
		t.Fatal(err)
	}

	// Incorrect height at 500
	grpSig500, bClaimsBin500, err := GenerateSnapshotData(1, 500, bnSigners, n, state.MasterPublicKey[:], false)
	if err != nil {
		t.Fatal(err)
	}

	// Incorrect chanid 2 at 1024
	grpSigChain2, bClaimsBinChain2, err := GenerateSnapshotData(2, 1024, bnSigners, n, state.MasterPublicKey[:], false)
	if err != nil {
		t.Fatal(err)
	}

	// Incorrect signature in a valid block at 1024
	grpSigIncorrect, bClaimsBinIncorrect, err := GenerateSnapshotData(2, 1024, bnSigners, n, state.MasterPublicKey[:], true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf(`
		export const validSnapshot1024: Snapshot = {
			GroupSignature: "0x%x",
			BClaims: "0x%x"
		};

		export const validSnapshot2048: Snapshot = {
			GroupSignature: "0x%x",
			BClaims: "0x%x"
		};

		export const invalidSnapshot500: Snapshot = {
			GroupSignature: "0x%x",
			BClaims: "0x%x"
		};

		export const invalidSnapshotChainID2: Snapshot = {
			GroupSignature: "0x%x",
			BClaims: "0x%x"
		};

		export const invalidSnapshotIncorrectSig: Snapshot = {
			GroupSignature: "0x%x",
			BClaims: "0x%x"
		};`,
		grpSig1024, bClaimsBin1024,
		grpSig2048, bClaimsBin2048,
		grpSig500, bClaimsBin500,
		grpSigChain2, bClaimsBinChain2,
		grpSigIncorrect, bClaimsBinIncorrect,
	)

}

func GenerateSnapshotData(chainID uint32, height uint32, bnSigners []*crypto.BNGroupSigner, n int, mpkI []*big.Int, fakeSig bool) ([]byte, []byte, error) {
	bclaims := &objs.BClaims{
		ChainID:    chainID,
		Height:     height,
		TxCount:    0,
		PrevBlock:  crypto.Hasher([]byte("")),
		TxRoot:     crypto.Hasher([]byte("")),
		StateRoot:  crypto.Hasher([]byte("")),
		HeaderRoot: crypto.Hasher([]byte("")),
	}

	blockHash, err := bclaims.BlockHash()
	if err != nil {
		return nil, nil, err
	}

	bClaimsBin, err := bclaims.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}

	grpsig := []byte{}
	if fakeSig {
		grpsig, err = bnSigners[0].Sign(blockHash)
		if err != nil {
			return nil, nil, err
		}
	} else {
		grpsig, err = GenerateBlockSignature(bnSigners, n, blockHash, mpkI)
		if err != nil {
			return nil, nil, err
		}
	}

	bnVal := &crypto.BNGroupValidator{}
	_, err = bnVal.Validate(blockHash, grpsig)
	if err != nil {
		return nil, nil, err
	}

	return grpsig, bClaimsBin, nil
}

func GenerateBlockSignature(bnSigners []*crypto.BNGroupSigner, n int, blockHash []byte, mpkI []*big.Int) ([]byte, error) {
	sigs := [][]byte{}
	groupShares := [][]byte{}
	for idx := 0; idx < n; idx++ {
		sig, err := bnSigners[idx].Sign(blockHash)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Sig: %x\n", sig)
		sigs = append(sigs, sig)
		pkShare, err := bnSigners[idx].PubkeyShare()
		if err != nil {
			return nil, err
		}
		groupShares = append(groupShares, pkShare)
		fmt.Printf("Pkshare: %x\n", pkShare)
	}
	s := new(crypto.BNGroupSigner)
	mpk, err := bn256A.MarshalBigIntSlice(mpkI)
	err = s.SetGroupPubk(mpk)
	if err != nil {
		return nil, err
	}

	// Finally submit signature
	grpsig, err := s.Aggregate(sigs, groupShares)
	if err != nil {
		return nil, err
	}
	return grpsig, nil

}

// We complete everything correctly, but we do not complete in time
func TestCompletion_Group_1_Bad1(t *testing.T) {
	numValidators := 6
	fixture := setupEthereum(t, numValidators)
	suite := StartFromGPKjPhase(t, fixture, []int{}, []int{}, 100)
	ctx := context.Background()

	dkgState, err := state.GetDkgState(suite.DKGStatesDbs[0])
	assert.Nil(t, err)
	tests.AdvanceTo(suite.Eth, dkgState.PhaseStart+dkgState.PhaseLength)

	task := suite.CompletionTasks[0]
	err = task.Initialize(ctx, nil, suite.DKGStatesDbs[0], fixture.Logger, suite.Eth, fixture.Contracts, "CompletionTask", "task-id", nil)
	assert.Nil(t, err)

	err = task.Prepare(ctx)
	assert.Nil(t, err)

	// Advance to completion submission phase; note we did *not* submit MPK
	tests.AdvanceTo(suite.Eth, task.Start+dkgState.PhaseLength)

	// Do MPK Submission task
	txn, err := task.Execute(ctx)
	assert.NotNil(t, err)
	assert.Nil(t, txn)
}

func TestCompletion_Group_1_Bad2(t *testing.T) {
	task := dkg.NewCompletionTask(1, 100)
	db := mocks.NewTestDB()
	log := logging.GetLogger("test").WithField("test", "test")

	err := task.Initialize(context.Background(), nil, db, log, nil, nil, "", "", nil)
	assert.Nil(t, err)

	taskErr := task.Prepare(context.Background())
	assert.NotNil(t, taskErr)
	assert.False(t, taskErr.IsRecoverable())
}
