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
 *     Initial: 2017/08/09       Zhang Zizhao
 */

package errcode

const (
	//PutIn
	CreateSucceed                = 0x0
	ErrCartPutInInvalidParams = 0x1
	ErrCartPutInProductNotFound  = 0x2
	ErrCartPutInDatabase         = 0x3

	//Delete
	CartDeleteSucceed             = 0x0
	ErrCartDeleteInvalidParams = 0x1
	ErrCartDeleteProductNotFound  = 0x3

	// CartsDelete
	CartsDeleteSucceed = 0x0
	ErrCartsDeleteErrInvalidParams = 0x1

	//Alter
	AlterCartSucceed            = 0x0
	ErrAlterCartInvalidParams   = 0x1
	ErrAlterCartProductNotFound = 0x2

	//Browse
	BrowseCartSucceed      = 0x0
	ErrBrowseInvalidParams = 0x1
	ErrBrowseCartNotFound  = 0x2
)
