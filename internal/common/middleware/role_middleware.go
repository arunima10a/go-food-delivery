package middleware

import (
	"net/http"
	"github.com/golang-jwt/jwt/v5"
	"fmt"

	"github.com/labstack/echo/v4"
)
func RequireRole(requiredRole string) echo.MiddlewareFunc  {
	return func(next echo.HandlerFunc) echo.HandlerFunc{
		return func( c echo.Context) error {
			userValue := c.Get("user")
			if userValue == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid token")
			}
			token, ok := userValue.(*jwt.Token)
				if !ok {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format")
				}
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					return echo.NewHTTPError(http.StatusUnauthorized,"Invalid token claims")
				}
				role, ok := claims["role"].(string)
				if !ok || role != requiredRole{
					return echo.NewHTTPError(http.StatusForbidden, "Role not found in token")
				}
				fmt.Printf("SECURITY CHECK: User has role [%s], but we need [%s]\n", role, requiredRole) 

				if role != requiredRole {
					return echo.NewHTTPError(http.StatusForbidden, "You do not have permission (Admin only)")
				}
				return next(c)
			}
		}
	}
	
