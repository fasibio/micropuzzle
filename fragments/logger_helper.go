package fragments

import "go.uber.org/zap"

func getLoggerWithUserInfo(logs *zap.SugaredLogger, user WebSocketUser) *zap.SugaredLogger {
	return logs.With("streamid", user.Id, "address", user.RemoteAddr)
}
