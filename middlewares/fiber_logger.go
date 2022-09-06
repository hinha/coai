package middlewares

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/hinha/coai/config"
	"github.com/hinha/coai/logger"
	"github.com/hinha/coai/utils"
	"github.com/hinha/coai/utils/security"
)

// Config defines the config for middleware
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Logger defines zap logger instance
	Logger *zap.Logger

	// AppConfig defines config by yaml
	AppConfig *config.Config
}

// New creates a new middleware handler
func NewLogger(config Config) fiber.Handler {
	var (
		errPadding  = 15
		start, stop time.Time
		once        sync.Once
		errHandler  fiber.ErrorHandler
	)

	if config.AppConfig == nil {
		panic("must required config")
	}

	cipher := security.NewCipher(config.AppConfig.Server.SecretKey)
	return func(c *fiber.Ctx) error {
		if config.Next != nil && config.Next(c) {
			return c.Next()
		}

		once.Do(func() {
			errHandler = c.App().Config().ErrorHandler
			stack := c.App().Stack()
			for m := range stack {
				for r := range stack[m] {
					if len(stack[m][r].Path) > errPadding {
						errPadding = len(stack[m][r].Path)
					}
				}
			}
		})

		start = time.Now()
		chainErr := c.Next()

		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		stop = time.Now()

		var fields []zap.Field
		fields = append(fields,
			zap.Namespace("context"),
			zap.String("pid", strconv.Itoa(os.Getpid())),
			zap.String("time", stop.Sub(start).String()),
		)
		if config.AppConfig.Log.Encrypted {
			log, err := logEncrypted(cipher, c)
			if err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
			fields = append(fields, log...)
		} else {
			fields = append(fields,
				zap.Object("response", logger.Resp(c.Response())),
				zap.Object("request", logger.Req(c)),
			)
		}

		if u := c.Locals("userId"); u != nil {
			fields = append(fields, zap.Uint("userId", u.(uint)))
		}

		formatErr := ""
		if chainErr != nil {
			formatErr = chainErr.Error()
			fields = append(fields, zap.String("error", formatErr))
			config.Logger.With(fields...).Error(formatErr)

			return nil
		}

		if c.Response().StatusCode() < 200 || c.Response().StatusCode() > 299 {
			config.Logger.With(fields...).Warn(c.Route().Name)
		} else {
			config.Logger.With(fields...).Info(c.Route().Name)
		}

		return nil
	}
}

func logEncrypted(cipher utils.Cipher, c *fiber.Ctx) ([]zap.Field, error) {
	bRequest, _ := json.Marshal(logger.Req(c))
	request, err := cipher.EncryptText(string(bRequest))
	if err != nil {
		return nil, err
	}

	bResponse, _ := json.Marshal(logger.Resp(c.Response()))
	response, err := cipher.EncryptText(string(bResponse))
	if err != nil {
		return nil, err
	}

	return []zap.Field{
		zap.String("request", request),
		zap.String("response", response),
	}, nil
}
