package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const HeapLimit = 5000

// Each Transaction stored as an object inside our Priority Queue
type TransactionItem struct {
	txhash     string
	gas        int
	feepercent float64
	signature  string
	priority   float64
	index      int
}

//
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// container/heap lets us define our own functions for the Priority Queue
type PriorityQueue []*TransactionItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*TransactionItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Strip each data item of their field name, i.e. TxHash=7B32CDE34 -> 7B32CDE34
func removeEqualsText(fieldStr string, strToRemove string) string {
	res := strings.ReplaceAll(fieldStr, strToRemove, "")
	return res
}

// Push every line of transactions.txt into PriorityQueue
func pushToPriorityQueue(file *os.File, pq *PriorityQueue) error {
	scanner := bufio.NewScanner(file)
	idx := 0
	// For every line in transactions.txt, parse each field, insert fields into TransactionItem obj, push obj to Priority Queue
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, " ")

		// Handle not enough fields being passed in
		if len(items) < 4 {
			return errors.New("Not enough fields in transaction")
		}

		// Convert string forms of Gas and FeePerGas to int and float
		feePercStr := removeEqualsText(items[2], "FeePerGas=")

		gasInt, err := strconv.Atoi(removeEqualsText(items[1], "Gas="))
		if err != nil {
			return err
		}
		feePercentInt, err := strconv.ParseFloat(feePercStr, 64)
		if err != nil {
			return err
		}

		// Error check negative numbers
		if gasInt < 0 || feePercentInt < 0 {
			return errors.New("Negative Number in transaction")
		}

		// Create TransactionItem obj for current line, also calculating priority as Gas * FeePerGas
		transactionitem := &TransactionItem{
			txhash:     removeEqualsText(items[0], "TxHash="),
			gas:        gasInt,
			feepercent: feePercentInt,
			signature:  removeEqualsText(items[3], "Signature="),
			priority:   (float64(gasInt) * feePercentInt),
			index:      idx,
		}
		idx++
		heap.Push(pq, transactionitem)

		// If we have reached our HeapLimit of 5000, pop the item with least priority off the Queue
		if pq.Len() > HeapLimit {
			transactionitem = heap.Pop(pq).(*TransactionItem)
		}
	}
	return nil
}

// Prints every object in PriorityQueue to prioritized-transactions.txt in order of their priority (Gas * FeePerGas)
func transactionsToOutFile(pq *PriorityQueue) {
	f, err := os.Create("prioritized-transactions.txt")
	check(err)

	defer f.Close()

	for pq.Len() > 0 {
		transactionitem := heap.Pop(pq).(*TransactionItem)

		// To view the priority, adding transactionitem.priority to the Sprintf will show priorities top down in ascending order
		outStr := fmt.Sprintf("TxHash=%s Gas=%d FeePerGas=%f Signature=%s\n", transactionitem.txhash, transactionitem.gas, transactionitem.feepercent, transactionitem.signature)
		_, err := f.WriteString(outStr)
		check(err)
	}
}

func main() {
	// Open transactions.txt for parsing, if we cannot open return error
	file, err := os.Open("transactions.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Create our Priority Queue
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Push every line in transactions.txt into Priority Queue (up to 5000 transactions)
	err = pushToPriorityQueue(file, &pq)

	// On invalid input, exit program. Can take a few line out to keep program running on bad input
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Print all transactions out to prioritized-transactions.txt
	transactionsToOutFile(&pq)

	fmt.Println("Done!, logged to prioritized-transactions.txt")
}
