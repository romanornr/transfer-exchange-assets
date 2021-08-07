// Copyright (c) 2021 Romano
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package main

import (
	_ "context"
	"github.com/romanornr/toolkit/engineManager"
	"github.com/romanornr/toolkit/webserver"
	"github.com/thrasher-corp/gocryptotrader/gctscript"
	gctlog "github.com/thrasher-corp/gocryptotrader/log"
	"github.com/thrasher-corp/gocryptotrader/signaler"
)

func init() {
	go gctscript.Setup()
}

func main() {
	go func() {
		engineManager.StartMainEngine()
		go webserver.New()
	}()
	interrupt := signaler.WaitForInterrupt()
	gctlog.Infof(gctlog.Global, "Captured %v, shutdown requested.\n", interrupt)
}
