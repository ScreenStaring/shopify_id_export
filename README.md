# Shopify ID Export

Dump Shopify product and variant IDs —along with other identifiers— to a CSV or JSON file.

When data munging in Shopify you'll often need to (cross-)reference product IDs, variant IDs, and SKUs.
Shopify does not make product nor variant ID available via their product export.

Other times one is provided a product spreadsheet with SKUs but the target system relies
on product/variant IDs.

This program can help.


## Installation

Download the version for your platform on the [releases page](https://github.com/screenstaring/shopify_id_export/releases).
Windows, macOS/OS X, and GNU/Linux are supported.

## Usage

### Credentials

To use the command you'll need access to the Shopify store you want to export from.

If you have access to the store via Shopify Admin, this can be done by [generating private app API credentials](https://shopify.dev/tutorials/generate-api-credentials).
Once obtained they can be specified to the `shopify_id_export` command:

```
shopify_id_export -k my-app-key -p my-app-password shop
```

If the store has your app installed, you can use the credentials generated when the shop installed your app:

```
shopify_id_export -t shop-token shop
```

In both cases credentials can be set via environment variables. See [Options](#options).

### Options

```
shopify_id_export [hjv] [-k key] [-p password] [-r root-property] [-t token] shop

Options
-h --help                display this help message
-j --json                output dump as JSON
-k --key       key       Shopify API key; defaults to the SHOPIFY_API_KEY environment variable
-p --password  password  Shopify API password; defaults to the SHOPIFY_API_PASSWORD environment variable
-r --json-root property  use property as the top-level property for each JSON object
-t --token     token     Shopify API token; defaults to the SHOPIFY_API_TOKEN environment variable
-v --version             display version information

By default data is output to a CSV file.

Valid properties for the --json-root option are: product_id, product_title, handle, variant_id, sku
```

### Output to CSV

This is the default output format. A file named `shop.csv` will be created in the current directory. `shop` will replaced by the shop's name.

### Output to JSON

To output a JSON file use the `-j` option:
```
shopify_id_export -j -t shop-token shop
```

This will create a file named `shop.json` in the current directory. `shop` will be replaced by the shop's name.

If you're cross-referencing IDs it may be useful to set the root property for the JSON object output for each product/variant.

This will output each object with the variant's SKU as the root:
```
shopify_id_export -j -r sku -t shop-token shop
```

Valid properties for the `-r`/`--json-root` option are: `product_id`, `product_title`, `handle`, `variant_id`, `sku`.

## License

Released under the MIT License: http://www.opensource.org/licenses/MIT

---

Made by [ScreenStaring](http://screenstaring.com)
