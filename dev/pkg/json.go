package pkg

import (
	"net/http"

	"github.com/pocketbase/pocketbase/core"
)

func SendSuccessResponse(c *core.RequestEvent, data any) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "สำเร็จ",
		"data":    data,
	})
}
