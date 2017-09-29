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
 *     Initial: 2017/07/20        Yusan Kurban
 */

package handler

import (
	"errors"

	"github.com/labstack/echo"
	jwt "github.com/dgrijalva/jwt-go"
	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"github.com/jinzhu/gorm"
	"ShopApi/models"
	"time"
	"ShopApi/config"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func MustLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get(general.JWTHEADERCODE)
		if tokenString == "" {
			err := errors.New("User Must Login.")
			log.Logger.Error("[ERROR] MustLogin:", err)
			return general.NewErrorWithMessage(errcode.ErrMustLogin, err.Error())
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(c.Get(general.JWTTokenKey).(string)), nil
		})
		if err != nil {
			log.Logger.Error("[ERROR] MustLogin:", err)
			return general.NewErrorWithMessage(errcode.ErrMustLogin, err.Error())
		}
		claims := token.Claims.(jwt.MapClaims)
		log.Logger.Info("Uid:%d", claims["uid"])
		c.Set(general.SessionUserID, uint64(claims[general.SessionUserID].(float64)))

		return next(c)
	}
}
/**
 * 微信登陆  玩法是，先微信客户端调我，传给我一个code 我再拿code 传给微信,确认是微信登陆
 */
// WeAppLogin 微信小程序登录
func WeAppLogin(c echo.Context) {
	//获取服务端传过来的code
	code := c.FormValue("code")
	if code =="" {
		err := errors.New("微信验证登陆code不能为空.")
		log.Logger.Error("[ERROR] 微信验证登陆code不能为空 :", err)
		return general.NewErrorWithMessage(errcode.ErrLoginInvalidParams, err.Error())
	}
	appID         := config.WeAppConfig.AppID
	secret        := config.WeAppConfig.Secret
	CodeToSessURL := config.WeAppConfig.CodeToSessURL
	CodeToSessURL  = strings.Replace(CodeToSessURL, "{appid}",  appID,  -1)
	CodeToSessURL  = strings.Replace(CodeToSessURL, "{secret}", secret, -1)
	CodeToSessURL  = strings.Replace(CodeToSessURL, "{code}",   code,   -1)
	resp, err := http.Get(CodeToSessURL)
	if err != nil {
		log.Logger.Error("[ERROR] 微信发起code验证出错:", err)
		return general.NewErrorWithMessage(errcode.ErrLoginInvalidParams, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err := errors.New("微信验证登陆code不能为空.")
		log.Logger.Error("[ERROR] 微信发起返回status !=200:", err)
		return general.NewErrorWithMessage(errcode.ErrLoginInvalidParams, err.Error())
	}
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Logger.Error("[ERROR] 微信发起返回json:", err)
		return general.NewErrorWithMessage(errcode.ErrLoginInvalidParams, err.Error())
	}

	if _, ok := data["session_key"]; !ok {
		err := errors.New("微信验证登陆session_key 不存在.")
		log.Logger.Error("[ERROR] 微信验证登陆session_key 不存在", err)
		return general.NewErrorWithMessage(errcode.ErrLoginInvalidParams, err.Error())
	}
	var openID string
	var sessionKey string
	openID     = data["openid"].(string)
	sessionKey = data["session_key"].(string)


	session   := ctx.Session()
	session.Set("weAppOpenID",     openID)
	session.Set("weAppSessionKey", sessionKey)

	resData := iris.Map{}
	resData[config.ServerConfig.SessionID] = session.ID()
	ctx.JSON(iris.StatusOK, iris.Map{
		"errNo" : model.ErrorCode.SUCCESS,
		"msg"   : "success",
		"data"  : resData,
	})

	/*

	var openID string
	var sessionKey string
	openID     = data["openid"].(string)
	sessionKey = data["session_key"].(string)
	session   := ctx.Session()
	session.Set("weAppOpenID",     openID)
	session.Set("weAppSessionKey", sessionKey)

	resData := iris.Map{}
	resData[config.ServerConfig.SessionID] = session.ID()
	ctx.JSON(iris.StatusOK, iris.Map{
		"errNo" : model.ErrorCode.SUCCESS,
		"msg"   : "success",
		"data"  : resData,
	})
	*/
}
/**
 * 普通 的 登陆
 */
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
	token :=jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims[general.SessionUserID] = userID
	claims[general.SessionExp] = time.Now().Add(time.Minute * 5).Unix()
	tokenStr, err := token.SignedString([]byte(c.Get(general.JWTTokenKey).(string)))
	//session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	//session.Set(general.SessionUserID, userID)
	if err != nil {
		err = errors.New("Login json web token make error.")
		log.Logger.Error("[ERROR] Login Login:", err)
		return general.NewErrorWithMessage(errcode.ErrLoginInvalidPassword, err.Error())
	}
	log.Logger.Info("[SUCCEED] Login: User ID %d", userID)
	log.Logger.Info("[SUCCEED] Login: User Token %s", tokenStr)

	return c.JSON(errcode.LoginSucceed, general.RespToken{Code:errcode.LoginSucceed,Token:tokenStr})
}