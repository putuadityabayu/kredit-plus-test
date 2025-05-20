/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package response

const (
	ErrMissingAuthorization = "MISSING_AUTHORIZATION"
	ErrJwt                  = "JWT_ERROR"

	MsgMissingAuthorization = "Missing authorization token"
	MsgExpiredToken         = "The token provided has expired"
	MsgInvalidToken         = "The token provided is invalid"
)

func Authorization(httpCode int, code string, msg string, err ...error) ErrorResponse {
	return NewError(httpCode, code, msg, nil, err...)
}
