package captcha

import (
	"github.com/afocus/captcha"
	"github.com/penguinn/penguin/component/config"
	"image/color"
)

var (
	captchaIns *captcha.Captcha
)

func init() {
	path := config.GetString("captcha.path")
	captchaIns = captcha.New()
	err := captchaIns.SetFont(path)
	if err != nil {
		panic(err)
	}
	captchaIns.SetSize(128, 64)
	captchaIns.SetDisturbance(captcha.NORMAL)
	captchaIns.SetFrontColor(color.RGBA{255, 255, 255, 255})
	captchaIns.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
}

func GetImgAndStr() (img *captcha.Image, str string) {
	img, str = captchaIns.Create(4, captcha.ALL)
	return
}
