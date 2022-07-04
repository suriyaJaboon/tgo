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
	"tgo/store/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func TestLc_LC(t *testing.T) {
	t.Run("http status ok", func(t *testing.T) {
		var f = fiber.New()

		ctrl := gomock.NewController(t)

		mlc := mocks.NewMockLC(ctrl)
		mlc.EXPECT().LCX().Return([]*store.Lcg{
			{
				ID:        primitive.NewObjectID(),
				Name:      "001",
				Age:       1,
				CreatedAt: time.Now(),
			},
			{
				ID:        primitive.NewObjectID(),
				Name:      "002",
				Age:       2,
				CreatedAt: time.Now(),
			},
		}, nil)

		NewLC(f, mlc)
		res, err := f.Test(httptest.NewRequest(http.MethodGet, "/lc", nil))
		if err != nil {
			t.Errorf("http request service lc shoud be error %v", err)
		}

		defer func() {
			if err = res.Body.Close(); err != nil {
				t.Skip("http request skip close body", err)
			}
		}()

		bt, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("http respose read all body error: %v", err)
		}

		var ls []store.Lcg
		if err = json.Unmarshal(bt, &ls); err != nil {
			assert.Errorf(t, err, "http respose read all body error: %v", err)
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, 2, len(ls))
	})

	t.Run("http internal server error", func(t *testing.T) {
		var f = fiber.New()

		ctrl := gomock.NewController(t)

		mlc := mocks.NewMockLC(ctrl)
		mlc.EXPECT().LCX().Return([]*store.Lcg{}, mongo.ErrNoDocuments)

		NewLC(f, mlc)
		res, err := f.Test(httptest.NewRequest(http.MethodGet, "/lc", nil))
		if err != nil {
			t.Errorf("http request service lc shoud be error %v", err)
		}

		defer func() {
			if err = res.Body.Close(); err != nil {
				t.Skip("http request skip close body", err)
			}
		}()
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func TestLc_LCByID(t *testing.T) {
	var id = primitive.NewObjectID()

	t.Run("http status ok", func(t *testing.T) {
		var f = fiber.New()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mlc := mocks.NewMockLC(ctrl)
		mlc.EXPECT().LCXByID(gomock.Any()).Return(&store.Lcg{
			ID:        id,
			Name:      "001",
			Age:       1,
			CreatedAt: time.Now(),
		}, nil)

		NewLC(f, mlc)
		res, err := f.Test(httptest.NewRequest(http.MethodGet, "/lc/"+id.Hex(), nil))
		if err != nil {
			t.Errorf("http request service lc shoud be error %v", err)
		}

		defer func() {
			if err = res.Body.Close(); err != nil {
				t.Skip("http request skip close body", err)
			}
		}()

		bt, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("http respose read all body error: %v", err)
		}

		var ls store.Lcg
		if err = json.Unmarshal(bt, &ls); err != nil {
			assert.Errorf(t, err, "http respose read all body error: %v", err)
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, id.Hex(), ls.ID.Hex())
	})

	t.Run("http invalid not accept table", func(t *testing.T) {
		var f = fiber.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mlc := mocks.NewMockLC(ctrl)

		NewLC(f, mlc)
		res, err := f.Test(httptest.NewRequest(http.MethodGet, "/lc/id", nil))
		if err != nil {
			t.Errorf("http request service lc shoud be error %v", err)
		}

		defer func() {
			if err = res.Body.Close(); err != nil {
				t.Skip("http request skip close body", err)
			}
		}()

		assert.Equal(t, http.StatusNotAcceptable, res.StatusCode)
	})

	t.Run("http internal server error", func(t *testing.T) {
		var f = fiber.New()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mlc := mocks.NewMockLC(ctrl)
		mlc.EXPECT().LCXByID(gomock.Any()).Return(nil, mongo.ErrNoDocuments)

		NewLC(f, mlc)
		res, err := f.Test(httptest.NewRequest(http.MethodGet, "/lc/"+id.Hex(), nil))
		if err != nil {
			t.Errorf("http request service lc shoud be error %v", err)
		}

		defer func() {
			if err = res.Body.Close(); err != nil {
				t.Skip("http request skip close body", err)
			}
		}()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
