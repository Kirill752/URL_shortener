package tests

import (
	"fmt"
	"net/url"
	"testing"
	"urlShotener/internal/http-server/handlers/url/del"
	"urlShotener/internal/http-server/handlers/url/save"
	"urlShotener/internal/lib/api"
	"urlShotener/internal/lib/random"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())
	e.POST("/url/save").WithJSON(save.Request{
		URL:   gofakeit.URL(),
		Alias: random.CreateRandomString(10),
	}).WithBasicAuth("admin", "12345").Expect().Status(200).JSON().Object().ContainsKey("alias")
}

func TestURLShortener_SaveRedirectDelete(t *testing.T) {
	// TODO: придумать больше тест кейсов
	testCases := []struct {
		name  string
		url   string
		alias string
		err   string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "not URL",
			alias: gofakeit.Word(),
			err:   "field URL is not valid URL",
		},
		{
			name:  "Error in URL",
			url:   "https://",
			alias: gofakeit.Word(),
			err:   "field URL is not valid URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Отправляем запрос на сохранение URL
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}
			e := httpexpect.Default(t, u.String())
			resp := e.POST("/url/save").WithJSON(save.Request{
				URL:   tc.url,
				Alias: tc.alias,
			}).WithBasicAuth("admin", "12345").Expect().Status(200).JSON().Object()

			// Если поле ошибки в тесте не пустое, то проверяем его равенство полученной ошибке
			if tc.err != "" {
				resp.Value("error").String().IsEqual(tc.err)
				return
			}
			// Провиряем алиас
			alias := tc.alias
			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}
			// Проверяем Redirect
			r := url.URL{
				Scheme: "http",
				Host:   host,
				Path:   alias,
			}
			redirectionResult, err := api.GetRedirect(r.String())
			require.NoError(t, err)
			require.Equal(t, tc.url, redirectionResult)
			// Проверяем Delete
			reqDel := e.DELETE("/url/delete").WithJSON(del.Request{
				Alias: alias,
			}).WithBasicAuth("admin", "12345").Expect().Status(200).JSON().Object()
			reqDel.Value("status").String().IsEqual("OK")
			// Проверяем, что редирект не происходит
			redirectionResult, err = api.GetRedirect(r.String())
			require.Equal(t, err, fmt.Errorf("%s: %w: %d", "api.GetRedirect", api.ErrInvalidStatusCode, 200))
		})
	}
}
