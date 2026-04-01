package utils

import (
	"fmt"
	"strconv"

	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

// GetUserID returns user ID from Fiber context as string (supports int, int64, string)
func GetUserID(c frameworkhttp.Context) (string, bool) {
	val := c.Locals("userID")
	switch v := val.(type) {
	case string:
		return v, true
	case int:
		return fmt.Sprintf("%d", v), true
	case int64:
		return fmt.Sprintf("%d", v), true
	default:
		return "", false
	}
}

// GetUserIDInt returns user ID as int from Fiber context (fallback 0)
func GetUserIDInt(c frameworkhttp.Context) (int, bool) {
	val := c.Locals("userID")
	switch v := val.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case string:
		i, err := strconv.Atoi(v)
		if err == nil {
			return i, true
		}
	}
	return 0, false
}

func GetDeptIDInt(c frameworkhttp.Context) (int, bool) {
	val := c.Locals("deptID")
	switch v := val.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case string:
		i, err := strconv.Atoi(v)
		if err == nil {
			return i, true
		}
	}
	return 0, false
}

// GetUserRole returns the user's role from context (assumed string)
func GetUserRole(c frameworkhttp.Context) string {
	if role, ok := c.Locals("role").(string); ok {
		return role
	}
	return ""
}

// GetUserEmail returns the user's email from context (assumed string)
func GetUserEmail(c frameworkhttp.Context) string {
	if email, ok := c.Locals("email").(string); ok {
		return email
	}
	return ""
}

// GetUserWithPermission returns both user ID and permission from context
func GetUserWithPermission(c frameworkhttp.Context) (string, string) {
	userID, _ := GetUserID(c)
	perm := ""
	if p, ok := c.Locals("permission").(string); ok {
		perm = p
	}
	return userID, perm
}
