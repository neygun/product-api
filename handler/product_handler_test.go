package handler

import (
	"chi-demo/model"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetOne(t *testing.T) {
	type args struct {
		givenID       string
		givenRequest  string
		expStatusCode int
		expResponse   model.Response
	}

	tcs := map[string]args{
		"err - cannot convert id": {
			givenID:       "id",
			expStatusCode: http.StatusBadRequest,
			expResponse: model.Response{
				Code:        http.StatusBadRequest,
				Description: "Cannot convert id to integer",
			},
		},
		"err - empty authorization_code": {
			givenID:       "id",
			givenRequest:  `{"redirect_uri":"uob://app"}`,
			expStatusCode: http.StatusBadRequest,
			expResponse: testhelper.BuiErrMessageResponse(
				exception.AuthorizeBadRequest.CodeString(),
				exception.AuthorizeBadRequest.Error(),
			),
		},
		"err - empty redirect_uri": {
			givenID:       "iam_id",
			givenRequest:  `{"authorization_code":"abczyxjkhz"}`,
			expStatusCode: http.StatusBadRequest,
			expResponse: testhelper.BuiErrMessageResponse(
				exception.AuthorizeBadRequest.CodeString(),
				exception.AuthorizeBadRequest.Error(),
			),
		},
		"err - UserLoggedNotFound": {
			expStatusCode: http.StatusUnauthorized,
			expResponse: testhelper.BuiErrMessageResponse(
				"PS_10100",
				"User Logged Not Found",
			),
		},
		"err - unable to authorize": {
			givenID:       "iam_id",
			givenRequest:  `{"authorization_code":"abczyxjkhz","redirect_uri":"uob://app"}`,
			expStatusCode: http.StatusUnprocessableEntity,
			expResponse: testhelper.BuiErrMessageResponse(
				"PS_20300",
				"unable to authorize",
			),
			mockAuthorize: mockAuthorize{
				wantCall: true,
				mockIn: authorization.AuthorizeInput{
					AuthorizationCode: "abczyxjkhz",
					RedirectURI:       "uob://app",
					IamID:             "iam_id",
				},
				mockErr: custom_error.NewBusinessLogicError("unable to authorize", 0, 3),
			},
		},
		"success": {
			givenID:       "iam_id",
			givenRequest:  `{"authorization_code":"abczyxjkhz","redirect_uri":"uob://app"}`,
			expStatusCode: http.StatusOK,
			mockAuthorize: mockAuthorize{
				wantCall: true,
				mockIn: authorization.AuthorizeInput{
					AuthorizationCode: "abczyxjkhz",
					RedirectURI:       "uob://app",
					IamID:             "iam_id",
				},
				mockErr: nil,
			},
		},
	}
	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Mock
			mockAuthCtrl := new(authorization.MockController)
			if tc.mockAuthorize.wantCall {
				mockAuthCtrl.On("Authorize", mock.Anything, tc.mockAuthorize.mockIn).
					Return(tc.mockAuthorize.mockErr)
			}

			// Given
			req := httptest.NewRequest(http.MethodGet, "/products/{id}", strings.NewReader(tc.givenRequest))
			routeCtx := chi.NewRouteContext()
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
			res := httptest.NewRecorder()
			if tc.givenID != "" {
				ctx = iam.SetUserProfileInContextForTestingOnly(ctx, tc.givenID, iam.UserMetadata{}, []string{}, []iam.RBAC{})
			}
			req = req.WithContext(ctx)

			// When
			instance := New(mockAuthCtrl)
			handler := instance.Authorize()
			handler.ServeHTTP(res, req)

			// Then
			require.Equal(t, tc.expStatusCode, res.Code)
			if tc.expResponse != "" {
				require.JSONEq(t, tc.expResponse, res.Body.String())
			}
			mockAuthCtrl.AssertExpectations(t)
		})
	}
}
