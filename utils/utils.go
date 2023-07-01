package utils

var transactionID int32 = 0

func NewTransactionID() int32 {
	transactionID++
	return transactionID
}
