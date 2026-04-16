package embed

import (
	"github.com/apex/log"
	"github.com/google/uuid" // go get github.com/google/uuid
	"github.com/labstack/echo/v5"
)

// ApexLoggerMiddleware injects a scoped logger into the request context
func ApexLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		req := c.Request()
		res := c.Response()

		// 1. Get or Generate Request ID
		rid := req.Header.Get(echo.HeaderXRequestID)
		if rid == "" {
			rid = uuid.New().String()
		}
		res.Header().Set(echo.HeaderXRequestID, rid)

		// 2. Create a logger entry with the Request ID
		entry := log.WithField("request_id", rid)

		// 3. Put the logger into the standard context.Context
		ctx := log.NewContext(req.Context(), entry)
		c.SetRequest(req.WithContext(ctx))

		return next(c)
	}
}

// Logger is a helper function to extract the logger from the echo context
func Logger(c *echo.Context) *log.Entry {
	l, ok := log.FromContext(c.Request().Context()).(*log.Entry)
	if !ok {
		// Fallback to default logger if middleware didn't run
		return log.WithField("request_id", "unknown")
	}
	return l
}
