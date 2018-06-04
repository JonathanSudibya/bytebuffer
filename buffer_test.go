package bytebuffer_test

import (
	"testing"

	"github.com/jonathansudibya/bytebuffer"
)

func assertError(t *testing.T) {
	if r := recover(); r != nil {
		t.Errorf("panic code \n%s", r)
	}
}

func assertNoError(t *testing.T) {
	if r := recover(); r == nil {
		t.Error("code does not produce error")
	}
}

func TestPoolCreation(t *testing.T) {
	tcs := []struct {
		name   string
		config *bytebuffer.Config
		len    int
		cap    int
	}{
		{"empty config", nil, 0, 1024},
		{"custom capacity 2048", &bytebuffer.Config{Cap: 2048}, 0, 2048},
		{"custom capacity 512", &bytebuffer.Config{Cap: 512}, 0, 512},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			defer assertError(t)
			pool := bytebuffer.NewPool(tc.config)

			b := pool.Get()

			if b.Len() != tc.len {
				t.Errorf("byte array length invalid, got %d expected %d", b.Len(), tc.len)
			}

			if b.Cap() != tc.cap {
				t.Errorf("byte array capacity invalid, got %d expected %d", b.Cap(), tc.cap)
			}

			b.Free()
		})
	}
}

func TestWriteByte(t *testing.T) {
	tcs := []struct {
		name     string
		input    byte
		expected byte
	}{
		{"single byte char", byte('a'), byte('a')},
		{"single byte int", byte('2'), byte('2')},
	}

	pool := bytebuffer.NewPool(nil)

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			defer assertError(t)
			b := pool.Get()
			defer b.Free()

			b.WriteByte(tc.input)

			result := b.Bytes()[0]

			if result != tc.expected {
				t.Errorf("result byte invalid, got %x expected %x", result, tc.expected)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tcs := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{"empty byte array", []byte(""), []byte("")},
		{"fill byte array", []byte("hello world"), []byte("hello world")},
	}

	pool := bytebuffer.NewPool(nil)

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			defer assertError(t)
			b := pool.Get()
			defer b.Free()

			b.Write(tc.input)

			result := b.Bytes()

			if result == nil {
				t.Errorf("byte array is nil")
			}

			for i, res := range result {
				if res != tc.expected[i] {
					t.Errorf("invalid byte array, got %s expected %s", result, tc.expected)
				}
			}
		})
	}
}

func TestWriteString(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"fill string", "Hello World", "Hello World"},
	}

	pool := bytebuffer.NewPool(nil)

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			defer assertError(t)
			b := pool.Get()
			defer b.Free()

			b.WriteString(tc.input)

			result := b.String()

			if result != tc.expected {
				t.Errorf("invalid byte array, got %s expected %s", result, tc.expected)
			}
		})
	}
}
