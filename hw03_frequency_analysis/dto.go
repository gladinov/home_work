package hw03frequencyanalysis

import "errors"

var (
	ErrWordIsEmpty         = errors.New("word is empty")
	ErrSingleDashIsNotWord = errors.New("single dash is not a valid word")
	ErrEmptyWordAfterTrim  = errors.New("word is empty after trim")
)
