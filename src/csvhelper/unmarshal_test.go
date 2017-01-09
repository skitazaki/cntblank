package csvhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

type TSPrefecture struct {
	Code string
	Name string
}

type TSDrug struct {
	Code    string
	Name    null.String
	Generic null.Bool
	Status  null.Int
	Price   null.Float
}

func TestUnmarshal(t *testing.T) {
	var p TSPrefecture
	err := Unmarshal([]string{"都道府県コード", "都道府県"}, &p, true)
	require.Nil(t, err, "Unmarshal returns error: %v", err)
	assert.Equal(t, "都道府県コード", p.Code, "Expected code is 都道府県コード, got %s", p.Code)
	assert.Equal(t, "都道府県", p.Name, "Expected name is 都道府県, got %s", p.Name)
}

func TestUnmarshalNull(t *testing.T) {
	var err error
	var p TSDrug
	// Pattern 1
	err = Unmarshal([]string{"123456789", "サンプル", "1", "9", "5.6"}, &p, true)
	require.Nil(t, err, "Unmarshal returns error: %v", err)
	assert.Equal(t, "123456789", p.Code, "Code is unexpected")
	assert.Equal(t, true, p.Name.Valid, "Name validity is unexpected: %v", p.Name)
	assert.Equal(t, "サンプル", p.Name.String, "Name value is unexpected: %v", p.Name)
	assert.Equal(t, true, p.Generic.Valid, "Generic validity is unexpected: %v", p.Generic)
	assert.Equal(t, true, p.Generic.Bool, "Generic value is unexpected: %v", p.Generic)
	assert.Equal(t, true, p.Status.Valid, "Status validity is unexpected: %v", p.Status)
	assert.Equal(t, int64(9), p.Status.Int64, "Status value is unexpected: %v", p.Status)
	assert.Equal(t, true, p.Price.Valid, "Price validity is unexpected: %v", p.Price)
	assert.Equal(t, 5.6, p.Price.Float64, "Price value is unexpected: %v", p.Price)
	// Pattern 2
	err = Unmarshal([]string{"123456789", "", "", "", ""}, &p, true)
	require.Nil(t, err, "Unmarshal returns error: %v", err)
	assert.Equal(t, "123456789", p.Code, "Code is unexpected")
	assert.Equal(t, false, p.Name.Valid, "Name validity is unexpected: %v", p.Name)
	assert.Equal(t, "", p.Name.String, "Name value is unexpected: %v", p.Name)
	assert.Equal(t, false, p.Generic.Valid, "Generic validity is unexpected: %v", p.Generic)
	assert.Equal(t, false, p.Generic.Bool, "Generic value is unexpected: %v", p.Generic)
	assert.Equal(t, false, p.Status.Valid, "Status validity is unexpected: %v", p.Status)
	assert.Equal(t, int64(0), p.Status.Int64, "Status value is unexpected: %v", p.Status)
	assert.Equal(t, false, p.Price.Valid, "Price validity is unexpected: %v", p.Price)
	assert.Equal(t, float64(0), p.Price.Float64, "Price value is unexpected: %v", p.Price)
}
