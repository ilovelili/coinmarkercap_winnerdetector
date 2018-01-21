package main

import (
	"bytes"
	"config"
	"context"
	"html/template"
	"os"
	"path"
	"strconv"
	"utils"

	"github.com/ilovelili/go_coinmarketcap"

	. "github.com/ahmetb/go-linq"
)

// Domain gainer or loser
type Domain int

const (
	// Gainer gainer
	Gainer Domain = 0
	// Loser loser
	Loser Domain = 1
)

var (
	conf               = config.GetConfig()
	gainerMetaFileName = "gainer.meta"
	loserMetaFileName  = "loser.meta"
)

type mailItem struct {
	Tickers []*coinmarketcap.Ticker
}

func main() {
	// step1. get all tickers
	client, _ := coinmarketcap.NewClient("", "")
	ctx := context.Background()
	tickers, err := client.GetTickers(ctx)
	if err != nil {
		panic(err)
	}

	// step2. write mail if there are candidates
	gainers := getCandidates(Gainer, tickers)
	if len(gainers) > 0 {
		writeLast(gainerMetaFileName, gainers)
		sendMail("[AgentSmith]: Soaring in one hour", &mailItem{Tickers: gainers})
	}

	losers := getCandidates(Loser, tickers)
	if len(losers) > 0 {
		writeLast(loserMetaFileName, losers)
		sendMail("[AgentSmith]: Falling in one hour", &mailItem{Tickers: losers})
	}
}

func getCandidates(domain Domain, tickers []*coinmarketcap.Ticker) []*coinmarketcap.Ticker {
	var candidates []*coinmarketcap.Ticker
	thresholdFilter := getTickerThresholdFilter(domain)
	mailFilter := getTickerMailSentFilter(domain)

	From(tickers).WhereT(thresholdFilter).Take(3).WhereT(mailFilter).ToSlice(&candidates)
	return candidates
}

func getTickerThresholdFilter(domain Domain) func(*coinmarketcap.Ticker) bool {
	var filter func(*coinmarketcap.Ticker) bool

	switch domain {
	case Gainer:
		filter = func(ticker *coinmarketcap.Ticker) bool {
			percentchange, _ := strconv.ParseFloat(ticker.PercentChangeOneHour, 64)
			return percentchange >= conf.Max
		}
	case Loser:
		filter = func(ticker *coinmarketcap.Ticker) bool {
			percentchange, _ := strconv.ParseFloat(ticker.PercentChangeOneHour, 64)
			return percentchange <= conf.Min
		}
	default:
		filter = func(ticker *coinmarketcap.Ticker) bool {
			return true
		}
	}

	return filter
}

func getTickerMailSentFilter(domain Domain) func(*coinmarketcap.Ticker) bool {
	var metafile string
	switch domain {
	case Gainer:
		metafile = gainerMetaFileName
	case Loser:
		metafile = loserMetaFileName
	}

	lines, _ := readLast(metafile)

	return func(ticker *coinmarketcap.Ticker) bool {
		id := ticker.ID
		for _, line := range lines {
			if line == id {
				return false
			}
		}
		return true
	}
}

func readLast(filename string) ([]string, error) {
	path := path.Join(utils.ResolveOutputDir(), filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	return utils.ReadFile(path)
}

func writeLast(filename string, items []*coinmarketcap.Ticker) error {
	var lines []string
	From(items).SelectT(func(ticker *coinmarketcap.Ticker) string {
		return ticker.ID
	}).ToSlice(&lines)

	path := path.Join(utils.ResolveOutputDir(), filename)
	return utils.WriteFile(path, lines)
}

func sendMail(subject string, mailitem *mailItem) error {
	t := template.Must(template.New("email.templ").ParseFiles("template/email.templ"))
	var tpl bytes.Buffer
	err := t.Execute(&tpl, mailitem)
	if err != nil {
		return err
	}

	return utils.SendMail(conf, subject, tpl.String())
}
