package helpers

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zamurabims/QA_avito/task2/internal/client"
	"github.com/zamurabims/QA_avito/task2/internal/models"
)

var uuidRegexp = regexp.MustCompile(
	`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
)

func RequireOK(t testing.TB, resp *client.Response, err error) {
	t.Helper()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode,
		"unexpected status, body: %s", resp.Body)
}

func RequireBadRequest(t testing.TB, resp *client.Response, err error) {
	t.Helper()
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode,
		"expected 400, got %d, body: %s", resp.StatusCode, resp.Body)
}

func RequireNotFound(t testing.TB, resp *client.Response, err error) {
	t.Helper()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode,
		"expected 404, got %d, body: %s", resp.StatusCode, resp.Body)
}

func RequireValidUUID(t testing.TB, id string) {
	t.Helper()
	require.Regexp(t, uuidRegexp, id, "expected valid UUID, got: %s", id)
}

func RequireValidItem(t testing.TB, item models.Item) {
	t.Helper()
	RequireValidUUID(t, item.ID)
	require.NotEmpty(t, item.Name)
	require.NotEmpty(t, item.CreatedAt)
	require.GreaterOrEqual(t, item.Price, 0)
	require.GreaterOrEqual(t, item.SellerID, 0)
}

func RequireValidStatistics(t testing.TB, s models.Statistics) {
	t.Helper()
	require.GreaterOrEqual(t, s.Likes, 0)
	require.GreaterOrEqual(t, s.ViewCount, 0)
	require.GreaterOrEqual(t, s.Contacts, 0)
}

func RequireContentTypeJSON(t testing.TB, resp *client.Response) {
	t.Helper()
	require.Contains(t, resp.Headers.Get("Content-Type"), "application/json")
}
