package model

// Increment 自增
type Increment struct {
	ID int64
}

func GetIncrement() (int64, error) {
	id := &Increment{}
	err := DB.Get(IncrementTable, IncrementKey, id)
	if err != nil && err != ISNULL {
		return 0, err
	}
	if err == ISNULL {
		id.ID = 1
	} else {
		id.ID++
	}
	err = DB.Set(IncrementTable, IncrementKey, id)
	if err != nil {
		return 0, err
	}
	return id.ID, nil
}

func ResetIncrement() error {
	id := &Increment{0}
	return DB.Set(IncrementTable, IncrementKey, id)
}
