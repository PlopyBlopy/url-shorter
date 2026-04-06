package handlers

import (
	"context"
	"testing"
	"time"

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

type MockUrlExpiredRemover struct {
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
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockUrlExpiredRemover) RemoveExpired(before time.Time, ctx context.Context) error {
	args := m.Called(before, ctx)
	return args.Error(0)
}

func TestAddUrlUsecase(t *testing.T) {
	type fields struct {
		generator     *MockShortUrlGenerator
		urladdGetter  *MockUrlAddGetter
		counterGetter *MockCounterGetter
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
				f.urladdGetter.On("GetShortUrl", "orig", mock.AnythingOfType("context.backgroundCtx")).Return("short", nil).Once()
			},
			want:    "short",
			wantErr: nil,
			args:    args{"orig", context.Background()},
		},
		{
			name: "a non-existing URL, success get counter, generate short url and added",
			setupMocks: func(f fields) {
				f.urladdGetter.On("GetShortUrl", "orig", mock.AnythingOfType("context.backgroundCtx")).Return("", domain.ErrURLSNotFound).Once()
				f.counterGetter.On("GetCounter", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(0), nil)
				f.generator.On("GenerateShortUrl", uint64(0)).Return("short").Once()
				f.urladdGetter.On("AddUrl", mock.MatchedBy(func(u domain.Url) bool {
					return u.OrigUrl == "orig" && u.ShortUrl == "short"
				}), mock.AnythingOfType("context.backgroundCtx")).Return(nil).Once()
			},
			want:    "short",
			wantErr: nil,
			args:    args{"orig", context.Background()},
		},
		{
			name: "a non-existing URL, success get counter, generate short url and unsuccessfully added",
			setupMocks: func(f fields) {
				f.urladdGetter.On("GetShortUrl", "orig", mock.AnythingOfType("context.backgroundCtx")).Return("", domain.ErrURLSNotFound).Once()
				f.counterGetter.On("GetCounter", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(0), nil).Once()
				f.generator.On("GenerateShortUrl", uint64(0)).Return("short").Once()
				f.urladdGetter.On("AddUrl", mock.MatchedBy(func(u domain.Url) bool {
					return u.OrigUrl == "orig" && u.ShortUrl == "short"
				}), mock.AnythingOfType("context.backgroundCtx")).Return(domain.ErrURLSNoAdded).Once()
			},
			want:    "",
			wantErr: domain.ErrURLSNoAdded,
			args:    args{"orig", context.Background()},
		},
		{
			name: "a non-existing URL, unsuccessfully get counter",
			setupMocks: func(f fields) {
				f.urladdGetter.On("GetShortUrl", "orig", mock.AnythingOfType("context.backgroundCtx")).Return("", domain.ErrURLSNotFound).Once()
				f.counterGetter.On("GetCounter", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(0), domain.ErrCounterNotFound).Once()
			},
			want:    "",
			wantErr: domain.ErrCounterNotFound,
			args:    args{"orig", context.Background()},
		},
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

func TestDeleteExpiredUsecase(t *testing.T) {
	type fields struct {
		expiredRemover *MockUrlExpiredRemover
	}

	type args struct {
		before time.Time
		ctx    context.Context
	}

	tests := []struct {
		name       string
		setupMocks func(f fields)
		wantErr    error
		args       args
	}{
		{
			name: "",
			setupMocks: func(f fields) {
				f.expiredRemover.On("RemoveExpired", time.Now().UTC().Truncate(time.Microsecond), mock.Anything).Return(nil).Once()
			},
			wantErr: nil,
			args:    args{time.Now().UTC().Truncate(time.Microsecond), context.Background()},
		},
	}

	for _, tt := range tests {
		// Arrange
		expiredRemover := new(MockUrlExpiredRemover)
		tt.setupMocks(fields{expiredRemover})

		// Act
		u := urls.DeleteExpiredUsecase(expiredRemover)

		err := u(tt.args.before.String(), tt.args.ctx)

		// Assert
		assert.ErrorIs(t, tt.wantErr, err)

		expiredRemover.AssertExpectations(t)
	}
}
