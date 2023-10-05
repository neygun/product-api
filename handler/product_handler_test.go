package handler

import (
	"chi-demo/model"
	"chi-demo/service"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetOne(t *testing.T) {
	type mockGetOneService struct {
		expCall bool
		output  model.Product
		err     error
	}

	type args struct {
		givenID           string
		mockGetOneService mockGetOneService
		expStatusCode     int
		expResponse       string
	}

	tcs := map[string]args{
		"success": {
			givenID: "1",
			mockGetOneService: mockGetOneService{
				expCall: true,
				output: model.Product{
					ID:    1,
					Name:  "test",
					Price: 1,
				},
			},
			expStatusCode: http.StatusOK,
			expResponse: ToJsonString(model.Product{
				ID:    1,
				Name:  "test",
				Price: 1,
			}),
		},
		"err - cannot convert id": {
			givenID:       "abc",
			expStatusCode: http.StatusBadRequest,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusBadRequest,
				Description: "Cannot convert id to integer",
			}),
		},
		"err - invalid id": {
			givenID:       "-1",
			expStatusCode: http.StatusBadRequest,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusBadRequest,
				Description: "Invalid id",
			}),
		},
		"service error": {
			givenID: "1",
			mockGetOneService: mockGetOneService{
				expCall: true,
				err:     errors.New("test"),
			},
			expStatusCode: http.StatusInternalServerError,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusInternalServerError,
				Description: "Internal Server Error",
			}),
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			req := httptest.NewRequest(http.MethodGet, "/products", nil)
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", tc.givenID)
			// req.Header.Set("Content-Type", "application/json")
			// bearer := "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjN9.PZLMJBT9OIVG2qgp9hQr685oVYFgRgWpcSPmNcw6y7M"
			// req.Header.Add("Authorization", bearer)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
			res := httptest.NewRecorder()

			req = req.WithContext(ctx)

			mockProductService := service.NewMockProductService(t)

			// When
			id, _ := strconv.ParseInt(tc.givenID, 10, 64)
			if tc.mockGetOneService.expCall {
				mockProductService.ExpectedCalls = []*mock.Call{
					mockProductService.On("GetOne", ctx, id).Return(tc.mockGetOneService.output, tc.mockGetOneService.err),
				}
			}
			instance := New(mockProductService)
			handler := instance.GetOne()
			handler.ServeHTTP(res, req)

			// Then
			require.Equal(t, tc.expStatusCode, res.Code)
			if tc.expResponse != "" {
				require.JSONEq(t, tc.expResponse, res.Body.String())
			}
		})
	}
}

func TestHandler_GetProducts(t *testing.T) {
	type mockGetAllService struct {
		expCall bool
		output  []model.Product
		err     error
	}

	type args struct {
		mockGetAllService mockGetAllService
		expStatusCode     int
		expResponse       string
	}

	tcs := map[string]args{
		"success": {
			mockGetAllService: mockGetAllService{
				expCall: true,
				output: []model.Product{
					{
						ID:    1,
						Name:  "test1",
						Price: 1,
					},
					{
						ID:    2,
						Name:  "test2",
						Price: 2,
					},
				},
			},
			expStatusCode: http.StatusOK,
			expResponse: ToJsonString([]model.Product{
				{
					ID:    1,
					Name:  "test1",
					Price: 1,
				},
				{
					ID:    2,
					Name:  "test2",
					Price: 2,
				},
			}),
		},
		"empty": {
			mockGetAllService: mockGetAllService{
				expCall: true,
				output:  nil,
			},
			expStatusCode: http.StatusOK,
			expResponse:   ToJsonString(nil),
		},
		"service error": {
			mockGetAllService: mockGetAllService{
				expCall: true,
				err:     errors.New("test"),
			},
			expStatusCode: http.StatusInternalServerError,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusInternalServerError,
				Description: "Internal Server Error",
			}),
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			req := httptest.NewRequest(http.MethodGet, "/products", nil)
			routeCtx := chi.NewRouteContext()
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
			res := httptest.NewRecorder()

			req = req.WithContext(ctx)

			mockProductService := service.NewMockProductService(t)

			// When
			if tc.mockGetAllService.expCall {
				mockProductService.ExpectedCalls = []*mock.Call{
					mockProductService.On("GetAll", ctx).Return(tc.mockGetAllService.output, tc.mockGetAllService.err),
				}
			}
			instance := New(mockProductService)
			handler := instance.GetProducts()
			handler.ServeHTTP(res, req)

			// Then
			require.Equal(t, tc.expStatusCode, res.Code)
			if tc.expResponse != "" {
				require.JSONEq(t, tc.expResponse, res.Body.String())
			}
		})
	}
}

func TestHandler_CreateProduct(t *testing.T) {
	type mockCreateService struct {
		expCall bool
		err     error
	}

	type args struct {
		givenRequest      string
		mockCreateService mockCreateService
		expStatusCode     int
		expResponse       string
	}

	tcs := map[string]args{
		"success": {
			givenRequest: `{"name":"test","price":1}`,
			mockCreateService: mockCreateService{
				expCall: true,
			},
			expStatusCode: http.StatusOK,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusOK,
				Description: "Product created",
			}),
		},
		"err - invalid product": {
			expStatusCode: http.StatusBadRequest,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusBadRequest,
				Description: "Invalid product",
			}),
		},
		"err - missing field": {
			givenRequest:  `{"name":"test"}`,
			expStatusCode: http.StatusBadRequest,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusBadRequest,
				Description: "Missing field",
			}),
		},
		"err - invalid price": {
			givenRequest:  `{"name":"test","price":-1}`,
			expStatusCode: http.StatusBadRequest,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusBadRequest,
				Description: "Invalid price",
			}),
		},
		"service error": {
			givenRequest: `{"name":"test","price":1}`,
			mockCreateService: mockCreateService{
				expCall: true,
				err:     errors.New("test"),
			},
			expStatusCode: http.StatusInternalServerError,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusInternalServerError,
				Description: "Internal Server Error",
			}),
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(tc.givenRequest))
			routeCtx := chi.NewRouteContext()
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
			res := httptest.NewRecorder()

			req = req.WithContext(ctx)

			mockProductService := service.NewMockProductService(t)

			// When
			if tc.mockCreateService.expCall {
				// i := ToStruct(tc.givenRequest)
				// product := i.(model.Product)
				product := model.Product{
					Name:  "test",
					Price: 1,
				}

				mockProductService.ExpectedCalls = []*mock.Call{
					mockProductService.On("Create", ctx, product).Return(tc.mockCreateService.err),
				}
			}
			instance := New(mockProductService)
			handler := instance.CreateProduct()
			handler.ServeHTTP(res, req)

			// Then
			require.Equal(t, tc.expStatusCode, res.Code)
			if tc.expResponse != "" {
				require.JSONEq(t, tc.expResponse, res.Body.String())
			}
		})
	}
}

func TestHandler_DeleteProduct(t *testing.T) {
	type mockDeleteService struct {
		expCall bool
		err     error
	}

	type args struct {
		givenID           string
		mockDeleteService mockDeleteService
		expStatusCode     int
		expResponse       string
	}

	tcs := map[string]args{
		"success": {
			givenID: "1",
			mockDeleteService: mockDeleteService{
				expCall: true,
			},
			expStatusCode: http.StatusOK,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusOK,
				Description: "Product deleted",
			}),
		},
		"err - cannot convert id": {
			givenID:       "abc",
			expStatusCode: http.StatusBadRequest,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusBadRequest,
				Description: "Cannot convert id to integer",
			}),
		},
		"err - invalid id": {
			givenID:       "-1",
			expStatusCode: http.StatusBadRequest,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusBadRequest,
				Description: "Invalid id",
			}),
		},
		"service error": {
			givenID: "1",
			mockDeleteService: mockDeleteService{
				expCall: true,
				err:     errors.New("test"),
			},
			expStatusCode: http.StatusInternalServerError,
			expResponse: ToJsonString(model.Response{
				Code:        http.StatusInternalServerError,
				Description: "Internal Server Error",
			}),
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			req := httptest.NewRequest(http.MethodDelete, "/products", nil)
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", tc.givenID)
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
			res := httptest.NewRecorder()

			req = req.WithContext(ctx)

			mockProductService := service.NewMockProductService(t)

			// When
			if tc.mockDeleteService.expCall {
				id, _ := strconv.ParseInt(tc.givenID, 10, 64)
				mockProductService.ExpectedCalls = []*mock.Call{
					mockProductService.On("Delete", ctx, id).Return(tc.mockDeleteService.err),
				}
			}
			instance := New(mockProductService)
			handler := instance.DeleteProduct()
			handler.ServeHTTP(res, req)

			// Then
			require.Equal(t, tc.expStatusCode, res.Code)
			if tc.expResponse != "" {
				require.JSONEq(t, tc.expResponse, res.Body.String())
			}
		})
	}
}
