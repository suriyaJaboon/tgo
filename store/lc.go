package store

//go:generate mockgen -source ./lc.go LC -destination ./mocks/lc_mock.go

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LC interface {
	LCX() ([]*Lcg, error)
	LCXByID(id primitive.ObjectID) (*Lcg, error)
	CreateLCX(dto *LcgDto) (*Lcg, error)
}

type lc struct {
	c   *mongo.Collection
	ctx context.Context
}

type (
	LcgDto struct {
		Name string `bson:"name" json:"name"`
		Age  int    `bson:"age" json:"age"`
	}

	Lcg struct {
		ID        primitive.ObjectID `bson:"_id" json:"id"`
		Name      string             `bson:"name" json:"name"`
		Age       int                `bson:"age" json:"age"`
		CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	}
)

func NewLC(mlc *mongo.Collection) LC {
	return &lc{c: mlc, ctx: context.Background()}
}

func (l lc) LCX() ([]*Lcg, error) {
	cur, err := l.c.Find(l.ctx, bson.D{}, &options.FindOptions{Sort: bson.D{{Key: "created_at", Value: -1}}})
	if err != nil {
		return nil, err
	}

	defer func() {
		err = cur.Close(l.ctx)
	}()

	var lgs []*Lcg
	if err = cur.All(l.ctx, &lgs); err != nil {
		return nil, err
	}

	return lgs, nil
}

func (l lc) LCXByID(id primitive.ObjectID) (*Lcg, error) {
	var lcg Lcg
	err := l.c.FindOne(l.ctx, bson.M{"_id": id}).Decode(&lcg)
	if err != nil {
		return nil, err
	}

	return &lcg, nil
}

func (l lc) CreateLCX(dto *LcgDto) (*Lcg, error) {
	var lcg = &Lcg{
		ID:        primitive.NewObjectID(),
		Name:      dto.Name,
		Age:       dto.Age,
		CreatedAt: time.Now(),
	}

	_, err := l.c.InsertOne(l.ctx, lcg)
	if err != nil {
		return nil, err
	}

	return lcg, nil
}
