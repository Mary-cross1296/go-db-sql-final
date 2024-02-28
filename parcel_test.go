package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	assert.True(t, number > 0)

	parcel.Number = number

	row, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, parcel, row)

	err = store.Delete(number)
	require.NoError(t, err)

	_, err = store.Get(number)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	assert.True(t, number > 0)

	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)

	row, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, newAddress, row.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	number, err := store.Add(parcel)
	require.NoError(t, err)
	assert.True(t, number > 0)

	newAStatus := ParcelStatusDelivered
	err = store.SetStatus(number, newAStatus)
	require.NoError(t, err)

	row, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, newAStatus, row.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	parcels := []Parcel{
		parcel,
		parcel,
		parcel,
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client

	for i := 0; i < len(parcels); i++ {
		number, err := store.Add(parcels[i])
		require.NoError(t, err)
		assert.True(t, number > 0)
		parcels[i].Number = number
		parcelMap[number] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, 2, len(storedParcels))

	for _, parcel := range storedParcels {
		require.Contains(t, parcelMap, parcel.Number)
		parcelMapValue := parcelMap[parcel.Number]

		assert.Equal(t, parcelMapValue.Client, parcel.Client)
		assert.Equal(t, parcelMapValue.Status, parcel.Status)
		assert.Equal(t, parcelMapValue.Address, parcel.Address)
		assert.Equal(t, parcelMapValue.CreatedAt, parcel.CreatedAt)
	}
}
