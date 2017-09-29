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
 *      Initial: 2017/07/21         Ai Hao
 *      Modify : 2017/07/21         Zhu Yaqiang
 *      Modify : 2017/08/10         Yu Yi
 *      Modify : 2017/07/21         Ma Chao
 *      Modify : 2017/08/10         Li Zebang
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

func CreateProduct(c echo.Context) error {
	var (
		err     error
		product models.CreateProduct
	)

	if err = c.Bind(&product); err != nil {
		log.Logger.Error("[ERROR] CreateProduct Bind", err)

		return general.NewErrorWithMessage(errcode.ErrCreateProductInvalidParams, err.Error())
	}

	if err = c.Validate(product); err != nil {
		log.Logger.Error("[ERROR] CreateProduct Validate", err)

		return general.NewErrorWithMessage(errcode.ErrCreateProductInvalidParams, err.Error())
	}

	err = models.ProductService.CreateProduct(&product)
	if err != nil {
		log.Logger.Error("[ERROR] CreateProduct CreateProduct:", err)

		return general.NewErrorWithMessage(errcode.ErrCreateProductDatabase, err.Error())
	}

	log.Logger.Info("[SUCCEED] CreateProduct: Name %s", product.Name)

	return c.JSON(errcode.CreateProductSucceed, general.NewMessage(errcode.CreateProductSucceed))
}

func GetProductList(c echo.Context) error {
	var (
		err    error
		header *[]models.ProductList
		list   *[]models.ProductList
	)

	header, err = models.ProductService.GetProductHeader()
	if err != nil {
		log.Logger.Error("[ERROR] GetProductList GetProductHeader", err)

		return general.NewErrorWithMessage(errcode.ErrGetListDatabase, err.Error())
	}

	list, err = models.ProductService.GetProductList()
	if err != nil {
		log.Logger.Error("[ERROR] GetProductList GetProductList", err)

		return general.NewErrorWithMessage(errcode.ErrGetListDatabase, err.Error())
	}

	log.Logger.Info("[SUCCEED] GetProductList %v")

	return c.JSON(errcode.GetListSucceed, general.NewMessageForProductList(errcode.GetListSucceed, header, list))
}

func GetProductListByCategory(c echo.Context) error {
	var (
		err      error
		category models.ProductCategory
		list     *[]models.ProductList
	)

	if err = c.Bind(&category); err != nil {
		log.Logger.Error("[ERROR] GetProductList Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrGetProductListByCategoryInvalidParams, err.Error())
	}

	if err = c.Validate(category); err != nil {
		log.Logger.Error("[ERROR] GetProductList Validate", err)

		return general.NewErrorWithMessage(errcode.ErrGetProductListByCategoryInvalidParams, err.Error())
	}

	pageStart := utility.Paging(category.Page, category.PageSize)
	list, err = models.ProductService.GetProductByCategory(category.Category, pageStart, category.PageSize)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Logger.Error("[ERROR] GetProductByCategory GetProductByCategory:", err)

			return general.NewErrorWithMessage(errcode.ErrGetProductListByCategoryNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] GetProductByCategory GetProductByCategory:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] GetProductByCategory %v")

	return c.JSON(errcode.GetProductListByCategorySucceed, general.NewMessageWithData(errcode.GetProductListByCategorySucceed, list))
}

func GetProInfo(c echo.Context) error {
	var (
		err         error
		productID   *models.ProductID
		productInfo *models.ProductInfo
	)

	if err = c.Bind(&productID); err != nil {
		log.Logger.Error("[ERROR] GetProInfo Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrGetProInfoInvalidParams, err.Error())
	}

	if err = c.Validate(productID); err != nil {
		log.Logger.Error("[ERROR] GetProInfo Validate", err)

		return general.NewErrorWithMessage(errcode.ErrGetProductListByCategoryInvalidParams, err.Error())
	}

	productInfo, err = models.ProductService.GetProInfo(productID.ID)

	if err != nil {
		if err != nil && !strings.Contains(err.Error(), "not found") {
			log.Logger.Error("[ERROR] GetUserInfo GetUserAvatar:", err)

			return general.NewErrorWithMessage(errcode.ErrGetProInfoNotFound, err.Error())
		}

		log.Logger.Error("[ERROR] GetProInfo GetProInfo:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] GetProInfo %v")

	return c.JSON(errcode.GetProInfoSucceed, general.NewMessageWithData(errcode.GetProInfoSucceed, productInfo))
}

func ChangeProStatus(c echo.Context) error {
	var (
		err error
		cps models.ChangeProStatus
	)

	if err = c.Bind(&cps); err != nil {
		log.Logger.Error("[ERROR] ChangeProStatus Bind:", err)

		return general.NewErrorWithMessage(errcode.ErrChangeProStatusInvalidParams, err.Error())
	}

	if err = c.Validate(cps); err != nil {
		log.Logger.Error("[ERROR] Product Status Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrChangeProStatusInvalidParams, err.Error())
	}

	if cps.Status != general.ProductOnSale && cps.Status != general.ProductUnSale {
		err = errors.New("Invalid Product Status")

		log.Logger.Error("[ERROR] ChangeProStatus:", err)

		return general.NewErrorWithMessage(errcode.ErrChangeProStatusInvalidParams, err.Error())
	}

	err = models.ProductService.ChangeProStatus(&cps)
	if err != nil {
		log.Logger.Error("[ERROR] ChangeProStatus ChangeProStatus:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] ChangeProStatus")

	return c.JSON(errcode.ChangeStatusSucceed, general.NewMessage(errcode.ChangeStatusSucceed))
}

func ChangeCategory(c echo.Context) error {
	var (
		err error
		cc  *models.ChangeCategory
	)

	if err = c.Bind(&cc); err != nil {
		log.Logger.Error("[ERROR] ChangeCategories Bind with error:", err)

		return general.NewErrorWithMessage(errcode.ErrCategoryInvalidParams, err.Error())
	}

	if err = c.Validate(cc); err != nil {
		log.Logger.Error("[ERROR] Categories Validate:", err)

		return general.NewErrorWithMessage(errcode.ErrCategoryInvalidParams, err.Error())
	}

	err = models.ProductService.ChangeCategory(cc)
	if err != nil {

		log.Logger.Error("[ERROR] Categories change with error:", err)

		return general.NewErrorWithMessage(errcode.ErrMysql, err.Error())
	}

	log.Logger.Info("[SUCCEED] Change Categories: Category %s", cc.Category)

	return c.JSON(errcode.ChangeCategorySucceed, general.NewMessage(errcode.ChangeCategorySucceed))
}

func GetMyPage(c echo.Context) error {
	var (
		err    error
		list   *[]models.ProductList
	)

	list, err = models.ProductService.GetMyPage()
	if err != nil {
		log.Logger.Error("[ERROR] GetMyPage ", err)

		return general.NewErrorWithMessage(errcode.ErrGetListDatabase, err.Error())
	}

	log.Logger.Info("[SUCCEED] GetMyPage %v")

	return c.JSON(errcode.GetListSucceed, general.NewMessageWithData(errcode.GetListSucceed, list))
}
