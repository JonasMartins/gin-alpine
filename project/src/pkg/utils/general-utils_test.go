package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneralUtilsUsecases(t *testing.T) {
	description := "If the general-utils scenarios are working correctly"
	defer func() {
		log.Printf("Test: %s\n", description)
		log.Println("Deferred tearing down.")
	}()

	t.Run("should PickUniqueRandomNumbers return success", func(t *testing.T) {
		start := 10
		end := 20
		n := 5
		result, err := PickUniqueRandomNumbers(start, end, n)
		assert.NoError(t, err)
		assert.Len(t, result, n)
		seen := make(map[int]bool)
		for _, v := range result {
			assert.GreaterOrEqual(t, v, start)
			assert.LessOrEqual(t, v, end)

			// ensure uniqueness
			assert.False(t, seen[v], "number repeated: %d", v)
			seen[v] = true
		}
	})

	t.Run("should PickUniqueRandomNumbers return full range", func(t *testing.T) {
		start := 1
		end := 5
		n := 5
		result, err := PickUniqueRandomNumbers(start, end, n)
		assert.NoError(t, err)
		assert.Len(t, result, n)
		seen := make(map[int]bool)
		for _, v := range result {
			seen[v] = true
		}
		for i := start; i <= end; i++ {
			assert.True(t, seen[i], "missing number %d", i)
		}
	})

	t.Run("should PickUniqueRandomNumbers error too large", func(t *testing.T) {
		start := 1
		end := 5
		n := 6
		result, err := PickUniqueRandomNumbers(start, end, n)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should PickUniqueRandomNumbers error with wrong start and finish", func(t *testing.T) {
		result, err := PickUniqueRandomNumbers(10, 5, 2)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
