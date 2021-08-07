// Copyright (c) 2021 Romano
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package webserver

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/engine"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"net/http"
	"strings"
	"time"
)

type DepositResponse struct {
	*Asset
	Time time.Time `json:"time"`
}

func (d DepositResponse) Render(w http.ResponseWriter, r *http.Request) error {
	d.Time = time.Now()
	return nil
}

func DepositHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "deposit.html", nil)
	if err != nil {
		logrus.Errorf("error template: %s\n", err)
	}
}

func getDepositAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	assetInfo, ok := ctx.Value("depositInfo").(*Asset)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		render.Render(w, r, ErrNotFound)
		return
	}
	depositInfo := DepositResponse{assetInfo, time.Now()}
	render.Render(w, r, depositInfo)
	return
}

func DepositAddressCtx(next http.Handler) http.Handler {
	var err error
	assetInfo := new(Asset)

	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		exchangeNameReq := chi.URLParam(request, "exchange")
		exchange := engine.Bot.GetExchangeByName(exchangeNameReq)
		assetInfo.Exchange = exchange.GetName()
		assetInfo.Code = currency.NewCode(strings.ToUpper(chi.URLParam(request, "asset")))
		assetInfo.Address, err = exchange.GetDepositAddress(assetInfo.Code, "")
		if err != nil {
			logrus.Errorf("failed to get address: %s\n", err)
			render.Render(w, request, ErrInvalidRequest(err))
			return
		}

		portfolio := engine.Bot.Portfolio.GetPortfolioByExchange(exchange.GetName())

		for p := range portfolio {
			if p.Match(assetInfo.Code.Upper()) {
				assetInfo.Item = *p.Item
				assetInfo.AssocChain = p.Item.AssocChain
			}
		}

		accountInfo, err := exchange.FetchAccountInfo(asset.Spot)
		if err != nil {
			logrus.Errorf("failed to fetch account info: %s\n", err)
			render.Render(w, request, ErrInvalidRequest(err))
			return
		}

		if len(accountInfo.Accounts) >= 0 {
			for _, i := range accountInfo.Accounts[0].Currencies {
				if i.CurrencyName == assetInfo.Code {
					assetInfo.Balance = fmt.Sprintf("%f", i.TotalValue)
					engine.Bot.Portfolio.AddExchangeAddress(exchange.GetName(), assetInfo.Code, i.TotalValue)
				}
			}
		}
		ctx := context.WithValue(request.Context(), "depositInfo", assetInfo)
		next.ServeHTTP(w, request.WithContext(ctx))
	})
}
