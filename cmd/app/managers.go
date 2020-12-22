package app

import (
	"encoding/json"
	"github.com/bahrom656/crud/pkg/customers"
	"github.com/bahrom656/crud/pkg/managers"
	"github.com/bahrom656/crud/pkg/security"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

//handleManagerRegistration
func (s *Server) handleManagerRegistration(writer http.ResponseWriter, request *http.Request) {
	var regM *managers.Registration
	MAuthentication(request.Context(), writer)

	err := json.NewDecoder(request.Body).Decode(&regM)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	item, err := s.managersSvs.Register(request.Context(), regM)
	if err != nil {
		errorWriter(writer, http.StatusInternalServerError, err)
		return
	}

	data, err := json.Marshal(item)
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

//handleManagerGetToken
func (s *Server) handleManagerGetToken(writer http.ResponseWriter, request *http.Request) {
	var item *managers.Auth

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

//handleManagerGetProducts
func (s *Server) handleManagerGetProducts(writer http.ResponseWriter, request *http.Request) {
	items, err := s.managersSvs.GetProducts(request.Context())
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

//handleManagerRemoveProductsByID
func (s *Server) handleManagerRemoveProductsByID(writer http.ResponseWriter, request *http.Request) {
	MAuthentication(request.Context(), writer)
	idP, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	productID, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	err = s.managersSvs.RemoveProductsByID(request.Context(), productID)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

}

//handleManagerRemoveCustomerByID
func (s *Server) handleManagerRemoveCustomerByID(writer http.ResponseWriter, request *http.Request) {
	idP, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	customerID, err := strconv.ParseInt(idP, 10, 64)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	err = s.managersSvs.RemoveCustomerByID(request.Context(), customerID)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}
}

//handleManagerGetCustomer
func (s *Server) handleManagerGetCustomer(writer http.ResponseWriter, request *http.Request) {
	items, err := s.managersSvs.GetCustomer(request.Context())
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

//handleManagerChangeCustomer
func (s *Server) handleManagerChangeCustomer(writer http.ResponseWriter, request *http.Request) {
	item := &customers.Customer{}

	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	customer, err := s.managersSvs.ManagerChangeCustomer(request.Context(), item)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	data, err := json.Marshal(customer)
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

//handleManagerChangeProduct
func (s *Server) handleManagerChangeProduct(writer http.ResponseWriter, request *http.Request) {
	item := &managers.Product{}

	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	product, err := s.managersSvs.ManagerChangeProduct(request.Context(), item)
	if err != nil {
		errorWriter(writer, http.StatusBadRequest, err)
		return
	}

	data, err := json.Marshal(product)
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

func (s *Server) handleManagerGetSales(http.ResponseWriter, *http.Request) {

}

func (s *Server) handleManagerMakeSales(http.ResponseWriter, *http.Request) {

}
