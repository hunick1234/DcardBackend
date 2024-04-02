package ad

import (
	"strconv"
	"testing"
)

var testAdCase = []struct {
	ad       AD
	expected bool
	err      string
}{
	{
		AD{
			Title:   "test",
			EndAt:   "2020-12-12T12:12:12Z",
			StartAt: "2020-12-12T12:12:12Z",
		},
		true,
		"",
	},
	{
		AD{
			Title:   "test_1",
			EndAt:   "2020-12-12T12:12:12Z",
			StartAt: "2020-12-22T12:12:12Z",
		},
		false,
		"start time should be less than end time",
	},
	{
		AD{
			Title: "test_2",
			EndAt: "2020-12-12T12:12:12Z",
		},
		false,
		"need start time",
	},
	{
		AD{
			Title:   "test_3",
			StartAt: "2020-12-12T12:12:12Z",
		},
		false,
		"need end time",
	},
}

var testConditionsCase = []struct {
	conditions Conditions
	expected   bool
	err        string
}{
	{
		Conditions{
			AgeStart: 20,
			AgeEnd:   30,
			Gender:   []string{"m", "f"},
			Country:  []string{"US"},
			Platform: []string{"IOS", "Android"},
		},
		true,
		"",
	},
	{
		Conditions{
			Country: []string{"USA", "Canada", "Mexico"},
		},
		false,
		"invalid country",
	},
	{
		Conditions{
			Gender: []string{"male", "female", "unknown"},
		},
		false,
		"invalid gender",
	},
	{
		Conditions{
			AgeStart: -1,
			AgeEnd:   101,
		},
		false,
		"age start should be greater than 0 and age end should be less than 100",
	},
	{
		Conditions{
			AgeStart: -1,
		},
		false,
		"age start should be greater than 0",
	},
	{
		Conditions{
			AgeEnd: 10122243,
		},
		false,
		"age end should be less than 100",
	},
}

var testAdQueryCase = []struct {
	query    AdQuery
	expected bool
	err      string
}{
	{
		AdQuery{
			Offset:   0,
			Limit:    10,
			Age:      20,
			Gender:   []string{},
			Country:  []string{"US"},
			Platform: []string{"Android", "ios"},
		},
		true,
		"",
	},
	{
		AdQuery{
			Offset: -1,
		},
		false,
		"invalid offset",
	},
}

func TestIsValidConditions(t *testing.T) {
	for i, test := range testConditionsCase {
		testName := "Conditions test case " + strconv.Itoa(i)
		t.Run(testName, func(t *testing.T) {
			if test.conditions.IsValidConditions() != test.expected {
				if test.expected {
					t.Errorf("Test case %s: shouldn't err", testName)
				}

				if !test.expected {
					t.Errorf("Test case %s: should err", testName)
				}
			}
		})
	}
}

func TestIsValidAdQuery(t *testing.T) {
	for i, test := range testAdQueryCase {
		testName := "Query test case " + strconv.Itoa(i)
		t.Run(testName, func(t *testing.T) {
			if test.query.IsValidAdQuery() != test.expected {
				if test.expected {
					t.Errorf("Test case %s: shouldn't err", testName)
				}

				if !test.expected {
					t.Errorf("Test case %s: should err", testName)
				}
			}
		})
	}
}

func TestIsValidAd(t *testing.T) {
	for i, test := range testAdCase {
		testName := "Ad test case " + strconv.Itoa(i)
		t.Run(testName, func(t *testing.T) {
			if test.ad.IsValidAd() != test.expected {
				if test.expected {
					t.Errorf("Test case %s: shouldn't err", testName)
				}

				if !test.expected {
					t.Errorf("Test case %s: should err", testName)
				}
			}
		})
	}
}
