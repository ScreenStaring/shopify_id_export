package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	shopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/screenstaring/shopify_id_export/exportformat"
)

const version = "v0.0.5"
const shopifyFields = "id,title,product_type,handle,variants"

const usage = `shopify_id_export [hjv] [-k key] [-p password] [-r root-property] [-t token] shop

Options
-h --help                display this help message
-j --json                output dump as JSON
-k --key       key       Shopify API key; defaults to the SHOPIFY_API_KEY environment variable
-p --password  password  Shopify API password; defaults to the SHOPIFY_API_PASSWORD environment variable
-r --json-root property  use property as the top-level property for each JSON object
--timeout      integer   set Shopify client timeout (default: 10 seconds)
-s --size      integer   page size to use when retrieving products (default: 250)
-t --token     token     Shopify API token; defaults to the SHOPIFY_API_TOKEN environment variable
-v --version             display version information
--verbose                output Shopify API request/response (default: false)

By default data is output to a CSV file.

Valid properties for the --json-root option are: %s
`

type dumper interface {
	Dump(shopify.Product) error
	Close() error
}

func exitFailure(error string, code int) {
	fmt.Fprintln(os.Stderr, error)
	os.Exit(code)
}

func dumpProducts(client *shopify.Client, dumper dumper, pageSize int) error {
	listOptions := shopify.ListOptions{
		Fields: shopifyFields,
		Limit: pageSize,
	}

	for {
		products, pages, err := client.Product.ListWithPagination(
			shopify.ProductListOptions{
				ListOptions: listOptions,
			},
		)

		if err != nil {
			return fmt.Errorf("Failed retrieve products: "+err.Error())
		}

		for _, product := range products {
			err = dumper.Dump(product)
			if err != nil {
				return fmt.Errorf("Failed saving products: "+err.Error())
			}
		}

		if pages.NextPageOptions == nil {
			break
		}

		listOptions.PageInfo = pages.NextPageOptions.PageInfo
	}

	return nil
}

func main() {
	var key, password, token string
	var asJSON, showHelp, showVersion, verbose bool
	var jsonRoot string
	var timeout int64
	var pageSize int
	var options []shopify.Option
	var client *shopify.Client

	flag.Usage = func() {
		exitFailure(fmt.Sprintf(usage, strings.Join(exportformat.JSONRootProperties, ", ")), 2)
	}

	flag.BoolVar(&showHelp, "h", false, "")
	flag.BoolVar(&showHelp, "help", false, "")
	flag.BoolVar(&asJSON, "j", false, "")
	flag.BoolVar(&asJSON, "json", false, "")
	flag.StringVar(&key, "k", os.Getenv("SHOPIFY_API_KEY"), "")
	flag.StringVar(&key, "key", os.Getenv("SHOPIFY_API_KEY"), "")
	flag.StringVar(&password, "p", os.Getenv("SHOPIFY_API_PASSWORD"), "")
	flag.StringVar(&password, "password", os.Getenv("SHOPIFY_API_PASSWORD"), "")
	flag.StringVar(&jsonRoot, "r", "", "")
	flag.StringVar(&jsonRoot, "json-root", "", "")
	flag.IntVar(&pageSize, "s", 250, "")
	flag.IntVar(&pageSize, "size", 250, "")
	flag.StringVar(&token, "t", os.Getenv("SHOPIFY_API_TOKEN"), "")
	flag.StringVar(&token, "token", os.Getenv("SHOPIFY_API_TOKEN"), "")
	flag.Int64Var(&timeout, "timeout", -1, "")
	flag.BoolVar(&showVersion, "v", false, "")
	flag.BoolVar(&showVersion, "version", false, "")
	flag.BoolVar(&verbose, "verbose", false, "")

	flag.Parse()
	argv := flag.Args()

	if showVersion {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if showHelp || len(argv) == 0 {
		flag.Usage()
	}

	var dumper dumper
	var err error

	if asJSON {
		dumper, err = exportformat.NewJSON(argv[0], jsonRoot)
	} else {
		dumper, err = exportformat.NewCSV(argv[0])
	}

	if err != nil {
		exitFailure(err.Error(), 1)
	}

	app := shopify.App{ApiKey: key, Password: password}

	if timeout > -1 {
		options = append(
			options,
			shopify.WithHTTPClient(
				&http.Client{
					Timeout: time.Second * time.Duration(timeout),
				},
			),
		)
	}

	if verbose {
		options = append(options, shopify.WithLogger(&shopify.LeveledLogger{Level: shopify.LevelDebug}))
	}

	client = shopify.NewClient(app, argv[0], token, options...)

	err = dumpProducts(client, dumper, pageSize)
	dumpErr := dumper.Close()

	if err != nil {
		exitFailure(err.Error(), 1)
	}

	if dumpErr != nil {
		exitFailure(dumpErr.Error(), 1)
	}
}
