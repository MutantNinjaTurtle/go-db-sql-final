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
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)
	parcelNew, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, parcel.Client, parcelNew.Client)
	assert.Equal(t, parcel.Address, parcelNew.Address)
	assert.Equal(t, parcel.Status, parcelNew.Status)
	assert.Equal(t, parcel.CreatedAt, parcelNew.CreatedAt)
	assert.Equal(t, parcel.Number, parcelNew.Number)
	err = store.Delete(id)
	require.NoError(t, err)
	_, err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err, nil)
	parcelNewAdress, err := store.Get(id)
	require.NoError(t, err)
	assert.NotEqual(t, parcel.Address, parcelNewAdress.Address)
	assert.Equal(t, parcelNewAdress.Address, newAddress)
	assert.Equal(t, parcel.Client, parcelNewAdress.Client)
	assert.Equal(t, parcel.Status, parcelNewAdress.Status)
	assert.Equal(t, parcel.CreatedAt, parcelNewAdress.CreatedAt)
	assert.Equal(t, parcel.Number, parcelNewAdress.Number)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db")
	store := NewParcelStore(db)
	parcel := getTestParcel()
	id, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, id)
	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)
	parcelNewStatus, err := store.Get(id)
	require.NoError(t, err)
	assert.NotEqual(t, parcel.Status, parcelNewStatus.Status)
	assert.Equal(t, parcelNewStatus.Status, ParcelStatusSent)
	assert.Equal(t, parcel.Client, parcelNewStatus.Client)
	assert.Equal(t, parcel.Address, parcelNewStatus.Address)
	assert.Equal(t, parcel.CreatedAt, parcelNewStatus.CreatedAt)
	assert.Equal(t, parcel.Number, parcelNewStatus.Number)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		assert.NotEmpty(t, id)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, len(storedParcels), len(parcelMap))
	for _, parcel := range storedParcels {
		assert.Equal(t, parcel, parcelMap[parcel.Number])
		assert.Equal(t, parcel.Address, parcelMap[parcel.Number].Address)
		assert.Equal(t, parcel.Client, parcelMap[parcel.Number].Client)
		assert.Equal(t, parcel.Status, parcelMap[parcel.Number].Status)
		assert.Equal(t, parcel.CreatedAt, parcelMap[parcel.Number].CreatedAt)
		assert.Equal(t, parcel.Number, parcelMap[parcel.Number].Number)

	}
}
