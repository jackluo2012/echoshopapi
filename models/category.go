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
 */

package models

import (
	"time"

	"ShopApi/general"
	"ShopApi/orm"
)

type CategoryServiceProvider struct {
}

var CategoryService *CategoryServiceProvider = &CategoryServiceProvider{}

type Category struct {
	ID      uint64    `sql:"auto_increment;primary_key;",json:"id"`
	Name    string    `json:"name"`
	PID     uint64    `gorm:"column:pid" json:"pid"`
	Status  uint64    `json:"status"`
	Created time.Time `json:"created"`
}

type CreateCategory struct {
	Name string `json:"name" validate:"required,alphanumunicode"`
	PID  uint64 `json:"pid"`
}

type GetCategory struct {
	PID uint64 `json:"pid"`
}

type CategoryGet struct {
	Name string `json:"name"`
}

func (Category) TableName() string {
	return "category"
}

func (csp *CategoryServiceProvider) CreateCategory(createCategory CreateCategory) error {
	category := &Category{
		Name:    createCategory.Name,
		PID:     createCategory.PID,
		Status:  general.CategoryOnUse,
		Created: time.Now(),
	}

	db := orm.Conn

	return db.Create(&category).Error
}

func (csp *CategoryServiceProvider) CheckPID(pid uint64) (err error) {
	var (
		category Category
	)

	db := orm.Conn

	return db.Where("id =? ", pid).First(&category).Error
}

func (csp *CategoryServiceProvider) GetCategory() (*[]CategoryGet, error) {
	var (
		err          error
		categories   []Category
		categoryList []CategoryGet
	)

	db := orm.Conn

	err = db.Where("status = ?", general.CategoryOnUse).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		categoryGet := CategoryGet{Name: category.Name}
		categoryList = append(categoryList, categoryGet)
	}

	return &categoryList, nil
}
