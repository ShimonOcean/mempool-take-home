// Couple of quick tests for checking mostly input stuff, testing with 0 input, negative numbers, short/cutoff input

package main

import (
	"container/heap"
	"os"
	"testing"
)

func TestRemoveEqualsText(t *testing.T) {
	if removeEqualsText("TxHash=40E10C7CF56A738C0B8AD4EE30EA8008C7B2334B3ADA195083F8CB18BD3911A0", "TxHash=") != "40E10C7CF56A738C0B8AD4EE30EA8008C7B2334B3ADA195083F8CB18BD3911A0" {
		t.Error("removeEqualsTest 1 failed")
	}
	if removeEqualsText("TxHash=TXHASH2899DC689106C8FCEA3E24E4AFFC597D2B4E701F99EB8CD909217D323F", "TxHash=") != "TXHASH2899DC689106C8FCEA3E24E4AFFC597D2B4E701F99EB8CD909217D323F" {
		t.Error("removeEqualsTest 2 failed")
	}
	if removeEqualsText("TxHash=", "TxHash=") != "" {
		t.Error("removeEqualsTest 3 failed")
	}
	if removeEqualsText("16633D0A25ECA886F100A34BA5C43366732836E6E7B140159298C71CF78309F9", "TxHash=") != "16633D0A25ECA886F100A34BA5C43366732836E6E7B140159298C71CF78309F9" {
		t.Error("removeEqualsTest 4 failed")
	}
}

func TestPushToPriorityQueue(t *testing.T) {
	// Test 1, Testing with Gas = 0 and FeePerGas = 0
	f, _ := os.Create("test_pq.txt")
	_, err := f.WriteString(`TxHash=1B0B5C0EE1E167DA2DEC3533E79BA67B949D50B281F2673D3D90D5D1DA8CC8ED Gas=102000 FeePerGas=0.5612550629386129 Signature=7AEB155062299FD2B67D4A75AF1EEC70FAC3F518AAFFE966CC16F7A28E73EFFC804357C49DC0FAE2CA618DC6D25EFCD4BF20BCF9B2C162B7501215AAE897AC64
TxHash=1B0B5C0EE1E167DA2DEC3533E79BA67B949D50B281F2673D3D90D5D1DA8CC8ED Gas=102000 FeePerGas=0 Signature=7AEB155062299FD2B67D4A75AF1EEC70FAC3F518AAFFE966CC16F7A28E73EFFC804357C49DC0FAE2CA618DC6D25EFCD4BF20BCF9B2C162B7501215AAE897AC64
TxHash=1B0B5C0EE1E167DA2DEC3533E79BA67B949D50B281F2673D3D90D5D1DA8CC8ED Gas=0 FeePerGas=0 Signature=7AEB155062299FD2B67D4A75AF1EEC70FAC3F518AAFFE966CC16F7A28E73EFFC804357C49DC0FAE2CA618DC6D25EFCD4BF20BCF9B2C162B7501215AAE897AC64`)
	if err != nil {
		os.Exit(1)
	}

	file, _ := os.Open("test_pq.txt")
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	err = pushToPriorityQueue(file, &pq)
	if err != nil {
		t.Error("pushToPriorityQueue 1 failed, expected to work")
	}

	// Test 2, Just TxHash, should return error, not enough fields
	f, _ = os.Create("test_pq.txt")
	_, err = f.WriteString(`TxHash=1B0B5C0EE1E167DA2DEC3533E79BA67B949D50B281F2673D3D90D5D1DA8CC8ED`)
	if err != nil {
		os.Exit(1)
	}
	file, _ = os.Open("test_pq.txt")
	pq = make(PriorityQueue, 0)
	heap.Init(&pq)
	err = pushToPriorityQueue(file, &pq)
	if err == nil {
		t.Error("pushToPriorityQueue 2 did not return error as expected")
	}
	transactionsToOutFile(&pq)

	// Test 3, Negative Numbers, should return error negative numbers
	f, _ = os.Create("test_pq.txt")
	_, err = f.WriteString(`TxHash=1B0B5C0EE1E167DA2DEC3533E79BA67B949D50B281F2673D3D90D5D1DA8CC8ED Gas=102000 FeePerGas=-0.5612550629386129 Signature=7AEB155062299FD2B67D4A75AF1EEC70FAC3F518AAFFE966CC16F7A28E73EFFC804357C49DC0FAE2CA618DC6D25EFCD4BF20BCF9B2C162B7501215AAE897AC64
TxHash=1B0B5C0EE1E167DA2DEC3533E79BA67B949D50B281F2673D3D90D5D1DA8CC8ED Gas=-102000 FeePerGas=0.5612550629386129 Signature=7AEB155062299FD2B67D4A75AF1EEC70FAC3F518AAFFE966CC16F7A28E73EFFC804357C49DC0FAE2CA618DC6D25EFCD4BF20BCF9B2C162B7501215AAE897AC64\n`)
	if err != nil {
		os.Exit(1)
	}

	file, _ = os.Open("test_pq.txt")

	pq = make(PriorityQueue, 0)
	heap.Init(&pq)
	err = pushToPriorityQueue(file, &pq)
	if err == nil {
		t.Error("pushToPriorityQueue 3 did not return error as expected")
	}

}
