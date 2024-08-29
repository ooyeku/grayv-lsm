package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	cases := []struct {
		password string
	}{
		{"password123"},
		{"simplepassword"},
		{"verylongpasswordthatexceedslengthlimit"},
		{""},
	}

	for _, tc := range cases {
		_, err := HashPassword(tc.password)
		if err != nil {
			t.Errorf("HashPassword(%q) => unexpected error %v", tc.password, err)
			continue
		}
		if _, err := HashPassword(tc.password); err != nil {
			t.Errorf("HashPassword(%q) => unexpected error %v", tc.password, err)
		}
	}
}

func TestCheckPasswordHash(t *testing.T) {
	cases := []struct {
		password string
	}{
		{"password123"},
		{"simplepassword"},
		{"verylongpasswordthatexceedslengthlimit"},
		{""},
	}

	for _, tc := range cases {
		hash, err := HashPassword(tc.password)
		if err != nil {
			t.Errorf("HashPassword(%q) => unexpected error %v", tc.password, err)
		}

		if valid := CheckPasswordHash(tc.password, hash); !valid {
			t.Errorf("CheckPasswordHash(%q, %q) => got %v, want %v", tc.password, hash, valid, true)
		}

		if valid := CheckPasswordHash("wrongpassword", hash); valid {
			t.Errorf("CheckPasswordHash(%q, %q) => got %v, want %v", "wrongpassword", hash, valid, false)
		}
	}
}
