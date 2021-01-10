package crawlers_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"equal_dark_crawler/crawlers"
)

var _ = Describe("Disturbia", func() {
	var disturbia crawlers.Disturbia

	var makeTestServer = func(body []byte) *httptest.Server {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("Content-Type", "text/html")
			rw.Write(body)
		}))
		return server
	}

	Describe("IsValidProductsPage", func() {
		Context("When is products page", func() {
			It("Should returns true", func() {
				server := makeTestServer([]byte(disturbiaProductsPageDocument))
				defer server.Close()

				doc, _ := goquery.NewDocument(server.URL)
				actual := disturbia.IsValidProductsPage(doc)

				Expect(actual).To(BeTrue())
			})
		})

		Context("When is not products page", func() {
			It("Should returns false", func() {
				server := makeTestServer([]byte(disturbiaSaleProductPageDocument))
				defer server.Close()

				doc, _ := goquery.NewDocument(server.URL)
				actual := disturbia.IsValidProductsPage(doc)

				Expect(actual).To(BeFalse())
			})
		})
	})

	Describe("GetProductsURL", func() {
		It("Should returns url array", func() {
			server := makeTestServer([]byte(disturbiaProductsPageDocument))
			defer server.Close()

			doc, _ := goquery.NewDocument(server.URL)
			actual := disturbia.GetProductsURL(doc)

			Expect(actual[0]).To(Equal("https://www.disturbia.co.uk/products/womens-all-tops/Blaze-Jumper"))
		})
	})

	Describe("GetProductName", func() {
		It("Should returns product name", func() {
			server := makeTestServer([]byte(disturbiaSaleProductPageDocument))
			defer server.Close()

			doc, _ := goquery.NewDocument(server.URL)
			actual := disturbia.GetProductName(doc)

			Expect(actual).To(Equal("Infernal Eternity Lace Up Vest"))
		})
	})

	Describe("GetProductCurrency", func() {
		Context("When is not sale product", func() {
			It("Should returns product currency", func() {
				server := makeTestServer([]byte(disturbiaNotSaleProductPageDocument))
				defer server.Close()

				doc, _ := goquery.NewDocument(server.URL)
				actual := disturbia.GetProductCurrency(doc)

				Expect(actual).To(Equal("GBP"))
			})
		})

		Context("When is sale product", func() {
			It("Should returns product currency", func() {
				server := makeTestServer([]byte(disturbiaSaleProductPageDocument))
				defer server.Close()

				doc, _ := goquery.NewDocument(server.URL)
				actual := disturbia.GetProductCurrency(doc)

				Expect(actual).To(Equal("GBP"))
			})
		})
	})

	Describe("GetProductPrice & GetProductSalePrice", func() {
		Context("When is not sale product", func() {
			It("Should returns equal price and sale price", func() {
				server := makeTestServer([]byte(disturbiaNotSaleProductPageDocument))
				defer server.Close()

				doc, _ := goquery.NewDocument(server.URL)
				price := disturbia.GetProductPrice(doc)
				salePrice := disturbia.GetProductSalePrice(doc)

				Expect(price).To(Equal(float64(46)))
				Expect(salePrice).To(Equal(float64(46)))
			})
		})

		Context("When is sale product", func() {
			It("Should returns different price and sale price", func() {
				server := makeTestServer([]byte(disturbiaSaleProductPageDocument))
				defer server.Close()

				doc, _ := goquery.NewDocument(server.URL)
				price := disturbia.GetProductPrice(doc)
				salePrice := disturbia.GetProductSalePrice(doc)

				Expect(price).To(Equal(float64(28)))
				Expect(salePrice).To(Equal(float64(16)))
			})
		})
	})
})
