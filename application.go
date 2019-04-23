/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush captcha plugin]
 * DependOn cookies [plugins]
 */

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

// Plugin for gin
func (c *Captcha) Plugin() bulrush.PNRet {
	defSecret := "123abc#@%"
	return func(cfg *bulrush.Config, router *gin.RouterGroup, httpProxy *gin.Engine) {
		router.Use(func(ctx *gin.Context) {
			if data, error := ctx.Cookie("captcha"); error == nil && data != "" {
				decData := decrypt([]byte(data), Some(c.Secret, defSecret).(string))
				dataStr := string(decData)
				ctx.Set("captcha", dataStr)
			}
			ctx.Next()
		})
		router.GET(Some(c.URLPrefix, "/captcha").(string), func(ctx *gin.Context) {
			if height, err := strconv.Atoi(ctx.Query("height")); err != nil && height != 0 {
				c.Config.Height = height
			}
			if width, err := strconv.Atoi(ctx.Query("width")); err != nil && width != 0 {
				c.Config.Width = width
			}
			idKey, captcha := base64Captcha.GenerateCaptcha("", c.Config)
			encryptData := encrypt([]byte(idKey), Some(c.Secret, defSecret).(string))
			base64 := base64Captcha.CaptchaWriteToBase64Encoding(captcha)
			ctx.SetCookie("captcha", string(encryptData), Some(c.Period, 60).(int), "/", "", false, false)
			ctx.String(http.StatusOK, base64)
		})
	}
}
