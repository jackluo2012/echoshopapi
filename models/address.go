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
 *     Initial: 2017/07/18        Li Zebang
 *     Modify : 2017/07/20        Yu Yi
 *     Modify : 2017/07/20        Yang Zhengtian
 *     Modify : 2017/07/28        Li Zebang
 */

package models

import (
	"time"

	"ShopApi/general"
	"ShopApi/orm"
	"ShopApi/utility"
)

type AddressServiceProvider struct {
}

var AddressService *AddressServiceProvider = &AddressServiceProvider{}

type Address struct {
	ID        string    `sql:"primary_key;" json:"id"`
	UserID    uint64    `gorm:"column:userid" json:"userid"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Area      string    `json:"area"`
	Address   string    `json:"address"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	IsDefault uint8     `gorm:"column:isdefault" json:"isdefault" `
}

type AddressJSON struct {
	ID        string `json:"id" validate:"required"`
	UserID    uint64 `json:"userid"`
	Name      string `json:"receiver" validate:"required,alphanumunicode"`
	Phone     string `json:"phone" validate:"required,numeric,len=11"`
	Area      string `json:"area" validate:"required"`
	Address   string `json:"detailAdress" validate:"required"`
	IsDefault bool   `json:"default"`
}

type AddressID struct {
	ID string `json:"id" validate:"required"`
}

func (Address) TableName() string {
	return "address"
}

func (asp *AddressServiceProvider) AddAddress(addAddress *AddressJSON) error {
	var (
		err     error
		address Address
	)

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	if addAddress.IsDefault {
		updateToNotDefault := map[string]uint8{"isdefault": general.AddressNotDefault}

		err = tx.Model(&address).Where("userid = ?", addAddress.UserID).Update(updateToNotDefault).Limit(1).Error
	}

	address = Address{
		ID:        addAddress.ID,
		UserID:    addAddress.UserID,
		Name:      addAddress.Name,
		Phone:     addAddress.Phone,
		Area:      addAddress.Area,
		Address:   addAddress.Address,
		Created:   time.Now(),
		Updated:   time.Now(),
		IsDefault: utility.BoolToUint8(addAddress.IsDefault),
	}

	err = tx.Create(&address).Error

	return err
}

func (asp *AddressServiceProvider) ChangeAddress(changeAddress *AddressJSON, userID uint64) error {
	var (
		err     error
		address Address
	)

	updater := map[string]interface{}{
		"name":      changeAddress.Name,
		"phone":     changeAddress.Phone,
		"area":      changeAddress.Area,
		"address":   changeAddress.Address,
		"updated":   time.Now(),
		"isdefault": utility.BoolToUint8(changeAddress.IsDefault),
	}

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	if changeAddress.IsDefault {
		updateToNotDefault := map[string]uint8{"isdefault": general.AddressNotDefault}

		err = tx.Model(&address).Where("userid = ?", userID).Update(updateToNotDefault).Limit(1).Error
	}

	err = tx.Model(&address).Where("id = ?", changeAddress.ID).Update(updater).Limit(1).Error

	return err
}

func (asp *AddressServiceProvider) GetAddressByUserID(userID uint64) (*[]AddressJSON, error) {
	var (
		err         error
		address     []Address
		addressList []AddressJSON
	)

	db := orm.Conn

	err = db.Where("userid = ?", userID).Find(&address).Error
	if err != nil {
		return &addressList, err
	}

	for _, addr := range address {
		addressGet := AddressJSON{
			ID:        addr.ID,
			Name:      addr.Name,
			Phone:     addr.Phone,
			Area:      addr.Area,
			Address:   addr.Address,
			IsDefault: utility.Uint8ToBool(addr.IsDefault),
		}
		addressList = append(addressList, addressGet)
	}

	return &addressList, nil
}

func (asp *AddressServiceProvider) AlterAddress(alterAddress *AddressID, userID uint64) (err error) {
	var (
		address Address
	)

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	updateToNotDefault := map[string]interface{}{"isdefault": general.AddressNotDefault}

	err = tx.Model(&address).Where("userid = ?", userID).Update(updateToNotDefault).Limit(1).Error
	if err != nil {
		return err
	}

	updaterToDefault := map[string]uint8{"isdefault": general.AddressDefault}

	err = tx.Model(&address).Where("id = ?", alterAddress.ID).Update(updaterToDefault).Limit(1).Error
	if err != nil {
		return err
	}

	return nil
}

func (asp *AddressServiceProvider) DeleteAddress(deleteAddress *AddressID) error {
	var (
		address Address
	)

	address.ID = deleteAddress.ID

	db := orm.Conn

	return db.Delete(&address).Error
}

func (asp *AddressServiceProvider) FindAddress(ID string, userID uint64) error {
	var (
		address Address
	)

	db := orm.Conn

	return db.Where("id = ? And userid = ?", ID, userID).First(&address).Error
}
