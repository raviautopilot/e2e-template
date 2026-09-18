package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString generates a random alphanumeric string of the specified length.
// If length <= 0, an empty string is returned.
func GenerateRandomString(length int) string {
	if length <= 0 {
		return ""
	}
	b := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))
	for i := range b {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			// Fallback in unlikely crypto/rand failure
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
			continue
		}
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// GenerateRandomStringwithTimeStamp generates a random string prepended with current timestamp
// in the format: YYYYMMDD_HHMMSS_<randomString>.
// Matches exact casing requested by user.
func GenerateRandomStringwithTimeStamp(length int) string {
	timestamp := time.Now().Format("20060102_150405")
	if length <= 0 {
		return timestamp
	}
	return fmt.Sprintf("%s_%s", timestamp, GenerateRandomString(length))
}

// GenerateRandomStringWithTimeStamp is an idiomatic PascalCase alias for GenerateRandomStringwithTimeStamp.
func GenerateRandomStringWithTimeStamp(length int) string {
	return GenerateRandomStringwithTimeStamp(length)
}

// GenerateRandomEmail generates a unique random email for testing purposes.
// Optional domain defaults to "example.com".
func GenerateRandomEmail(domain ...string) string {
	d := "example.com"
	if len(domain) > 0 && domain[0] != "" {
		d = domain[0]
	}
	return fmt.Sprintf("user_%d_%s@%s", time.Now().UnixNano(), GenerateRandomString(6), d)
}

// GenerateRandomInt returns a random integer in the range [min, max].
func GenerateRandomInt(min, max int) int {
	if min >= max {
		return min
	}
	diff := big.NewInt(int64(max - min + 1))
	n, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return min
	}
	return min + int(n.Int64())
}

// GenerateRandomPhone generates a mock 10-digit phone number.
func GenerateRandomPhone() string {
	return fmt.Sprintf("+1%03d%03d%04d", GenerateRandomInt(200, 999), GenerateRandomInt(100, 999), GenerateRandomInt(1000, 9999))
}
