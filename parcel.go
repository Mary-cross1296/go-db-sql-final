package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	response, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		fmt.Printf("Add s.db.Exec error %v \n", err)
		return 0, err
	}
	lastId, err := response.LastInsertId()
	if err != nil {
		fmt.Printf("Add response.LastInsertId error %v \n", err)
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	return int(lastId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT client, status, address, created_at FROM parcel WHERE id = :id",
		sql.Named("id", number))
	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		fmt.Printf("Get row.Scan error %v \n", err)
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT client, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		fmt.Printf("GetByClient s.db.Query error %v \n", err)
		return nil, err
	}
	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Printf("GetByClient rows.Scan err %v \n", err)
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE id := number",
		sql.Named("status", status),
		sql.Named("id", number))
	if err != nil {
		fmt.Printf("SetStatus s.db.Exec error %v \n", err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE id := number",
		sql.Named("address", address),
		sql.Named("id", number))
	if err != nil {
		fmt.Printf("SetAddress s.db.Exec error %v \n", err)
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE status = :status AND id = :number",
		sql.Named("status", "registered"),
		sql.Named("id", number))
	if err != nil {
		fmt.Printf("Delete s.db.Exec error %v \n", err)
		return err
	}
	return nil
}
