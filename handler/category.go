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
 *     Initial: 2017/07/21        Yang Zhengtian
 *     Modify : 2017/07/21        Li Zebang
 *     Modify : 2017/07/29        Li Zebang
 */

package handler

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
)

func CreateCategory(c echo.Context) error {
	var (
		err            error
		createCategory models.CreateCategory
	)

	if err = c.Bind(&createCategory); err != nil {
		log.Logger.Error("[ERROR] CreateCategory Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	if err = c.Validate(createCategory); err != nil {
		log.Logger.Error("[ERROR] CreateCategory Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidParams, err.Error())
	}

	if createCategory.PID != 0 {
		err = models.CategoryService.CheckPID(createCategory.PID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Logger.Error("[ERROR] CreateCategory CheckPID: PID Not Found", err)

				return general.NewErrorWithMessage(errcode.ErrNotFound, err.Error())
			}
			log.Logger.Error("[ERROR] CreateCategory CheckPID: MySQL ERROR", err)

			return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
		}
	}

	err = models.CategoryService.CreateCategory(createCategory)
	if err != nil {
		log.Logger.Error("[ERROR] CreateCategory CreateCategory:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, nil)
}

func GetCategory(c echo.Context) error {
	var (
		err          error
		categoryList *[]models.CategoryGet
	)

	categoryList, err = models.CategoryService.GetCategory()
	if err != nil {
		log.Logger.Error("[ERROR] GetCategory GetCategory: MySQL ERROR", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	if len(*categoryList) == 0 {
		err = errors.New("[ERROR] Categories Not Found")

		log.Logger.Error("[ERROR] GetCategory GetCategory: ", err)

		return general.NewErrorWithMessage(errcode.ErrNotFound, err.Error())
	}

	return c.JSON(errcode.ErrSucceed, *categoryList)
}
