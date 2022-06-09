package books

type Book struct {
	UUID   string `json:"uuid"` // field name in dynamodb !case sensitive
	Title  string `json:"title"`
	Author string `json:"author"`
	Editor string `json:"editor"`
}
