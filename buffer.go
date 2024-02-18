package main

import (
	"math/rand"
	"strconv"
	"time"
)

var RawCommand = map[string]string{
	"EOL":          "\r\n",
	"ESC":          "\x1b",
	"UEL":          "\x1b%-12345X",
	"DELIMITER":    "DELIMITER",
	"RAND_NUM":     rand_num(),
	"RAND_NUM_PCL": rand_num_PCL(),
	"PS_HEADER":    "@PJL ENTER LANGUAGE = POSTSCRIPT\n%!\n",
	"PS_IOHACK":    "/print {(%stdout) (w) file dup 3 2 roll writestring flushfile} def\n/== {128 string cvs print (\\n) print} def\n",
	"PCL_HEADER":   "@PJL ENTER LANGUAGE = PCL\r\n\x1b",
}

func rand_num() string {
	rand.Seed(time.Now().UnixNano())
	num_ := rand.Intn(55536) + 10000 //PS & PJL token: 10000-65535
	num := strconv.Itoa(num_)
	return num
}

func rand_num_PCL() string {
	rand.Seed(time.Now().UnixNano())
	num_ := rand.Intn(32513) - 32768 //PCL token: -256..-32767
	num := strconv.Itoa(num_)
	return num
}
