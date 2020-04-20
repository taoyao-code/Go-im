// example of HTTP server that uses the captcha package.
package util

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
)

//configJsonBody json request body.
type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

var store = base64Captcha.DefaultMemStore

func GenerateCaptchaHandler(w http.ResponseWriter, r *http.Request) {
	var param configJsonBody
	param.Id = uuid.New().String()
	param.DriverDigit = base64Captcha.DefaultDriverDigit
	driver := param.DriverDigit
	cap := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cap.Generate()
	body := map[string]interface{}{"code": 0, "data": b64s, "id": id, "msg": "success"}
	if err != nil {
		body = map[string]interface{}{"code": -1, "msg": err.Error()}
	}
	ret, err := json.Marshal(body)
	if err != nil {
		log.Printf(err.Error())
	}
	// 3.输出
	w.Write(ret)
}
func CaptchaVerifyHandle(UUID, Code string) error {
	if !store.Verify(UUID, Code, true) {
		return errors.New("验证码错误")
	}
	return nil
}
