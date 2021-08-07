package webserver

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/engine"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/ftx"
	"net/http"
	"strconv"
	"strings"
)

func (withdraw WithdrawResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (history WithdrawHistoryResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type WithdrawHistoryResponse struct {
	History []exchange.WithdrawalHistory `json:"history"`
}

type WithdrawResponse struct {
	TransactionData interface{} `json:"transactionData"`
	Error           FtxErr      `json:"error"`
}

type FtxErr struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (f FtxErr) Render(w http.ResponseWriter, r *http.Request) error {
	f.Success = false
	return nil
}

func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "withdraw.html", nil)
	if err != nil {
		logrus.Errorf("error template: %s\n", err)
	}
}

func getWithdrawalInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	withdraw, ok := ctx.Value("withdraw").(WithdrawResponse)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
		logrus.Errorf("failed withdraw response CTX: %v", withdraw)
		withdraw.Error.Render(w, r)
		return
	}
	render.Render(w, r, withdraw)
	return
}

// get withdrawal history from exchange from an asset
func getWithdrawHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	withdrawHistory, ok := ctx.Value("history").([]exchange.WithdrawalHistory)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
		render.Render(w, r, ErrNotFound)
		return
	}
	history := WithdrawHistoryResponse{
		History: withdrawHistory,
	}
	render.Render(w, r, history)
	return
}

// TODO fix FTX ERRO[0006] failed fetch history: not yet implemented
// only works for binance
func withdrawHistoryCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		var ctx context.Context
		exchangeNameReq := chi.URLParam(request, "exchange")
		exchangeEngine := engine.Bot.GetExchangeByName(exchangeNameReq)
		code := currency.NewCode(strings.ToUpper(chi.URLParam(request, "asset")))
		history, err := exchangeEngine.GetWithdrawalsHistory(code)
		if err != nil {
			logrus.Errorf("failed fetch history: %s", err)
		}
		ctx = context.WithValue(request.Context(), "history", history)
		next.ServeHTTP(w, request.WithContext(ctx))
	})
}

func WithdrawCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		var err error
		assetInfo := new(Asset)
		var ftxBase *ftx.FTX
		var ctx context.Context

		exchangeNameReq := chi.URLParam(request, "exchange")
		logrus.Info("received withdrawCtx request")
		//asset := chi.URLParam(request, "asset")
		destinationAddress := chi.URLParam(request, "destinationAddress")
		sizeReq := chi.URLParam(request, "size")
		size, err := strconv.ParseFloat(sizeReq, 64)
		if err != nil {
			logrus.Errorf("") // 3.14159265
		}

		assetInfo.Code = currency.NewCode(strings.ToUpper(chi.URLParam(request, "asset")))
		exchangeEngine := engine.Bot.GetExchangeByName(exchangeNameReq)

		var withdrawResponse = WithdrawResponse{
			TransactionData: nil,
			Error:           FtxErr{Success: true},
		}

		switch chi.URLParam(request, "exchange") {
		case "ftx":
			b := exchangeEngine.GetBase()
			ftxBase = &ftx.FTX{Base: *b}
			withdrawResponse.TransactionData, err = ftxBase.Withdraw(assetInfo.Code.String(), destinationAddress, "", "", assetInfo.Code.String(), size)
			if err != nil {
				withdrawResponse.Error = FtxErr{Success: false, Message: fmt.Sprintf("failed to withdraw %s %s: %s", assetInfo.Code, sizeReq, err.Error())}
				logrus.Errorf("withdraw success %t: %s\n", withdrawResponse.Error.Success, withdrawResponse.Error.Message)
			}
		case "binance":
			//exchangeEngine = engine.Bot.GetExchangeByName(exchangeNameReq)
			//b := exchangeEngine.GetBase()
			//base := &binance.Binance{Base: *b}
		default:
			logrus.Errorf("exchange not supported: %s", err)
		}
		logrus.Infof("withdraw: %v\n", withdrawResponse.TransactionData)
		ctx = context.WithValue(request.Context(), "withdraw", withdrawResponse)
		next.ServeHTTP(w, request.WithContext(ctx))
	})
}
