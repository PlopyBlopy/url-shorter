///go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/PlopyBlopy/url-shorter/internal"
	"github.com/PlopyBlopy/url-shorter/internal/adapters"
	"github.com/PlopyBlopy/url-shorter/internal/api"
	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/PlopyBlopy/url-shorter/tests/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlersV1(t *testing.T) {
	pgContainerCtx := context.Background()

	// postgres container
	pgContainer, err := testdata.NewPostgresTestcontainer(pgContainerCtx)
	require.NoError(t, err)

	// config
	config, err := testdata.NewHTTPConfig()
	require.NoError(t, err)

	// http test suite
	testSuiteHttp, err := testdata.NewTestSuiteHTTP(config.Domain, config.Port, pgContainer, pgContainerCtx)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := pgContainer.Terminate(pgContainerCtx)
		assert.NoError(t, err)
	})

	// base url, http client

	t.Run("add_urls", func(t *testing.T) {
		// Arrange
		require := require.New(t)
		assert := assert.New(t)

		ctx := context.Background()

		err := testSuiteHttp.SetupTestPg(ctx)
		require.NoError(err)

		rep := adapters.NewRepository(testSuiteHttp.Db)
		g, err := internal.NewGenerator()
		require.NoError(err)

		err, shutdown := testSuiteHttp.SetupTestHTTP(api.NewRouter(1, g, rep), ctx)
		require.NoError(err)
		t.Cleanup(func() {
			err := shutdown()
			assert.NoError(err)
		})

		expectedUrl, err := testdata.GetUrl(1, 0)
		require.NoError(err)

		// Act
		data, err := json.Marshal(map[string]string{"url": expectedUrl.OrigUrl})
		require.NoError(err)
		req, err := http.NewRequest(http.MethodPost, testSuiteHttp.BaseUrl+"/", bytes.NewBuffer(data))
		require.NoError(err)
		req.Header.Set("Content-Type", "application/json")

		respShortUrl, err := testSuiteHttp.Client.Do(req)
		require.NoError(err)
		t.Cleanup(func() {
			err := respShortUrl.Body.Close()
			assert.NoError(err)
		})

		bodyShortUrl, err := io.ReadAll(respShortUrl.Body)
		require.NoError(err)

		var shortUrl domain.ShortUrl
		err = json.Unmarshal(bodyShortUrl, &shortUrl)
		require.NoError(err)

		actualUrl, err := rep.GetUrl(shortUrl.ShortUrl, ctx)
		require.NoError(err)

		// Assert
		assert.NotEqual(domain.Url{}, actualUrl)
		assert.Equal(expectedUrl.ShortUrl, actualUrl.ShortUrl)
	})

	t.Run("get_url", func(t *testing.T) {
		// Arrange
		require := require.New(t)
		assert := assert.New(t)

		ctx := context.Background()

		err := testSuiteHttp.SetupTestPg(ctx)
		require.NoError(err)

		rep := adapters.NewRepository(testSuiteHttp.Db)
		g, err := internal.NewGenerator()
		require.NoError(err)

		err, shutdown := testSuiteHttp.SetupTestHTTP(api.NewRouter(1, g, rep), ctx)
		require.NoError(err)
		t.Cleanup(func() {
			err := shutdown()
			assert.NoError(err)
		})

		expectedUrl, err := testdata.GetUrl(0, 1)
		require.NoError(err)

		// Act
		err = rep.AddUrl(expectedUrl, ctx)
		require.NoError(err)

		data, err := json.Marshal(map[string]string{"url": expectedUrl.ShortUrl})
		require.NoError(err)
		req, err := http.NewRequest(http.MethodGet, testSuiteHttp.BaseUrl+"/url", bytes.NewBuffer(data))
		require.NoError(err)
		req.Header.Set("Content-Type", "application/json")

		respUrl, err := testSuiteHttp.Client.Do(req)
		require.NoError(err)
		t.Cleanup(func() {
			respUrl.Body.Close()
		})

		bodyUrl, err := io.ReadAll(respUrl.Body)
		require.NoError(err)

		var actualUrl domain.Url
		err = json.Unmarshal(bodyUrl, &actualUrl)
		require.NoError(err)

		// Assert
		assert.NotEqual(domain.Url{}, actualUrl)
		assert.Equal(expectedUrl, actualUrl)
	})

	t.Run("get_urls", func(t *testing.T) {
		// Arrange
		require := require.New(t)
		assert := assert.New(t)

		ctx := context.Background()

		err := testSuiteHttp.SetupTestPg(ctx)
		require.NoError(err)

		rep := adapters.NewRepository(testSuiteHttp.Db)
		g, err := internal.NewGenerator()
		require.NoError(err)

		err, shutdown := testSuiteHttp.SetupTestHTTP(api.NewRouter(1, g, rep), ctx)
		require.NoError(err)
		t.Cleanup(func() {
			err := shutdown()
			assert.NoError(err)
		})

		expectedUrls, err := testdata.GetUrls(10, 0, 10)
		require.NoError(err)

		tests := []struct {
			name         string
			limit        int
			expectedUrls []domain.Url
		}{
			{"get 1", 1, expectedUrls[:1]},
			{"get 5", 5, expectedUrls[:5]},
			{"get 10", 10, expectedUrls[:]},
		}

		// Act
		err = rep.AddUrls(expectedUrls, ctx)
		require.NoError(err)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				params := url.Values{}
				params.Add("limit", strconv.Itoa(tt.limit))

				req, err := http.NewRequest(http.MethodGet, testSuiteHttp.BaseUrl+"/urls?"+params.Encode(), nil)
				require.NoError(err)

				resp, err := testSuiteHttp.Client.Do(req)
				require.NoError(err)
				t.Cleanup(func() {
					err := resp.Body.Close()
					if err != nil {
						assert.NoError(err)
					}
				})

				body, err := io.ReadAll(resp.Body)
				require.NoError(err)

				var actualUrls []domain.Url
				err = json.Unmarshal(body, &actualUrls)
				require.NoError(err)

				//Assert
				assert.Equal(tt.expectedUrls, actualUrls[:tt.limit])
			})
		}
	})
}
