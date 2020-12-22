package app

import (
	"encoding/json"
	"github.com/bahrom656/crud/cmd/app/middleware"
	"github.com/bahrom656/crud/pkg/customers"
	"github.com/bahrom656/crud/pkg/security"
	"log"
	"net/http"
)

//handleCustomerRegistration
func (s *Server) handleCustomerRegistration(writer http.ResponseWriter, request *http.Request) {
	var item *customers.Registration
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	saved, err := s.customersSvs.Register(request.Context(), item)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	data, err := json.Marshal(saved)
	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

//handleCustomerGetToken
func (s *Server) handleCustomerGetToken(writer http.ResponseWriter, request *http.Request) {
	var item *customers.Auth
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	token, err := s.customersSvs.Token(request.Context(), item.Login, item.Password)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	data, err := json.Marshal(&security.Token{Token: token})
	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

//handleCustomerGetProducts
func (s *Server) handleCustomerGetProducts(writer http.ResponseWriter, request *http.Request) {
	items, err := s.customersSvs.GetProducts(request.Context())
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

//handleCustomerGetPurchases
func (s *Server) handleCustomerGetPurchases(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}
	if id == 0 {
		errorWriter(writer, http.StatusForbidden, err)
		return
	}

	items, err := s.customersSvs.Purchases(request.Context(), id)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

//handleCustomerMakePurchases
func (s *Server) handleCustomerMakePurchases(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}
	if id == 0 {
		errorWriter(writer, http.StatusForbidden, err)
		return
	}
	items, err := s.customersSvs.Purchases(request.Context(), id)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	data, err := json.Marshal(items)
	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}
