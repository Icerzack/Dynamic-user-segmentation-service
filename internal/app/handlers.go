package app

import (
	"avito-backend-internship/internal/pkg/model"
	"encoding/json"
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
				return
			}
			return
		}

		err = s.service.InsertSegmentIntoDatabase(s.ctx, s.db, segment)
		if err != nil {
			fmt.Println("Error adding segment to database:", err.Error())
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
				return
			}
			return
		}

		ok, err := s.service.DeleteSegmentFromDatabase(s.ctx, s.db, segment)
		if err != nil {
			fmt.Println("Error deleting segment from database:", err.Error())
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
				return
			}
			return
		}

		addNotExist, deleteNotExist, err := s.service.ModifyUsersSegmentsInDatabase(s.ctx, s.db, usersSegment)
		if err != nil {
			fmt.Println("Error modifying user-segment in database:", err.Error())
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
				return
			}
			return
		}

		segments, err := s.service.GetUserSegmentsFromDatabase(s.ctx, s.db, usersSegment)
		if err != nil {
			fmt.Println("Error reading users segments from database:", err.Error())
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
