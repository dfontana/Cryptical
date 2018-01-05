package bot

import (
	"errors"
	"time"
)

// Portfolio describes the amount of currency and its USD value.
type Portfolio struct {
	CurrencyPair string
	History      []PfEntry
}

// PfEntry describes the portfolio at a single moment in time.
type PfEntry struct {
	Crypto       float64
	Pair         float64
	NetWorth     float64
	CurrentValue float64
	Time         time.Time
}

// Initialize is used to prepare an empty portfolio with its initial contents
// of cryptocurrency and pair amounts. The current values of 1 Cryptocurrency
// to pairing must be supplied as well.
func (p *Portfolio) Initialize(cryptoAmt, pairAmt, currentVal float64) error {
	if len(p.History) > 0 {
		return errors.New("cannot reinitalize an initalize portfolio")
	}
	e := PfEntry{
		Crypto:       cryptoAmt,
		Pair:         pairAmt,
		NetWorth:     pairAmt + cryptoAmt*currentVal,
		CurrentValue: currentVal,
		Time:         time.Now(),
	}
	p.History = append(p.History, e)
	return nil
}

// Update the portfolio with the amount of cryptocurrency bought or sold (Sign matters),
// at the current value for that currency against its pair (ie 1 ETH = $800 USD). This
// assumes the currency was purchase or sold against the portfolio's balances, and thus
// will update the pair value based on the provided information. Returns the new entry.
func (p *Portfolio) Update(cryptoDiff, currentVal float64) (PfEntry, error) {
	if len(p.History) == 0 {
		return PfEntry{}, errors.New("must initialize portfolio before updating")
	}

	entry, _ := p.Latest()
	pAmt := entry.Pair - cryptoDiff*currentVal
	cAmt := entry.Crypto + cryptoDiff

	e := PfEntry{
		Crypto:       cAmt,
		Pair:         pAmt,
		NetWorth:     pAmt + cAmt*currentVal,
		CurrentValue: currentVal,
		Time:         time.Now(),
	}
	p.History = append(p.History, e)
	return e, nil
}

// Latest returns the latest status of your portfolio.
func (p *Portfolio) Latest() (PfEntry, error) {
	if len(p.History) == 0 {
		return PfEntry{}, errors.New("Portfolio is empty")
	}
	return p.History[len(p.History)-1], nil
}
