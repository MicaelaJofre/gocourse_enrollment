package enrollment

import (
	"context"
	"errors"

	courseSdk "github.com/MicaelaJofre/go_course_sdk/course"
	userSdk "github.com/MicaelaJofre/go_course_sdk/user"
	"github.com/MicaelaJofre/go_lib_response/response"
	"github.com/MicaelaJofre/gocourse_meta/meta"
)

type Controller func(ctx context.Context, request interface{}) (interface{}, error)

type Endpoint struct {
		Create Controller
		GetAll Controller
		Update Controller
}

type CreateReq struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
	}

type GetAllReq struct {
		UserID   string
		CourseID string
		Limit    int
		Page     int
	}

type UpdateReq struct {
		ID     string
		Status *string `json:"status"`
	}
type Config struct {
		LimPageDef string
	}


func MakeEndpoint(s Service, config Config) Endpoint {
	return Endpoint{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		req := request.(GetAllReq)

		filters := Filters{
			UserID:   req.UserID,
			CourseID: req.CourseID,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		enrollments, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", enrollments, meta), nil
	}
} 

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReq)

		if req.UserID == "" {
			return nil, response.BadRequest(ErrUserIDRequired.Error())
		}

		if req.CourseID == "" {
			return nil, response.BadRequest(ErrCourseIDRequired.Error())
		}

		enroll, err := s.Create(ctx, req.UserID, req.CourseID)
		if err != nil {

			if errors.As(err, &userSdk.ErrNotFound{}) ||
				errors.As(err, &courseSdk.ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", enroll, nil), nil

	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReq)

		if req.Status != nil && *req.Status == "" {
			return nil, response.BadRequest(ErrStatusRequired.Error())
		}

		if err := s.Update(ctx, req.ID, req.Status); err != nil {

			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}

			if errors.As(err, &ErrInvalidStatus{}) {
				return nil, response.BadRequest(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}
