package authen

import 
	(
		"os"
		"github.com/labstack/echo/v4"
	)

func AuthMiddleware(username, password string, c echo.Context) (bool, error) {
	if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
		return true, nil
	}
	return false, nil
}