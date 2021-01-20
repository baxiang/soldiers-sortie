package service

import (
	golorem "github.com/drhodes/golorem"
)
type Service interface {
	//
	Word(min,max int)string
	// generate a sentence with at least min words and at most max words.
	Sentence(min, max int) string

	// generate a paragraph with at least min sentences and at most max sentences.
	Paragraph(min, max int) string
}

type LoremService struct {

}

// Implement service functions
func (LoremService) Word(min, max int) string {
	return golorem.Word(min, max)
}

func (LoremService) Sentence(min, max int) string {
	return golorem.Sentence(min, max)
}

func (LoremService) Paragraph(min, max int) string {
	return golorem.Paragraph(min, max)
}