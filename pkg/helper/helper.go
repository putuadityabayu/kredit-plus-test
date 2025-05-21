/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package helper

import "context"

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
