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
 *     Initial: 2017/07/21       Li Zebang
 *     Modify : 2017/07/21       Zhang Zizhao
 *	   Modify : 2017/07/21       Ai Hao
 *     Modify : 2017/07/21       Ma Chao
 */

package handler

import (
	"errors"

	"github.com/labstack/echo"

	"ShopApi/general"
	"ShopApi/general/errcode"
	"ShopApi/log"
	"ShopApi/models"
	"ShopApi/utility"
	"github.com/jinzhu/gorm"
)

func CreateOrder(c echo.Context) error {
	var (
		order models.CreateOrder
		err   error

	)

	if err = c.Bind(&order); err != nil {
		log.Logger.Error("[ERROR] Create crash with error:", err)

		return general.NewErrorWithMessage(errcode.ErrCreateOrderInvalidParams, err.Error())
	}

	if err = c.Validate(order); err != nil {
		log.Logger.Error("[ERROR] CreateOrder Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrCreateOrderInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	UserID := session.Get(general.SessionUserID).(uint64)

	err = models.OrderService.CreateOrder(UserID, order)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] CreateOrder: Address doesn't exist", err)

			return general.NewErrorWithMessage(errcode.ErrAddressNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] Mysql error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	err = models.CartsService.CartsDelete(&models.CartsDeleted, UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] CreateOrder: Carts doesn't exist", err)

			return general.NewErrorWithMessage(errcode.ErrAddressNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] Mysql error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] CartsDelete %v")
	log.Logger.Info("[SUCCEED] CreateOrder %v")

	return c.JSON(errcode.ErrCreateOrderSucceed, general.NewMessage(errcode.ErrCreateOrderSucceed))
}

func GetOrders(c echo.Context) error {
	var (
		err       error
		getOrders models.GetOrders
		orders    *[]models.OrdersGet
	)

	if err = c.Bind(&getOrders); err != nil {
		log.Logger.Error("[ERROR] Bind with error:", err)

		return general.NewErrorWithMessage(errcode.ErrGetOrdersInvalidParams, err.Error())
	}

	if err = c.Validate(getOrders); err != nil {
		log.Logger.Error("[ERROR] GetOrders Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrGetOrdersInvalidParams, err.Error())
	}

	if getOrders.Status != general.OrderUnfinished && getOrders.Status != general.OrderFinished && getOrders.Status != general.OrderGetAll && getOrders.Status != general.OrderCanceled {
		err = errors.New("[ERROR] Invalid Orders Status")

		log.Logger.Error("[ERROR] Error:", err)

		return general.NewErrorWithMessage(errcode.ErrInvalidOrdersStatus, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	getOrders.UserID = session.Get(general.SessionUserID).(uint64)

	pageStart := utility.Paging(getOrders.Page, getOrders.PageSize)

	orders, err = models.OrderService.GetOrders(&getOrders, pageStart)
	if err != nil {
		log.Logger.Error("[ERROR] Mysql error in GetOrders Function:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	if len(*orders) == 0 {
		err = errors.New("[ERROR] Orders Not Found")

		log.Logger.Error("[ERROR] Error:", err)

		return general.NewErrorWithMessage(errcode.ErrNotFound, err.Error())
	}

	return c.JSON(errcode.ErrGetOrdersSucceed, orders)
}

func GetOneOrder(c echo.Context) error {
	var (
		err    error
		order  models.GetOne
		OutPut []models.OrmOrders
	)

	if err = c.Bind(&order); err != nil {
		log.Logger.Error("[ERROR] Bind with error:", err)

		return general.NewErrorWithMessage(errcode.ErrGetOrderInvalidParams, err.Error())
	}

	session := utility.GlobalSessions.SessionStart(c.Response().Writer, c.Request())
	UserID := session.Get(general.SessionUserID).(uint64)

	OutPut, err = models.OrderService.GetOneOrder(UserID, order.ID)
	if err != nil {
		log.Logger.Error("[ERROR] GetOneOrder with error:", err)

		return general.NewErrorWithMessage(errcode.ErrNotFound, err.Error())
	}

	return c.JSON(errcode.ErrGetOrderSucceed, general.NewMessageWithData(errcode.GetProInfoSucceed, OutPut))
}

func ChangeStatus(c echo.Context) error {
	var (
		err error
		st  models.ChangeStatus
	)

	if err = c.Bind(&st); err != nil {
		log.Logger.Error("[ERROR] Input order status with error:", err)

		return general.NewErrorWithMessage(errcode.ErrChangeOrderInvalidParams, err.Error())
	}

	if st.Status != general.OrderFinished && st.Status != general.OrderUnfinished && st.Status != general.OrderCanceled && st.Status != general.OrderPaid && st.Status != general.OrderUnpaid {
		err = errors.New("[ERROR] Status InExistent")
		log.Logger.Error("", err)

		return general.NewErrorWithMessage(errcode.ErrChangeOrderInvalidParams, err.Error())
	}
	err = models.OrderService.ChangeStatus(st.OrderID, st.Status)
	if err != nil {
		log.Logger.Error("[ERROR] Change status with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	return c.JSON(errcode.ErrChangeOrderSucceed, general.NewMessage(errcode.ErrChangeOrderSucceed))
}
