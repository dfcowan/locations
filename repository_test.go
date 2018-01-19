package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveUser(t *testing.T) {
	migrate()

	tkn := accessToken{
		AccessToken:  "ar5j7SDvII051ed2DPa2wVyxvQOEABg3b6pGh_b0XV5BeE15Uke47SA817LmhQT6",
		RefreshToken: "C4Jpg9MbIwKfYNK_QimBexYI08c4lrO2pe2qemMehJoJhQYj4s0l2twiTG4Ug355",
		UserID:       3,
	}

	firstDate := "20151216"

	err := createUser(tkn, firstDate)
	require.Nil(t, err)

	usr, uss, err := loadUser(tkn.UserID)
	require.Nil(t, err)

	assert.Equal(t, tkn.UserID, usr.UserID)
	assert.Equal(t, tkn.AccessToken, usr.AccessToken)
	assert.Equal(t, tkn.RefreshToken, usr.RefreshToken)
	assert.Equal(t, tkn.UserID, uss.UserID)
	assert.Equal(t, firstDate, uss.StartDate)
}

func TestSaveBreadcrumbs(t *testing.T) {
	migrate()

	id := 3
	bcs := []breadcrumb{
		{
			Coordinate: coordinate{
				Lat: 1,
				Lon: 1,
			},
			Time: "20151204",
		},
		{
			Coordinate: coordinate{
				Lat: 2,
				Lon: 2,
			},
			Time: "20151214",
		},
	}

	err := saveBreadcrumbs(id, bcs)

	assert.Nil(t, err)
}

func TestLoadCounts(t *testing.T) {
	migrate()

	cnts, err := loadCounts(56200834001288605, "20151216", "20180118")

	require.Nil(t, err)

	assert.NotZero(t, len(cnts))
}
