package endpoints

type IsPalRequest struct {
	Word string `json:"word"`
}

type ReverseRequest struct {
	Word string `json:"word"`
}