package executor

import (
	"context"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/alicenet/alicenet/consensus/db"
	"github.com/alicenet/alicenet/layer1/executor/tasks"
	"github.com/alicenet/alicenet/layer1/transaction"
	"github.com/alicenet/alicenet/logging"
	"github.com/alicenet/alicenet/test/mocks"
	mockrequire "github.com/derision-test/go-mockgen/testutil/require"
	"github.com/dgraph-io/badger/v2"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func getTaskProcessor(t *testing.T) (*TasksProcessor, *mocks.MockClient, *db.Database, *taskResponseChan, *mocks.MockWatcher) {
	db := mocks.NewTestDB()
	client := mocks.NewMockClient()
	client.ExtractTransactionSenderFunc.SetDefaultReturn(common.Address{}, nil)
	client.GetTxMaxStaleBlocksFunc.SetDefaultReturn(10)
	hdr := &types.Header{
		Number: big.NewInt(1),
	}
	client.GetHeaderByNumberFunc.SetDefaultReturn(hdr, nil)
	client.GetBlockBaseFeeAndSuggestedGasTipFunc.SetDefaultReturn(big.NewInt(100), big.NewInt(1), nil)
	client.GetDefaultAccountFunc.SetDefaultReturn(accounts.Account{Address: common.Address{}})
	client.GetTransactionByHashFunc.SetDefaultReturn(nil, false, nil)
	client.GetFinalityDelayFunc.SetDefaultReturn(10)

	logger := logging.GetLogger("test")

	txWatcher := mocks.NewMockWatcher()
	taskProcessor, err := NewTaskProcessor(txWatcher, db, logger.WithField("Component", "schedule"))
	assert.Nil(t, err)

	taskRespChan := &taskResponseChan{trChan: make(chan tasks.Response, 100)}
	return taskProcessor, client, db, taskRespChan, txWatcher
}

func Test_TaskProcessor_HappyPath(t *testing.T) {
	processor, client, db, taskRespChan, txWatcher := getTaskProcessor(t)
	defer taskRespChan.close()

	receipt := &types.Receipt{
		Status:      types.ReceiptStatusSuccessful,
		BlockNumber: big.NewInt(20),
	}

	receiptResponse := mocks.NewMockReceiptResponse()
	receiptResponse.IsReadyFunc.SetDefaultReturn(true)
	receiptResponse.GetReceiptBlockingFunc.SetDefaultReturn(receipt, nil)

	txWatcher.SubscribeFunc.SetDefaultReturn(receiptResponse, nil)

	client.GetTransactionReceiptFunc.SetDefaultReturn(receipt, nil)

	task := mocks.NewMockTask()
	task.PrepareFunc.SetDefaultReturn(nil)
	txn := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     []byte{52, 66, 175, 92},
	})
	task.ExecuteFunc.SetDefaultReturn(txn, nil)
	task.ShouldExecuteFunc.SetDefaultReturn(true, nil)
	task.GetLoggerFunc.SetDefaultReturn(processor.logger)

	mainCtx := context.Background()
	processor.ProcessTask(mainCtx, task, "", "123", db, processor.logger, client, taskRespChan)

	mockrequire.CalledOnce(t, task.PrepareFunc)
	mockrequire.CalledOnce(t, task.ExecuteFunc)
	mockrequire.CalledOnceWith(t, task.FinishFunc, mockrequire.Values(nil))
	assert.Emptyf(t, processor.TxsBackup, "Expected transactions to be empty")
}

func Test_TaskProcessor_TaskErrorRecoverable(t *testing.T) {
	processor, client, db, taskRespChan, txWatcher := getTaskProcessor(t)
	defer taskRespChan.close()

	receipt := &types.Receipt{
		Status:      types.ReceiptStatusSuccessful,
		BlockNumber: big.NewInt(20),
	}

	receiptResponse := mocks.NewMockReceiptResponse()
	receiptResponse.IsReadyFunc.SetDefaultReturn(true)
	receiptResponse.GetReceiptBlockingFunc.SetDefaultReturn(receipt, nil)

	txWatcher.SubscribeFunc.SetDefaultReturn(receiptResponse, nil)

	client.GetTransactionReceiptFunc.SetDefaultReturn(receipt, nil)

	taskErr := tasks.NewTaskErr("Recoverable error", true)
	task := mocks.NewMockTask()
	task.PrepareFunc.PushReturn(taskErr)
	task.PrepareFunc.PushReturn(nil)
	txn := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     []byte{52, 66, 175, 92},
	})
	task.ExecuteFunc.SetDefaultReturn(txn, nil)
	task.ShouldExecuteFunc.SetDefaultReturn(true, nil)
	task.GetLoggerFunc.SetDefaultReturn(processor.logger)

	mainCtx := context.Background()
	processor.ProcessTask(mainCtx, task, "", "123", db, processor.logger, client, taskRespChan)

	mockrequire.CalledN(t, task.PrepareFunc, 2)
	mockrequire.CalledOnce(t, task.ExecuteFunc)
	mockrequire.CalledOnceWith(t, task.FinishFunc, mockrequire.Values(nil))
	assert.Emptyf(t, processor.TxsBackup, "Expected transactions to be empty")
}

func Test_TaskProcessor_UnrecoverableError(t *testing.T) {
	processor, client, db, taskRespChan, txWatcher := getTaskProcessor(t)
	defer taskRespChan.close()

	receipt := &types.Receipt{
		Status:      types.ReceiptStatusSuccessful,
		BlockNumber: big.NewInt(20),
	}

	receiptResponse := mocks.NewMockReceiptResponse()
	receiptResponse.IsReadyFunc.SetDefaultReturn(true)
	receiptResponse.GetReceiptBlockingFunc.SetDefaultReturn(receipt, nil)

	txWatcher.SubscribeFunc.SetDefaultReturn(receiptResponse, nil)

	client.GetTransactionReceiptFunc.SetDefaultReturn(receipt, nil)

	task := mocks.NewMockTask()
	taskErr := tasks.NewTaskErr("Unrecoverable error", false)

	task.PrepareFunc.SetDefaultReturn(taskErr)
	mainCtx := context.Background()
	processor.ProcessTask(mainCtx, task, "", "123", db, processor.logger, client, taskRespChan)

	mockrequire.CalledOnce(t, task.PrepareFunc)
	mockrequire.NotCalled(t, task.ExecuteFunc)
	mockrequire.CalledOnceWith(t, task.FinishFunc, mockrequire.Values(taskErr))
	assert.Emptyf(t, processor.TxsBackup, "Expected transactions to be empty")
}

func Test_TaskProcessor_TaskInTasksManagerTransactions(t *testing.T) {
	processor, client, db, taskRespChan, txWatcher := getTaskProcessor(t)
	defer taskRespChan.close()

	receipt := &types.Receipt{
		Status:      types.ReceiptStatusSuccessful,
		BlockNumber: big.NewInt(20),
	}

	receiptResponse := mocks.NewMockReceiptResponse()
	receiptResponse.IsReadyFunc.SetDefaultReturn(true)
	receiptResponse.GetReceiptBlockingFunc.SetDefaultReturn(receipt, nil)

	txWatcher.SubscribeFunc.SetDefaultReturn(receiptResponse, nil)

	client.GetTransactionReceiptFunc.SetDefaultReturn(receipt, nil)

	task := mocks.NewMockTask()
	task.PrepareFunc.SetDefaultReturn(nil)
	txn := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     []byte{52, 66, 175, 92},
	})
	task.ShouldExecuteFunc.SetDefaultReturn(true, nil)
	task.GetLoggerFunc.SetDefaultReturn(processor.logger)
	taskId := "123"
	task.GetIdFunc.SetDefaultReturn(taskId)

	mainCtx := context.Background()
	processor.TxsBackup[task.GetId()] = txn
	processor.ProcessTask(mainCtx, task, "", taskId, db, processor.logger, client, taskRespChan)

	mockrequire.NotCalled(t, task.PrepareFunc)
	mockrequire.NotCalled(t, task.ExecuteFunc)
	assert.Emptyf(t, processor.TxsBackup, "Expected transactions to be empty")
}

func Test_TaskProcessor_ExecuteWithErrors(t *testing.T) {
	processor, client, db, taskRespChan, txWatcher := getTaskProcessor(t)
	defer taskRespChan.close()

	receipt := &types.Receipt{
		Status:      types.ReceiptStatusSuccessful,
		BlockNumber: big.NewInt(20),
	}

	receiptResponse := mocks.NewMockReceiptResponse()
	receiptResponse.IsReadyFunc.SetDefaultReturn(true)
	receiptResponse.GetReceiptBlockingFunc.SetDefaultReturn(receipt, nil)

	txWatcher.SubscribeFunc.SetDefaultReturn(receiptResponse, nil)

	client.GetTransactionReceiptFunc.SetDefaultReturn(receipt, nil)

	task := mocks.NewMockTask()
	task.PrepareFunc.SetDefaultReturn(nil)
	task.ExecuteFunc.PushReturn(nil, tasks.NewTaskErr("Recoverable error", true))
	unrecoverableError := tasks.NewTaskErr("Unrecoverable error", false)
	task.ExecuteFunc.PushReturn(nil, unrecoverableError)
	task.ShouldExecuteFunc.SetDefaultReturn(true, nil)
	task.GetLoggerFunc.SetDefaultReturn(processor.logger)

	mainCtx := context.Background()
	processor.ProcessTask(mainCtx, task, "", "123", db, processor.logger, client, taskRespChan)

	mockrequire.CalledOnce(t, task.PrepareFunc)
	mockrequire.CalledN(t, task.ExecuteFunc, 2)
	mockrequire.CalledOnceWith(t, task.FinishFunc, mockrequire.Values(unrecoverableError))
	assert.Emptyf(t, processor.TxsBackup, "Expected transactions to be empty")
}

func Test_TaskProcessor_ReceiptWithErrorAndFailure(t *testing.T) {
	processor, client, db, taskRespChan, txWatcher := getTaskProcessor(t)
	defer taskRespChan.close()

	receipt := &types.Receipt{
		Status:      types.ReceiptStatusFailed,
		BlockNumber: big.NewInt(20),
	}

	receiptResponse := mocks.NewMockReceiptResponse()

	task := mocks.NewMockTask()
	task.PrepareFunc.SetDefaultReturn(nil)
	txn := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     []byte{52, 66, 175, 92},
	})
	task.ShouldExecuteFunc.PushReturn(true, nil)
	task.ExecuteFunc.SetDefaultReturn(txn, nil)
	task.GetLoggerFunc.SetDefaultReturn(processor.logger)

	receiptResponse.IsReadyFunc.PushReturn(false)
	task.ShouldExecuteFunc.PushReturn(true, nil)

	receiptResponse.IsReadyFunc.PushReturn(true)
	receiptResponse.GetReceiptBlockingFunc.PushReturn(nil, errors.New("got error while getting receipt"))
	task.ShouldExecuteFunc.PushReturn(true, nil)
	txWatcher.SubscribeFunc.PushReturn(receiptResponse, nil)

	receiptResponse.IsReadyFunc.PushReturn(true)
	receiptResponse.GetReceiptBlockingFunc.PushReturn(receipt, nil)
	task.ShouldExecuteFunc.PushReturn(true, nil)
	txWatcher.SubscribeFunc.PushReturn(receiptResponse, nil)

	receiptResponse.IsReadyFunc.PushReturn(false)
	txWatcher.SubscribeFunc.PushReturn(receiptResponse, nil)
	task.ShouldExecuteFunc.PushReturn(false, nil)

	mainCtx := context.Background()
	processor.ProcessTask(mainCtx, task, "", "123", db, processor.logger, client, taskRespChan)

	mockrequire.CalledOnce(t, task.PrepareFunc)
	mockrequire.CalledN(t, task.ExecuteFunc, 3)
	mockrequire.CalledN(t, receiptResponse.IsReadyFunc, 4)
	mockrequire.CalledN(t, receiptResponse.GetReceiptBlockingFunc, 2)
	mockrequire.CalledN(t, task.ShouldExecuteFunc, 5)

	mockrequire.CalledOnceWith(t, task.FinishFunc, mockrequire.Values(nil))
	assert.Emptyf(t, processor.TxsBackup, "Expected transactions to be empty")
}

func Test_TaskProcessor_RecoveringTaskProcessor(t *testing.T) {
	dir, err := ioutil.TempDir("", "db-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	}()
	opts := badger.DefaultOptions(dir)
	rawDB, err := badger.Open(opts)
	if err != nil {
		t.Fatal(err)
	}
	defer rawDB.Close()

	db := &db.Database{}
	db.Init(rawDB)

	client := mocks.NewMockClient()
	client.ExtractTransactionSenderFunc.SetDefaultReturn(common.Address{}, nil)
	client.GetTxMaxStaleBlocksFunc.SetDefaultReturn(10)
	hdr := &types.Header{
		Number: big.NewInt(1),
	}
	client.GetHeaderByNumberFunc.SetDefaultReturn(hdr, nil)
	client.GetBlockBaseFeeAndSuggestedGasTipFunc.SetDefaultReturn(big.NewInt(100), big.NewInt(1), nil)
	client.GetDefaultAccountFunc.SetDefaultReturn(accounts.Account{Address: common.Address{}})
	client.GetTransactionByHashFunc.SetDefaultReturn(nil, false, nil)
	client.GetFinalityDelayFunc.SetDefaultReturn(10)

	logger := logging.GetLogger("test")

	txWatcher := mocks.NewMockWatcher()
	processor, err := NewTaskProcessor(txWatcher, db, logger.WithField("Component", "schedule"))
	assert.Nil(t, err)

	taskRespChan := &taskResponseChan{trChan: make(chan tasks.Response, 100)}
	defer taskRespChan.close()

	receipt := &types.Receipt{
		Status:      types.ReceiptStatusSuccessful,
		BlockNumber: big.NewInt(20),
	}

	receiptResponse := mocks.NewMockReceiptResponse()
	receiptResponse.IsReadyFunc.SetDefaultReturn(true)
	receiptResponse.GetReceiptBlockingFunc.PushReturn(nil, &transaction.ErrTransactionStale{})
	receiptResponse.GetReceiptBlockingFunc.PushReturn(receipt, nil)

	txWatcher.SubscribeFunc.SetDefaultReturn(receiptResponse, nil)

	client.GetTransactionReceiptFunc.SetDefaultReturn(receipt, nil)

	task := mocks.NewMockTask()
	task.PrepareFunc.SetDefaultReturn(nil)
	txn := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		Value:    big.NewInt(1),
		Gas:      1,
		GasPrice: big.NewInt(1),
		Data:     []byte{52, 66, 175, 92},
	})
	task.ExecuteFunc.SetDefaultReturn(txn, nil)
	task.ShouldExecuteFunc.SetDefaultReturn(true, nil)
	task.GetLoggerFunc.SetDefaultReturn(processor.logger)

	mainCtx := context.Background()
	processor.ProcessTask(mainCtx, task, "", "123", db, processor.logger, client, taskRespChan)

	assert.Equalf(t, 1, len(processor.TxsBackup), "Expected one transaction (stale status)")
	processor, err = NewTaskProcessor(txWatcher, db, logger.WithField("Component", "schedule"))
	assert.Nil(t, err)
	processor.ProcessTask(mainCtx, task, "", "123", db, processor.logger, client, taskRespChan)

	mockrequire.CalledOnce(t, task.PrepareFunc)
	mockrequire.CalledOnce(t, task.ExecuteFunc)
	mockrequire.CalledOnceWith(t, task.FinishFunc, mockrequire.Values(nil))
	assert.Emptyf(t, processor.TxsBackup, "Expected transactions to be empty")

}
