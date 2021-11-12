package v2_transaction_db

var store = syncMap{m: make(map[string]string)}

func Put(key string, value string) error {
	store.Lock()
	defer store.Unlock()
	store.m[key] = value

	return nil
}

func Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	defer store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func Delete(key string) error {
	store.Lock()
	defer store.Unlock()
	delete(store.m, key)

	return nil
}