/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package response

const (
	ErrUnauthorized = "UNAUTHORIZED"

	MsgMissingAuthorization = "Missing authorization token"
	MsgInvalidToken         = "The token provided is invalid"
	MsgLoginRequired        = "Authentication required. Please provide a valid token"
)

func Authorization(httpCode int, code string, msg string, err ...error) ErrorResponse {
	return NewError(httpCode, code, msg, nil, err...)
}
