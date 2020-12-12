package sigurlx

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/drsigned/gos"
	"github.com/drsigned/sigurlx/pkg/categorize"
	"github.com/drsigned/sigurlx/pkg/paramscan"
	"github.com/valyala/fasthttp"
)

// Runner is a
type Runner struct {
	Options    *Options
	Categories categorize.Categories
	Params     []paramscan.Params
	Client     *fasthttp.Client
}

// Results is a
type Results struct {
	URL           string             `json:"url,omitempty"`
	Category      string             `json:"category,omitempty"`
	StatusCode    int                `json:"status_code,omitempty"`
	ContentType   string             `json:"content_type,omitempty"`
	ContentLength int                `json:"content_length,omitempty"`
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

	runner.Categories.MEDIA, err = newRegex(`(?m).*?\.(jpg|jpeg|png|ico|svg|gif|webp|mp3|mp4|woff|woff2|ttf|eot)(\?.*?|)$`)
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

	// Client
	runner.Client = &fasthttp.Client{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return runner, nil
}

// Process is a
func (runner *Runner) Process(URL string) (results Results, err error) {
	parsedURL, err := gos.ParseURL(URL)
	if err != nil {
		return results, err
	}

	results.URL = parsedURL.URL.String()

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
		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()

		defer func() {
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(res)
		}()

		req.SetRequestURI(URL)
		req.Header.Add("UserAgent", runner.Options.UserAgent)
		req.Header.Add("Connection", "close")

		if err := runner.Client.Do(req, res); err != nil {
			return results, err
		}

		results.StatusCode = res.StatusCode()
		results.ContentType = strings.Split(string(res.Header.ContentType()), ";")[0]
		results.ContentLength = res.Header.ContentLength()
	}

	return results, nil
}
