package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/account", NewHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", NewHandlerFunc(s.handleAccount))
	log.Print("Server is Running....")
	http.ListenAndServe(s.listenAddr, router)
}

// handle is prefix for handler
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		params := mux.Vars(r)
		if len(params) == 0 {
			return s.handleGetAccounts(w, r)
		} else {
			return s.handleGetAccount(w, r)
		}
	} else if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	} else if r.Method == "PUT" {
		return s.handleUpdateAccount(w, r)
	} else if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id, err := ConvertID(params["id"])
	if err != nil {
		return err
	}
	acc, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, acc)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accountReq := &CreateAccountReq{}
	if err := json.NewDecoder(r.Body).Decode(accountReq); err != nil {
		return err
	}
	account := NewAccount(accountReq.FirstName, accountReq.LastName)
	s.store.CreateAccount(account)
	return WriteJson(w, http.StatusAccepted, account)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	accountUpdateDetails := CreateAccountReq{}
	id, err := ConvertID(params["id"])
	if err != nil {
		return err
	}
	json.NewDecoder(r.Body).Decode(&accountUpdateDetails)
	account, err := s.store.UpdateAccount(id, &accountUpdateDetails)
	if err != nil {
		return err
	}
	WriteJson(w, http.StatusOK, account)
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id, err := ConvertID(params["id"])
	if err != nil {
		return err
	}
	s.store.DeleteAccount(id)
	return WriteJson(w, http.StatusOK, map[string]int{"deleted": id})

}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, accounts)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

func NewHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, 400, APIError{Error: err.Error()})
		}
	}
}

func ConvertID(idStr string) (int, error) {
	accid, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id %s", idStr)
	}
	return accid, nil
}
