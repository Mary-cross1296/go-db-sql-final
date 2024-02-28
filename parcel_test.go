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
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
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
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	number, err := store.Add(parcel)
	require.NoError(t, err)
	assert.True(t, number > 0)

	parcel.Number = number

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	row, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, parcel, row)

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	err = store.Delete(number)
	require.NoError(t, err)

	// проверьте, что посылку больше нельзя получить из БД
	_, err = store.Get(number)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	number, err := store.Add(parcel)
	require.NoError(t, err)
	assert.True(t, number > 0)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	row, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, newAddress, row.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	number, err := store.Add(parcel)
	require.NoError(t, err)
	assert.True(t, number > 0)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	newAStatus := ParcelStatusDelivered
	err = store.SetStatus(number, newAStatus)
	require.NoError(t, err)

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	row, err := store.Get(number)
	require.NoError(t, err)
	assert.Equal(t, newAStatus, row.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// Создаем новый массив с объектами parcel
	parcels := []Parcel{
		parcel,
		parcel,
		parcel,
	}
	// Создаем мапу ключ-число : значение-объект parcel
	parcelMap := map[int]Parcel{}

	// задаём двум из трех посылок одинаковый идентификатор клиента, чтобы в дальнейшем проверить правильность выборки, что не вернулось ничего лишнего
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		number, err := store.Add(parcels[i]) // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		require.NoError(t, err)
		assert.True(t, number > 0)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = number

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[number] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	require.NoError(t, err)
	// убедитесь, что количество полученных посылок совпадает с количеством тех, что присвоены клиенту с идентификатором, сохранённым в переменной client
	assert.Equal(t, 2, len(storedParcels))
	// check

	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		require.Contains(t, parcelMap, parcel.Number)
		// убедитесь, что значения полей полученных посылок заполнены верно
		parcelMapValue := parcelMap[parcel.Number]

		assert.Equal(t, parcelMapValue.Client, parcel.Client)
		assert.Equal(t, parcelMapValue.Status, parcel.Status)
		assert.Equal(t, parcelMapValue.Address, parcel.Address)
		assert.Equal(t, parcelMapValue.CreatedAt, parcel.CreatedAt)
	}
}
