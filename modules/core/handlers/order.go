// Package handlers contains HTTP handlers for the core module of nutrix.
//
// The handlers in this package are used to handle incoming HTTP requests for
// the core module of nutrix. They interact with the services package, which
// contains the business logic of the core module.
//
// The handlers in this package create a RESTful API for the core module of
// nutrix. The API endpoints are documented using the Swagger specification.
// Each handler function is responsible for processing HTTP requests, calling
// the appropriate service methods, and returning HTTP responses.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/modules/core/dto"
	"github.com/elmawardy/nutrix/modules/core/models"
	"github.com/elmawardy/nutrix/modules/core/services"
)

// PayUnpaidOrder returns a HTTP handler function to pay an unpaid order.
func PayUnpaidOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		orderId := r.URL.Query().Get("id")
		if orderId == "" {
			http.Error(w, "id query string is required", http.StatusBadRequest)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err := orderService.PayUnpaidOrder(orderId)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

// GetUnpaidOrders returns a HTTP handler function to get all unpaid orders.
func GetUnpaidOrders(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		unpaidOrders, err := orderService.GetUnpaidOrders()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Failed to get unpaid orders", http.StatusInternalServerError)
			return
		}

		response := struct {
			Orders []models.Order `json:"orders"`
		}{
			Orders: unpaidOrders,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Failed to marshal unpaid orders response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	}
}

// CancelOrder returns a HTTP handler function to cancel an order.
func CancelOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id quert string is required", http.StatusBadRequest)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err := orderService.CancelOrder(id)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// GetStashedOrders returns a HTTP handler function to get all orders in stash.
func GetStashedOrders(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		stashed_orders, err := orderService.GetStashedOrders()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(stashed_orders)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

// OrderRemoveFromStash returns a HTTP handler function to remove an order from the stash.
func OrderRemoveFromStash(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var order_remove_stash_request dto.OrderRemoveStashRequest
		err := decoder.Decode(&order_remove_stash_request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err = orderService.RemoveStashedOrder(order_remove_stash_request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// OrderStash returns a HTTP handler function to stash an order.
func OrderStash(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var order_stash_request dto.OrderStashRequest
		err := decoder.Decode(&order_stash_request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		order, err := orderService.StashOrder(order_stash_request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := dto.OrderStashResponse{
			Order: order,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

// FinishOrder returns a HTTP handler function to finish an order.
func FinishOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var finish_order_request dto.FinishOrderRequest
		err := decoder.Decode(&finish_order_request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err = orderService.FinishOrder(finish_order_request)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		display_id, err := orderService.GetOrderDisplayId()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msg := models.WebsocketOrderFinishServerMessage{
			OrderId: display_id,
			WebsocketTopicServerMessage: models.WebsocketTopicServerMessage{
				Type:      "topic_message",
				TopicName: "order_finished",
				Severity:  "info",
			},
		}

		msgJson, err := json.Marshal(msg)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		notifications_svc, err := services.SpawnNotificationSingletonSvc("melody", logger, config)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		notifications_svc.SendToTopic("order_finished", string(msgJson))
	}
}

// SubmitOrder returns a HTTP handler function to submit an order.
func SubmitOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var order models.Order
		err := decoder.Decode(&order)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		err = orderService.SubmitOrder(order)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

}

// GetOrders returns a HTTP handler function to retrieve a list of orders.
// to use pagination, send a "first" and "rows" query string
// to select all rows, send a "rows" query string with value -1
// to filter for orders that contains a specific display_id, just send a display_id query string
func GetOrders(config config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := services.GetOrdersParameters{}

		orderDisplayId := r.URL.Query().Get("display_id")
		if orderDisplayId != "" {
			params.OrderDisplayIdContains = orderDisplayId
		}

		first := r.URL.Query().Get("first")
		if i, err := strconv.Atoi(first); err == nil {
			params.First = i
		}

		rows := r.URL.Query().Get("rows")
		if i, err := strconv.Atoi(rows); err == nil {
			if i != -1 {
				params.Rows = i
			}
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		orders, total_records, err := orderService.GetOrders(params)

		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := struct {
			Orders       []models.Order `json:"orders"`
			TotalRecords int64          `json:"total_records"`
		}{
			Orders:       orders,
			TotalRecords: total_records,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}

// StartOrder returns a HTTP handler function to start an order.
func StartOrder(config config.Config, logger logger.ILogger, settings models.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var order_start_request dto.OrderStartRequest
		err := decoder.Decode(&order_start_request)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		orderService := services.OrderService{
			Logger:   logger,
			Config:   config,
			Settings: settings,
		}

		err = orderService.StartOrder(order_start_request)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)

			response := struct {
				Data string `json:"body"`
			}{
				Data: err.Error(),
			}

			json_response, err := json.Marshal(response)
			if err != nil {
				logger.Error(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write(json_response)
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

// GetOrder returns a HTTP handler function to retrieve an order.
func GetOrder(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id quert string is required", http.StatusBadRequest)
			return
		}

		orderService := services.OrderService{
			Logger: logger,
			Config: config,
		}

		order, err := orderService.GetOrder(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonOrder, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonOrder)
	}
}
