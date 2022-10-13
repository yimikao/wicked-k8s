package main

import (
	"github.com/yimikao/wicked-k8s/conf"
)

func main() {
	conf.A = "added"
	conf.Do()
}
