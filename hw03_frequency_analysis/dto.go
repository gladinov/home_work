package hw03frequencyanalysis

import "errors"

var (
	ErrWordIsEmpty         error = errors.New("word is empty")
	ErrSingleDashIsNotWord error = errors.New("single dash is not a valid word")
	ErrEmptyWordAfterTrim  error = errors.New("word is empty after trim")
)
