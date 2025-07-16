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

func f7() error {
	_, err := func() (int, error) {
		return 1, errSentinel
	}()
	if err == nil {
		fmt.Println("ok")
	}
	return err
}

var errSentinel = errors.New("sentinel")

func f8() error {
	_, err := func() (int, error) {
		return 1, errSentinel
	}()
	if err != nil {
		return fmt.Errorf("error: %v", err) // want "returning sentinel .* error"
	}
	return nil
}

func f9() error {
	_, err := func() (int, error) {
		return 1, errSentinel
	}()
	if err != nil {
		return fmt.Errorf("error") // want "returning sentinel .* error"
	}
	return nil
}

func f10() {
	if _, err := func() (int, error) {
		return 1, errSentinel
	}(); nil != err {
		fmt.Println("ok")
	}
}

func f11() (error, int) {
	return nil, 0
}

func f12() {
	err, _ := f11()
	if err != nil {
		return
	}
}

func f13() error {
	_, err := func() (int, error) {
		return 1, errSentinel
	}()
	if err != nil {
		return fmt.Errorf("error: %w", errors.New("another error")) // want "returning sentinel .* error"
	}
	return nil
}

func f14() error {
	_, err := func() (int, error) {
		return 1, errSentinel
	}()
	if err != nil {
		return fmt.Errorf("error: %w") // want "returning sentinel .* error"
	}
	return nil
}

func anotherError() error {
	return errors.New("another")
}

func f15() error {
	if anotherError() != nil {
		return ErrNotFound
	}
	return nil
}

func f16() (int, string) {
	err := errors.New("foo")
	if err != nil {
		return 1, "not an error"
	}
	return 0, ""
}

func f17() error {
	err := errors.New("original error")
	if err != nil {
		return fmt.Errorf("an error occurred: %w", err)
	}
	return nil
}