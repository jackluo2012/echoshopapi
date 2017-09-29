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
 */

package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"

	"ShopApi/log"
	"ShopApi/orm"
	"ShopApi/server/router"
	jwt "github.com/dgrijalva/jwt-go"
	"ShopApi/general"
)

var (
	server *echo.Echo
)

type JwtCustomClaims struct {
	UserID   uint64 `sql:"primary_key" gorm:"column:userid" json:"userid"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Sex      uint8  `json:"sex"`
	jwt.StandardClaims
}

func startServer() {
	server = echo.New()
	server.Use(middleware.CORS())
	server.Use(middleware.Recover())
	server.Use(middleware.Logger())
	//加入jwt
	/*
	server.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:configuration.tokenKey,
		TokenLookup: "header:x-access-token",

		keyFunc:func() (interface{}, error){
			return nil,nil
		},
	}))*/

	server.HTTPErrorHandler = general.EchoRestfulErrorHandler
	server.Validator = general.NewEchoValidator()
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error{
			c.Set(general.JWTTokenKey,configuration.tokenKey)
			c.Set("configure",configuration)
			return  next(c)
		}
	})
	router.InitRouter(server)
	log.Logger.Info("Router already init %v")
	log.Logger.Fatal(server.Start(configuration.address))
}

func init() {
	readConfiguration()
	initMysql()
//	InitMetal()
	startServer()
}

func initMysql() {
	user := configuration.mysqlUser
	pass := configuration.mysqlPass
	url := configuration.mysqlHost
	port := configuration.mysqlPort
	sqlName := configuration.mysqlDb

	conf := fmt.Sprintf(user + ":" + pass + "@" + "tcp(" + url + port + ")/" + sqlName + "?charset=utf8&parseTime=True&loc=Local")

	orm.InitOrm(conf)
}

func InitMetal() {
	var err error
	url := configuration.MgoUrl

	orm.MDSession, err = mgo.DialWithTimeout(url, time.Second)

	if err != nil {
		panic(err)
	}

	log.Logger.Info("the MongoDB of %s connected!", orm.MD)

	orm.MDSession.SetMode(mgo.Monotonic, true)
}
