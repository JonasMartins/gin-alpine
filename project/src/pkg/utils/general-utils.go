// Package utils
package utils

import (
	"bytes"
	"errors"
	"log"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// GetResponseRecorder simulates a responseWriter for testing
func GetResponseRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

// ExecuteRequest creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func ExecuteRequest(req *http.Request, router http.Handler) *httptest.ResponseRecorder {
	req.RemoteAddr = "localhost:3000"
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

// CheckResponseCode is a simple utility to check the response code
// of the response
func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func GenerateNewUUID() ([]byte, error) {
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(newUUID, []byte("\n")), nil
}

func FatalResult(s string, err error) {
	log.Fatalf("%s %v", s, err)
}

func ValidatePassword(pass string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}

func ConvertIntSliceToInt32Slice(s *[]int) *[]int32 {
	res := []int32{}
	for _, x := range *s {
		res = append(res, int32(x))
	}
	return &res
}

func GetEnvDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// BallPairs extract all the unique combinations
// from a slice of int
func BallPairs(balls []int) [][2]int {
	var pairs [][2]int
	for i := range balls {
		for j := i + 1; j < len(balls); j++ {
			pairs = append(pairs, [2]int{balls[i], balls[j]})
		}
	}
	return pairs
}

func PickUniqueRandomNumbers(start, end, pickNumber int) ([]int, error) {
	if start > end {
		return nil, errors.New("start must be less than or equal to end")
	}

	rangeSize := end - start + 1
	if pickNumber > rangeSize {
		return nil, errors.New("pickNumber is greater than the available range")
	}

	// Create slice with all numbers in range
	numbers := make([]int, rangeSize)
	for i := 0; i < rangeSize; i++ {
		numbers[i] = start + i
	}

	// Shuffle using modern global rand (auto-seeded)
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	return numbers[:pickNumber], nil
}

func CalculateDaysAmount(date time.Time) int {
	today := time.Now()

	// Normalize both to midnight
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

	diff := today.Sub(date).Hours() / 24

	if diff < 0 {
		return 0
	}

	return int(diff)
}
