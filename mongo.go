package backendgis

import (
	"context"
	"os"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetConnectionMongo(MongoString, dbname string) *mongo.Database {
	MongoInfo := atdb.DBInfo{
		DBString: os.Getenv(MongoString),
		DBName:   dbname,
	}
	conn := atdb.MongoConnect(MongoInfo)
	return conn
}

func GetAllData(MongoConnect *mongo.Database, colname string) []GeoJson {
	data := atdb.GetAllDoc[[]GeoJson](MongoConnect, colname)
	return data
}

func InsertDataGeojson(MongoConn *mongo.Database, colname string, coordinate [][]float64, name, volume, tipe string) (InsertedID interface{}) {
	req := new(LonLatProperties)
	req.Type = tipe
	req.Coordinates = coordinate
	req.Name = name
	req.Volume = volume

	ins := atdb.InsertOneDoc(MongoConn, colname, req)
	return ins
}

func UpdateDataGeojson(MongoConn *mongo.Database, colname, name, newVolume, newTipe string) error {
	// Filter berdasarkan nama
	filter := bson.M{"name": name}

	// Update data yang akan diubah
	update := bson.M{
		"$set": bson.M{
			"volume": newVolume,
			"tipe":   newTipe,
		},
	}

	// Mencoba untuk mengupdate dokumen
	_, err := MongoConn.Collection(colname).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func DeleteDataGeojson(Mongoenv, dbname string, ctx context.Context, val LonLatProperties) (DeletedId interface{}) {
	conn := GetConnectionMongo(Mongoenv, dbname)
	filter := bson.D{{"volume", val.Volume}}
	res, err := conn.Collection("lonlatpost").DeleteOne(ctx, filter)
	if err != nil {
		return "Gagal Delete"
	}
	return res
}

func IsExist(Tokenstr, PublicKey string) bool {
	id := watoken.DecodeGetId(PublicKey, Tokenstr)
	if id == "" {
		return false
	}
	return true
}

func GetCoordinateNear(MongoConn *mongo.Database, colname string, coordinate []float64) (result []GeoJson, err error) {
	filter := bson.M{"geometry.coordinates": bson.M{
		"$near": bson.M{
			"$geometry": bson.M{
				"type":        "LineString",
				"coordinates": coordinate,
			},
		},
	}}
	curr, err := MongoConn.Collection(colname).Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer curr.Close(context.Background())
	err = curr.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
