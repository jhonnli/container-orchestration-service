package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	//"github.com/jhonnli/container-orchestration-service/filter"
	"github.com/json-iterator/go"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"net/http"
)
import zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"

var zh1 = zh.New()
var uni = ut.New(zh1, zh1)
var trans, _ = uni.GetTranslator("zh")
var Validate = validator.New()

const OK string = "ok"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	zh_translations.RegisterDefaultTranslations(Validate, trans)
}

func GetZhError(err error) string {
	errs := err.(validator.ValidationErrors)
	msgMap := errs.Translate(trans)
	result := ""
	if msgMap != nil && len(msgMap) > 0 {
		idx := 0
		for _, v := range msgMap {
			result += v
			idx++
			if idx <= len(msgMap)-1 {
				result += ","
			}
		}

	}
	return result
}

type BoolResult struct {
	Result bool `json:"result"`
}

type StringResult struct {
	Result string `json:"result"`
}

type Result struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GenFailureResult(code string, message string) Result {
	return Result{
		Code:    code,
		Message: message,
	}
}

func GenBoolResult() Result {
	return Result{
		Code:    OK,
		Message: "",
		Data:    BoolResult{true},
	}
}

func GenStringResult(data string) Result {
	return Result{
		Code:    OK,
		Message: "",
		Data:    StringResult{data},
	}
}

func GenParamErrorResult(code string) Result {
	return Result{
		Code:    code,
		Message: "参数异常",
	}
}

func GenSuccessResult(data interface{}) Result {
	return Result{
		Code:    OK,
		Message: "",
		Data:    data,
	}
}

func GetStringParam(ctx *gin.Context, key string) string {
	return ctx.Param(key)
}

func GetHeaderParam(ctx *gin.Context, key string) string {
	return ctx.Request.Header.Get(key)
}

func GetJSONBody(ctx *gin.Context, res interface{}) error {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println("获取 %s body报错，原因:%s\n", GetHeaderParam(ctx, "Referer"))
		return err
	}
	return json.Unmarshal(body, res)
}

func ParamIsEmpty(ctx *gin.Context, code string, params ...string) bool {
	var flag = false
	for _, item := range params {
		if item == "" {
			flag = true
			break
		}
	}
	if flag {
		ctx.JSON(http.StatusBadRequest, GenParamErrorResult(code))
		return true
	}
	return false
}

func AddFilter(rg *gin.RouterGroup) {
	//rg.Use(filter.ClientAuth())
	//rg.Use(filter.JWTMiddleware())
	//rg.Use(filter.CORSMiddleware())
}

func AddApiFilter(rg *gin.RouterGroup) {
	//rg.Use(filter.ClientAuth())
	//rg.Use(filter.CORSMiddleware())
}
