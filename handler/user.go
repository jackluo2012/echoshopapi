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
 *     Initial: 2017/07/18        Yusan Kurban
 *	   Modify : 2017/07/19		  Ai Hao
 *	   Modify : 2017/07/20        Zhang Zizhao
 *     Modify : 2017/07/21        Xu Haosheng
 *	   Modify : 2017/07/21        Yang Zhengtian
 *     Modify : 2017/07/21        Ma Chao
 */

package handler

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
	"ShopApi/utility"
)

func Register(c echo.Context) error {
	var (
		err      error
		register models.Register
	)

	if err = c.Bind(&register); err != nil {
		log.Logger.Error("[ERROR] Register Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrRegisterInvalidParams, err.Error())
	}

	if err = c.Validate(register); err != nil {
		log.Logger.Error("[ERROR] Register Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrRegisterInvalidParams, err.Error())
	}

	match := utility.IsValidPhone(*register.Mobile)
	if !match {
		log.Logger.Error("[ERROR] Register IsValidPhone: InvalidPhone", err)

		return general.NewErrorWithMessage(errcode.ErrRegisterInvalidParams, err.Error())
	}

	err = models.UserService.Register(register.Mobile, register.Pass)
	if err != nil {
		if strings.Contains(err.Error(), general.DuplicateEntry) {
			log.Logger.Error("[ERROR] Register Register: Mobile Duplicate", err)

			return general.NewErrorWithMessage(errcode.ErrRegisterUserDuplicate, err.Error())
		}

		log.Logger.Error("[ERROR] Register Register: Mysql Error", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] Register: Mobile %s", *register.Mobile)

	return c.JSON(errcode.RegisterSucceed, general.NewMessage(errcode.RegisterSucceed))
}
/*
func Login(c echo.Context) error {
	var (
		err   error
		login models.Login
	)

	if err = c.Bind(&login); err != nil {
		log.Logger.Error("[ERROR] Login Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrLoginInvalidParams, err.Error())
	}

	if err = c.Validate(login); err != nil {
		log.Logger.Error("[ERROR] Login Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrLoginInvalidParams, err.Error())
	}

	flag, userID, err := models.UserService.Login(login.Mobile, login.Pass)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] Login Login: User doesn't exist", err)

			return general.NewErrorWithMessage(errcode.ErrLoginUserNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] Login Login: Mysql Error", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	if !flag {
		err = errors.New("Login mobile and password not match.")

		log.Logger.Error("[ERROR] Login Login:", err)

		return general.NewErrorWithMessage(errcode.ErrLoginInvalidPassword, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	session.Set(general.SessionUserID, userID)

	log.Logger.Info("[SUCCEED] Login: User ID %d", userID)

	return c.JSON(errcode.LoginSucceed, general.NewMessage(errcode.LoginSucceed))
}
*/
func Logout(c echo.Context) error {
	var (
		err error
	)

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID := session.Get(general.SessionUserID)

	err = session.Delete(general.SessionUserID)
	if err != nil {
		log.Logger.Error("[ERROR] Logout:", err)

		return general.NewErrorWithMessage(errcode.ErrLogout, err.Error())
	}

	log.Logger.Info("[SUCCEED] Logout: User ID %d", userID)

	return c.JSON(errcode.LogoutSucceed, general.NewMessage(errcode.LogoutSucceed))
}

func GetUserInfo(c echo.Context) error {
	var (
		err    error
		output = new(models.UserGet)
//		avatar = new(models.UserAvatar)
	)

	//session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	//userID := session.Get(general.SessionUserID).(uint64)

	userID := c.Get(general.SessionUserID).(uint64) //uint64(1)

	output, err = models.UserService.GetUserInfo(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] GetUserInfo : User Information Not Found", err)

			return general.NewErrorWithMessage(errcode.ErrGetUserInfoInvalidParams, err.Error())
		}

		log.Logger.Error("[ERROR] GetUserInfo : Mysql Error", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}
	/*
	avatar, err = models.UserService.GetUserAvatar(userID)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		log.Logger.Error("[ERROR] GetUserInfo GetUserAvatar: Mongo Error", err)

		return general.NewErrorWithMessage(errcode.ErrMongo, err.Error())
	}

	output.Avatar = avatar.Avatar
	*/
	log.Logger.Info("[SUCCEED] GetUserInfo: User ID %d", userID)

	return c.JSON(errcode.GetUserInfoSucceed, general.NewMessageWithData(errcode.GetUserInfoSucceed, *output))
}

func ChangeUserInfo(c echo.Context) error {
	var (
		err  error
		info models.ChangeUserInfo
	)

	if err = c.Bind(&info); err != nil {
		log.Logger.Error("[ERROR] ChangeUserInfo Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrChangeUserInfoInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID := session.Get(general.SessionUserID).(uint64)

	err = models.UserService.ChangeUserInfo(&info, userID)
	if err != nil {
		log.Logger.Error("[ERROR] ChangeUserInfo ChangeUserInfo:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] ChangeUserInfo: User ID %d", userID)

	return c.JSON(errcode.ChangeUserInfoSucceed, general.NewMessage(errcode.ChangeUserInfoSucceed))
}

func ChangeUserAvatar(c echo.Context) error {
	var (
		err    error
		avatar models.UserAvatar
	)

	if err = c.Bind(&avatar); err != nil {
		log.Logger.Error("[ERROR] ChangeUserAvatar Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrChangeUserAvatarInvalidParams, err.Error())
	}

	if err = c.Validate(&avatar); err != nil {
		log.Logger.Error("[ERROR] ChangeUserAvatar Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrChangeUserAvatarInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	avatar.UserID = session.Get(general.SessionUserID).(uint64)

	err = models.UserService.ChangeUserAvatar(&avatar)
	if err != nil {
		log.Logger.Error("[ERROR] ChangeUserAvatar ChangeUserAvatar:", err)

		return general.NewErrorWithMessage(errcode.ErrMongo, err.Error())
	}

	log.Logger.Info("[SUCCEED] ChangeUserAvatar: User ID %d", avatar.UserID)

	return c.JSON(errcode.ChangeUserAvatarSucceed, general.NewMessage(errcode.ChangeUserAvatarSucceed))
}

func ChangePhone(c echo.Context) error {
	var (
		err         error
		changePhone models.ChangePhone
	)
	if err = c.Bind(&changePhone); err != nil {
		log.Logger.Error("[ERROR] ChangePhone Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrChangePhoneInvalidParams, err.Error())
	}

	if err = c.Validate(changePhone); err != nil {
		log.Logger.Error("[ERROR] ChangePhone Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrChangePhoneInvalidParams, err.Error())
	}

	match := utility.IsValidPhone(changePhone.Phone)
	if !match {
		log.Logger.Error("[ERROR] ChangePhone IsValidPhone: Invalid Phone", err)

		return general.NewErrorWithMessage(errcode.ErrChangePhoneInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID := session.Get(general.SessionUserID).(uint64)

	err = models.UserService.ChangePhone(userID, changePhone.Phone)
	if err != nil {
		if strings.Contains(err.Error(), general.DuplicateEntry) {
			log.Logger.Error("[ERROR] ChangePhone ChangePhone: Mobile Duplicate", err)

			return general.NewErrorWithMessage(errcode.ErrChangePhoneDuplicate, err.Error())
		}

		log.Logger.Error("[ERROR] ChangePhone ChangePhone: Mysql Error", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	err = session.Delete(general.SessionUserID)
	if err != nil {
		log.Logger.Error("[ERROR] ChangePhone Delete:", err)

		return general.NewErrorWithMessage(errcode.ErrLogout, err.Error())
	}

	log.Logger.Info("[SUCCEED] ChangePhone: User ID %d", userID)

	return c.JSON(errcode.ChangePhoneSucceed, general.NewMessage(errcode.ChangePhoneSucceed))
}

func ChangePassword(c echo.Context) error {
	var (
		changePassword models.ChangePassword
		userID         uint64
		err            error
	)

	if err = c.Bind(&changePassword); err != nil {
		log.Logger.Error("[ERROR] ChangePassword Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrChangePasswordInvalidParams, err.Error())
	}

	if err = c.Validate(changePassword); err != nil {
		log.Logger.Error("[ERROR] ChangePassword Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrChangePasswordInvalidParams, err.Error())
	}

	if *changePassword.Password == *changePassword.NewPass {
		err = errors.New("The password hasn't change.")

		log.Logger.Error("[ERROR] ChangePassword:", err)

		return general.NewErrorWithMessage(errcode.ErrChangePasswordInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	userID = session.Get(general.SessionUserID).(uint64)

	ok, err := models.UserService.ChangePassword(&changePassword, userID)
	if err != nil {
		log.Logger.Error("[ERROR] ChangePassword ChangePassword: Mysql Error", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	if !ok {
		err = errors.New("Password is wrong.")

		log.Logger.Error("[ERROR]", err)

		return general.NewErrorWithMessage(errcode.ErrChangePasswordInvalidParams, err.Error())
	}

	err = session.Delete(general.SessionUserID)
	if err != nil {
		log.Logger.Error("[ERROR] ChangePhone Delete:", err)

		return general.NewErrorWithMessage(errcode.ErrLogout, err.Error())
	}

	log.Logger.Info("[SUCCEED] ChangePassword: User ID %d", userID)

	return c.JSON(errcode.ChangePasswordSucceed, general.NewMessage(errcode.ChangePasswordSucceed))
}
