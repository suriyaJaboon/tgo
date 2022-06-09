package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"tgo/store"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func Test_lc_LC(t *testing.T) {
	var f = fiber.New()

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	var requestTesting = func(t *testing.T, mt *mtest.T, statusCode int) {
		NewLC(f, store.NewLC(mt.Coll))

		res, err := f.Test(httptest.NewRequest(http.MethodGet, "/lc", nil))
		if err != nil {
			t.Errorf("http request service lc shoud be error %v", err)
		}

		defer func() {
			if err = res.Body.Close(); err != nil {
				t.Skip("http request skip close body", err)
			}
		}()

		if res.StatusCode != statusCode {
			t.Errorf("http request service lc expected status code %v result: %v", statusCode, res.StatusCode)
		}

		bt, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("http respose read all body error: %v", err)
		}

		if res.StatusCode == 200 {
			var lgc []*store.Lcg
			if err = json.Unmarshal(bt, &lgc); err != nil {
				t.Errorf("http response unmarshal body error: %v", err)
			}
			if len(lgc) != 2 {
				t.Errorf("http response expected 2 slices  but got %v", len(lgc))
			}
		}
		if res.StatusCode == 500 {
			var er *fiber.Error
			if err = json.Unmarshal(bt, &er); err != nil {
				t.Errorf("http response unmarshal body error: %v", err)
			}
			fmt.Println(string(bt), res.StatusCode, er)
		}
	}

	mt.Run("http response success", func(mt *mtest.T) {
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "lc.lc", mtest.FirstBatch, bson.D{
			{"_id", id1},
			{"name", "lcx1"},
			{"age", 1},
			{"created_at", time.Now()},
		})
		second := mtest.CreateCursorResponse(1, "lc.lc", mtest.NextBatch, bson.D{
			{"_id", id2},
			{"name", "lcx2"},
			{"age", 2},
			{"created_at", time.Now()},
		})
		killCursors := mtest.CreateCursorResponse(0, "lc.lc", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		requestTesting(t, mt, http.StatusOK)
	})
	mt.Run("http response error", func(mt *mtest.T) {
		id1 := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "lc.lc", mtest.FirstBatch, bson.D{
			{"_ids", id1},
			{"names", "lcx1"},
			{"ages", 1},
			{"created_ats", time.Now()},
		})
		killCursors := mtest.CreateCursorResponse(0, "lc.lc", mtest.NextBatch)
		mt.AddMockResponses(first, killCursors)

		requestTesting(t, mt, http.StatusInternalServerError)
	})
}
