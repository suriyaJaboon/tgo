package store

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

//var noClientOpts = mtest.NewOptions().CreateClient(false)

func TestNewLC(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tests := []struct {
		name string
		want *lc
	}{
		// TODO: Add test cases.
		{
			name: "New create lc",
			want: &lc{c: mt.Coll, ctx: context.Background()},
		},
		{
			name: "New create lc set of nil",
			want: &lc{ctx: context.Background()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLC(mt.Coll); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLc_CreateLCX(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	//var noClientOpts = mtest.NewOptions().CreateClient(false)

	mt.RunOpts("insert one", mtest.NewOptions().ClientType(mtest.Mock), func(mt *mtest.T) {
		// DatabaseName
		mt.Run("success", func(mt *mtest.T) {
			// CollectionName
			id := primitive.NewObjectID()
			mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "id", Value: id}))
			cl, err := NewLC(mt.Coll).CreateLCX(&LcgDto{
				Name: "lcx",
				Age:  1,
			})
			assert.Nil(t, err)
			assert.NotEqual(t, primitive.NilObjectID, cl.ID)
			assert.Equal(t, "lcx", cl.Name)
			assert.Equal(t, 1, cl.Age)
		})
		mt.Run("error", func(t *mtest.T) {
			// CollectionName
			mt.AddMockResponses(bson.D{{"ok", 0}})
			cl, err := NewLC(mt.Coll).CreateLCX(&LcgDto{})
			assert.NotNil(t, err)
			assert.Nil(t, cl)
		})
	})
}

func BenchmarkNewLC(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	mt := mtest.New(&testing.T{}, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for i := 0; i < b.N; i++ {
		_ = NewLC(mt.Coll)
	}
}

func Test_lc_LCX(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.RunOpts("find", mtest.NewOptions().ClientType(mtest.Mock), func(mt *mtest.T) {
		// DatabaseName
		mt.Run("success", func(mt *mtest.T) {
			// CollectionName
			id1 := primitive.NewObjectID()
			id2 := primitive.NewObjectID()

			first := mtest.CreateCursorResponse(1, "lc.lc", mtest.FirstBatch, bson.D{
				{"_id", id1},
				{"name", "lcx1"},
				{"age", 1},
				//{"created_at", time.Now()},
			})
			second := mtest.CreateCursorResponse(1, "lc.lc", mtest.NextBatch, bson.D{
				{"_id", id2},
				{"name", "lcx2"},
				{"age", 2},
				//{"created_at", time.Now()},
			})
			killCursors := mtest.CreateCursorResponse(0, "lc.lc", mtest.NextBatch)
			mt.AddMockResponses(first, second, killCursors)

			nlc := NewLC(mt.Coll)
			lcs, err := nlc.LCX()
			assert.Nil(t, err)
			assert.Equal(t, []*Lcg{
				{ID: id1, Name: "lcx1", Age: 1},
				{ID: id2, Name: "lcx2", Age: 2},
			}, lcs)
		})

		mt.Run("error-context", func(t *mtest.T) {
			// CollectionName
			mt.AddMockResponses(bson.D{{"ok", 0}})
			l := lc{
				c:   mt.Coll,
				ctx: nil,
			}

			lcs, err := l.LCX()
			assert.NotNil(t, err)
			assert.Nil(t, lcs)
		})

		mt.Run("error-cursor", func(t *mtest.T) {
			// CollectionName
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "lc.lc", mtest.FirstBatch, bson.D{
				{"_ids", primitive.NewObjectID()},
				{"names", "lcx1"},
				{"ages", 1},
				//{"created_at", time.Now()},
			}))
			lcs, err := NewLC(mt.Coll).LCX()
			assert.NotNil(t, err)
			assert.Nil(t, lcs)
		})
	})
}

func BenchmarkLc_LCX(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	mt := mtest.New(&testing.T{}, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		id1 := primitive.NewObjectID()
		id2 := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "lc.lc", mtest.FirstBatch, bson.D{
			{"_id", id1},
			{"name", "lcx1"},
			{"age", 1},
			//{"created_at", time.Now()},
		})
		second := mtest.CreateCursorResponse(1, "lc.lc", mtest.NextBatch, bson.D{
			{"_id", id2},
			{"name", "lcx2"},
			{"age", 2},
			//{"created_at", time.Now()},
		})
		killCursors := mtest.CreateCursorResponse(0, "lc.lc", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		nlc := NewLC(mt.Coll)
		for i := 0; i < b.N; i++ {
			_, err := nlc.LCX()
			if err != nil {
				b.Errorf("Benchmark find lcx error: %v", err)
			}
		}
	})
}
