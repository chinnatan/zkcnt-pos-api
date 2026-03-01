package auth

import (
	"errors"
	"net/http"
	auth_dto "zkcnt-pos-api/internal/auth/dto"
	"zkcnt-pos-api/pkg"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type AuthHandler struct {
	app *pocketbase.PocketBase
}

// NewAuthHandler สร้าง instance ของ handler เพื่อใช้ app ร่วมกัน
func NewAuthHandler(app *pocketbase.PocketBase) *AuthHandler {
	return &AuthHandler{app: app}
}

// Login จะทำการลงชื่อเข้าใช้งาน
func (h *AuthHandler) Login(ctx *core.RequestEvent) error {
	return ctx.String(http.StatusOK, "Login success")
}

// SignUp จะทำการสมัครสมาชิก
func (h *AuthHandler) SignUp(ctx *core.RequestEvent) error {
	// ดึงข้อมูลจาก request
	var request auth_dto.SignUpDTO
	if err := ctx.BindBody(&request); err != nil {
		return ctx.BadRequestError("Invalid request body", err)
	}

	// ตรวจสอบความถูกต้องของข้อมูล
	if request.Email != "" && request.Password != "" && request.ConfirmPassword != "" && request.Name != "" && request.Role != "" && request.StoreName != "" && request.StoreAddress != "" {
		return ctx.BadRequestError("Invalid request body", errors.New("invalid request body"))
	}

	if err := h.app.RunInTransaction(func(txApp core.App) error {
		// Step 1: หา collection users
		users, err := txApp.FindCollectionByNameOrId("users")
		if err != nil {
			h.app.Logger().Error("ไม่สามารถหา collection users ได้", "error", err)
			return ctx.InternalServerError("ไม่สามารถสร้างบัญชีได้ กรุณาติดต่อผู้ดูแลระบบ", err)
		}

		// Step 2: สร้าง record ใหม่ใน collection users
		record := core.NewRecord(users)
		record.Set("email", request.Email)
		record.Set("password", request.Password)
		record.Set("confirmPassword", request.ConfirmPassword)
		record.Set("name", request.Name)
		record.Set("role", "admin")

		// Step 3: บันทึก record ใหม่
		if err := txApp.Save(record); err != nil {
			h.app.Logger().Error("ไม่สามารถบันทึก record ใน collection users ได้", "error", err)
			return ctx.InternalServerError("ไม่สามารถสร้างบัญชีได้ กรุณาติดต่อผู้ดูแลระบบ", err)
		}

		// Step 4: สร้าง record ใหม่ใน collection stores
		stores, err := txApp.FindCollectionByNameOrId("stores")
		if err != nil {
			h.app.Logger().Error("ไม่สามารถหา collection stores ได้", "error", err)
			return ctx.InternalServerError("ไม่สามารถสร้างบัญชีได้ กรุณาติดต่อผู้ดูแลระบบ", err)
		}

		storeRecord := core.NewRecord(stores)
		storeRecord.Set("name", request.StoreName)
		storeRecord.Set("address", request.StoreAddress)
		storeRecord.Set("created_id", record.Id)

		// Step 5: บันทึก record ใหม่
		if err := txApp.Save(storeRecord); err != nil {
			h.app.Logger().Error("ไม่สามารถบันทึก record ใน collection stores ได้", "error", err)
			return ctx.InternalServerError("ไม่สามารถสร้างบัญชีได้ กรุณาติดต่อผู้ดูแลระบบ", err)
		}

		return pkg.SendSuccessResponse(ctx, auth_dto.User{
			ID:           record.Id,
			Email:        record.GetString("email"),
			Name:         record.GetString("name"),
			Role:         record.GetString("role"),
			StoreID:      storeRecord.Id,
			StoreName:    storeRecord.GetString("name"),
			StoreAddress: storeRecord.GetString("address"),
		})
	}); err != nil {
		h.app.Logger().Error("ไม่สามารถสร้างบัญชีได้ กรุณาติดต่อผู้ดูแลระบบ", "error", err)
		return err
	}

	return nil
}
