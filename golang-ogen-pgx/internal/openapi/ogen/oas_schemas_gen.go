// Code generated by ogen, DO NOT EDIT.

package ogen

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *ErrorResponseStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

// Ref: #/components/schemas/CreateItemRequest
type CreateItemRequest struct {
	Data ItemIn `json:"data"`
}

// GetData returns the value of Data.
func (s *CreateItemRequest) GetData() ItemIn {
	return s.Data
}

// SetData sets the value of Data.
func (s *CreateItemRequest) SetData(val ItemIn) {
	s.Data = val
}

// Ref: #/components/schemas/CreateItemResponse
type CreateItemResponse struct {
	Data Item                   `json:"data"`
	Meta CreateItemResponseMeta `json:"meta"`
}

// GetData returns the value of Data.
func (s *CreateItemResponse) GetData() Item {
	return s.Data
}

// GetMeta returns the value of Meta.
func (s *CreateItemResponse) GetMeta() CreateItemResponseMeta {
	return s.Meta
}

// SetData sets the value of Data.
func (s *CreateItemResponse) SetData(val Item) {
	s.Data = val
}

// SetMeta sets the value of Meta.
func (s *CreateItemResponse) SetMeta(val CreateItemResponseMeta) {
	s.Meta = val
}

// Ref: #/components/schemas/CreateItemResponseMeta
type CreateItemResponseMeta struct {
	Created bool `json:"created"`
}

// GetCreated returns the value of Created.
func (s *CreateItemResponseMeta) GetCreated() bool {
	return s.Created
}

// SetCreated sets the value of Created.
func (s *CreateItemResponseMeta) SetCreated(val bool) {
	s.Created = val
}

// Ref: #/components/schemas/ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetError returns the value of Error.
func (s *ErrorResponse) GetError() string {
	return s.Error
}

// SetError sets the value of Error.
func (s *ErrorResponse) SetError(val string) {
	s.Error = val
}

// ErrorResponseStatusCode wraps ErrorResponse with StatusCode.
type ErrorResponseStatusCode struct {
	StatusCode int
	Response   ErrorResponse
}

// GetStatusCode returns the value of StatusCode.
func (s *ErrorResponseStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *ErrorResponseStatusCode) GetResponse() ErrorResponse {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *ErrorResponseStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *ErrorResponseStatusCode) SetResponse(val ErrorResponse) {
	s.Response = val
}

// Ref: #/components/schemas/GetItemResponse
type GetItemResponse struct {
	Data Item                `json:"data"`
	Meta GetItemResponseMeta `json:"meta"`
}

// GetData returns the value of Data.
func (s *GetItemResponse) GetData() Item {
	return s.Data
}

// GetMeta returns the value of Meta.
func (s *GetItemResponse) GetMeta() GetItemResponseMeta {
	return s.Meta
}

// SetData sets the value of Data.
func (s *GetItemResponse) SetData(val Item) {
	s.Data = val
}

// SetMeta sets the value of Meta.
func (s *GetItemResponse) SetMeta(val GetItemResponseMeta) {
	s.Meta = val
}

type GetItemResponseMeta struct{}

// Ref: #/components/schemas/GetItemsResponse
type GetItemsResponse struct {
	Data []Item               `json:"data"`
	Meta GetItemsResponseMeta `json:"meta"`
}

// GetData returns the value of Data.
func (s *GetItemsResponse) GetData() []Item {
	return s.Data
}

// GetMeta returns the value of Meta.
func (s *GetItemsResponse) GetMeta() GetItemsResponseMeta {
	return s.Meta
}

// SetData sets the value of Data.
func (s *GetItemsResponse) SetData(val []Item) {
	s.Data = val
}

// SetMeta sets the value of Meta.
func (s *GetItemsResponse) SetMeta(val GetItemsResponseMeta) {
	s.Meta = val
}

type GetItemsResponseMeta struct{}

// Ref: #/components/schemas/Item
type Item struct {
	ID        int64     `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Price     float32   `json:"price"`
}

// GetID returns the value of ID.
func (s *Item) GetID() int64 {
	return s.ID
}

// GetUUID returns the value of UUID.
func (s *Item) GetUUID() uuid.UUID {
	return s.UUID
}

// GetCreatedAt returns the value of CreatedAt.
func (s *Item) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetName returns the value of Name.
func (s *Item) GetName() string {
	return s.Name
}

// GetPrice returns the value of Price.
func (s *Item) GetPrice() float32 {
	return s.Price
}

// SetID sets the value of ID.
func (s *Item) SetID(val int64) {
	s.ID = val
}

// SetUUID sets the value of UUID.
func (s *Item) SetUUID(val uuid.UUID) {
	s.UUID = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *Item) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

// SetName sets the value of Name.
func (s *Item) SetName(val string) {
	s.Name = val
}

// SetPrice sets the value of Price.
func (s *Item) SetPrice(val float32) {
	s.Price = val
}

// Ref: #/components/schemas/ItemIn
type ItemIn struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

// GetName returns the value of Name.
func (s *ItemIn) GetName() string {
	return s.Name
}

// GetPrice returns the value of Price.
func (s *ItemIn) GetPrice() float32 {
	return s.Price
}

// SetName sets the value of Name.
func (s *ItemIn) SetName(val string) {
	s.Name = val
}

// SetPrice sets the value of Price.
func (s *ItemIn) SetPrice(val float32) {
	s.Price = val
}

type PingGetOK struct {
	Message string `json:"message"`
}

// GetMessage returns the value of Message.
func (s *PingGetOK) GetMessage() string {
	return s.Message
}

// SetMessage sets the value of Message.
func (s *PingGetOK) SetMessage(val string) {
	s.Message = val
}
