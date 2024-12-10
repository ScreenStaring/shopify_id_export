package gql

import (
	"context"
	shopify "github.com/bold-commerce/go-shopify/v4"
)

const productsQuery = `
  query findProducts($after: String) {
    products(first:250 after:$after) {
      pageInfo {
	hasNextPage
	endCursor
      }
      edges {
	node {
          legacyResourceId
          handle
          productType
	  title
	  variants(first:100) {
	    edges {
	      node {
                barcode
		legacyResourceId
                sku
		title
	      }
	    }
	  }
	}
      }
    }
  }
`

type productResponse struct {
	Products Products `json:"products"`
}

type Products struct {
	Edges []ProductEdge `json:"edges"`
	PageInfo PageInfo `json:"pageInfo"`
}

type ProductEdge struct {
	Node Product `json:"node"`
}

type PageInfo struct {
	HasNextPage bool `json:"hasNextPage"`
	Cursor string `json:"endCursor"`
}

type Product struct {
	Handle   string   `json:"handle"`
	ID       string   `json:"legacyResourceId"`
	Title    string   `json:"title"`
	ProductType string  `json:"productType"`
	Variants Variants `json:"variants"`
}

type Variants struct {
	Edges []VariantEdge `json:"edges"`
}

type VariantEdge struct {
	Node Variant `json:"node"`
}

type Variant struct {
	Barcode string `json:"barcode"`
	ID    string `json:"legacyResourceId"`
	ProductID string `json:"productId"`
	Sku string `json:"sku"`
	Title string `json:"title"`
}

func FindProducts(client *shopify.Client, options map[string]interface{}) (Products, error) {
	response := productResponse{}

	err := client.GraphQL.Query(context.Background(), productsQuery, options, &response)

	return response.Products, err
}
