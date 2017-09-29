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
 *     Modify: 2017/07/19         Yang Zhengtian
 *     Modify: 2017/07/20         Yang Zhengtain
 */

package router

import (
	"github.com/labstack/echo"

	"ShopApi/handler"
)

func InitRouter(server *echo.Echo) {
	if server == nil {
		panic("[InitRouter], server couldn't be nil")
	}

	// user
	server.POST("/api/v1/user/register", handler.Register)
	server.POST("/api/v1/user/login", handler.Login)
	server.GET("/api/v1/user/logout", handler.Logout, handler.MustLogin)
	server.GET("/api/v1/user/getinfo", handler.GetUserInfo, handler.MustLogin)
	server.POST("/api/v1/user/changeavatar", handler.ChangeUserAvatar, handler.MustLogin)
	server.POST("/api/v1/user/changeinfo", handler.ChangeUserInfo, handler.MustLogin)
	server.POST("/api/v1/user/changephone", handler.ChangePhone, handler.MustLogin)
	server.POST("/api/v1/user/changepass", handler.ChangePassword, handler.MustLogin)

	// address
	server.POST("/api/v1/address/add", handler.AddAddress, handler.MustLogin)
	server.POST("/api/v1/address/change", handler.ChangeAddress, handler.MustLogin)
	server.GET("/api/v1/address/get", handler.GetAddress, handler.MustLogin)
	server.POST("/api/v1/address/alter", handler.AlterDefault, handler.MustLogin)
	server.POST("/api/v1/address/delete", handler.DeleteAddress, handler.MustLogin)

	// products
	server.POST("/api/v1/product/create", handler.CreateProduct)
	server.GET("/api/v1/product/gethomepage", handler.GetProductList)
	server.POST("/api/v1/product/getlistbycategory", handler.GetProductListByCategory)
	server.POST("/api/v1/product/getinfo", handler.GetProInfo)
	server.POST("/api/v1/product/changestatus", handler.ChangeProStatus)
	server.POST("/api/v1/product/changecate", handler.ChangeCategory)
	server.GET("/api/v1/product/getmypage",handler.GetMyPage)

	// orders
	server.POST("/api/v1/orders/create", handler.CreateOrder, handler.MustLogin)
	server.POST("/api/v1/orders/getone", handler.GetOneOrder, handler.MustLogin)
	server.POST("/api/v1/orders/changestatus", handler.ChangeStatus)
	server.POST("/api/v1/orders/get", handler.GetOrders, handler.MustLogin)

	// category
	server.POST("/api/v1/category/create", handler.CreateCategory)
	server.GET("/api/v1/category/get", handler.GetCategory)

	// carts
	server.POST("/api/v1/carts/create", handler.CreateCarts, handler.MustLogin)
	server.POST("/api/v1/carts/delete", handler.CartsDelete, handler.MustLogin)
	server.POST("/api/v1/carts/alter", handler.AlterCartPro, handler.MustLogin)
	server.GET("/api/v1/carts/getlist", handler.CartsBrowse, handler.MustLogin)
}
