package utils

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

var (
	ErrorMethodNotAllowed        = "method Not allowed"
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToMarshalItem     = "failed to marshal item"

	ErrorCouldNotFetchItem  = "could not fetch item"
	ErrorCouldNotFetchItems = "could not fetch records"
	ErrorCouldNotDeleteItem = "could not delete item"
	ErrorCouldNotPutItem    = "could not put item"
	ErrorCouldNotUpdateItem = "could not update item"

	ErrorInvalidBookData   = "invalid book data"
	ErrorBookDoesNotExists = "Book doesn't exist"
)
