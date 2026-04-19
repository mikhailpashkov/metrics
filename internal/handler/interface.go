package handler

import (
	"net/http"

	"go.uber.org/zap"
)

type MHandler interface {
	GetLogger() *zap.Logger // fixme: шляпа какая-то)) может через context можно как-то это пробрасывать?
	GetUrlPattern() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
