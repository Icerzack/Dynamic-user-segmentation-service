package app

import (
	"avito-backend-internship/internal/pkg/db"
	"avito-backend-internship/internal/pkg/model"
	db_service "avito-backend-internship/internal/pkg/service/db"
	history_service "avito-backend-internship/internal/pkg/service/history"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	cfg            *Config
	ctx            context.Context
	db             db.DBops
	dbService      db_service.Service
	historyService history_service.Service
}

func NewServer(ctx context.Context, cfg *Config, db db.DBops, dbService db_service.Service, historyService history_service.Service) *Server {
	return &Server{
		cfg:            cfg,
		ctx:            ctx,
		db:             db,
		dbService:      dbService,
		historyService: historyService,
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()

	s.initHandlers(mux)

	fmt.Println("Server started on port:", s.cfg.Port)

	err := http.ListenAndServe(":"+s.cfg.Port, mux)

	if err != nil {
		fmt.Println("Failed to start server on port:", s.cfg.Port, "ERROR:", err.Error())
	}
}

func (s *Server) initHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", s.mainHandler)
	mux.HandleFunc("/create-segment", s.createSegmentHandler)
	mux.HandleFunc("/delete-segment", s.deleteSegmentHandler)
	mux.HandleFunc("/user-in-segment", s.userInSegmentHandler)
	mux.HandleFunc("/get-user-segments", s.getUserSegmentsHandler)
	mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs/"))))
}

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
		var segment model.SegmentRequest
		err := json.NewDecoder(request.Body).Decode(&segment)
		if err != nil || segment.Title == nil {
			if err != nil {
				fmt.Println("Error decoding JSON:", err.Error())
			}
			var output model.SegmentResponse
			output.Status = "Invalid JSON"

			writer.WriteHeader(http.StatusBadRequest)
			writer.Header().Set("Content-Type", "application/json")
			result, err := json.Marshal(output)
			if err != nil {
				fmt.Println("Error marshalling output JSON:", err.Error())
				return
			}

			_, err = writer.Write(result)
			if err != nil {
				fmt.Println("Error while writing response:", err.Error())
			}
			return
		}

		err = s.dbService.InsertSegmentIntoDatabase(s.ctx, s.db, segment)
		if err != nil {
			fmt.Println("Error adding segment to db_service:", err.Error())
			return
		}

		var output model.SegmentResponse
		output.Status = "Success"

		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		result, err := json.Marshal(output)
		if err != nil {
			fmt.Println("Error marshalling output JSON:", err.Error())
			return
		}

		_, err = writer.Write(result)
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
		var segment model.SegmentRequest
		err := json.NewDecoder(request.Body).Decode(&segment)
		if err != nil || segment.Title == nil {
			if err != nil {
				fmt.Println("Error decoding JSON:", err.Error())
			}
			var output model.SegmentResponse
			output.Status = "Invalid JSON"

			writer.WriteHeader(http.StatusBadRequest)
			writer.Header().Set("Content-Type", "application/json")
			result, err := json.Marshal(output)
			if err != nil {
				fmt.Println("Error marshalling output JSON:", err.Error())
				return
			}

			_, err = writer.Write(result)
			if err != nil {
				fmt.Println("Error while writing response:", err.Error())
			}
			return
		}

		ok, err := s.dbService.DeleteSegmentFromDatabase(s.ctx, s.db, segment)
		if err != nil {
			fmt.Println("Error deleting segment from db_service:", err.Error())
			return
		}

		var output model.SegmentResponse
		if ok {
			output.Status = "Success"
		} else {
			output.Status = "Failed: no such element to delete"
		}

		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		result, err := json.Marshal(output)
		if err != nil {
			fmt.Println("Error marshalling output JSON:", err.Error())
			return
		}

		_, err = writer.Write(result)
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
	if request.Method == http.MethodPut {
		var usersSegment model.UserSegmentRequest
		err := json.NewDecoder(request.Body).Decode(&usersSegment)
		if err != nil || usersSegment.UserID == nil || (usersSegment.UserID != nil && usersSegment.SegmentsTitlesToAdd == nil && usersSegment.SegmentsTitlesToDelete == nil) {
			if err != nil {
				fmt.Println("Error decoding JSON:", err.Error())
			}
			var output model.SegmentResponse
			output.Status = "Invalid JSON"

			writer.WriteHeader(http.StatusBadRequest)
			writer.Header().Set("Content-Type", "application/json")
			result, err := json.Marshal(output)
			if err != nil {
				fmt.Println("Error marshalling output JSON:", err.Error())
				return
			}

			_, err = writer.Write(result)
			if err != nil {
				fmt.Println("Error while writing response:", err.Error())
			}
			return
		}

		addNotExist, deleteNotExist, err := s.dbService.ModifyUsersSegmentsInDatabase(s.ctx, s.db, usersSegment, s.historyService)
		if err != nil {
			fmt.Println("Error modifying user-segment in db_service:", err.Error())
			return
		}

		var output model.UserSegmentResponse
		output.Status = "Success"
		output.SegmentsTitlesNotExistAdd = addNotExist
		output.SegmentsTitlesNotExistDelete = deleteNotExist
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		result, err := json.Marshal(output)
		if err != nil {
			fmt.Println("Error marshalling output JSON:", err.Error())
			return
		}

		_, err = writer.Write(result)
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
		var usersSegment model.UserSegmentRequest
		err := json.NewDecoder(request.Body).Decode(&usersSegment)
		if err != nil || usersSegment.UserID == nil {
			if err != nil {
				fmt.Println("Error decoding JSON:", err.Error())
			}
			var output model.SegmentResponse
			output.Status = "Invalid JSON"

			writer.WriteHeader(http.StatusBadRequest)
			writer.Header().Set("Content-Type", "application/json")
			result, err := json.Marshal(output)
			if err != nil {
				fmt.Println("Error marshalling output JSON:", err.Error())
				return
			}

			_, err = writer.Write(result)
			if err != nil {
				fmt.Println("Error while writing response:", err.Error())
			}
			return
		}

		segments, err := s.dbService.GetUserSegmentsFromDatabase(s.ctx, s.db, usersSegment)
		if err != nil {
			fmt.Println("Error reading users segments from db_service:", err.Error())
			return
		}

		var output model.UserSegmentResponse
		if len(segments) == 0 {
			output.Status = "Failure: this user is not in any segment"
		} else {
			output.Status = "Success"
			output.Segments = segments
		}
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		result, err := json.Marshal(output)
		if err != nil {
			fmt.Println("Error marshalling output JSON:", err.Error())
			return
		}

		_, err = writer.Write(result)
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
