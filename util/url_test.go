package util

import (
	"fmt"
	"testing"
)

func TestGetPreAndFileName(t *testing.T) {
	fmt.Println(GetPreAndFileName("/tracinglogger/logrus.go"))
	fmt.Println(GetPreAndFileName("tracinglogger/logrus.go"))
	fmt.Println(GetPreAndFileName("/logrus.go"))
	fmt.Println(GetPreAndFileName("logrus.go"))
	fmt.Println(GetPreAndFileName(""))
	fmt.Println(GetPreAndFileName("/Users/liumin/GolandProjects/jaeger/tracinglogger/logrus.go"))
}
