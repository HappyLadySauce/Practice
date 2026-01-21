package common

import "log"

func DealError(err error) {
	if err != nil {
		log.Panicf("err: %v\n", err)
	}
}