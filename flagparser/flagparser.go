package flagparser

import (
	"flag"
	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/dispatch"
	"github.com/thrasher-corp/gocryptotrader/engine"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
	_ "github.com/thrasher-corp/gocryptotrader/gctscript"
	gctscriptVM "github.com/thrasher-corp/gocryptotrader/gctscript/vm"
	"github.com/thrasher-corp/gocryptotrader/portfolio/withdraw"
	"runtime"
	"time"
)

// Handle flags

func flagParse() engine.Settings {
	//versionFlag := flag.Bool("version", false, "retrieves current GoCryptoTrader version")
	var settings engine.Settings
	// Core settings
	flag.StringVar(&settings.ConfigFile, "config", config.DefaultFilePath(), "config file to load")
	flag.StringVar(&settings.DataDir, "datadir", common.GetDefaultDataDir(runtime.GOOS), "default data directory for GoCryptoTrader files")
	flag.IntVar(&settings.GoMaxProcs, "gomaxprocs", runtime.GOMAXPROCS(-1), "sets the runtime GOMAXPROCS value")
	flag.BoolVar(&settings.EnableDryRun, "dryrun", false, "dry runs bot, doesn't save config file")
	flag.BoolVar(&settings.EnableAllExchanges, "enableallexchanges", false, "enables all exchanges")
	flag.BoolVar(&settings.EnableAllPairs, "enableallpairs", false, "enables all pairs for enabled exchanges")
	flag.BoolVar(&settings.EnablePortfolioManager, "portfoliomanager", true, "enables the portfolio manager")
	flag.DurationVar(&settings.PortfolioManagerDelay, "portfoliomanagerdelay", time.Duration(0), "sets the portfolio managers sleep delay between updates")
	flag.BoolVar(&settings.EnableGRPC, "grpc", true, "enables the grpc server")
	flag.BoolVar(&settings.EnableGRPCProxy, "grpcproxy", false, "enables the grpc proxy server")
	flag.BoolVar(&settings.EnableWebsocketRPC, "websocketrpc", true, "enables the websocket RPC server")
	flag.BoolVar(&settings.EnableDeprecatedRPC, "deprecatedrpc", true, "enables the deprecated RPC server")
	flag.BoolVar(&settings.EnableCommsRelayer, "enablecommsrelayer", true, "enables available communications relayer")
	flag.BoolVar(&settings.Verbose, "verbose", false, "increases logging verbosity for GoCryptoTrader")
	flag.BoolVar(&settings.EnableExchangeSyncManager, "syncmanager", true, "enables to exchange sync manager")
	flag.BoolVar(&settings.EnableWebsocketRoutine, "websocketroutine", true, "enables the websocket routine for all loaded exchanges")
	flag.BoolVar(&settings.EnableCoinmarketcapAnalysis, "coinmarketcap", false, "overrides config and runs currency analysis")
	flag.BoolVar(&settings.EnableEventManager, "eventmanager", true, "enables the event manager")
	flag.BoolVar(&settings.EnableOrderManager, "ordermanager", true, "enables the order manager")
	flag.BoolVar(&settings.EnableDepositAddressManager, "depositaddressmanager", true, "enables the deposit address manager")
	flag.BoolVar(&settings.EnableConnectivityMonitor, "connectivitymonitor", true, "enables the connectivity monitor")
	flag.BoolVar(&settings.EnableDatabaseManager, "databasemanager", true, "enables database manager")
	flag.BoolVar(&settings.EnableGCTScriptManager, "gctscriptmanager", true, "enables gctscript manager")
	flag.DurationVar(&settings.EventManagerDelay, "eventmanagerdelay", time.Duration(0), "sets the event managers sleep delay between event checking")
	flag.BoolVar(&settings.EnableNTPClient, "ntpclient", true, "enables the NTP client to check system clock drift")
	flag.BoolVar(&settings.EnableDispatcher, "dispatch", true, "enables the dispatch system")
	flag.IntVar(&settings.DispatchMaxWorkerAmount, "dispatchworkers", dispatch.DefaultMaxWorkers, "sets the dispatch package max worker generation limit")
	flag.IntVar(&settings.DispatchJobsLimit, "dispatchjobslimit", dispatch.DefaultJobsLimit, "sets the dispatch package max jobs limit")

	// Exchange syncer settings
	flag.BoolVar(&settings.EnableTickerSyncing, "tickersync", true, "enables ticker syncing for all enabled exchanges")
	flag.BoolVar(&settings.EnableOrderbookSyncing, "orderbooksync", true, "enables orderbook syncing for all enabled exchanges")
	flag.BoolVar(&settings.EnableTradeSyncing, "tradesync", false, "enables trade syncing for all enabled exchanges")
	flag.IntVar(&settings.SyncWorkers, "syncworkers", engine.DefaultSyncerWorkers, "the amount of workers (goroutines) to use for syncing exchange data")
	flag.BoolVar(&settings.SyncContinuously, "synccontinuously", true, "whether to sync exchange data continuously (ticker, orderbook and trade history info")
	flag.DurationVar(&settings.SyncTimeout, "synctimeout", engine.DefaultSyncerTimeout,
		"the amount of time before the syncer will switch from one protocol to the other (e.g. from REST to websocket)")

	// Forex provider settings
	flag.BoolVar(&settings.EnableCurrencyConverter, "currencyconverter", false, "overrides config and sets up foreign exchange Currency Converter")
	flag.BoolVar(&settings.EnableCurrencyLayer, "currencylayer", false, "overrides config and sets up foreign exchange Currency Layer")
	flag.BoolVar(&settings.EnableFixer, "fixer", false, "overrides config and sets up foreign exchange Fixer.io")
	flag.BoolVar(&settings.EnableOpenExchangeRates, "openexchangerates", false, "overrides config and sets up foreign exchange Open Exchange Rates")

	// Exchange tuning settings
	flag.BoolVar(&settings.EnableExchangeAutoPairUpdates, "exchangeautopairupdates", true, "enables automatic available currency pair updates for supported exchanges")
	flag.BoolVar(&settings.DisableExchangeAutoPairUpdates, "exchangedisableautopairupdates", false, "disables exchange auto pair updates")
	flag.BoolVar(&settings.EnableExchangeWebsocketSupport, "exchangewebsocketsupport", false, "enables Websocket support for exchanges")
	flag.BoolVar(&settings.EnableExchangeRESTSupport, "exchangerestsupport", true, "enables REST support for exchanges")
	flag.BoolVar(&settings.EnableExchangeVerbose, "exchangeverbose", false, "increases exchange logging verbosity")
	flag.BoolVar(&settings.ExchangePurgeCredentials, "exchangepurgecredentials", false, "purges the stored exchange API credentials")
	flag.BoolVar(&settings.EnableExchangeHTTPRateLimiter, "ratelimiter", true, "enables the rate limiter for HTTP requests")
	flag.IntVar(&settings.MaxHTTPRequestJobsLimit, "requestjobslimit", int(request.DefaultMaxRequestJobs), "sets the max amount of jobs the HTTP request package stores")
	flag.IntVar(&settings.RequestMaxRetryAttempts, "httpmaxretryattempts", request.DefaultMaxRetryAttempts, "sets the number of retry attempts after a retryable HTTP failure")
	flag.DurationVar(&settings.HTTPTimeout, "httptimeout", time.Duration(0), "sets the HTTP timeout value for HTTP requests")
	flag.StringVar(&settings.HTTPUserAgent, "httpuseragent", "", "sets the HTTP user agent")
	flag.StringVar(&settings.HTTPProxy, "httpproxy", "", "sets the HTTP proxy server")
	flag.BoolVar(&settings.EnableExchangeHTTPDebugging, "exchangehttpdebugging", false, "sets the exchanges HTTP debugging")

	// Common tuning settings
	flag.DurationVar(&settings.GlobalHTTPTimeout, "globalhttptimeout", time.Duration(0), "sets common HTTP timeout value for HTTP requests")
	flag.StringVar(&settings.GlobalHTTPUserAgent, "globalhttpuseragent", "", "sets the common HTTP client's user agent")
	flag.StringVar(&settings.GlobalHTTPProxy, "globalhttpproxy", "", "sets the common HTTP client's proxy server")

	// GCTScript tuning settings
	flag.UintVar(&settings.MaxVirtualMachines, "maxvirtualmachines", uint(gctscriptVM.DefaultMaxVirtualMachines), "set max virtual machines that can load")

	// Withdraw Cache tuning settings
	flag.Uint64Var(&settings.WithdrawCacheSize, "withdrawcachesize", withdraw.CacheSize, "set cache size for withdrawal requests")

	//if *versionFlag {
	//	fmt.Print(core.Version(false))
	//	os.Exit(0)
	//}
	flag.Parse()
	return settings
}

func DefaultEngineSettings() (engine.Settings, map[string]bool) {
	//go gctscript.Setup()
	settings := flagParse()

	//var err error
	settings.CheckParamInteraction = true

	// collect flags
	flagSet := make(map[string]bool)

	// Stores the set flags
	flag.Visit(func(f *flag.Flag) { flagSet[f.Name] = true })
	if !flagSet["config"] {
		// If config file is not explicitly set, fall back to default path resolution
		//settings.ConfigFile = "/home/romano/.gocryptotrader/config.json"
		settings.ConfigFile = ""
	}
	return settings, flagSet
}
