package backendgis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	model "github.com/angkringankuy/backendyak"
	"github.com/petapedia/peda"
	"github.com/whatsauth/watoken"
)

func GCHandlerFunc(publickey, Mongostring, dbname, colname string, r *http.Request) string {
	resp := new(model.Credential)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		resp.Status = false
		resp.Message = "Header Login Not Exist"
	} else {
		existing := IsExist(tokenlogin, os.Getenv(publickey))
		if !existing {
			resp.Status = false
			resp.Message = "Kamu kayaknya belum punya akun"
		} else {
			koneksyen := GetConnectionMongo(Mongostring, dbname)
			datageo := GetAllData(koneksyen, colname)
			jsoncihuy, _ := json.Marshal(datageo)
			resp.Status = true
			resp.Message = "Data Berhasil diambil"
			resp.Token = string(jsoncihuy)
		}
	}
	return ReturnStringStruct(resp)
}

func GCFPostCoordinate(Mongostring, Publickey, dbname, colname string, r *http.Request) string {
	req := new(Credents)
	conn := GetConnectionMongo(Mongostring, dbname)
	resp := new(LonLatProperties)
	err := json.NewDecoder(r.Body).Decode(&resp)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = strconv.Itoa(http.StatusNotFound)
		req.Message = "Header Login Not Exist"
	} else {
		existing := IsExist(tokenlogin, os.Getenv(Publickey))
		if !existing {
			req.Status = strconv.Itoa(http.StatusNotFound)
			req.Message = "Kamu kayaknya belum punya akun"
		} else {
			if err != nil {
				req.Status = strconv.Itoa(http.StatusNotFound)
				req.Message = "error parsing application/json: " + err.Error()
			} else {
				req.Status = strconv.Itoa(http.StatusOK)
				Ins := InsertDataGeojson(conn, colname,
					resp.Coordinates,
					resp.Name,
					resp.Volume,
					resp.Type)
				req.Message = fmt.Sprintf("%v:%v", "Berhasil Input data", Ins)
			}
		}
	}
	return ReturnStringStruct(req)
}

func ReturnStringStruct(Data any) string {
	jsonee, _ := json.Marshal(Data)
	return string(jsonee)
}

func GCFUpdateNameGeojson(publickey, Mongostring, dbname, colname string, r *http.Request) string {
	req := new(Credents)
	resp := new(LonLatProperties)
	conn := GetConnectionMongo(Mongostring, dbname)
	err := json.NewDecoder(r.Body).Decode(&resp)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = strconv.Itoa(http.StatusNotFound)
		req.Message = "Header Login Not Exist"
	} else {
		existing := IsExist(tokenlogin, os.Getenv(publickey))
		if !existing {
			req.Status = strconv.Itoa(http.StatusNotFound)
			req.Message = "Kamu kayaknya belum punya akun"
		} else {
			if err != nil {
				req.Status = strconv.Itoa(http.StatusNotFound)
				req.Message = "error parsing application/json: " + err.Error()
			} else {
				req.Status = strconv.Itoa(http.StatusOK)
				Ins := UpdateDataGeojson(conn, colname,
					resp.Name,
					resp.Volume,
					resp.Type)
				req.Message = fmt.Sprintf("%v:%v", "Berhasil Update data", Ins)
			}
		}

	}
	return ReturnStringStruct(req)
}

func GCFDeleteDataGeojson(publickey, Mongostring, dbname, colname string, r *http.Request) string {
	req := new(Credents)
	resp := new(LonLatProperties)
	err := json.NewDecoder(r.Body).Decode(&resp)
	tokenlogin := r.Header.Get("Login")
	if tokenlogin == "" {
		req.Status = strconv.Itoa(http.StatusNotFound)
		req.Message = "Header Login Not Exist"
	} else {
		existing := IsExist(tokenlogin, os.Getenv(publickey))
		if !existing {
			req.Status = strconv.Itoa(http.StatusNotFound)
			req.Message = "Kamu kayaknya belum punya akun"
		} else {
			if err != nil {
				req.Status = strconv.Itoa(http.StatusNotFound)
				req.Message = "error parsing application/json: " + err.Error()
			} else {
				req.Status = strconv.Itoa(http.StatusOK)
				Ins := DeleteDataGeojson(Mongostring, dbname, context.Background(),
					LonLatProperties{
						Type:   resp.Type,
						Name:   resp.Name,
						Volume: resp.Volume,
					})
				req.Message = fmt.Sprintf("%v:%v", "Berhasil Hapus data", Ins)
			}
		}
	}
	return ReturnStringStruct(req)
}

func Login(Privatekey, MongoEnv, dbname, Colname string, r *http.Request) string {
	var resp model.Credential
	mconn := model.SetConnection(MongoEnv, dbname)
	var datauser peda.User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		if peda.IsPasswordValid(mconn, Colname, datauser) {
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(Privatekey))
			if err != nil {
				resp.Message = "Gagal Encode Token : " + err.Error()
			} else {
				resp.Status = true
				resp.Message = "Selamat Datang"
				resp.Token = tokenstring
			}
		} else {
			resp.Message = "Password Salah"
		}
	}
	return ReturnStringStruct(resp)
}
