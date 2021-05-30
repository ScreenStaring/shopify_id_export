package exportformat

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	shopify "github.com/bold-commerce/go-shopify"
)

type JSON struct {
	out      *os.File
	root     string
	products []shopify.Product
}

var JSONRootProperties = []string{
	"barcode",
	"product_id",
	"product_title",
	"handle",
	"variant_id",
	"sku",
}

func isValidJSONRootProperty(name string) bool {
	for _, property := range JSONRootProperties {
		if name == property {
			return true
		}
	}

	return false
}

func productMap(product shopify.Product) map[string]interface{} {
	record := make(map[string]interface{})

	record["product_id"] = strconv.FormatInt(product.ID, 10)
	record["handle"] = product.Handle
	record["product_title"] = product.Title
	record["product_type"] = product.ProductType

	return record
}

func variantMap(variant shopify.Variant) map[string]string {
	record := make(map[string]string)

	record["barcode"] = variant.Barcode
	record["variant_id"] = strconv.FormatInt(variant.ID, 10)
	record["variant_title"] = variant.Title
	record["sku"] = variant.Sku

	return record
}

func (j *JSON) formatForOutput() interface{} {
	if j.root == "variant_id" || j.root == "sku" || j.root == "barcode" {
		return j.formatWithVariantRoot()
	} else {
		return j.formatWithProduct()
	}
}

func (j *JSON) formatWithVariantRoot() map[string]interface{} {
	output := make(map[string]interface{})

	for _, product := range j.products {
		for _, variant := range product.Variants {
			record := variantMap(variant)
			key := record[j.root]
			if len(key) == 0 {
				continue
			}

			for key, value := range productMap(product) {
				svalue, ok := value.(string)
				if !ok {
					panic(fmt.Sprintf("Cannot convert product property '%s' to string for product '%s'", key, product.Title))
				}

				record[key] = svalue
			}

			output[key] = record
		}

	}

	return output
}

func (j *JSON) formatWithProduct() interface{} {
	if len(j.root) > 0 {
		return j.formatWithProductRoot()
	}

	var output []map[string]interface{}

	for _, product := range j.products {
		record := productMap(product)

		var variants []map[string]string
		for _, variant := range product.Variants {
			variants = append(variants, variantMap(variant))
		}

		record["variants"] = variants
		output = append(output, record)
	}

	return output
}

func (j *JSON) formatWithProductRoot() map[string]interface{} {
	output := make(map[string]interface{})

	for _, product := range j.products {
		record := productMap(product)

		var variants []map[string]string
		for _, variant := range product.Variants {
			variants = append(variants, variantMap(variant))
		}

		record["variants"] = variants

		key, ok := record[j.root].(string)
		if !ok {
			panic(fmt.Sprintf("Cannot convert JSON root property '%s' to string for product '%s'", j.root, product.Title))
		}

		output[key] = record
	}

	return output
}

func NewJSON(shop string, jsonRoot string) (*JSON, error) {
	j := new(JSON)

	if len(jsonRoot) > 0 && !isValidJSONRootProperty(jsonRoot) {
		return nil, fmt.Errorf("Invalid JSON root property: %s", jsonRoot)
	}

	out, err := os.Create(shop+".json")
	if err != nil {
		return nil, fmt.Errorf("Failed to create JSON file: %s", err)
	}

	j.root = jsonRoot
	j.out = out

	return j, nil
}

func (j *JSON) Dump(product shopify.Product) error {
	j.products = append(j.products, product)

	return nil
}

func (j *JSON) Close() error {
	var out []byte
	var n int
	var err error

	defer j.out.Close()

	out, err = json.Marshal(j.formatForOutput())
	if err != nil {
		return err
	}

	n, err = j.out.Write(out)
	if err != nil {
		return err
	}

	if n != len(out) {
		return fmt.Errorf("Was only able to write %d/%d bytes to JSON file", n, len(out))
	}

	return nil
}
