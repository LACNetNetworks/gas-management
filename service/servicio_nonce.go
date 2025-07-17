package service

import (
	"math/big"
)

// BigIntQueue implementa una cola FIFO de *big.Int
// para almacenar nonces o cualquier valor numérico.
type BigIntQueue []*big.Int

// NewQueue crea e inicializa una nueva cola.
func NewQueue() *BigIntQueue {
	q := BigIntQueue{}
	return &q
}

// Enqueue agrega un elemento al final de la cola.
func (q *BigIntQueue) Enqueue(val *big.Int) {
	*q = append(*q, val)
}

// Dequeue saca el primer elemento de la cola y lo devuelve.
// Si la cola está vacía, devuelve nil.
func (q *BigIntQueue) Dequeue() *big.Int {
	if len(*q) == 0 {
		return nil
	}
	val := (*q)[0]
	*q = (*q)[1:]
	return val
}
