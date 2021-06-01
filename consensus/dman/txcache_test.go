package dman

import (
	"testing"

	"github.com/MadBase/MadNet/interfaces"
)

type testingTxMarshaller struct {
}

func (tm *testingTxMarshaller) MarshalTx(tx interfaces.Transaction) ([]byte, error) {
	return tx.TxHash()
}

func (tm *testingTxMarshaller) UnmarshalTx(tx []byte) (interfaces.Transaction, error) {
	return makeTx(tx), nil
}

func makeCache() *txCache {
	txc := &txCache{}
	t := &testingTxMarshaller{}
	txc.Init(t)
	return txc
}

func makeTx(h []byte) interfaces.Transaction {
	return new(testingTransaction).setHashForTest(h)
}

func Test_txCache_Get(t *testing.T) {
	type args struct {
		txHsh []byte
	}

	txc := makeCache()
	tx := makeTx(nil)
	h, _ := tx.TxHash()
	txc.Add(1, tx)
	if !txc.Contains(h) {
		t.Fatal("test failed")
	}
	_, ok := txc.Get(h)
	if !ok {
		t.Fatal("not found in get")
	}

}

func Test_txCache_GetHeight(t *testing.T) {
	txc := makeCache()
	tx1 := makeTx([]byte("bar"))
	tx2 := makeTx([]byte("foo"))
	h1, _ := tx1.TxHash()
	h2, _ := tx2.TxHash()
	txc.Add(1, tx1)
	txc.Add(2, tx2)
	txs1, _ := txc.GetHeight(1)
	if len(txs1) != 1 {
		for hash, rh := range txc.rcache {
			t.Logf("%v %s", rh, hash)
		}
		t.Fatalf("1: not found in get, got len %v", len(txs1))
	}
	h1t, _ := txs1[0].TxHash()
	if string(h1t) != string(h1) {
		t.Fatalf("1: bad hash: %s vs %s", h1t, h1)
	}
	txs2, _ := txc.GetHeight(2)
	if len(txs1) != 1 {
		t.Fatal("2: not found in get")
	}
	h2t, _ := txs2[0].TxHash()
	if string(h2t) != string(h2) {
		t.Fatal("2: bad hash")
	}
	txs3, _ := txc.GetHeight(3)
	if len(txs3) != 0 {
		t.Fatal("3: found in get")
	}

}

func Test_txCache_Del(t *testing.T) {
	txc := makeCache()
	tx := makeTx(nil)
	h, _ := tx.TxHash()
	txc.Add(1, tx)
	txc.Del(h)
}

func Test_txCache_DropBeforeHeight(t *testing.T) {
	txc := makeCache()
	tx1 := makeTx([]byte("bar"))
	tx2 := makeTx([]byte("foo"))
	txc.Add(1, tx1)
	txc.Add(2, tx2)
	txc.DropBeforeHeight(1)
	txc.DropBeforeHeight(257)
}
