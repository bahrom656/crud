package app

import (
	"encoding/json"
	"github.com/bahrom656/crud/cmd/app/middleware"
	"github.com/bahrom656/crud/pkg/customers"
	"github.com/bahrom656/crud/pkg/managers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	mux          *mux.Router
	customersSvs *customers.Service
	managersSvs  *managers.Service
}

const (
	POST   = "POST"
	GET    = "GET"
	DELETE = "DELETE"
)

func NewServer(mux *mux.Router, customersSvs *customers.Service) *Server {
	return &Server{mux: mux, customersSvs: customersSvs}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	customersAuthenticateMd := middleware.Authenticate(s.customersSvs.IDByToken)
	managersAuthenticateMd := middleware.Authenticate(s.managersSvs.IDByToken)

	//customers
	customersSubrouter := s.mux.PathPrefix("/api/customers").Subrouter()
	customersSubrouter.Use(customersAuthenticateMd)
	customersSubrouter.HandleFunc("", s.handleCustomerRegistration).Methods(POST)
	customersSubrouter.HandleFunc("/token", s.handleCustomerGetToken).Methods(POST)
	customersSubrouter.HandleFunc("/products", s.handleCustomerGetProducts).Methods(GET)
	customersSubrouter.HandleFunc("/purchases", s.handleCustomerGetPurchases).Methods(GET)
	customersSubrouter.HandleFunc("/purchases", s.handleCustomerMakePurchases).Methods(POST)

	//managers
	managersSubroter := s.mux.PathPrefix("/api/managers").Subrouter()
	managersSubroter.Use(managersAuthenticateMd)
	managersSubroter.HandleFunc("", s.handleManagerRegistration).Methods(POST)
	managersSubroter.HandleFunc("/token", s.handleManagerGetToken).Methods(POST)
	managersSubroter.HandleFunc("/sales", s.handleManagerGetSales).Methods(GET)
	managersSubroter.HandleFunc("/sales", s.handleManagerMakeSales).Methods(POST)
	managersSubroter.HandleFunc("/products", s.handleManagerGetProducts).Methods(GET)
	managersSubroter.HandleFunc("/products", s.handleManagerChangeProducts).Methods(POST)
	managersSubroter.HandleFunc("/products/{id}", s.handleManagerRemoveProductByID).Methods(DELETE)
	managersSubroter.HandleFunc("/customer", s.handleManagerGetCustomers).Methods(GET)
	managersSubroter.HandleFunc("/customer", s.handleManagerChangeCustomer).Methods(POST)
	managersSubroter.HandleFunc("/customer/{id]", s.handleManagerRemoveCustomerByID).Methods(DELETE)
}

func respondJSON(w http.ResponseWriter, iData interface{}) {

	data, err := json.Marshal(iData)
	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}
func errorWriter(writer http.ResponseWriter, statusError int, err error) {
	log.Print(err)
	http.Error(writer, http.StatusText(statusError), statusError)
}
