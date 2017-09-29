/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/05/14        Li Zebang
 */

package errcode

const (
	// Register
	RegisterSucceed          = 0x0
	ErrRegisterInvalidParams = 0x1
	ErrRegisterUserDuplicate = 0x2

	// Login
	LoginSucceed            = 0x0
	ErrLoginInvalidParams   = 0x1
	ErrLoginUserNotFound    = 0x2
	ErrLoginInvalidPassword = 0x3

	// Logout
	LogoutSucceed = 0x1111
	ErrLogout     = 0x1

	// GetUserInfo
	GetUserInfoSucceed          = 0x0
	ErrGetUserInfoInvalidParams = 0x1

	// ChangeUserInfo
	ChangeUserInfoSucceed          = 0x0
	ErrChangeUserInfoInvalidParams = 0x1

	// ChangeUserAvatar
	ChangeUserAvatarSucceed          = 0x0
	ErrChangeUserAvatarInvalidParams = 0x1

	// ChangePhone
	ChangePhoneSucceed          = 0x0
	ErrChangePhoneInvalidParams = 0x1
	ErrChangePhoneDuplicate     = 0x2

	// ChangePassword
	ChangePasswordSucceed          = 0x0
	ErrChangePasswordInvalidParams = 0x1
)
