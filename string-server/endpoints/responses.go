package endpoints

type IsPalResponse struct {
	Message bool `json:"message"`
}

type ReverseResponse struct {
	Word string `json:"reversed_word"`
}