package main
import(
"os"
"database/sql"
"net/http/httptest"
"testing"
"encoding/json"
"net/http"
_ "github.com/go-sql-driver/mysql"
)
func openTestDB(t *testing.T)*sql.DB{
	dsn := os.Getenv("RECORD_STORE_DSN")
	if dsn==""{
		t.Fatalf("RECORD_STORE_DSN")
	}
	db,err := sql.Open("mysql",dsn)
	if err != nil{
		 t.Fatalf("failed to open database: %v", err)
	}
	err = db.Ping()
	if err != nil{
		db.Close()
		t.Fatalf("failed to ping database: %v", err)
	}
	return db
}
func TestRecordByIDsucceed(t *testing.T){
	db := openTestDB(t)
	defer db.Close()
	mux := registerRoutes(db)
	request := httptest.NewRequest("GET","/records/1",nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder,request)
	response:=recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK{
		t.Fatalf("expected:statusok,get:%v",response.StatusCode)
	}
	var record Record
	err :=json.NewDecoder(response.Body).Decode(&record)
	if err != nil{
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if record.ID!=1 || record.Title!="Abbey Road" || record.Artist!= "The Beatles"{
		t.Fatalf("expect id:1 got:%d,expect title:Abbey Road,get :%s,expect Artist:The Beatles,got:%s",record.ID,record.Title,record.Artist)
	}
}
func TestRecordByIDInvalidID(t *testing.T) {
    db := openTestDB(t)
	defer db.Close()
	mux := registerRoutes(db)
	request:=httptest.NewRequest("GET","/records/abc",nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder,request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest{
		t.Fatalf("expected:StatusBadRequest,get:%v",response.StatusCode)
	}
}

func TestRecordByIDNotFound(t *testing.T) {
    db := openTestDB(t)
	defer db.Close()
	mux:=registerRoutes(db)
	request := httptest.NewRequest("GET","/records/999999",nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder,request)
	response := recorder.Result()
	defer response.Body.Close()
	if response.StatusCode!=http.StatusNotFound{
		t.Fatalf("expect:StatusNotFound,get:%v",response.StatusCode)
	}
}
func TestRecordByIDDatabaseError(t *testing.T) {
    db := openTestDB(t)

	mux:=registerRoutes(db)
	db.Close()
	request := httptest.NewRequest("GET","/records/1",nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder,request)
	response:=recorder.Result()
	defer response.Body.Close()
	if response.StatusCode != http.StatusInternalServerError{
		t.Fatalf("expect:internalservererror,get:%d",response.StatusCode)
	}
	if response.Header.Get("Content-Type")!="application/json"{
		t.Fatalf("expected application/json, got %s",
    response.Header.Get("Content-Type"))
	}
}