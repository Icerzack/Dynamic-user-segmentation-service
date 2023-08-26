package server

import (
	"fmt"
	"net/http"
)

// mainHandler нужен, чтобы отлавливать несуществующие url и контролировать их поведение
func (s *Server) mainHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		if request.URL.Path != "/" {
			writer.WriteHeader(http.StatusNotFound)
			_, err := writer.Write([]byte("404 Not found"))
			if err != nil {
				fmt.Println("Error while writing response:", err.Error())
				return
			}
			return
		}
		writer.WriteHeader(http.StatusOK)
		_, err := writer.Write([]byte("it's good"))
		if err != nil {
			fmt.Println("Error while writing response:", err.Error())
			return
		}
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := writer.Write([]byte("405 Method not allowed"))
		if err != nil {
			fmt.Println("Error while writing response:", err.Error())
			return
		}
	}
}

func (s *Server) createSegmentHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		writer.WriteHeader(http.StatusOK)
		_, err := writer.Write([]byte("add"))
		if err != nil {
			fmt.Println("Error while writing response:", err.Error())
			return
		}
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := writer.Write([]byte("405 Method not allowed"))
		if err != nil {
			fmt.Println("Error while writing response:", err.Error())
			return
		}
	}
}

func (s *Server) deleteSegmentHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodDelete {
		writer.WriteHeader(http.StatusOK)
		_, err := writer.Write([]byte("delete"))
		if err != nil {
			fmt.Println("Error while writing response:", err.Error())
			return
		}
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := writer.Write([]byte("405 Method not allowed"))
		if err != nil {
			fmt.Println("Error while writing response:", err.Error())
			return
		}
	}
}

func (s *Server) userInSegmentHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		writer.WriteHeader(http.StatusOK)
		_, err := writer.Write([]byte("modify"))
		if err != nil {
			fmt.Println("Error while writing response:", err.Error())
			return
		}
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := writer.Write([]byte("405 Method not allowed"))
		if err != nil {
			fmt.Println("Error while writing response:", err.Error())
			return
		}
	}
}

func (s *Server) getUserSegmentsHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		writer.WriteHeader(http.StatusOK)
		_, err := writer.Write([]byte("read"))
		if err != nil {
			fmt.Println("Error while writing response:", err.Error())
			return
		}
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := writer.Write([]byte("405 Method not allowed"))
		if err != nil {
			fmt.Println("Error while writing response:", err.Error())
			return
		}
	}
}
