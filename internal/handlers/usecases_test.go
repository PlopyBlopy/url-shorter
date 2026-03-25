package handlers

import (
	"context"
	"testing"

	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/PlopyBlopy/url-shorter/internal/handlers/urls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockShortUrlGenerator struct {
	mock.Mock
}

type MockUrlAddGetter struct {
	mock.Mock
}

type MockCounterGetter struct {
	mock.Mock
}

func (m *MockShortUrlGenerator) GenerateShortUrl(counter uint64) string {
	args := m.Called(counter)
	return args.String(0)
}

func (m *MockUrlAddGetter) AddUrl(url domain.Url, ctx context.Context) error {
	args := m.Called(url, ctx)
	return args.Error(0)
}
func (m *MockUrlAddGetter) GetShortUrl(origUrl string, ctx context.Context) (string, error) {
	args := m.Called(origUrl, ctx)
	return args.String(0), args.Error(1)
}

func (m *MockCounterGetter) GetCounter(ctx context.Context) (uint64, error) {
	args := m.Called()
	return args.Get(0).(uint64), args.Error(1)
}

func TestAddUrlUsecase(t *testing.T) {
	type fields struct {
		generator  *MockShortUrlGenerator
		urlRep     *MockUrlAddGetter
		counterRep *MockCounterGetter
	}

	type args struct {
		origUrl string
		ctx     context.Context
	}

	tests := []struct {
		name       string
		setupMocks func(f fields)
		want       string
		wantErr    error
		args       args
	}{
		{
			name: "success existing URL",
			setupMocks: func(f fields) {
				f.urlRep.On("GetShortUrl", "orig", mock.Anything).Return("short", nil).Once()
			},
			want:    "short",
			wantErr: nil,
			args:    args{"orig", context.Background()},
		},
		{
			name: "a non-existing URL, success generate short url and added",
			setupMocks: func(f fields) {
				f.urlRep.On("GetShortUrl", "orig", mock.Anything).Return("", domain.ErrURLSNotFound).Once()
				f.generator.On("GenerateShortUrl").Return("short").Once()
				f.urlRep.On("AddUrl", mock.MatchedBy(func(u domain.Url) bool {
					return u.OrigUrl == "orig" && u.ShortUrl == "short"
				}), mock.Anything).Return(nil).Once()
			},
			want:    "short",
			wantErr: nil,
			args:    args{"orig", context.Background()},
		},
		// {
		// 	name: "a non-existing URL, success generate short url and unsuccessfully added",
		// },
		// {
		// 	name: "other db error",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			generator := new(MockShortUrlGenerator)
			urlRep := new(MockUrlAddGetter)
			counterRep := new(MockCounterGetter)
			tt.setupMocks(fields{generator, urlRep, counterRep})

			// Act
			u := urls.AddUrlUsecase(generator, urlRep, counterRep)
			got, err := u(tt.args.origUrl, tt.args.ctx)

			// Assert
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)

			generator.AssertExpectations(t)
			urlRep.AssertExpectations(t)
			counterRep.AssertExpectations(t)
		})
	}
}
