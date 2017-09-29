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
 *     Modify : 2017/07/21        Xu Haosheng
 *     Modify : 2017/07/20	      Zhang Zizhao
 *     Modify : 2017/07/21        Yang Zhengtian
 *     Modify : 2017/07/19        Ma Chao
 *     Modify : 2017/08/10        Li Zebang
 *     Modify : 2017/08/11        Yu Yi
 */

package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"ShopApi/general"
	"ShopApi/orm"
	"ShopApi/utility"
)

type UserServiceProvider struct {
}

var UserService *UserServiceProvider = &UserServiceProvider{}

type User struct {
	UserID   uint64    `sql:"auto_increment;primary_key;" gorm:"column:id" json:"userid"`
	Password string    `json:"password" validate:"required,alphanum,min=6,max=30"`
	Name     string    `json:"name"`
	Status   uint16    `json:"status"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type UserInfo struct {
	UserID   uint64 `sql:"primary_key" gorm:"column:userid" json:"userid"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Sex      uint8  `json:"sex"`
}

type UserAvatar struct {
	UserID uint64 `bson:"_id,omitempty" json:"id"`
	Avatar string `bson:"avatar" json:"avatar" validate:"required"`
}

type UserGet struct {
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"name"`
	Sex      uint8  `json:"sex"`
}

type Register struct {
	Mobile *string `json:"mobile" validate:"required,numeric,len=11"`
	Pass   *string `json:"password" validate:"required,alphanum,min=6,max=64"`
}

type Login struct {
	Mobile *string `json:"mobile" validate:"required,numeric,len=11"`
	Pass   *string `json:"password" validate:"required,alphanum,min=6,max=64"`
}

type ChangeUserInfo struct {
	Nickname string `json:"name"`
	Sex      uint8  `json:"sex"`
}

type ChangePhone struct {
	Phone string `json:"phone" validate:"required,numeric,len=11"`
}

type ChangePassword struct {
	Password *string `json:"password" validate:"required,alphanum,min=6,max=64"`
	NewPass  *string `json:"newpassword" validate:"required,alphanum,min=6,max=64"`
}

func (User) TableName() string {
	return "user"
}

func (UserInfo) TableName() string {
	return "userinfo"
}

func (us *UserServiceProvider) Register(name, pass *string) error {
	hashedPass, err := utility.GenerateHash(*pass)
	if err != nil {
		return err
	}

	u := User{
		Name:     *name,
		Password: string(hashedPass),
		Status:   general.UserActive,
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	db := orm.Conn

	tx := db.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Create(&u).Error
	if err != nil {
		return err
	}

	info := UserInfo{
		UserID: u.UserID,
		Phone:  *name,
		Sex:    general.Man,
	}

	err = tx.Create(&info).Error
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (us *UserServiceProvider) Login(name, pass *string) (bool, uint64, error) {
	var (
		u   User
		err error
	)

	db := orm.Conn
	err = db.Where("name = ?", *name).First(&u).Error
	if err != nil {
		return false, 0, err
	}

	if !utility.CompareHash([]byte(u.Password), *pass) {
		return false, 0, nil
	}

	return true, u.UserID, nil
}

func (us *UserServiceProvider) GetUserInfo(UserID uint64) (*UserGet, error) {
	var (
		err error
		ui  UserInfo
		ug  UserGet
	)

	db := orm.Conn
	err = db.Where("userid = ?", UserID).First(&ui).Error
	if err != nil {
		return &ug, err
	}

	ug = UserGet{
		Phone:    ui.Phone,
		Nickname: ui.Nickname,
		Sex:      ui.Sex,
	}

	return &ug, nil
}

func (us *UserServiceProvider) GetUserAvatar(userID uint64) (*UserAvatar, error) {
	var (
		err    error
		avatar UserAvatar
	)

	collection := orm.MDSession.DB(orm.MD).C("useravatar")
	orm.MDSession.Refresh()
	err = collection.Find(bson.M{"_id": userID}).One(&avatar)

	return &avatar, err
}

func (us *UserServiceProvider) ChangeUserInfo(info *ChangeUserInfo, userID uint64) error {
	var (
		con     UserInfo
		updater = make(map[string]interface{})
	)

	if info.Nickname != "" {
		updater["nickname"] = info.Nickname
	}

	if info.Sex != general.Sex {
		updater["sex"] = info.Sex
	}

	db := orm.Conn
	err := db.Model(&con).Where("userid =? ", userID).Update(updater).Limit(1).Error

	return err
}

func (us *UserServiceProvider) ChangeUserAvatar(avatar *UserAvatar) error {
	collection := orm.MDSession.DB(orm.MD).C("useravatar")
	orm.MDSession.Refresh()
	_, err := collection.Upsert(bson.M{"_id": avatar.UserID}, avatar)

	return err
}

func (us *UserServiceProvider) ChangePhone(userID uint64, phone string) error {
	var (
		err  error
		user User
		info UserInfo
	)

	changeUser := map[string]string{"name": phone}
	changeInfo := map[string]string{"phone": phone}

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Model(&user).Where("id = ?", userID).Update(changeUser).Limit(1).Error
	if err != nil {
		return err
	}

	err = tx.Model(&info).Where("userid = ?", userID).Update(changeInfo).Limit(1).Error
	if err != nil {
		return err
	}

	return err
}

func (us *UserServiceProvider) ChangePassword(changePassword *ChangePassword, id uint64) (bool, error) {
	var (
		user User
		err  error
	)

	tx := orm.Conn.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()

	err = tx.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return false, err
	}

	if !utility.CompareHash([]byte(user.Password), *changePassword.Password) {
		return false, nil
	}

	hashPass, err := utility.GenerateHash(*changePassword.NewPass)

	updater := map[string]interface{}{"password": hashPass}

	err = tx.Model(&user).Where("id =? ", id).Update(updater).Limit(1).Error

	return true, err
}
