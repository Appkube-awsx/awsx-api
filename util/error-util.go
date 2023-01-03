package util

import (
	"awsx-api/log"
	"fmt"
)

func CommonError(err error) error {
	fmt.Println("Error: ", err)
	log.Fatal(err)
	return err
}
func DashboardError(err error, UID string) error {
	fmt.Println("UID: "+UID+", Error: ", err)
	log.Fatal("Error: UID:\n%s\n%s", UID, err)
	return err
}
