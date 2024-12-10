package exportformat

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/screenstaring/shopify_id_export/gql"
)

type CSV struct {
	out           *csv.Writer
	file          *os.File
	header        []string
	headerWritten bool
}

func NewCSV(shop string) (*CSV, error) {
	c := new(CSV)

	file, err := os.Create(shop + ".csv")
	if err != nil {
		return nil, fmt.Errorf("Failed to create CSV file: %s", err)
	}

	c.file = file
	c.out = csv.NewWriter(c.file)
	c.header = []string{
		"Product ID",
		"Product Title",
		"Product Type",
		"Variant ID",
		"Variant Title",
		"SKU",
		"Barcode",
		"Handle",
	}

	return c, nil
}

func (c *CSV) Dump(product gql.Product) error {
	var err error

	if !c.headerWritten {
		c.out.Write(c.header)
		c.headerWritten = true
	}

	for _, edge := range product.Variants.Edges {
		variant := edge.Node
		row := []string{
			product.ID,
			product.Title,
			product.ProductType,
			variant.ID,
			variant.Title,
			variant.Sku,
			variant.Barcode,
			product.Handle,
		}

		err = c.out.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CSV) Close() error {
	defer c.file.Close()

	c.out.Flush()

	err := c.out.Error()
	if err != nil {
		return err
	}

	return nil
}
