//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package main

import (
	"reflect"
	"runtime"

	"github.com/sfmunoz/logit"
	ex1 "github.com/sfmunoz/logit/cmd/logit-demo/example1"
	ex2 "github.com/sfmunoz/logit/cmd/logit-demo/example2"
	ex3 "github.com/sfmunoz/logit/cmd/logit-demo/example3"
	ex4 "github.com/sfmunoz/logit/cmd/logit-demo/example4"
)

func funcName(fn any) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func main() {
	var log = logit.Logit()
	examples := []func(){ex1.Run, ex2.Run, ex3.Run, ex4.Run}
	for _, f := range examples {
		fName := funcName(f)
		log.Info("================ " + fName + " ================")
		f()
		log.Info("---------------- " + fName + " ----------------")
	}
}
