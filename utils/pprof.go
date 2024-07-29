package utils

import (
	_ "net/http/pprof"
)

//func StartUpPprof(cfg *config.Config) {
//	if cfg.EnablePprof {
//		go func() {
//			fmt.Println(http.ListenAndServe(cfg.PprofAddr, nil))
//		}()
//	}
//}
