package utils

import (
	"github.com/labstack/echo/v4"
	"github.com/shopicano/shopicano-backend/errors"
	"github.com/shopicano/shopicano-backend/models"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	PasswordCost = 11

	StoreID         = "store_id"
	StorePermission = "store_permission"
	UserID          = "user_id"
	UserPermission  = "user_permission"
)

func GeneratePassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	return string(bytes), err
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func ParseBearerToken(ctx echo.Context) (string, error) {
	bearer := ctx.Request().Header.Get("Authorization")
	bearerWithToken := strings.Split(bearer, " ")

	if len(bearerWithToken) != 2 {
		return "", errors.NewError("Bearer token not found")
	}
	return bearerWithToken[1], nil
}

func IsStoreStaff(ctx echo.Context) bool {
	perm := ctx.Get(StorePermission)
	return ctx.Get(StoreID) != nil && perm != nil && (perm.(models.Permission) == models.ManagerPerm || perm.(models.Permission) == models.AdminPerm)
}

func GetStoreID(ctx echo.Context) string {
	return ctx.Get(StoreID).(string)
}
