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
 *     Initial: 2017/07/21         Ai Hao
 *     Modify : 2017/07/21         Zhu Yaqiang
 *     Modify : 2017/08/10         Yu Yi
 *     Modify : 2017/07/21         Ma chao
 *     Modify : 2017/08/10         Li Zebang
 */

package models

import (
	"fmt"
	"time"

	"ShopApi/general"
	"ShopApi/orm"
	"gopkg.in/mgo.v2/bson"
)

type ProductServiceProvider struct {
}

var ProductService *ProductServiceProvider = &ProductServiceProvider{}

type Product struct {
	ID        uint64    `sql:"auto_increment;primary_key;" gorm:"column:id" json:"id"`
	Name      string    `json:"name"`
	TotalSale uint64    `gorm:"column:totalsale" json:"totalsale"`
	Category  uint64    `json:"categories"`
	Price     float64   `json:"price"`
	Detail    string    `json:"detail"`
	Status    uint8     `json:"status"`
	Created   time.Time `json:"created"`
}

type ProductImages struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Class     uint8         `bson:"class" json:"class"`
	ProductID uint64        `bson:"productid" json:"productid"`
	Image     string        `bson:"image" json:"image"`
}

type ProductSize struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ProductID uint64        `bson:"productid" json:"productid"`
	Size      string        `bson:"size" json:"size"`
}

type ProductColor struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ProductID uint64        `bson:"productid" json:"productid"`
	Color     string        `bson:"color" json:"color"`
}

type CreateProduct struct {
	Name         string   `json:"name" validate:"required"`
	Avatar       string   `json:"avatar"`
	Images       []string `json:"images"`
	DetailImages []string `json:"detailimages"`
	Category     uint64   `json:"category"`
	Price        float64  `json:"price"`
	Size         []string `json:"size" validate:"required"`
	Color        []string `json:"color" validate:"required"`
	Detail       string   `json:"detail" validate:"required"`
}

type ProductList struct {
	ID       uint64  `json:"id"`
	Name     string  `json:"title"`
	Avatar   string  `json:"img"`
	Category uint64  `json:"category"`
	Price    float64 `json:"price"`
}

type ProductCategory struct {
	Category uint64 `json:"category"`
	Page     uint64 `json:"page"`
	PageSize uint64 `json:"pagesize"`
}

type ProductID struct {
	ID uint64 `json:"id"`
}

type ProductInfo struct {
	Name         string   `json:"name"`
	Images       []string `json:"images"`
	DetailImages []string `json:"detailimages"`
	TotalSale    uint64   `json:"totalsale"`
	Category     uint64   `json:"category"`
	Price        float64  `json:"price"`
	Size         []string `json:"size"`
	Color        []string `json:"color"`
	Detail       string   `json:"detail"`
}

type ChangeProStatus struct {
	ID     uint64 `json:"id"`
	Status uint8  `json:"status"`
}

type ChangeCategory struct {
	ID       uint64 `json:"id"`
	Category uint64 `json:"category"`
}

func (Product) TableName() string {
	return "product"
}

func (ps *ProductServiceProvider) CreateProduct(create *CreateProduct) error {
	var (
		err     error
		product Product
	)

	product = Product{
		Name:     create.Name,
		Category: create.Category,
		Price:    create.Price,
		Detail:   create.Detail,
		Status:   general.ProductOnSale,
		Created:  time.Now(),
	}

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Create(&product).Error
	if err != nil {
		return err
	}

	err = tx.Where("name = ?", create.Name).First(&product).Error
	if err != nil {
		return err
	}

	err = AddProductImage(product.ID, create)
	if err != nil {
		return err
	}

	err = AddProductSize(product.ID, create)
	if err != nil {
		return err
	}

	err = AddProductColor(product.ID, create)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return err
}

func AddProductImage(productID uint64, create *CreateProduct) error {
	var (
		err    error
		image  ProductImages
		images []ProductImages
	)

	image = ProductImages{
		Class:     general.ProductAvatar,
		ProductID: productID,
		Image:     create.Avatar,
	}

	images = append(images, image)

	for _, img := range create.Images {
		image = ProductImages{
			Class:     general.ProductImage,
			ProductID: productID,
			Image:     img,
		}
		images = append(images, image)
	}

	for _, img := range create.DetailImages {
		image = ProductImages{
			Class:     general.ProductDetailImage,
			ProductID: productID,
			Image:     img,
		}
		images = append(images, image)
	}

	collection := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	for _, img := range images {
		err = collection.Insert(img)
		if err != nil {
			break
		}
	}

	return err
}

func AddProductSize(productID uint64, create *CreateProduct) error {
	var (
		err   error
		size  ProductSize
		sizes []ProductSize
	)

	for _, si := range create.Size {
		size = ProductSize{
			ProductID: productID,
			Size:      si,
		}
		sizes = append(sizes, size)
	}

	collection := orm.MDSession.DB(orm.MD).C("productsize")
	orm.MDSession.Refresh()

	for _, si := range sizes {
		err = collection.Insert(si)
		if err != nil {
			break
		}
	}

	return err
}

func AddProductColor(productID uint64, create *CreateProduct) error {
	var (
		err    error
		color  ProductColor
		colors []ProductColor
	)

	for _, co := range create.Color {
		color = ProductColor{
			ProductID: productID,
			Color:     co,
		}
		colors = append(colors, color)
	}

	collection := orm.MDSession.DB(orm.MD).C("productcolors")
	orm.MDSession.Refresh()

	for _, co := range colors {
		err = collection.Insert(co)
		if err != nil {
			break
		}
	}

	return err
}

func (ps *ProductServiceProvider) GetProductHeader() (*[]ProductList, error) {
	var (
		err     error
		product ProductList
		image   ProductImages
		header  []ProductList
	)

	db := orm.Conn

	sql := "SELECT * FROM product WHERE status = ? LIMIT 5 LOCK IN SHARE MODE"
	rows, err := db.Raw(sql, general.ProductOnSale).Rows()
	if err != nil {
		return &header, err
	}
	defer rows.Close()

	collection := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	for rows.Next() {
		db.ScanRows(rows, &product)
		err = collection.Find(bson.M{"productid": product.ID, "class": general.ProductImage}).One(&image)
		if err != nil {
			return &header, err
		}
		product.Avatar = image.Image
		header = append(header, product)
	}

	return &header, nil
}

func (ps *ProductServiceProvider) GetProductList() (*[]ProductList, error) {
	var (
		err     error
		product ProductList
		image   ProductImages
		list    []ProductList
	)

	db := orm.Conn

	sql := "SELECT * FROM product WHERE status = ? LIMIT 5, 6 LOCK IN SHARE MODE"
	rows, err := db.Raw(sql, general.ProductOnSale).Rows()
	if err != nil {
		return &list, err
	}
	defer rows.Close()

	collection := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	for rows.Next() {
		db.ScanRows(rows, &product)
		err = collection.Find(bson.M{"productid": product.ID, "class": general.ProductAvatar}).One(&image)
		if err != nil {
			return &list, err
		}
		product.Avatar = image.Image
		list = append(list, product)
	}

	return &list, nil
}

func (ps *ProductServiceProvider) GetProductByCategory(cate, pageStart, pageSize uint64) (*[]ProductList, error) {
	var (
		err     error
		product ProductList
		image   ProductImages
		list    []ProductList
	)

	db := orm.Conn

	sql := fmt.Sprintf("SELECT * FROM product WHERE category = ? AND status = ? LIMIT %d, %d LOCK IN SHARE MODE", pageStart, pageSize)

	rows, err := db.Raw(sql, cate, general.ProductOnSale).Rows()
	defer rows.Close()
	if err != nil {
		return &list, err
	}

	collection := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	for rows.Next() {
		db.ScanRows(rows, &product)
		err = collection.Find(bson.M{"productid": product.ID, "class": general.ProductAvatar}).One(&image)
		if err != nil {
			return &list, err
		}
		product.Avatar = image.Image
		list = append(list, product)
	}

	return &list, nil
}

func (ps *ProductServiceProvider) GetProInfo(id uint64) (*ProductInfo, error) {
	var (
		err     error
		product Product
		info    ProductInfo
		images  []ProductImages
		sizes   []ProductSize
		colors  []ProductColor
	)

	db := orm.Conn

	err = db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return &info, err
	}

	info = ProductInfo{
		Name:      product.Name,
		TotalSale: product.TotalSale,
		Category:  product.Category,
		Price:     product.Price,
		Detail:    product.Detail,
	}

	collection1 := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	err = collection1.Find(bson.M{"productid": id}).All(&images)
	if err != nil {
		return &info, err
	}

	for _, image := range images {
		switch image.Class {
		case general.ProductAvatar:
			continue
		case general.ProductImage:
			info.Images = append(info.Images, image.Image)
		case general.ProductDetailImage:
			info.DetailImages = append(info.DetailImages, image.Image)
		}
	}

	collection2 := orm.MDSession.DB(orm.MD).C("productsize")
	orm.MDSession.Refresh()

	err = collection2.Find(bson.M{"productid": id}).All(&sizes)
	if err != nil {
		return &info, err
	}

	for _, size := range sizes {
		info.Size = append(info.Size, size.Size)
	}

	collection3 := orm.MDSession.DB(orm.MD).C("productcolors")
	orm.MDSession.Refresh()

	err = collection3.Find(bson.M{"productid": id}).All(&colors)
	if err != nil {
		return &info, err
	}

	for _, color := range colors {
		info.Color = append(info.Color, color.Color)
	}

	return &info, nil
}

func (ps *ProductServiceProvider) ChangeProStatus(sta *ChangeProStatus) error {
	var (
		pro Product
	)

	updater := map[string]uint8{"status": sta.Status}

	db := orm.Conn

	return db.Model(&pro).Where("id = ?", sta.ID).Update(updater).Limit(1).Error
}

func (ps *ProductServiceProvider) ChangeCategory(cate *ChangeCategory) error {
	var (
		pro Product
	)

	db := orm.Conn
	err := db.Model(&pro).Where("id = ?", cate.ID).Update("category", cate.Category).Limit(1).Error

	return err
}

func (ps *ProductServiceProvider) GetMyPage() (*[]ProductList, error) {
	var (
		err        error
		product    ProductList
		image      ProductImages
		list       []ProductList
	)

	db := orm.Conn

	sql := "SELECT * FROM product WHERE status = ? LIMIT 6 LOCK IN SHARE MODE"
	rows, err := db.Raw(sql, general.ProductOnSale).Rows()
	if err != nil {
		return &list, err
	}
	defer rows.Close()

	collection := orm.MDSession.DB(orm.MD).C("productimage")
	orm.MDSession.Refresh()

	for rows.Next() {
		db.ScanRows(rows, &product)
		err = collection.Find(bson.M{"productid": product.ID, "class": general.ProductAvatar}).One(&image)
		if err != nil {
			return &list, err
		}
		product.Avatar = image.Image
		list = append(list, product)
	}

	return &list, nil
}
