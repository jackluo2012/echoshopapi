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
	"github.com/spf13/viper"
)

type shopServerConfig struct {
	address   string
	isDebug   bool
	corsHosts []string
	tokenKey  string
	mysqlHost string
	mysqlPort string
	mysqlUser string
	mysqlPass string
	mysqlDb   string
	mysqlSize int
	MgoUrl    string
}

var (
	configuration *shopServerConfig
)

func readConfiguration() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	configuration = &shopServerConfig{
		address:   viper.GetString("server.address"),
		isDebug:   viper.GetBool("server.debug"),
		corsHosts: viper.GetStringSlice("middleware.cors.hosts"),
		tokenKey:  viper.GetString("middleware.jwt.tokenkey"),
		mysqlHost: viper.GetString("mysql.host"),
		mysqlPort: viper.GetString("mysql.port"),
		mysqlUser: viper.GetString("mysql.user"),
		mysqlPass: viper.GetString("mysql.pass"),
		mysqlDb:   viper.GetString("mysql.db"),
		mysqlSize: viper.GetInt("mysql.size"),
		MgoUrl:    viper.GetString("mongodb.url"),
	}
}
