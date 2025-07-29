package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MicaelaJofre/go_lib_response/response"
	"github.com/MicaelaJofre/gocourse_enrollment/internal/enrollment"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewUserHTTpServer(ctx context.Context, endpoints enrollment.Endpoint) http.Handler {
	mux := http.NewServeMux()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	mux.Handle("/enrollments", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeStoreEnrollment, 
		encodeResponse,
		opts...,
	))
	mux.Handle("/enrollments", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllEnrollment,
		encodeResponse,
		opts...,
	))
	mux.Handle("/enrollments/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateEnrollment,
		encodeResponse,
		opts...,
	))
	return mux
}
func decodeStoreEnrollment(_ context.Context, r *http.Request) (interface{}, error) {
	// --- Validar el Método HTTP aquí ---
	if r.Method != http.MethodPost {
		return nil, response.BadRequest(fmt.Sprintf("invalid method: '%s'", r.Method))
	}

	var req enrollment.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return req, nil
}

func decodeGetAllEnrollment(_ context.Context, r *http.Request) (interface{}, error) {

	if r.Method != http.MethodGet {
		return nil, response.BadRequest(fmt.Sprintf("invalid method: '%s'", r.Method))
	}

	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := enrollment.GetAllReq{
		UserID:   v.Get("user_id"),
		CourseID: v.Get("course_id"),
		Limit:    limit,
		Page:     page,
	}

	return req, nil
}

func decodeUpdateEnrollment(_ context.Context, r *http.Request) (interface{}, error) {
	// --- Validar el Método HTTP aquí ---
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		return nil, response.BadRequest(fmt.Sprintf("invalid method: '%s'", r.Method))
	}

	var req enrollment.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}


	pathSegments := strings.Split(r.URL.Path, "/")
	var p string
	if len(pathSegments) > 2 { // Asegurarse de que hay al menos 3 partes (/, enrollment, ID)
		p = pathSegments[2] // El ID debería ser la tercera parte (índice 2)
	} else {
		return nil, response.BadRequest("ID de inscripción no proporcionado en la URL. Formato esperado: /enrollment/{id}")
	}
	// Verifica si el ID es vacío (podría pasar si la URL termina en /enrollment/)
	if p == "" {
		return nil, response.BadRequest("ID de inscripción vacío en la URL")
	}

		req.ID = p
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(resp)
}
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}