package util

import (
	"awsx-api/log"
)

func CommonError(err error) error {
	//fmt.Println("Error: ", err)
	log.Error("Error: ", err)
	return err
}
func DashboardError(err error, UID string) error {
	//fmt.Println("UID: "+UID+", Error: ", err)
	log.Error("Error: UID:\n%s\n%s", UID, err)
	return err
}
func Error(message string, err error) error {
	log.Error("%s. Error: %s ", message, err)
	return err
}
