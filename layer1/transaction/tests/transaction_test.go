package tests

import (
	"context"
	"fmt"
	"github.com/alicenet/alicenet/layer1/ethereum"
	"github.com/alicenet/alicenet/layer1/tests"
	"github.com/alicenet/alicenet/test/mocks"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"testing"
	"time"

	"github.com/alicenet/alicenet/layer1/transaction"
	"github.com/stretchr/testify/assert"
)

func Setup(t *testing.T, accounts int) (*tests.ClientFixture, *transaction.FrontWatcher) {
	fixture := setupEthereum(t, accounts)
	db := mocks.NewTestDB()
	watcher := transaction.WatcherFromNetwork(fixture.Client, db, false, 1*time.Second)

	return fixture, watcher
}

func TestSubscribeAndWaitForValidTx(t *testing.T) {
	numAccounts := 2
	fixture, watcher := Setup(t, numAccounts)
	eth := fixture.Client
	accounts := eth.GetKnownAccounts()
	assert.Equal(t, numAccounts, len(accounts))

	owner := accounts[0]
	user := accounts[1]

	ctx, cf := context.WithTimeout(context.Background(), 300*time.Second)
	defer cf()

	amount := big.NewInt(12_345)
	txn, err := ethereum.TransferEther(eth, fixture.Logger, owner.Address, user.Address, amount)
	assert.Nil(t, err)

	receipt, err := watcher.SubscribeAndWait(ctx, txn, nil)
	assert.Nil(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, txn.Hash(), receipt.TxHash)

	currentHeight, err := eth.GetCurrentHeight(ctx)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, currentHeight, receipt.BlockNumber.Uint64()+eth.GetFinalityDelay())

	_, isPending, err := eth.GetTransactionByHash(ctx, txn.Hash())
	assert.Nil(t, err)
	assert.False(t, isPending)

	mintTxnOpts, err := eth.GetTransactionOpts(ctx, owner)
	assert.Nil(t, err)
	mintTxnOpts.NoSend = false
	//mintTxnOpts.GasFeeCap = big.NewInt(1_000_000_000)
	mintTxnOpts.Value = amount

	mintTxn, err := ethereum.GetContracts().BToken().MintTo(mintTxnOpts, user.Address, big.NewInt(1))
	assert.Nil(t, err)
	assert.NotNil(t, mintTxn)

	txnRough := &types.DynamicFeeTx{}
	txnRough.ChainID = txn.ChainId()
	txnRough.To = txn.To()
	txnRough.GasFeeCap = new(big.Int).Mul(new(big.Int).SetInt64(2), txn.GasFeeCap())
	txnRough.GasTipCap = new(big.Int).Mul(new(big.Int).SetInt64(2), txn.GasTipCap())
	txnRough.Gas = txn.Gas()
	txnRough.Nonce = txn.Nonce() + 1
	txnRough.Value = txn.Value()
	txnRough.Data = txn.Data()

	<-time.After(2 * time.Second)
	fixture.Logger.Infof("New Gasfee: %v", txnRough.GasFeeCap.String())

	signer := types.NewLondonSigner(txnRough.ChainID)

	_, adminPk := tests.GetAdminAccount()
	signedTx, err := types.SignNewTx(adminPk, signer, txnRough)
	if err != nil {
		fixture.Logger.Errorf("signing error:%v", err)
	}
	err = eth.SendTransaction(ctx, signedTx)
	if err != nil {
		fixture.Logger.Errorf("sending error:%v", err)
	}

}

func TestSubscribeAndWaitForInvalidTxNotSigned(t *testing.T) {
	numAccounts := 2
	fixture, watcher := Setup(t, numAccounts)
	eth := fixture.Client

	accounts := eth.GetKnownAccounts()
	assert.Equal(t, numAccounts, len(accounts))

	owner := accounts[0]
	user := accounts[1]

	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	amount := big.NewInt(12_345)

	txOpts, err := eth.GetTransactionOpts(ctx, owner)
	txOpts.NoSend = true
	txOpts.Value = amount
	assert.Nil(t, err)

	//Creating tx but not sending it
	txn, err := ethereum.GetContracts().BToken().MintTo(txOpts, user.Address, big.NewInt(1))
	assert.Nil(t, err)

	txnRough := &types.DynamicFeeTx{}
	txnRough.ChainID = txn.ChainId()
	txnRough.To = txn.To()
	txnRough.GasFeeCap = new(big.Int).Mul(new(big.Int).SetInt64(2), txn.GasFeeCap())
	txnRough.GasTipCap = new(big.Int).Mul(new(big.Int).SetInt64(2), txn.GasTipCap())
	txnRough.Gas = txn.Gas()
	txnRough.Nonce = txn.Nonce() + 1
	txnRough.Value = txn.Value()
	txnRough.Data = txn.Data()

	txnNotSigned := types.NewTx(txnRough)

	receipt, err := watcher.SubscribeAndWait(ctx, txnNotSigned, nil)
	assert.NotNil(t, err)
	assert.Nil(t, receipt)
	_, ok := err.(*transaction.ErrInvalidMonitorRequest)
	assert.True(t, ok)
}

func TestSubscribeAndWaitForTxNotFound(t *testing.T) {
	numAccounts := 2
	fixture, watcher := Setup(t, numAccounts)
	eth := fixture.Client

	accounts := eth.GetKnownAccounts()
	assert.Equal(t, numAccounts, len(accounts))

	owner := accounts[0]
	user := accounts[1]

	ctx, cf := context.WithTimeout(context.Background(), 300*time.Second)
	defer cf()

	amount := big.NewInt(12_345)

	txOpts, err := eth.GetTransactionOpts(ctx, owner)
	txOpts.NoSend = true
	txOpts.Value = amount
	assert.Nil(t, err)

	//Creating tx but not sending it
	txn, err := ethereum.GetContracts().BToken().MintTo(txOpts, user.Address, big.NewInt(1))
	assert.Nil(t, err)

	txnRough := &types.DynamicFeeTx{}
	txnRough.ChainID = txn.ChainId()
	txnRough.To = txn.To()
	txnRough.GasFeeCap = new(big.Int).Mul(new(big.Int).SetInt64(2), txn.GasFeeCap())
	txnRough.GasTipCap = new(big.Int).Mul(new(big.Int).SetInt64(2), txn.GasTipCap())
	txnRough.Gas = txn.Gas()
	txnRough.Nonce = txn.Nonce() + 1
	txnRough.Value = txn.Value()
	txnRough.Data = txn.Data()

	signer := types.NewLondonSigner(txnRough.ChainID)

	_, adminPk := tests.GetAdminAccount()
	signedTx, err := types.SignNewTx(adminPk, signer, txnRough)
	if err != nil {
		fixture.Logger.Errorf("signing error:%v", err)
	}

	//hardhatEndpoint := "http://127.0.0.1:8545"
	//tests.SetBlockInterval(hardhatEndpoint, 50)

	receipt, err := watcher.SubscribeAndWait(ctx, signedTx, nil)
	assert.NotNil(t, err)
	assert.Nil(t, receipt)

	_, ok := err.(*transaction.ErrTxNotFound)
	assert.True(t, ok)
}

func TestSubscribeAndWaitForStaleTx(t *testing.T) {
	numAccounts := 2
	fixture, watcher := Setup(t, numAccounts)
	eth := fixture.Client

	// setting base fee to 10k GWei
	hardhatEndpoint := "http://127.0.0.1:8545"
	tests.SetNextBlockBaseFee(hardhatEndpoint, 10_000_000_000_000)

	accounts := eth.GetKnownAccounts()
	assert.Equal(t, numAccounts, len(accounts))

	owner := accounts[0]
	user := accounts[1]

	ctx, cf := context.WithTimeout(context.Background(), 300*time.Second)
	defer cf()

	txOpts, err := eth.GetTransactionOpts(ctx, owner)
	txOpts.GasFeeCap = big.NewInt(1_000_000_000_000)
	txOpts.Value = big.NewInt(12_345)
	assert.Nil(t, err)

	txn, err := ethereum.GetContracts().BToken().MintTo(txOpts, user.Address, big.NewInt(1))
	assert.Nil(t, err)
	assert.NotNil(t, txn)

	receipt, err := watcher.Subscribe(ctx, txn, nil)
	//tests.MineBlocks(hardhatEndpoint, 10)
	<-time.After(10 * time.Second)

	fmt.Printf("err - %v", err)
	assert.NotNil(t, err)
	assert.Nil(t, receipt)
	_, ok := err.(*transaction.ErrTransactionStale)
	assert.True(t, ok)
}

//func TestAutoSubscribeOfPendingTransaction(t *testing.T) {
//	finalityDelay := uint64(6)
//	numAccounts := 2
//	_, _ := Setup(finalityDelay, numAccounts, common.HexToAddress("0x0b1F9c2b7bED6Db83295c7B5158E3806d67eC5bc"))
//	assert.Nil(t, err)
//	defer eth.Close()
//
//	testutils.SetBlockInterval(t, eth, 500)
//	// setting base fee to 10k GWei
//	testutils.SetNextBlockBaseFee(t, eth, 10_000_000_000_000)
//
//	accounts := eth.GetKnownAccounts()
//	assert.Equal(t, numAccounts, len(accounts))
//
//	for _, acct := range accounts {
//		err := eth.UnlockAccount(acct)
//		assert.Nil(t, err)
//	}
//
//	owner := accounts[0]
//	user := accounts[1]
//
//	ctx, cf := context.WithTimeout(context.Background(), 300*time.Second)
//	defer cf()
//
//	amount := big.NewInt(12_345)
//
//	txOpts, err := eth.GetTransactionOpts(ctx, owner)
//	txOpts.GasFeeCap = big.NewInt(1_000_000_000)
//	txOpts.Value = amount
//	assert.Nil(t, err)
//
//	_, err = eth.Contracts().BToken().MintTo(txOpts, user.Address, big.NewInt(1))
//	assert.Nil(t, err)
//
//	// logger.Infof("%v", testutils.GetPendingBlock(t, eth))
//	// block, err := eth.GetBlockByNumber(ctx, big.NewInt(-1))
//	// assert.Nil(t, err)
//	// for _, txn := range block.Transactions() {
//	// 	logger.Print(txn.Hash().Hex())
//	// }
//
//	// receipt, err := watcher.SubscribeAndWait(ctx, txn)
//	// assert.NotNil(t, err)
//	// assert.Nil(t, receipt)
//	// _, ok := err.(*transaction.ErrTransactionStale)
//	// assert.True(t, ok)
//}
