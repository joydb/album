package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/penguinn/penguin/component/log"
	"net/http"
	"github.com/penguinn/album/components/captcha"
	"github.com/penguinn/album/models"
	"image/png"
	"bytes"
	"strconv"
)

type CheckController struct {
	PostImageAuth func(*gin.Context) `path:"/check/image-auth"`
	PostUserName func(*gin.Context) `path:"/check/username"`
}

func(CheckController) Name() string {
	return "CheckController"
}

func NewCheckController() CheckController {
	return CheckController{
		PostImageAuth:PostImageAuth,
		PostUserName:PostUserName,
	}
}

type ImageAuthRequest struct {
	Type    int  `form:"type" binding:"required"` 		//1：代表主界面
}

type PostUsernameRequest struct {
	Username  string  `form:"type" binding:"required"`
}

func PostImageAuth(c *gin.Context) {
	var imageAuthRequest ImageAuthRequest
	err := c.Bind(&imageAuthRequest)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, PARAM_ERROR_CODE))
		return
	}
	cookie, err := c.Request.Cookie("sess-token")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, PARAM_ERROR_CODE, "获取不到token"))
		return
	}
	token := cookie.Value
	img, str := captcha.GetImgAndStr()
	err = models.ImageAuth{}.Insert(token, str)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, SERVER_ERROR_CODE))
		return
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)

	c.Writer.Header().Del("Content-Type")
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))
	c.Data(http.StatusOK, "image/png", buf.Bytes())
}

func PostUserName(c *gin.Context) {
	var postUsernameRequest PostUsernameRequest
	err := c.Bind(&postUsernameRequest)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, PARAM_ERROR_CODE))
		return
	}
	isExist, err := models.User{}.ValidateUsername(postUsernameRequest.Username)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, SERVER_ERROR_CODE))
		return
	}
	if isExist{
		c.JSON(http.StatusOK, NewResponse(false, USERNAME_EXIST_CODE))
	}
	c.JSON(http.StatusOK, NewResponse(true, SUCCESS_CODE))
}