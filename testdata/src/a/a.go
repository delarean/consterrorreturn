package a

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

func f1() error {
	err := errors.New("original error")
	if err != nil {
		return ErrNotFound // want "returning sentinel"
	}
	return nil
}

func f2() error {
	err := errors.New("original error")
	if err != nil {
		return errors.New("failed to marshal event") // want "returning sentinel"
	}
	return nil
}

func f3() error {
	err := errors.New("original error")
	if err != nil {
		return fmt.Errorf("... %w", ErrNotFound) // want "returning sentinel"
	}
	return nil
} 

func f4() error {
	err := errors.New("original error")
	if err != nil {
		return err
	}
	return nil
} 

func f5() error {
	err := fmt.Errorf("original error: %w", ErrNotFound)
	if errors.Is(err, ErrNotFound) {
		return ErrNotFound
	}
	return nil
}

func f6() error {
	err := fmt.Errorf("original error: %w", ErrNotFound)
	if errors.Is(err, ErrNotFound) {
		return errors.New("new not found error")
	}
	return nil
}