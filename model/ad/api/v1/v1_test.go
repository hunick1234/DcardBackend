package v1_test

var testGetAdCase = []struct {
	url     string
	isValid bool
}{
	{"/api/v1/ads", true},
	{"/api/v1/ads?offset=0&limit=10&age=20", true},
	{"/api/v1/ads?offset=0&limit=10&age=", false},
	{"/api/v1/ads?offset=0&limit=10&age=20&gender=M", true},
	{"/api/v1/ads?offset=0&limit=10&age=20&gender=M&country=TW", true},
	{"/api/v1/ads?offset=0&limit=10&age=20&gender=M&country=JP&country=TW", true},
	{"/api/v1/ads?offset=0&limit=10&age=20&gender=M&country=TW&platform=IOS", true},
}

var TestCreateADCase = []struct {
	body    string
	isValid bool
}{
	{`{"title": "AD 1", "endAt": "2023-12-22T01:00:00.000Z"}`, false},
}
