package relay

import "github.com/sorcix/irc"

// StdLogger ...
type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

func log(logger StdLogger, s string) {
	if logger == nil {
		return
	}

	logger.Print(s)
}

func logf(logger StdLogger, s string, p ...interface{}) {
	if logger == nil {
		return
	}

	logger.Printf(s, p...)
}

func logMessage(logger StdLogger, message *irc.Message, out, s string) {
	logf(logger, "[%s]%v[%s] - %s ", message.Command, message.Params, out, s)
}

func logOutgoingMessage(logger StdLogger, message *irc.Message) {
	if message.EmptyTrailing {
		logMessage(logger, message, "OUT", "Sended")
	} else {
		logMessage(logger, message, "OUT", message.Trailing)
	}
}

func logIncomingMessage(logger StdLogger, message *irc.Message) {
	if message.EmptyTrailing {
		logMessage(logger, message, "INC", "Received")
	} else {
		logMessage(logger, message, "INC", message.Trailing)
	}
}
