package exchange

var ExchangeManager *exchangeManager

func init() {
	ExchangeManager = newExchangeManager()
}

type exchangeManager struct {
}

func newExchangeManager() *exchangeManager {
	return &exchangeManager{}
}

func (mgr *exchangeManager) RegisterExchange() {

}
