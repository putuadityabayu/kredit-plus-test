/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package helper

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// GetValueContext retrieves a value from context or returns the provided default.
func GetValueContext[T any](ctx context.Context, key any, def ...T) T {
	val := ctx.Value(key)
	if v, ok := val.(T); ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	var zero T
	return zero
}

// GenerateContractNumber generates a unique contract number based on a timestamp.
// Format: TRX-YYYYMMDDHHMMSSmmm-XXXX (TRX-DateHourMinuteSecondMillisecond-Random4Digit)
// Example: TRX-20250522214530123-8765
func GenerateContractNumber() string {
	// Get the current time
	now := time.Now()

	// Format date and time up to milliseconds
	// "20060102150405.000" is Go's reference layout for YYYYMMDDHHmmss.milliseconds
	// We will remove the dot to get only digits
	timestampStr := now.Format("20060102150405.000")
	timestampStr = timestampStr[:len(timestampStr)-4] + timestampStr[len(timestampStr)-3:] // Remove the dot

	// Add a 4-digit random component to reduce the risk of collisions
	// in the same millisecond or if the system has low time precision.
	// rand.Seed() is not needed in Go 1.20+ when using rand.IntN
	randomNumber := rand.Intn(10000) // Generates a number between 0 and 9999

	// Combine all components
	return fmt.Sprintf("TRX-%s-%04d", timestampStr, randomNumber)
}

// GetTransactionAmount helper to count admin fee, interest amount and installment amount from otr and tenor
//
// Real conditions will be based on company policy
func GetTransactionAmount(otr float64, tenor int) (interestAmount, adminFee float64, installmentAmount float64, totalAmount float64) {
	interestAmount = 2 * otr / 100              // Anggap 2%
	adminFee = (otr + interestAmount) * 1 / 100 // Anggap 1%
	totalAmount = otr + interestAmount + adminFee
	installmentAmount = totalAmount / float64(tenor)
	return
}
