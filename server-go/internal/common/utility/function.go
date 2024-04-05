package utility

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"server-go/internal/consts"
	"strings"
	"time"
)

// GetVerErr @Title 处理验证错误消息
func GetVerErr(err error) string {
	er, ok := err.(validator.ValidationErrors)
	if ok {
		field := er[0].Field()
		if field == er[0].StructField() {
			field = strings.ToLower(field[0:1]) + field[1:]
		}
		switch er[0].Tag() {
		case "required":
			return field + "不能为空"
		case "min":
			if er[0].Type().String() == "string" {
				return field + "不能小于" + er[0].Param() + "位"
			}
			return field + "不能小于" + er[0].Param()
		case "gte":
			if er[0].Type().String() == "string" {
				return field + "不能小于" + er[0].Param() + "位"
			}
			return field + "不能小于" + er[0].Param()
		}
		return field + "错误"
	} else {
		return err.Error()
		//return "参数格式错误"
	}
}

// InArray @Title 是否在数组中
func InArray(need interface{}, data interface{}) bool {
	if datas, ok := data.([]int); ok {
		for _, item := range datas {
			if item == need {
				return true
			}
		}
	}
	if datas, ok := data.([]uint); ok {
		for _, item := range datas {
			if item == need {
				return true
			}
		}
	}
	if datas, ok := data.([]string); ok {
		for _, item := range datas {
			if item == need {
				return true
			}
		}
	}
	return false
}

// CreateDateDir 根据当前日期来创建文件夹
func CreateDateDir(Path string) string {
	folderName := time.Now().Format("20060102")
	folderPath := filepath.Join(Path, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0755)
	}
	return folderName
}

// DownloadFormUrl @Title 下载文件
func DownloadFormUrl(src string, filename string) error {
	res, err := http.Get(src)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, bytes.NewReader(body)); err != nil {
		return err
	}
	return nil
}

// Utf8ToGbk @Title utf8转gbk
func Utf8ToGbk(source string) string {
	result, _, _ := transform.String(simplifiedchinese.GBK.NewEncoder(), source)
	return result
}

// GbkToUtf8 @Title gbk转utf8
func GbkToUtf8(source string) string {
	result, _, _ := transform.String(simplifiedchinese.GBK.NewDecoder(), source)
	return result
}

// GetDeviceType @Title 获取设备类型
func GetDeviceType(c *gin.Context) uint {
	switch strings.ToLower(c.GetHeader("Device-Type")) {
	case "android":
		return consts.DeviceTypeAndroid
	case "ios":
		return consts.DeviceTypeIos
	default:
		return consts.DeviceTypeWeb
	}
}

// IsPhone @Title 是否手机号
func IsPhone(phone string) bool {
	return regexp.MustCompile(`^1\d{10}$`).MatchString(phone)
}

// Ternary @Title 三元运算
func Ternary(condition bool, resultA, resultB interface{}) interface{} {
	if condition {
		return resultA
	}
	return resultB
}

var timeTemplates = []string{
	"2006-01-02 15:04:05", //常规类型
	"2006/01/02 15:04:05",
	"2006-01-02",
	"2006/01/02",
	"15:04:05",
}

/* 时间格式字符串转换 */
func TimeStringToGoTime(tm string) time.Time {
	for i := range timeTemplates {
		t, err := time.ParseInLocation(timeTemplates[i], tm, time.Local)
		if nil == err && !t.IsZero() {
			return t
		}
	}
	return time.Time{}
}
