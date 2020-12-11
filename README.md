# sigurlx

![made with go](https://img.shields.io/badge/made%20with-Go-0040ff.svg) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0040ff.svg) [![open issues](https://img.shields.io/github/issues-raw/drsigned/sigurlx.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigurlx/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/drsigned/sigurlx.svg?style=flat&color=0040ff)](https://github.com/drsigned/sigurlx/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/license-MIT-gray.svg?colorB=0040FF)](https://github.com/drsigned/sigurlx/blob/master/LICENSE) [![twitter](https://img.shields.io/badge/twitter-@drsigned-0040ff.svg)](https://twitter.com/drsigned)

sigurlx is a helper utility written in Go. It takes urls from stdin then:

* categorize them into:
    * endpoint
    * style {css}
    * js {js|json|xml|csv}
    * doc {pdf|doc|docx|xlsx}
    * media {jpg|jpeg|png|ico|svg|gif|webp|mp3|mp4|woff|woff2|tff}
* checks for risky parameter names and the risks commonly associated with them
* send HTTP request

## Resources

* [Installation](#installation)
    * [From Binary](#from-binary)
    * [From source](#from-source)
    * [From github](#from-github)
* [Usage](#usage)
* [Credits](#credits)
* [Contribution](#contribution)

## Installation

#### From Binary

You can download the pre-built binary for your platform from this repository's [releases](https://github.com/drsigned/sigurlx/releases/) page, extract, then move it to your `$PATH`and you're ready to go.

#### From Source

sigurlx requires **go1.14+** to install successfully. Run the following command to get the repo

```bash
$ GO111MODULE=on go get -u -v github.com/drsigned/sigurlx/cmd/sigurlx
```

#### From Github

```bash
$ git clone https://github.com/drsigned/sigurlx.git; cd sigurlx/cmd/sigurlx/; go build; mv sigurlx /usr/local/bin/; sigurlx -h
```

## Usage

To display help message for sigurlx use the `-h` flag:

```
$ sigurlx -h

     _                  _      
 ___(_) __ _ _   _ _ __| |_  __
/ __| |/ _` | | | | '__| \ \/ /
\__ \ | (_| | |_| | |  | |>  < 
|___/_|\__, |\__,_|_|  |_/_/\_\ v1.0.0
       |___/

Usage:
  sigurlx [OPTIONS]

TASK OPTIONS:
   -cat                       categorize urls
   -param-scan                scan url parameters
   -request                   send HTTP request

GENERAL OPTIONS:
   -t                         number of concurrent threads. (default: 50)
   -nc                        no color mode
   -s                         silent mode
   -v                         verbose mode

REQUEST OPTIONS (used with -request):
   -UA                        HTTP user agent

OUTPUT OPTIONS:
   -o                         output file

```

**NOTE:** The output is in JSON format. Use `jq` to process the results.

## Credits

The list of parameter names and the riskss associated with them is mainly created from the public work of various people of the community - inital list was obtained from [Somdev Sangwan](https://github.com/s0md3v)'s [Parth](https://github.com/s0md3v/Parth) .

## Contribution

[Issues](https://github.com/drsigned/sigurlx/issues) and [Pull Requests](https://github.com/drsigned/sigurlx/pulls) are welcome!