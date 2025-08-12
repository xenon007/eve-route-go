package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/tkhamez/eve-route-go/internal/graph"
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

// Graph loads the universe graph from MongoDB.
func (m *Mongo) Graph(ctx context.Context) (graph.Graph, error) {
	var g graph.Graph
	err := m.collection("graph").FindOne(ctx, bson.D{}).Decode(&g)
	if err == mongo.ErrNoDocuments {
		return graph.Graph{}, nil
	}
	return g, err
}

// SaveGraph stores the universe graph in MongoDB.
func (m *Mongo) SaveGraph(ctx context.Context, g graph.Graph) error {
	_, err := m.collection("graph").ReplaceOne(ctx, bson.D{}, g, options.Replace().SetUpsert(true))
	return err
}

// EnsureMongoConnection pings the database to check connection.
func (m *Mongo) EnsureMongoConnection(ctx context.Context) {
	if err := m.client.Ping(ctx, nil); err != nil {
		log.Printf("mongo ping error: %v", err)
	}
}
