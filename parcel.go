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
	response, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		fmt.Printf("Add s.db.Exec error %v \n", err)
		return 0, err
	}
	number, err := response.LastInsertId()
	if err != nil {
		fmt.Printf("Add response.LastInsertId error %v \n", err)
		return 0, err
	}
	return int(number), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		fmt.Printf("Get row.Scan error %v \n", err)
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		fmt.Printf("GetByClient s.db.Query error %v \n", err)
		return nil, err
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Printf("GetByClient rows.Scan err %v \n", err)
			return nil, err
		}
		res = append(res, p)
	}
	if err = rows.Err(); err != nil {
		fmt.Printf("GetByClient rows.Next error %v \n", err)
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		fmt.Printf("SetStatus s.db.Exec error %v \n", err)
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", "registered"))
	if err != nil {
		fmt.Printf("SetAddress s.db.Exec error %v \n", err)
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE status = :status AND number = :number",
		sql.Named("status", "registered"),
		sql.Named("number", number))
	if err != nil {
		fmt.Printf("Delete s.db.Exec error %v \n", err)
		return err
	}
	return nil
}
