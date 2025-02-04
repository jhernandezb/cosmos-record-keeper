package recordkeeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func TestUint64Keys(t *testing.T) {
	ctx, keeper, _ := mockRecordKeeper()
	assert.NotNil(t, ctx)
	assert.NotNil(t, keeper)

	type Record struct{}

	// adding
	record := Record{}
	id := keeper.Add(ctx, record)
	assert.Equal(t, uint64(1), id)

	// getting
	var expectedRecord Record
	err := keeper.Get(ctx, id, &expectedRecord)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord, record)

	// iteration
	err = keeper.Each(ctx, func(recordBytes []byte) bool {
		var r Record
		keeper.codec.MustUnmarshalBinaryLengthPrefixed(recordBytes, &r)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), id)
		return true
	})
	assert.NoError(t, err)
}

func TestStringKeys(t *testing.T) {
	ctx, keeper, _ := mockRecordKeeper()
	assert.NotNil(t, ctx)
	assert.NotNil(t, keeper)

	type Record struct{}

	// setter
	record := Record{}
	keeper.StringSet(ctx, "key1", record)

	// getter
	var expectedRecord Record
	keeper.StringGet(ctx, "key1", &expectedRecord)
	assert.Equal(t, expectedRecord, record)
}

func TestAssociationList(t *testing.T) {
	ctx, k, k2 := mockRecordKeeper()

	k.Push(ctx, k.storeKey, k2.storeKey, uint64(1), uint64(2))
	k.Map(ctx, k2.storeKey, uint64(2), func(id uint64) {
		assert.Equal(t, uint64(1), id)
	})
}

func mockRecordKeeper() (sdk.Context, RecordKeeper, RecordKeeper) {
	db := dbm.NewMemDB()

	storeKey := sdk.NewKVStoreKey("records")
	storeKey2 := sdk.NewKVStoreKey("records2")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storeKey2, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	keeper := NewRecordKeeper(storeKey, codec)
	keeper2 := NewRecordKeeper(storeKey2, codec)

	return ctx, keeper, keeper2
}
