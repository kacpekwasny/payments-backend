package funcs

import (
	"fmt"
	"time"
)

func Log(str string, values ...interface{}) {
	if str[0] == '\n' {
		fmt.Printf("\n"+time.Now().Format("02/01/2006 – 15:04:05 ")+str[1:]+"\n", values...)
	} else {
		fmt.Printf(time.Now().Format("02/01/2006 – 15:04:05 ")+str+"\n", values...)
	}
}

func formatTime2TimeStamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
