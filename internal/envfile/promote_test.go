package envfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromote_BasicAllKeys(t *testing.T) {
	src := EnvMap{"A": "1", "B": "2"}
	dst := EnvMap{"C": "3"}
	res, err := Promote(src, dst, PromoteOptions{})
	require.NoError(t, err)
	assert.Equal(t, "1", res.Merged["A"])
	assert.Equal(t, "2", res.Merged["B"])
	assert.NotContains(t, res.Merged, "C")
	assert.ElementsMatch(t, []string{"A", "B"}, res.Added)
}

func TestPromote_PreserveTarget(t *testing.T) {
	src := EnvMap{"A": "new"}
	dst := EnvMap{"B": "keep", "A": "old"}
	res, err := Promote(src, dst, PromoteOptions{PreserveTarget: true})
	require.NoError(t, err)
	assert.Equal(t, "new", res.Merged["A"])
	assert.Equal(t, "keep", res.Merged["B"])
	assert.ElementsMatch(t, []string{"A"}, res.Updated)
}

func TestPromote_IncludeKeys(t *testing.T) {
	src := EnvMap{"A": "1", "B": "2", "C": "3"}
	res, err := Promote(src, nil, PromoteOptions{IncludeKeys: []string{"A", "C"}})
	require.NoError(t, err)
	assert.Contains(t, res.Merged, "A")
	assert.NotContains(t, res.Merged, "B")
	assert.Contains(t, res.Merged, "C")
	assert.Contains(t, res.Skipped, "B")
}

func TestPromote_ExcludeKeys(t *testing.T) {
	src := EnvMap{"SECRET": "x", "PORT": "8080"}
	res, err := Promote(src, nil, PromoteOptions{ExcludeKeys: []string{"SECRET"}})
	require.NoError(t, err)
	assert.NotContains(t, res.Merged, "SECRET")
	assert.Equal(t, "8080", res.Merged["PORT"])
	assert.Contains(t, res.Skipped, "SECRET")
}

func TestPromote_NilSourceReturnsError(t *testing.T) {
	_, err := Promote(nil, EnvMap{}, PromoteOptions{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source env map must not be nil")
}

func TestPromote_EmptySource(t *testing.T) {
	res, err := Promote(EnvMap{}, EnvMap{"A": "1"}, PromoteOptions{PreserveTarget: true})
	require.NoError(t, err)
	assert.Equal(t, "1", res.Merged["A"])
	assert.Empty(t, res.Added)
	assert.Empty(t, res.Updated)
}
