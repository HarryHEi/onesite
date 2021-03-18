package log

import "go.uber.org/zap"

var _logger *zap.Logger

func InitLogger() (err error) {
	if _logger != nil {
		return nil
	}

	_logger, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}

func Sync() {
	_ = _logger.Sync()
}

func Debug(msg string, fields ...zap.Field) {
	_logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	_logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	_logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	_logger.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	_logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	_logger.Fatal(msg, fields...)
}
