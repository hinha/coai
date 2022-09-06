package logger

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap/zapcore"
)

func getAllowedHeaders() map[string]bool {
	return map[string]bool{
		"User-Agent": true,
		"X-Mobile":   true,
	}
}

type resp struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Body []byte `json:"body"`
}

func Resp(r *fasthttp.Response) *resp {
	return &resp{
		Code: r.StatusCode(),
		Type: bytes.NewBuffer(r.Header.ContentType()).String(),
		Body: r.Body(),
	}
}

func (r *resp) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("type", r.Type)
	enc.AddInt("code", r.Code)

	var body interface{}
	if err := json.Unmarshal(r.Body, &body); err != nil {
		return err
	}
	if err := enc.AddReflected("body", body); err != nil {
		return err
	}

	return nil
}

type req struct {
	Body     string     `json:"body"`
	Fullpath string     `json:"fullPath"`
	User     string     `json:"user"`
	IP       string     `json:"ip"`
	Method   string     `json:"method"`
	Route    string     `json:"router"`
	Headers  *headerbag `json:"header"`
}

func Req(c *fiber.Ctx) *req {
	reqq := c.Request()
	var body []byte
	buffer := new(bytes.Buffer)
	err := json.Compact(buffer, reqq.Body())
	if err != nil {
		body = reqq.Body()
	} else {
		body = buffer.Bytes()
	}

	headers := &headerbag{
		vals: make(map[string]string),
	}
	allowedHeaders := getAllowedHeaders()
	reqq.Header.VisitAll(func(key, val []byte) {
		k := bytes.NewBuffer(key).String()
		if _, exist := allowedHeaders[k]; exist {
			headers.vals[strings.ToLower(k)] = bytes.NewBuffer(val).String()
		}
	})

	var userEmail string
	if u := c.Locals("userEmail"); u != nil {
		userEmail = u.(string)
	}

	return &req{
		Body:     bytes.NewBuffer(body).String(),
		Fullpath: bytes.NewBuffer(reqq.RequestURI()).String(),
		Headers:  headers,
		IP:       c.IP(),
		Method:   c.Method(),
		Route:    c.Route().Path,
		User:     userEmail,
	}
}

func (r *req) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("fullPath", r.Fullpath)
	enc.AddString("ip", r.IP)
	enc.AddString("method", r.Method)
	enc.AddString("route", r.Route)

	if r.Body != "" {
		enc.AddString("body", r.Body)
	}

	if r.User != "" {
		enc.AddString("user", r.User)
	}

	err := enc.AddObject("headers", r.Headers)
	if err != nil {
		return err
	}

	return nil
}

type headerbag struct {
	vals map[string]string
}

func (h *headerbag) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range h.vals {
		enc.AddString(k, v)
	}

	return nil
}
