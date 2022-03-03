package errors

import log "github.com/sirupsen/logrus"

func checkError(msg string, er error) {
	if er != nil {
		log.Error(msg, er.Error())
	}
}

func checkErrorF(msg string, er error) {
	if er != nil {
		log.Fatalf(msg, er.Error())
	}
}

func ErrorRepository(msg string, er error) {
	checkError("[REPOSITORY] "+msg, er)
}
func ErrorRepositoryF(msg string, er error) {
	checkErrorF("[REPOSITORY] "+msg, er)
}

func ErrorConfig(msg string, er error) {
	checkError("[CONFIG] "+msg, er)
}
