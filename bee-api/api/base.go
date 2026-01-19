package api

import (
	"errors"
	"net/http"

	"gitee.com/stuinfer/bee-api/db"
	"gitee.com/stuinfer/bee-api/enum"
	"gitee.com/stuinfer/bee-api/model"
	"gitee.com/stuinfer/bee-api/proto"
	"gitee.com/stuinfer/bee-api/util"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type BaseApi struct {
}

func (api BaseApi) GetUserInfo(c *gin.Context) *model.BeeUser {
	data, _ := c.Get(enum.UserInfoKey)
	return data.(*model.BeeUser)
}
func (api BaseApi) GetUserInfoById(c *gin.Context) *model.BeeUser {
	// data, _ := c.Get(enum.UserInfoKey)
	if id, flag := c.Params.Get("id"); flag {
		var user model.BeeUser
		db.GetDB().First(&user, id)
		return &user
	}
	return nil
}

func (api BaseApi) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &proto.Response{
		Code: enum.ResCodeSuccess,
		Data: data,
		Msg:  "",
	})
}

func (api BaseApi) Fail(c *gin.Context, code enum.ResCode, msg string, data ...interface{}) {
	var resData interface{}
	if len(data) >= 1 {
		resData = data[0]
	}
	c.JSON(http.StatusOK, &proto.Response{
		Code: code,
		Data: resData,
		Msg:  msg,
	})
}

func (api BaseApi) Res(c *gin.Context, data interface{}, err error) {
	if err == nil {
		if lo.IsNil(data) || util.IsEmptyArrayOrSlice(data) || errors.Is(err, gorm.ErrRecordNotFound) {
			err = enum.ErrEmpty
		}
	}
	if err != nil {
		var e *enum.BussError
		ok := errors.As(err, &e)
		if ok {
			c.JSON(http.StatusOK, &proto.Response{
				Code: e.Code,
				Data: data,
				Msg:  e.Message,
			})
			return
		}
		c.JSON(http.StatusOK, &proto.Response{
			Code: enum.ResCodeFail,
			Data: data,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &proto.Response{
		Code: enum.ResCodeSuccess,
		Data: data,
		Msg:  "",
	})
}
