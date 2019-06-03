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

func TestStringRecordKeeper(t *testing.T) {
	ctx, keeper := mockStringKeeper()
	assert.NotNil(t, ctx)
	assert.NotNil(t, keeper)

	type Record struct{}

	// setter
	record := Record{}
	keeper.Set(ctx, "key1", record)

	// getter
	var expectedRecord Record
	err := keeper.Get(ctx, "key1", &expectedRecord)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord, record)
}

func TestUint64RecordKeeper(t *testing.T) {
	ctx, keeper := mockUint64Keeper()
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

func mockUint64Keeper() (sdk.Context, RecordKeeper) {
	db := dbm.NewMemDB()

	storeKey := sdk.NewKVStoreKey("records")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	keeper := NewRecordKeeper(storeKey, codec)

	return ctx, keeper
}

func mockStringKeeper() (sdk.Context, StringRecordKeeper) {
	db := dbm.NewMemDB()

	storeKey := sdk.NewKVStoreKey("records")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := codec.New()
	keeper := NewStringRecordKeeper(storeKey, codec)

	return ctx, keeper
}
