# Cosmos Record Keeper

A `uint64` indexed iterable type keeper for [Cosmos SDK](https://github.com/cosmos/cosmos-sdk).

## Example

Embed a `RecordKeeper` struct inside a `Keeper`.

i.e:
```go
    type Keeper struct {
        RecordKeeper
    }
```

Define an index field for a type:

```go
type Record struct {
	ID uint64
}
```

### Initialization

```go
	storeKey := sdk.NewKVStoreKey("records")
	keeper := Keeper{&NewRecordKeeper(storeKey)}
```

### Setting

```go
	record := Record{
		ID: keeper.NextID(ctx),
	}

	// marshal record
	recordBytes, err := json.Marshal(record)

	// set in kvstore
	keeper.Store(ctx).Set(
		keeper.IDKey(record.ID),
		recordBytes)
```

### Getting

```go
	idBytes := keeper.IDKey(record.ID)
	recordBytes := keeper.Store(ctx).Get(idBytes)
```
