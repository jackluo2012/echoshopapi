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
 *	   Modify : 2017/07/21		 Ai Hao
 *	   Modify : 2017/07/21		 Zhang Zizhao
 *     Modify : 2017/07/21       Ma Chao
 */

package models

import (
	"fmt"
	"time"

	"ShopApi/general"
	"ShopApi/orm"
	"gopkg.in/mgo.v2/bson"
)

type OrderServiceProvider struct {
}

var OrderService *OrderServiceProvider = &OrderServiceProvider{}

type Orders struct {
	ID         uint64    `sql:"auto_increment;primary_key;"json:"id"`
	UserID     uint64    `gorm:"column:userid" json:"userid"`
	AddressID  string    `gorm:"column:addressid" json:"addressid"`
	TotalPrice float64   `gorm:"column:totalprice" json:"totalprice"`
	PayWay     uint8     `gorm:"column:payway" json:"payway"`
	Freight    float64   `json:"freight"`
	Remark     string    `json:"remark"`
	Status     uint8     `json:"status"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

type OrderProduct struct {
	ID        uint64 `sql:"auto_increment;primary_key;"json:"id"`
	OrderID   uint64 `gorm:"column:orderid" json:"orderid"`
	ProductID uint64 `gorm:"column:productid" json:"productid"`
	Discount  uint8  `json:"discount"`
	Size      string `json:"size"`
	Count     uint64 `json:"count"`
	Color     string `json:"color"`
}

type OrmOrders struct {
	TotalPrice float64   `json:"totalprice" `
	Freight    float64   `json:"freight" `
	Discount   uint8     `json:"discount" `
	Name       string    `json:"name"validate:"required, alphaunicode,min=2,max=18"`
	Size       string    `json:"size" validate:"required,alphanum"`
	Count      uint64    `json:"count"`
	Color      string    `json:"color" validate:"required,alphanum"`
	Status     uint8     `json:"status"`
	PayWay     uint8     `json:"payway"`
	AddressID  string    `json:"addressid" validate:"required,numeric"`
	ProductID  uint64    `json:"productid"`
	Remark     string    `json:"remark" validate:"required,alphanum"`
	Created    time.Time `json:"created"`
	Avatar     string    `json:"avatar"`
}

type CreateOrder struct {
	AddressID    string  `json:"addressid" validate:"required"`
	TotalPrice   float64 `json:"totalprice"`
	Freight      float64 `json:"freight"`
	Remark       string  `json:"remark"`
	PayWay       uint8   `json:"payway"`
	OrderProduct []OrderPro
}

type OrderPro struct {
	ProductID uint64 `json:"productid"`
	OrderID   uint64 `json:"orderid" `
	Discount  uint8  `json:"discount"`
	Size      string `json:"size" validate:"required,alphanum"`
	Count     uint64 `json:"count"`
	Color     string `json:"color" validate:"required,alphanum"`
}

type GetOrders struct {
	UserID   uint64 `json:"userid"`
	Status   uint8  `json:"status"`
	Page     uint64 `json:"page"`
	PageSize uint64 `json:"pagesize"`
}

type GetOne struct {
	ID uint64 `json:"orderid" `
}

type OrdersGet struct {
	TotalPrice float64 `json:"totalprice"`
	Freight    float64 `json:"freight"`
	Remark     string  `json:"remark"`
	Status     uint8   `json:"status"`
}

type ChangeStatus struct {
	Status  uint8  `json:"status"`
	OrderID uint64 `json:"orderid"`
}

func (Orders) TableName() string {
	return "orders"
}

func (OrderProduct) TableName() string {
	return "orderproduct"
}
var CartsDeleted CartsDelete
func (osp *OrderServiceProvider) CreateOrder(UserID uint64, ord CreateOrder) (error) {
	var (
		err    error
	//	car    Cart
		orders Orders

	)

	db := orm.Conn

	order := Orders{
		UserID:     UserID,
		AddressID:  ord.AddressID,
		TotalPrice: ord.TotalPrice,
		Freight:    ord.Freight,
		Remark:     ord.Remark,
		Status:     general.OrderUnfinished,
		PayWay:     ord.PayWay,
		Created:    time.Now(),
		Updated:    time.Now(),
	}

	tx := db.Begin()

	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Create(&order).Error
	if err != nil {
		return err
	}

	err = tx.Where("userid = ? AND totalprice = ? AND addressid = ?", UserID, order.TotalPrice, order.AddressID).First(&orders).Error
	if err != nil {
		return err
	}

	for _, value := range ord.OrderProduct {
		OrderProduct := OrderProduct{
			OrderID:   orders.ID,
			ProductID: value.ProductID,
			Discount:  value.Discount,
			Size:      value.Size,
			Count:     value.Count,
			Color:     value.Color,
		}

		err = tx.Create(&OrderProduct).Error
		if err != nil {
			return err
		}

		add1 := CartDelete{
			ProductID:  OrderProduct.ProductID,
			Size:       OrderProduct.Size,
			Color:      OrderProduct.Color,
		}
		CartsDeleted.Data = append(CartsDeleted.Data, add1)
	}

	return err

}

func (osp *OrderServiceProvider) GetOrders(getOrders *GetOrders, pageStart uint64) (*[]OrdersGet, error) {
	var (
		sql        string
		order      Orders
		orderGet   OrdersGet
		ordersList []OrdersGet
	)

	db := orm.Conn

	if getOrders.Status == general.OrderUnfinished || getOrders.Status == general.OrderFinished || getOrders.Status == general.OrderCanceled {
		sql = fmt.Sprintf("SELECT * FROM orders WHERE userid = %d AND status = %d LIMIT %d, %d LOCK IN SHARE MODE", getOrders.UserID, getOrders.Status, pageStart, getOrders.PageSize)
	} else {
		sql = fmt.Sprintf("SELECT * FROM orders WHERE userid = %d LIMIT %d, %d LOCK IN SHARE MODE", getOrders.UserID, pageStart, getOrders.PageSize)
	}

	rows, err := db.Raw(sql).Rows()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		db.ScanRows(rows, &order)
		fmt.Println(order)
		orderGet = OrdersGet{
			TotalPrice: order.TotalPrice,
			Freight:    order.Freight,
			Remark:     order.Remark,
			Status:     order.Status,
		}
		ordersList = append(ordersList, orderGet)
	}

	return &ordersList, nil
}

func (osp *OrderServiceProvider) GetOneOrder(userID uint64, ID uint64) ([]OrmOrders, error) {
	var (
		err           error
		order         Orders
		OrderProduct  []OrderProduct
		getOrder      []OrmOrders
		productAvatar ProductImages
	)

	db := orm.Conn
	tx := db.Begin()

	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Where("id = ? AND userid = ?", ID, userID).First(&order).Error
	if err != nil {
		return getOrder, err
	}

	add1 := OrmOrders{
		TotalPrice: order.TotalPrice,
		Freight:    order.Freight,
		Status:     order.Status,
		Created:    order.Created,
		PayWay:     order.PayWay,
		AddressID:  order.AddressID,
		Remark:     order.Remark,
	}
	getOrder = append(getOrder, add1)

	err = tx.Where("orderid = ?", order.ID).Find(&OrderProduct).Error
	if err != nil {
		return getOrder, err
	}

	for _, v := range OrderProduct {
		add1 := OrmOrders{
			Discount:  v.Discount,
			Count:     v.Count,
			Size:      v.Size,
			Color:     v.Color,
			ProductID: v.ProductID,
		}
		getOrder = append(getOrder, add1)

		collection := orm.MDSession.DB(orm.MD).C("productimage")
		orm.MDSession.Refresh()

		err = collection.Find(bson.M{"productid": add1.ProductID, "class": general.ProductAvatar}).One(&productAvatar)
		if err != nil {
			return nil, err
		}

		add1.Avatar = productAvatar.Image
		getOrder = append(getOrder, add1)
	}

	return getOrder, nil
}

func (osp *OrderServiceProvider) ChangeStatus(OrderID uint64, status uint8) error {
	cha := Orders{
		Status: status,
	}

	updater := map[string]interface{}{"status": status}
	db := orm.Conn

	err := db.Model(&cha).Where("id = ?", OrderID).Update(updater).Limit(1).Error
	return err
}
