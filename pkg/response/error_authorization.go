/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package response

const (
	ErrUnauthorized = "UNAUTHORIZED"
	ErrForbidden    = "FORBIDDEN"

	MsgMissingAuthorization = "Missing authorization token"
	MsgInvalidToken         = "The token provided is invalid"
	MsgForbidden            = "You don't have permission to access this resource"
)

func Authorization(httpCode int, code string, msg string, err ...error) ErrorResponse {
	return NewError(httpCode, code, msg, nil, err...)
}
