/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package model

import "github.com/west2-online/DomTok/kitex_gen/model"

type PaymentTokenRequest struct {
	OrderID int64 `thrift:"orderID,1,required" frugal:"1,required,i64" json:"orderID"`
	UserID  int64 `thrift:"userID,2,required" frugal:"2,required,i64" json:"userID"`
}

type PaymentTokenResponse struct {
	Base           *model.BaseResp `thrift:"base,1" frugal:"1,default,model.BaseResp" json:"base"`
	PaymentToken   string          `thrift:"paymentToken,2,required" frugal:"2,required,string" json:"paymentToken"`
	ExpirationTime int64           `thrift:"expirationTime,3,required" frugal:"3,required,i64" json:"expirationTime"`
}

type PaymentRequest struct {
	OrderID      int64                 `thrift:"orderID,1,required" frugal:"1,required,i64" json:"orderID"`
	UserID       int64                 `thrift:"userID,2,required" frugal:"2,required,i64" json:"userID"`
	PaymentToken string                `thrift:"paymentToken,3,required" frugal:"3,required,string" json:"paymentToken"`
	CreditCard   *model.CreditCardInfo `thrift:"creditCard,4,required" frugal:"4,required,model.CreditCardInfo" json:"creditCard"`
	Description  *string               `thrift:"description,5,optional" frugal:"5,optional,string" json:"description,omitempty"`
}

type PaymentResponse struct {
	Base      *model.BaseResp `thrift:"base,1" frugal:"1,default,model.BaseResp" json:"base"`
	PaymentID int64           `thrift:"paymentID,2,required" frugal:"2,required,i64" json:"paymentID"`
	Status    int64           `thrift:"status,3,required" frugal:"3,required,i64" json:"status"`
}

type RefundTokenRequest struct {
	OrderID int64 `thrift:"orderID,1,required" frugal:"1,required,i64" json:"orderID"`
	UserID  int64 `thrift:"userID,2,required" frugal:"2,required,i64" json:"userID"`
}

type RefundTokenResponse struct {
	Base           *model.BaseResp `thrift:"base,1" frugal:"1,default,model.BaseResp" json:"base"`
	RefundToken    string          `thrift:"refundToken,2,required" frugal:"2,required,string" json:"refundToken"`
	ExpirationTime int64           `thrift:"expirationTime,3,required" frugal:"3,required,i64" json:"expirationTime"`
}
type RefundRequest struct {
	OrderID      int64   `thrift:"orderID,1,required" frugal:"1,required,i64" json:"orderID"`
	UserID       int64   `thrift:"userID,2,required" frugal:"2,required,i64" json:"userID"`
	RefundAmount float64 `thrift:"refundAmount,3,required" frugal:"3,required,double" json:"refundAmount"`
	RefundReason string  `thrift:"refundReason,4,required" frugal:"4,required,string" json:"refundReason"`
}

type RefundResponse struct {
	Base     *model.BaseResp `thrift:"base,1" frugal:"1,default,model.BaseResp" json:"base"`
	RefundID int64           `thrift:"refundID,2,required" frugal:"2,required,i64" json:"refundID"`
	Status   int64           `thrift:"status,3,required" frugal:"3,required,i64" json:"status"`
}

type Payment struct {
	Base      *model.BaseResp // 这里要导入哪一个包？
	PaymentID int64
	Status    int64
}
