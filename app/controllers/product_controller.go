package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"AntiqueGo/app/models"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func formatRupiah(amount float64) string {
    intAmount := int(amount) // Konversi float64 ke int
    strAmount := fmt.Sprintf("%d", intAmount)
    n := len(strAmount)
    if n <= 3 {
        return "Rp " + strAmount
    }

    var result strings.Builder
    for i := 0; i < n; i++ {
        if (n-i)%3 == 0 && i != 0 {
            result.WriteString(".")
        }
        result.WriteByte(strAmount[i])
    }

    return "Rp " + result.String()
}


func (s *Server) Products(w http.ResponseWriter,r *http.Request) {
	render:= render.New(render.Options{
		Layout:"layout",
		Extensions: []string{".html", ".tmpl"},
	})

	q:=r.URL.Query()
	page,_:=strconv.Atoi(q.Get("page"))
	if page <= 0 {
		page=1
	}

	perPage := 9

	searchQuery := q.Get("search")
	

	if searchQuery != "" {
        productModel := models.Product{}
        products, totalRows, err := productModel.SearchProducts(s.DB, searchQuery, perPage, page)
		if err!= nil {
            return 
        }

		for i := range products {
			priceFloat, _ := products[i].Price.Float64() // Mengonversi decimal.Decimal ke float64
			products[i].FormattedPrice = formatRupiah(priceFloat)
		}
		

		pagination,_:=GetPaginationLinks(s.AppConfig, PaginationParams{
			Path:	"products",
			TotalRows: int64(totalRows),
			PerPage: int64(perPage),
			CurrentPage: int64(page),
		})
	
		cartID := GetShoppingCartID(w, r)
		cart, _ := GetShoppingCart(s.DB, cartID)
		itemCount := len(cart.CartItems)
	
		_ = render.HTML(w,http.StatusOK, "products",map[string]interface{}{
			"products": products,
			"pagination":pagination,
			"user": s.CurrentUser(w,r),
			"itemCount": itemCount,
		})
		return

    } else {
		productModel :=models.Product{}
		products,totalRows,err := productModel.GetProducts(s.DB,perPage,page)
		if err!= nil {
            return 
        }

		for i := range products {
			priceFloat, _ := products[i].Price.Float64() // Mengonversi decimal.Decimal ke float64
			products[i].FormattedPrice = formatRupiah(priceFloat)
		}
		
		
		pagination,_:=GetPaginationLinks(s.AppConfig, PaginationParams{
			Path:	"products",
			TotalRows: int64(totalRows),
			PerPage: int64(perPage),
			CurrentPage: int64(page),
		})
	
		cartID := GetShoppingCartID(w, r)
		cart, _ := GetShoppingCart(s.DB, cartID)
		itemCount := len(cart.CartItems)
	
		_ = render.HTML(w,http.StatusOK, "products",map[string]interface{}{
			"products": products,
			"pagination":pagination,
			"user": s.CurrentUser(w,r),
			"itemCount": itemCount,
		})
		return

	}


}

func (s *Server) GetProductBySlug(w http.ResponseWriter, r *http.Request){
	render:= render.New(render.Options{
        Layout:"layout",
		Extensions: []string{".html", ".tmpl"},
    })

	vars:= mux.Vars(r)

	if vars["slug"]==""{
		return 
	}

	productModel:= models.Product{}
	product, err := productModel.FindBySlug(s.DB,vars["slug"])
	if err!= nil {
        return 
    }

	cartID := GetShoppingCartID(w, r)
	cart, _ := GetShoppingCart(s.DB, cartID)
	itemCount := len(cart.CartItems)

	_=render.HTML(w,http.StatusOK,"product",map[string]interface{}{
		"product": product,
		"success": GetFlash(w,r,"success"),
		"error": GetFlash(w,r,"error"),
		"user": s.CurrentUser(w,r),
		"itemCount": itemCount,
	})
}