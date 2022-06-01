package rsp

import (
	"github.com/dartagnanli/alpha/aerror"
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, err error) {
	var errResp *aerror.Error
	var ok bool

	if errResp, ok = err.(*aerror.Error); !ok {
		errResp = aerror.ErrInternalError().WithMessage(err.Error())
	}

	c.JSON(errResp.HTTPStatusCode, errResp)
}
