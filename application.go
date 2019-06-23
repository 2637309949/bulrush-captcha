// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package captcha

import (
	"net/http"
	"strconv"

	"github.com/2637309949/bulrush"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// Captcha provide Digit Captcha
type Captcha struct {
	bulrush.PNBase
	Period    int
	Secret    string
	Config    base64Captcha.ConfigDigit
	URLPrefix string
}

// New defined return Captcha with default params
func New() *Captcha {
	c := &Captcha{
		URLPrefix: "/captcha",
		Secret:    "123abc#@%",
		Period:    60,
	}
	c.Config = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}
	return c
}

// Plugin for gin
func (c *Captcha) Plugin() bulrush.PNRet {
	return func(cfg *bulrush.Config, router *gin.RouterGroup, httpProxy *gin.Engine) {
		router.Use(func(ctx *gin.Context) {
			if data, err := ctx.Cookie("captcha"); err == nil && data != "" {
				decData := decrypt([]byte(data), c.Secret)
				dataStr := string(decData)
				ctx.Set("captcha", dataStr)
			}
			ctx.Next()
		})
		router.GET(c.URLPrefix, func(ctx *gin.Context) {
			if height, err := strconv.Atoi(ctx.Query("height")); err != nil && height != 0 {
				c.Config.Height = height
			}
			if width, err := strconv.Atoi(ctx.Query("width")); err != nil && width != 0 {
				c.Config.Width = width
			}
			idKey, captcha := base64Captcha.GenerateCaptcha("", c.Config)
			encryptData := encrypt([]byte(idKey), c.Secret)
			base64 := base64Captcha.CaptchaWriteToBase64Encoding(captcha)
			ctx.SetCookie("captcha", string(encryptData), Some(c.Period, 60).(int), "/", "", false, false)
			ctx.String(http.StatusOK, base64)
		})
	}
}
