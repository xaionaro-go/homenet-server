package api

type Option interface {
	IsAPIOptionType()
}
type Options []Option

type option struct{}

func (opt option) IsAPIOptionType() {}

type optSetLogger struct {
	option

	logger logger
}

func (opt *optSetLogger) GetLogger() logger {
	return opt.logger
}
func (opt *optSetLogger) setLogger(log logger) *optSetLogger {
	opt.logger = log
	return opt
}

type optSetLoggerDebug struct {
	optSetLogger
}

func OptSetLoggerDebug(log logger) Option {
	result := &optSetLoggerDebug{}
	result.setLogger(log)
	return result
}

type optSetLoggerInfo struct {
	optSetLogger
}

func OptSetLoggerInfo(log logger) Option {
	result := &optSetLoggerInfo{}
	result.setLogger(log)
	return result
}
