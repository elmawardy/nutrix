package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
)

type JSONAPILinks struct {
	Self string `json:"self"`
}

type JSONAPIMeta struct {
	TotalRecords int `json:"total_records"`
	PageNumber   int `json:"page_number"`
	PageSize     int `json:"page_size"`
	PageCount    int `json:"page_count"`
}

type JSONApiOkResponse struct {
	Links JSONAPILinks `json:"links"`
	Data  interface{}  `json:"data"`
	Meta  JSONAPIMeta  `json:"meta"`
}

type JSONAPIResource struct {
	ID            string       `json:"id,omitempty"`
	Type          string       `json:"type"`
	Links         JSONAPILinks `json:"links"`
	Attributes    interface{}  `json:"attributes"`
	Relationships interface{}  `json:"relationships"`
}

// InsertCategory returns a HTTP handler function to insert a Category into the database.
func InsertCategory(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		request := struct {
			Category models.Category `json:"category"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		categoryService := services.CategoryService{
			Logger: logger,
			Config: config,
		}

		err = categoryService.InsertCategory(request.Category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// DeleteCategory returns a HTTP handler function to delete a Category from the database.
func DeleteCategory(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id query string is required", http.StatusBadRequest)
			return
		}

		categoryService := services.CategoryService{
			Logger: logger,
			Config: config,
		}

		err := categoryService.DeleteCategory(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// UpdateCategory returns a HTTP handler function to update a Category in the database.
func UpdateCategory(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body := struct {
			Category models.Category `json:"category"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		categoryService := services.CategoryService{
			Logger: logger,
			Config: config,
		}

		err = categoryService.UpdateCategory(body.Category)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// GetCategories returns a HTTP handler function to retrieve a list of Categories from the database.
func GetCategories(config config.Config, logger logger.ILogger, url_prefix string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		page_number, err := strconv.Atoi(r.URL.Query().Get("page[number]"))
		if err != nil {
			page_number = 1
		}

		page_size, err := strconv.Atoi(r.URL.Query().Get("page[size]"))
		if err != nil {
			page_size = 50
		}

		categoryService := services.CategoryService{
			Logger: logger,
			Config: config,
		}

		categories, err := categoryService.GetCategories(page_number, page_size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resources := make([]JSONAPIResource, len(categories))

		for i, category := range categories {

			products := make([]JSONAPIResource, len(category.Products))

			for i, product := range category.Products {
				products[i] = ConvertProductToJSONAPIResource(models.Product{
					Id:   product.Id,
					Name: product.Name,
				}, url_prefix)
			}

			resources[i] = JSONAPIResource{
				ID:    category.Id,
				Type:  "categories",
				Links: JSONAPILinks{Self: url_prefix + "/api/categories/" + category.Id},
				Attributes: struct {
					Name string `json:"name"`
				}{
					Name: category.Name,
				},
				Relationships: struct {
					Products struct {
						Data []JSONAPIResource `json:"data"`
						Meta JSONAPIMeta       `json:"meta"`
					} `json:"products"`
				}{
					Products: struct {
						Data []JSONAPIResource `json:"data"`
						Meta JSONAPIMeta       `json:"meta"`
					}{
						Data: products,
						Meta: JSONAPIMeta{
							TotalRecords: len(products),
						},
					},
				},
			}
		}

		response := JSONApiOkResponse{
			Data: resources,
			Meta: JSONAPIMeta{
				TotalRecords: len(categories),
			},
			Links: JSONAPILinks{
				Self: url_prefix + "/api/categories",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

func ConvertProductToJSONAPIResource(product models.Product, url_prefix string) JSONAPIResource {

	sub_products := make([]JSONAPIResource, len(product.SubProducts))

	for index, sub_product := range product.SubProducts {
		sub_products[index] = ConvertProductToJSONAPIResource(sub_product, url_prefix)
	}

	resource := JSONAPIResource{
		ID:    product.Id,
		Type:  "products",
		Links: JSONAPILinks{Self: url_prefix + "/api/products/" + product.Id},
		Attributes: struct {
			Name string `json:"name"`
		}{
			Name: product.Name,
		},
		Relationships: struct {
			SubProducts []JSONAPIResource `json:"subproducts"`
		}{
			SubProducts: sub_products,
		},
	}

	return resource
}
