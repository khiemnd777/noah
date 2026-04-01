package middleware

import (
	stdhttp "net/http"

	"github.com/khiemnd777/noah_framework/shared/app/client_error"
	"github.com/khiemnd777/noah_framework/shared/utils"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func RequireDepartmentMember(deptIDFromPathParam string) frameworkhttp.Handler {
	return func(c frameworkhttp.Context) error {
		userID, ok := utils.GetUserIDInt(c)
		if !ok || userID <= 0 {
			return client_error.ResponseError(c, stdhttp.StatusUnauthorized, nil, "unauthorized")
		}

		deptID, ok := utils.GetDeptIDInt(c)
		if !ok || deptID <= 0 {
			return client_error.ResponseError(c, stdhttp.StatusUnauthorized, nil, "unauthorized")
		}

		paramDeptID, err := utils.GetParamAsInt(c, deptIDFromPathParam)
		if err != nil || paramDeptID <= 0 {
			return client_error.ResponseError(c, stdhttp.StatusBadRequest, err, "invalid department id")
		}

		ok = paramDeptID == deptID

		if !ok {
			return client_error.ResponseError(c, stdhttp.StatusForbidden, nil, "forbidden: not a member of department")
		}
		return c.Next()
	}
}
