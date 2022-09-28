package User

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID          *primitive.ObjectID `json:"-"   bson:"_id,omitempty"`
	FirstName   string              `json:"first_name"   bson:"first_name"`
	LastName    string              `json:"last_name"   bson:"last_name"`
	Username    string              `json:"username"   bson:"username"`
	Password    string              `json:"password"   bson:"password"`
	Admin       bool                `json:"admin"   bson:"admin"`
	_isDocument bool
}

func (m *Model) GetCollectionName() string {
	//TODO implement me
	return "users"
}

func (m *Model) SetIsDocumented(b bool) {
	//TODO implement me
	m._isDocument = b
}

func (m *Model) GetIsDocumented() bool {
	//TODO implement me
	return m._isDocument
}

func (m *Model) GetID() interface{} {
	return m.ID
}

func (m *Model) SetID(id interface{}) {
	i := id.(primitive.ObjectID)
	m.ID = &i
}
