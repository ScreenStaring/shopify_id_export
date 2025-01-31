# Shopify ID Export

Dump Shopify product and variant IDs —along with other identifiers— to a CSV or JSON file.

When data munging in Shopify you'll often need to (cross-)reference product IDs, variant IDs, SKUs,
barcodes, etc... Shopify does not make product nor variant ID available via their product export.

Other times one is provided a product spreadsheet with SKUs but the target system relies
on product/variant IDs.

This program can help.


## Shopify REST API Removal

To account for Shopify's [Products and Variants API removal](https://shopify.dev/docs/apps/build/graphql/migrate/new-product-model)
v0.1.0 has been updated from the REST Admin API to the GraphQL Admin API. From testing (prior to Feb 1st, 2025, the REST API cutoff) this changed
resulted in a 4x increase in export time! 😱

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
-s --size      integer   page size to use when retrieving products (default: 250)
-t --token     token     Shopify API token; defaults to the SHOPIFY_API_TOKEN environment variable
--timeout      integer   set Shopify client timeout (default: 10 seconds)
-v --version             display version information
--verbose                output Shopify API request/response (default: false)

By default data is output to a CSV file.

Valid properties for the --json-root option are: barcode, product_id, product_title, handle, variant_id, sku
```

### Output to CSV

This is the default output format. A file named `shop.csv` will be created in the current directory. `shop` will replaced by the shop's name.

#### Combining CSVs

Once you have a CSV with the missing product info you can combine it with original CSV via the [`xsv`](https://github.com/BurntSushi/xsv) program.
(csvkit's [`csvjoin`](https://csvkit.readthedocs.io/en/latest/scripts/csvjoin.html) will work as well).

For example, if the original spreadsheet contains a column called `Product SKU` it can be combined with Shopify ID Export's spreadsheet via:
```
xsv join 'Product SKU' original-data.csv SKU shop.csv | xsv select '!SKU' > combined-data.csv
```

The 2nd command `xsv select '!SKU'` removes the `SKU` column as it's now redundant.

For more info and options see `xsv join --help`.

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

Valid properties for the `-r`/`--json-root` option are: `product_id`, `product_title`, `barcode`, `handle`, `variant_id`, `sku`.

## See Also

- [Shopify Development Tools](https://github.com/ScreenStaring/shopify-dev-tools) - Assists with the development and/or maintenance of Shopify apps and stores


## License

Released under the MIT License: http://www.opensource.org/licenses/MIT

---

Made by [ScreenStaring](http://screenstaring.com)
