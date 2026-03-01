package auth

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func RegisterRoutes(app *pocketbase.PocketBase, group *router.RouterGroup[*core.RequestEvent]) {
	// สร้าง Handler
	handler := NewAuthHandler(app)

	// ผูก Route กับ Method ใน Handler
	group.GET("/login", handler.Login)
	group.POST("/signup", handler.SignUp)
}
