package service

import (
	"chi-demo/model"
	"chi-demo/repository"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestController_GetProducts(t *testing.T) {
	type mockGetProductsRepo struct {
		expCall bool
		output  []model.Product
		err     error
	}

	tcs := map[string]struct {
		mockGetProductsRepo mockGetProductsRepo
		expRes              []model.Product
		expErr              error
	}{
		"success": {
			mockGetProductsRepo: mockGetProductsRepo{
				expCall: true,
				output: []model.Product{
					{
						ID:    1,
						Name:  "test",
						Price: 1,
					},
				},
			},
			expRes: []model.Product{
				{
					ID:    1,
					Name:  "test",
					Price: 1,
				},
			},
		},
		"error": {
			mockGetProductsRepo: mockGetProductsRepo{
				expCall: true,
				output:  nil,
				err:     errors.New("test"),
			},
			expRes: nil,
			expErr: errors.New("test"),
		},
		"empty": {
			mockGetProductsRepo: mockGetProductsRepo{
				expCall: true,
				output:  nil,
			},
			expRes: nil,
		},
	}

	for scenario, tc := range tcs {
		t.Run(scenario, func(t *testing.T) {
			// Given
			ctx := context.Background()
			// repo := new(repository.MockRegistry)
			mockProductRepo := repository.NewMockProductRepository(t)
			// repo.ExpectedCalls = []*mock.Call{
			// 	repo.On("Product").Return(mockProductRepo),
			// }

			// When
			if tc.mockGetProductsRepo.expCall {
				mockProductRepo.ExpectedCalls = []*mock.Call{
					mockProductRepo.On("GetAll", ctx).Return(tc.mockGetProductsRepo.output, tc.mockGetProductsRepo.err),
				}
			}

			serv := New(mockProductRepo)
			rs, err := serv.GetAll(ctx)

			// Then
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expRes, rs)
			}
		})
	}
}
