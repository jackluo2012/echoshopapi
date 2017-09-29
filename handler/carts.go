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
 *     Initial: 2017/07/21     Zhu Yaqiang
 *     Modify : 2017/07/22     Xu Haosheng
 *     Modify : 2017/07/23     Wang Ke
 *     Modify : 2017/07/24     Ma Chao
 *	   Modify : 2017/08/10     Zhang Zizhao
 *     Modify : 2017/08/12     Yu Yi
 */

package handler

import (
	"io"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
	"ShopApi/utility"
)

func CreateCarts(c echo.Context) error {
	var (
		err            error
		carts          models.CartPutIn
		ProInfo        *models.ProductInfo
	)

	if err = c.Bind(&carts); err != nil {
		log.Logger.Error("[ERROR] Bind with error:", err)

		return general.NewErrorWithMessage(errcode.ErrCartPutInInvalidParams, err.Error())
	}

	if err = c.Validate(&carts); err != nil {
		log.Logger.Error("[ERROR] Create Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrCartPutInInvalidParams, err.Error())
	}

	ProInfo, err = models.ProductService.GetProInfo(carts.ProductID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] CartsCreate: Product doesn't exist", err)

			return general.NewErrorWithMessage(errcode.ErrCartPutInProductNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] Get Information with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	if err != nil {
		log.Logger.Error("[ERROR] Get Information with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMongo, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID := session.Get(general.SessionUserID).(uint64)

	err = models.CartsService.CreateCarts(&carts, userID, ProInfo.Name, ProInfo.Price)
	if err != nil {
		log.Logger.Error("[ERROR] Mysql error with CartCreate:", err)

		return general.NewErrorWithMessage(errcode.ErrCartPutInDatabase, err.Error())
	}

	log.Logger.Info("[SUCCEED] CartsPutIn name:%s", ProInfo.Name)

	return c.JSON(errcode.CreateSucceed, general.NewMessage(errcode.CreateSucceed))
}

func CartsDelete(c echo.Context) error {
	var (
		err  error
		Data models.CartsDelete
	)

	if err = c.Bind(&Data); err != nil && err != io.EOF {
		log.Logger.Error("[ERROR] CartsDelete Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrCartsDeleteErrInvalidParams, err.Error())
	}

	if err = c.Validate(Data); err != nil {
		log.Logger.Error("[ERROR] CartDelete Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrCartsDeleteErrInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID := session.Get(general.SessionUserID).(uint64)

	err = models.CartsService.CartsDelete(&Data, userID)
	if err != nil {
		log.Logger.Error("[ERROR] CartsDelete CartsDelete:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] CartsDelete: UserID %d", userID)

	return c.JSON(errcode.CartsDeleteSucceed, general.NewMessage(errcode.CartsDeleteSucceed))
}

func AlterCartPro(c echo.Context) error {
	var (
		err         error
		cartProduct *models.CartPutIn
	)

	if err = c.Bind(&cartProduct); err != nil {
		log.Logger.Error("[ERROR] AlterCartPro Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrAlterCartInvalidParams, err.Error())
	}

	if err = c.Validate(cartProduct); err != nil {
		log.Logger.Error("[ERROR] AlterCartPro Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrAlterCartInvalidParams, err.Error())
	}

	err = models.CartsService.AlterCartPro(cartProduct)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] AlterCartPro: Product doesn't exist!", err)

			return general.NewErrorWithMessage(errcode.ErrAlterCartProductNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] AlterCartPro with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] AlterCartPro %v")

	return c.JSON(errcode.AlterCartSucceed, general.NewMessage(errcode.AlterCartSucceed))
}

func CartsBrowse(c echo.Context) error {
	var (
		err    error
		output *[]models.ConCarts
	)

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID := session.Get(general.SessionUserID).(uint64)

	output, err = models.CartsService.CartsBrowse(userID)
	if err != nil {
		log.Logger.Error("[ERROR] CartsBrowse", err)

		return general.NewErrorWithMessage(errcode.ErrBrowseCartNotFound, err.Error())
	}

	log.Logger.Info("[SUCCEED] CartsBrowse %v")

	return c.JSON(errcode.BrowseCartSucceed, general.NewMessageWithData(errcode.BrowseCartSucceed, output))
}
