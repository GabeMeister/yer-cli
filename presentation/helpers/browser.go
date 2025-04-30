package helpers

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetIntQueryParam(c echo.Context, name string) (int, error) {
	idParam := c.QueryParam(name)
	if idParam == "" {
		return 0, fmt.Errorf("unable to find query param: %s", name)
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, fmt.Errorf("invalid `%s` query param: %s", name, idParam)
	}

	return id, nil
}

func MustGetFormValue(c echo.Context, name string) string {
	val := c.FormValue(name)
	if val == "" {
		panic(fmt.Sprintf("\n\n\n*** MISSING FORM VALUE: %s ***\n\n\n", name))
	}

	return val
}
