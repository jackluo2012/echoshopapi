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
 *     Initial: 2017/07/19       Li Zebang
 *     Modify : 2017/07/20       Yu Yi
 *     Modify : 2017/07/20       Yang Zhengtian
 *     Modify : 2017/07/27       Li Zebang
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
	"ShopApi/utility"
)

func AddAddress(c echo.Context) error {
	var (
		err        error
		addAddress models.AddressJSON
	)

	if err = c.Bind(&addAddress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrAddAddressInvalidParams, err.Error())
	}

	if err = c.Validate(addAddress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrAddAddressInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	addAddress.UserID = session.Get(general.SessionUserID).(uint64)

	err = models.AddressService.AddAddress(&addAddress)
	if err != nil {
		log.Logger.Error("[ERROR] AddAddress AddAddress:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] AddAddress: UserID %d", addAddress.UserID)

	return c.JSON(errcode.ErrSucceed, general.NewMessage(errcode.AddAddressSucceed))
}

func ChangeAddress(c echo.Context) error {
	var (
		err           error
		changeAddress models.AddressJSON
	)

	if err = c.Bind(&changeAddress); err != nil {
		log.Logger.Error("[ERROR] ChangeAddress Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrChangeAddressInvalidParams, err.Error())
	}

	if err = c.Validate(changeAddress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrChangeAddressInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID := session.Get(general.SessionUserID).(uint64)

	err = models.AddressService.FindAddress(changeAddress.ID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] ChangeAddress FindAddressByAddressID: Not Found", err)

			return general.NewErrorWithMessage(errcode.ErrChangeAddressNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] ChangeAddress FindAddressByAddressID: MySQL ERROR", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	err = models.AddressService.ChangeAddress(&changeAddress, userID)
	if err != nil {
		log.Logger.Error("[ERROR] ChangeAddress ChangeAddress:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] ChangeAddress: UserID %d", userID)

	return c.JSON(errcode.ChangeAddressSucceed, general.NewMessage(errcode.ChangeAddressSucceed))
}

func GetAddress(c echo.Context) error {
	var (
		err         error
		userID      uint64
		addressList *[]models.AddressJSON
	)

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID = session.Get(general.SessionUserID).(uint64)

	addressList, err = models.AddressService.GetAddressByUserID(userID)
	if err != nil {
		log.Logger.Error("[ERROR] GetAddress GetAddressByUserID ", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	if len(*addressList) == 0 {
		err = errors.New("Address Not Found")

		log.Logger.Error("[ERROR] GetAddress GetAddressByUserID:", err)

		return general.NewErrorWithMessage(errcode.ErrGetAddressNotFound, err.Error())
	}

	log.Logger.Info("[SUCCEED] GetAddress: UserID %d", userID)

	return c.JSON(errcode.GetAddressSucceed, general.NewMessageWithData(errcode.GetAddressSucceed, *addressList))
}

func AlterDefault(c echo.Context) error {
	var (
		err          error
		alterAddress *models.AddressID
		userID       uint64
	)

	if err = c.Bind(&alterAddress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrAlterDefaultInvalidParams, err.Error())
	}

	if err = c.Validate(alterAddress); err != nil {
		log.Logger.Error("[ERROR] AddAddress Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrAlterDefaultInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID = session.Get(general.SessionUserID).(uint64)

	err = models.AddressService.FindAddress(alterAddress.ID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] AlterDefault FindAddressByAddressID: Not Found", err)

			return general.NewErrorWithMessage(errcode.ErrAlterDefaultNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] AlterDefault FindAddressByAddressID: MySQL ERROR", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	err = models.AddressService.AlterAddress(alterAddress, userID)
	if err != nil {
		log.Logger.Error("[ERROR] AlterDefault AlterAddress:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] AlterDefault: UserID %d", userID)

	return c.JSON(errcode.AlterDefaultSucceed, general.NewMessage(errcode.AlterDefaultSucceed))
}

func DeleteAddress(c echo.Context) error {
	var (
		err           error
		deleteAddress *models.AddressID
		userID        uint64
	)

	if err = c.Bind(&deleteAddress); err != nil {
		log.Logger.Error("[ERROR] DeleteAddress Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrDeleteAddressInvalidParams, err.Error())
	}

	if err = c.Validate(deleteAddress); err != nil {
		log.Logger.Error("[ERROR] DeleteAddress Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrDeleteAddressInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID = session.Get(general.SessionUserID).(uint64)

	err = models.AddressService.FindAddress(deleteAddress.ID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] DeleteAddress FindAddressByAddressID: Not Found", err)

			return general.NewErrorWithMessage(errcode.ErrDeleteAddressNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] DeleteAddress FindAddressByAddressID: MySQL ERROR", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	err = models.AddressService.DeleteAddress(deleteAddress)
	if err != nil {
		log.Logger.Debug(err.Error())
	}

	log.Logger.Info("[SUCCEED] DeleteAddress: UserID %d", userID)

	return c.JSON(errcode.DeleteAddressSucceed, general.NewMessage(errcode.DeleteAddressSucceed))
}
