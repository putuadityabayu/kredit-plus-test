/*
 * Copyright (c) - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Putu Aditya <aditya@portalnesia.com>
 */

package model

import "time"

const (
	TrxPENDING  = "pending"
	TrxAPPROVED = "approved"
	TrxREJECTED = "rejected"
)

type Transaction struct {
	ID                string    `gorm:"column:id;type:uuid;primarykey" json:"id"`
	ContractNumber    string    `gorm:"column:contract_number;type:varchar(255);not null;unique" json:"contract_number"`
	UserID            string    `gorm:"column:user_id;type:uuid;not null" json:"user_id"`
	OTR               float64   `gorm:"column:otr;type:decimal(10,2);not null" json:"otr"`
	AdminFee          float64   `gorm:"column:admin_fee;type:decimal(10,2);not null" json:"admin_fee"`
	InstallmentAmount float64   `gorm:"column:installment_amount;type:decimal(10,2);not null" json:"installment_amount"`
	InterestAmount    float64   `gorm:"column:interest_amount;type:decimal(10,2);not null" json:"interest_amount"`
	AssetName         string    `gorm:"column:asset_name;type:varchar(255);not null" json:"asset_name"`
	Tenor             int       `gorm:"column:tenor;type:int;not null" json:"tenor"`
	TransactionDate   time.Time `gorm:"column:transaction_date;type:timestamp;not null" json:"transaction_date"`
	Status            string    `gorm:"column:status;type:enum('pending', 'approved', 'rejected');not null" json:"status"`
	CreatedAt         time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;autoUpdateTime"`
}

func (Transaction) TableName() string {
	return "transactions"
}
