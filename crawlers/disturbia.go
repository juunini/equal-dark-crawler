package crawlers

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Disturbia crawler
type Disturbia struct{}

// IsValidProductsPage checks page
func (disturbia *Disturbia) IsValidProductsPage(doc *goquery.Document) bool {
	productsClass := doc.Find("body").First().HasClass("products")
	productsLength := doc.Find("ul.products > li").Length() != 0
	categorySelector := doc.Find(".category").Length() != 0

	return productsClass && productsLength && categorySelector
}

// GetProductsURL returns product url
func (disturbia *Disturbia) GetProductsURL(doc *goquery.Document) (productsURL []string) {
	products := doc.Find("ul.products > li")
	productsURL = products.Map(func(i int, product *goquery.Selection) string {
		href, _ := product.Find("a").First().Attr("href")
		return "https://www.disturbia.co.uk" + href
	})
	return
}

// GetProductName returns product name
func (disturbia *Disturbia) GetProductName(doc *goquery.Document) (name string) {
	nameSelector := doc.Find("h1")
	name = nameSelector.First().Text()
	return
}

// GetProductCurrency returns currency
func (disturbia *Disturbia) GetProductCurrency(doc *goquery.Document) (currency string) {
	priceSelector := doc.Find(".product .detail .price")
	priceStr := priceSelector.First().Text()
	currency = GetCurrencyFromText(priceStr)
	return
}

func (disturbia *Disturbia) isSaleProduct(priceSelector *goquery.Selection) bool {
	return priceSelector.HasClass("reduced")
}

func (disturbia *Disturbia) getProductPrice(doc *goquery.Document, sale bool) (price float64) {
	priceSelector := doc.Find(".product .detail .price")
	var priceStr string

	if disturbia.isSaleProduct(priceSelector) {
		priceSelectorText := priceSelector.First().Text()
		index := GetIntFromBool(sale)

		priceStr = strings.Split(priceSelectorText, "Now")[index]
	} else {
		priceStr = priceSelector.First().Text()
	}
	price = GetFloatFromText(priceStr)
	return
}

// GetProductPrice returns float price
func (disturbia *Disturbia) GetProductPrice(doc *goquery.Document) (price float64) {
	price = disturbia.getProductPrice(doc, false)
	return
}

// GetProductSalePrice returns float price
func (disturbia *Disturbia) GetProductSalePrice(doc *goquery.Document) (salePrice float64) {
	salePrice = disturbia.getProductPrice(doc, true)
	return
}
