package sigurlx

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/drsigned/sigurlx/pkg/categorize"
	"github.com/drsigned/sigurlx/pkg/paramscan"
)

// Runner is a
type Runner struct {
	Options    *Options
	Categories categorize.Categories
	Params     []paramscan.Params
	Client     *http.Client
}

// Results is a
type Results struct {
	URL           string             `json:"url,omitempty"`
	Category      string             `json:"category,omitempty"`
	StatusCode    int                `json:"status_code,omitempty"`
	ContentType   string             `json:"content_type,omitempty"`
	ContentLength int64              `json:"content_length,omitempty"`
	List          []string           `json:"params_list,omitempty"`
	Risky         []paramscan.Params `json:"risky_params,omitempty"`
}

// New is a
func New(options *Options) (runner Runner, err error) {
	// Options
	runner.Options = options

	// Regex
	runner.Categories.STYLE, err = newRegex(`(?m).*?\.(css)(\?.*?|)$`)
	if err != nil {
		return runner, err
	}

	runner.Categories.JS, err = newRegex(`(?m).*?\.(js|json|xml|csv)(\?.*?|)$`)
	if err != nil {
		return runner, err
	}

	runner.Categories.DOC, err = newRegex(`(?m).*?\.(pdf|xlsx|doc|docx|txt)(\?.*?|)$`)
	if err != nil {
		return runner, err
	}

	runner.Categories.MEDIA, err = newRegex(`(?m).*?\.(jpg|jpeg|png|ico|svg|gif|webp|mp3|mp4|woff|woff2|ttf|eot|tif|tiff)(\?.*?|)$`)
	if err != nil {
		return runner, err
	}

	// Params
	raw, err := ioutil.ReadFile(runner.Options.ParamsPath)
	if err != nil {
		return runner, err
	}

	if err = json.Unmarshal(raw, &runner.Params); err != nil {
		return runner, err
	}

	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(runner.Options.Timeout) * time.Second,
			KeepAlive: time.Second,
		}).DialContext,
	}

	if !runner.Options.VerifyTLS {
		tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	if runner.Options.Proxy != "" {
		if p, err := url.Parse(runner.Options.Proxy); err == nil {
			tr.Proxy = http.ProxyURL(p)
		}
	}

	runner.Client = &http.Client{
		Timeout:   time.Duration(runner.Options.Timeout) * time.Second,
		Transport: tr,
	}

	return runner, nil
}

// Process is a
func (runner *Runner) Process(URL string) (results Results, err error) {
	results.URL = URL

	_, err = url.Parse(results.URL)
	if err != nil {
		return results, err
	}

	// Categorize
	if runner.Options.Categorize || runner.Options.All {
		results.Category, err = categorize.Run(URL, runner.Categories)
		if err != nil {
			return results, err
		}
	}

	// Scan Parameters
	if runner.Options.ScanParam || runner.Options.All {
		results.List, results.Risky, err = paramscan.Run(URL, runner.Params)
		if err != nil {
			return results, err
		}
	}

	// Request
	if runner.Options.Request || runner.Options.All {
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			return results, err
		}

		req.Header.Set("User-Agent", runner.Options.UserAgent)

		res, err := runner.Client.Do(req)
		if err != nil {
			return results, err
		}

		defer res.Body.Close()

		results.StatusCode = res.StatusCode
		results.ContentType = strings.Split(res.Header.Get("Content-Type"), ";")[0]
		results.ContentLength = res.ContentLength
	}

	return results, nil
}
