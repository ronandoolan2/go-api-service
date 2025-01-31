package main

type mockDB struct {
	// We could track calls or store in-memory data here if needed
	execError error
}

func (m *mockDB) Init() error {
	return nil
}

func (m *mockDB) Exec(query string, args ...interface{}) error {
	return m.execError
}

func (m *mockDB) Ping() error {
	return nil
}
