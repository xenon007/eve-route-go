package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo implements Store using MongoDB.
type Mongo struct {
	client *mongo.Client
	dbName string
}

// NewMongo creates a new Mongo store.
func NewMongo(client *mongo.Client, dbName string) *Mongo {
	return &Mongo{client: client, dbName: dbName}
}

func (m *Mongo) collection(name string) *mongo.Collection {
	return m.client.Database(m.dbName).Collection(name)
}

// Ansiblexes loads Ansiblex gates from MongoDB.
func (m *Mongo) Ansiblexes(ctx context.Context) ([]Ansiblex, error) {
	cur, err := m.collection("ansiblex").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var res []Ansiblex
	for cur.Next(ctx) {
		var a Ansiblex
		if err := cur.Decode(&a); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, cur.Err()
}

// TemporaryConnections loads temporary connections from MongoDB.
func (m *Mongo) TemporaryConnections(ctx context.Context) ([]TemporaryConnection, error) {
	cur, err := m.collection("temporary_connections").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var res []TemporaryConnection
	for cur.Next(ctx) {
		var c TemporaryConnection
		if err := cur.Decode(&c); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, cur.Err()
}

// Systems loads capital systems from MongoDB.
func (m *Mongo) Systems(ctx context.Context) (map[int]System, error) {
	cur, err := m.collection("systems").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	systems := make(map[int]System)
	for cur.Next(ctx) {
		var s System
		if err := cur.Decode(&s); err != nil {
			return nil, err
		}
		systems[s.ID] = s
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return systems, nil
}

// CreateAnsiblex inserts a new Ansiblex gate.
func (m *Mongo) CreateAnsiblex(ctx context.Context, a Ansiblex) (int64, error) {
	if a.ID == 0 {
		a.ID = time.Now().UnixNano()
	}
	_, err := m.collection("ansiblex").InsertOne(ctx, a)
	return a.ID, err
}

// UpdateAnsiblex updates an existing Ansiblex gate.
func (m *Mongo) UpdateAnsiblex(ctx context.Context, a Ansiblex) error {
	_, err := m.collection("ansiblex").UpdateOne(ctx, bson.M{"id": a.ID}, bson.M{"$set": a})
	return err
}

// DeleteAnsiblex removes an Ansiblex gate.
func (m *Mongo) DeleteAnsiblex(ctx context.Context, id int64) error {
	_, err := m.collection("ansiblex").DeleteOne(ctx, bson.M{"id": id})
	return err
}

// CreateTemporaryConnection inserts a new temporary connection.
func (m *Mongo) CreateTemporaryConnection(ctx context.Context, c TemporaryConnection) (int64, error) {
	if c.ID == 0 {
		c.ID = time.Now().UnixNano()
	}
	_, err := m.collection("temporary_connections").InsertOne(ctx, c)
	return c.ID, err
}

// UpdateTemporaryConnection updates an existing temporary connection.
func (m *Mongo) UpdateTemporaryConnection(ctx context.Context, c TemporaryConnection) error {
	_, err := m.collection("temporary_connections").UpdateOne(ctx, bson.M{"id": c.ID}, bson.M{"$set": c})
	return err
}

// DeleteTemporaryConnection removes a temporary connection.
func (m *Mongo) DeleteTemporaryConnection(ctx context.Context, id int64) error {
	_, err := m.collection("temporary_connections").DeleteOne(ctx, bson.M{"id": id})
	return err
}

// EnsureMongoConnection pings the database to check connection.
func (m *Mongo) EnsureMongoConnection(ctx context.Context) {
	if err := m.client.Ping(ctx, nil); err != nil {
		log.Printf("mongo ping error: %v", err)
	}
}
