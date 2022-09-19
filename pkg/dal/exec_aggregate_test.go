package dal

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/stretchr/testify/require"
)

func TestStepAggregate(t *testing.T) {
	basicAttrs := []simpleAttribute{
		{ident: "k1"},
		{ident: "k2"},
		{ident: "v1"},
		{ident: "txt"},
	}

	type (
		testCase struct {
			name string

			group            []simpleAttribute
			outAttributes    []simpleAttribute
			sourceAttributes []simpleAttribute

			in  []simpleRow
			out []simpleRow

			f internalFilter
		}
	)

	baseBehavior := []testCase{
		// Basic behavior
		{
			name:             "basic one key group",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "g1", "v1": 10, "txt": "foo"},
				{"k1": "g1", "v1": 20, "txt": "fas"},
				{"k1": "g2", "v1": 15, "txt": "bar"},
			},

			out: []simpleRow{
				{"k1": "g1", "v1": float64(30)},
				{"k1": "g2", "v1": float64(15)},
			},

			f: internalFilter{orderBy: filter.SortExprSet{{Column: "k1"}}},
		},
		{
			name:             "basic one key group rename values",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident:  "key_one",
				source: "k1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "something_something",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "g1", "v1": 10, "txt": "foo"},
				{"k1": "g1", "v1": 20, "txt": "fas"},
				{"k1": "g2", "v1": 15, "txt": "bar"},
			},

			out: []simpleRow{
				{"key_one": "g1", "something_something": float64(30)},
				{"key_one": "g2", "something_something": float64(15)},
			},

			f: internalFilter{orderBy: filter.SortExprSet{{Column: "key_one"}}},
		},
		{
			name:             "basic multi key group",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(12)},
				{"k1": "a", "k2": "b", "v1": float64(6)},

				{"k1": "b", "k2": "a", "v1": float64(51)},
			},

			f: internalFilter{orderBy: filter.SortExprSet{{Column: "k1"}, {Column: "k2"}}},
		},
		{
			name:             "basic expr in value aggregation",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(add(v1, 2))",
			}},

			in: []simpleRow{
				{"k1": "g1", "v1": 10, "txt": "foo"},
				{"k1": "g1", "v1": 20, "txt": "fas"},
				{"k1": "g2", "v1": 15, "txt": "bar"},
			},

			out: []simpleRow{
				{"k1": "g1", "v1": float64(34)},
				{"k1": "g2", "v1": float64(17)},
			},

			f: internalFilter{orderBy: filter.SortExprSet{{Column: "k1"}}},
		},
	}

	filtering := []testCase{
		{
			name:             "filtering constraints single attr",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			},
			},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(12)},
				{"k1": "a", "k2": "b", "v1": float64(6)},
			},

			f: internalFilter{
				constraints: map[string][]any{"k1": {"a"}},
			},
		},
		{
			name:             "filtering constraints multiple attrs",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				// ---
				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "b", "v1": float64(6)},
			},

			f: internalFilter{
				constraints: map[string][]any{"k1": {"a"}, "k2": {"b"}},
			},
		},
		{
			name:             "filtering constraints single attr multiple options",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "b", "v1": 2, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 3, "txt": "fas"},
				{"k1": "c", "k2": "a", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(10)},
				{"k1": "a", "k2": "b", "v1": float64(2)},
				{"k1": "b", "k2": "a", "v1": float64(3)},
			},

			f: internalFilter{
				orderBy:     filter.SortExprSet{{Column: "k1"}, {Column: "k2"}},
				constraints: map[string][]any{"k1": {"a", "b"}},
			},
		},
		{
			name:             "filtering expression simple expression",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				// ---
				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(12)},
			},

			f: internalFilter{
				expression: "v1 > 10 && v1 < 20",
			},
		},
		{
			name:             "filtering expression check renamed aggregate",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "some_sum_value",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				// ---
				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "some_sum_value": float64(12)},
			},

			f: internalFilter{
				expression: "some_sum_value > 10 && some_sum_value < 20",
			},
		},
		{
			name:             "filtering expression constant true",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				// ---
				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(12)},
				{"k1": "a", "k2": "b", "v1": float64(6)},

				{"k1": "b", "k2": "a", "v1": float64(51)},
			},

			f: internalFilter{
				expression: "true",
				orderBy:    filter.SortExprSet{{Column: "k1"}, {Column: "k2"}},
			},
		},
		{
			name:             "filtering expression constant false",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				// ---
				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out: []simpleRow{},

			f: internalFilter{
				expression: "false",
			},
		},
	}

	sorting := []testCase{
		{
			name:             "sorting single key full key asc",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "v1": 2, "txt": "fas"},
				{"k1": "b", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "v1": float64(12)},
				{"k1": "b", "v1": float64(3)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "k1", Descending: false}},
			},
		},
		{
			name:             "sorting single aggregate full asc",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "some_sum_value",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "v1": 2, "txt": "fas"},
				{"k1": "b", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "b", "some_sum_value": float64(3)},
				{"k1": "a", "some_sum_value": float64(12)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "some_sum_value", Descending: false}},
			},
		},
		{
			name:             "sorting single aggregate full desc",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "some_sum_value",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "v1": 2, "txt": "fas"},
				{"k1": "b", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "some_sum_value": float64(12)},
				{"k1": "b", "some_sum_value": float64(3)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "some_sum_value", Descending: true}},
			},
		},
		{
			name:             "sorting single key full key dsc",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "v1": 2, "txt": "fas"},
				{"k1": "b", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "b", "v1": float64(3)},
				{"k1": "a", "v1": float64(12)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "k1", Descending: true}},
			},
		},
		{
			name:             "sorting multiple key full key asc",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "b", "v1": 2, "txt": "fas"},
				{"k1": "b", "k2": "c", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(10)},
				{"k1": "a", "k2": "b", "v1": float64(2)},
				{"k1": "b", "k2": "c", "v1": float64(3)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "k1", Descending: false}, {Column: "k2", Descending: false}},
			},
		},
		{
			name:             "sorting multiple key full key dsc",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "b", "v1": 2, "txt": "fas"},
				{"k1": "b", "k2": "c", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "b", "k2": "c", "v1": float64(3)},
				{"k1": "a", "k2": "b", "v1": float64(2)},
				{"k1": "a", "k2": "a", "v1": float64(10)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "k1", Descending: true}, {Column: "k2", Descending: true}},
			},
		},
		{
			name:             "sorting multiple key full key mixed",
			sourceAttributes: basicAttrs,
			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			}},

			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "b", "v1": 2, "txt": "fas"},
				{"k1": "b", "k2": "c", "v1": 3, "txt": "fas"},
			},

			out: []simpleRow{
				{"k1": "a", "k2": "b", "v1": float64(2)},
				{"k1": "a", "k2": "a", "v1": float64(10)},
				{"k1": "b", "k2": "c", "v1": float64(3)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "k1", Descending: false}, {Column: "k2", Descending: true}},
			},
		},
	}

	exprGroups := []testCase{
		{
			name: "expression as key year",
			sourceAttributes: []simpleAttribute{
				{ident: "dob"},
				{ident: "name"},
			},
			group: []simpleAttribute{{
				ident: "dob_y",
				expr:  "year(dob)",
			}},
			outAttributes: []simpleAttribute{{
				ident: "users",
				expr:  "count(name)",
			}},

			in: []simpleRow{
				{"dob": "2022-10-20T09:44:49Z", "name": "Ana"},
				{"dob": "2022-10-20T09:44:49Z", "name": "John"},
				{"dob": "2021-10-20T09:44:49Z", "name": "Jane"},
			},

			out: []simpleRow{
				{"dob_y": 2021, "users": float64(1)},
				{"dob_y": 2022, "users": float64(2)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "dob_y", Descending: false}},
			},
		},
		{
			name: "expression as key year with calc",
			sourceAttributes: []simpleAttribute{
				{ident: "dob"},
				{ident: "name"},
			},
			group: []simpleAttribute{{
				ident: "dob_y",
				expr:  "year(dob)/10",
			}},
			outAttributes: []simpleAttribute{{
				ident: "users",
				expr:  "count(name)",
			}},

			in: []simpleRow{
				{"dob": "2022-10-20T09:44:49Z", "name": "Ana"},
				{"dob": "2022-10-20T09:44:49Z", "name": "John"},
				{"dob": "2021-10-20T09:44:49Z", "name": "Jane"},
			},

			out: []simpleRow{
				{"dob_y": 202.1, "users": float64(1)},
				{"dob_y": 202.2, "users": float64(2)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "dob_y", Descending: false}},
			},
		},
		{
			name: "same group expression",
			sourceAttributes: []simpleAttribute{
				{ident: "name"},
			},
			group: []simpleAttribute{{
				ident: "d",
				// @note will only run for a year then will need to be changed
				expr: "year(now())",
			}},
			outAttributes: []simpleAttribute{{
				ident: "users",
				expr:  "count(name)",
			}},

			in: []simpleRow{
				{"name": "Ana"},
				{"name": "John"},
				{"name": "Jane"},
			},

			out: []simpleRow{
				{"d": 2022, "users": float64(3)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "d", Descending: false}},
			},
		},
		{
			name: "same group constant",
			sourceAttributes: []simpleAttribute{
				{ident: "name"},
			},
			group: []simpleAttribute{{
				ident: "d",
				expr:  "'a'",
			}},
			outAttributes: []simpleAttribute{{
				ident: "users",
				expr:  "count(name)",
			}},

			in: []simpleRow{
				{"name": "Ana"},
				{"name": "John"},
				{"name": "Jane"},
			},

			out: []simpleRow{
				{"d": "a", "users": float64(3)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "d", Descending: false}},
			},
		},
		{
			name: "expression as key concatenated",
			sourceAttributes: []simpleAttribute{
				{ident: "dob"},
				{ident: "name"},
			},
			group: []simpleAttribute{{
				ident: "dob_y",
				expr:  "concat(string(year(dob)), '-', string(month(dob)))",
			}},
			outAttributes: []simpleAttribute{{
				ident: "users",
				expr:  "count(name)",
			}},

			in: []simpleRow{
				{"dob": "2022-10-20T09:44:49Z", "name": "Ana"},
				{"dob": "2022-10-20T09:44:49Z", "name": "John"},
				{"dob": "2021-10-20T09:44:49Z", "name": "Jane"},
			},

			out: []simpleRow{
				{"dob_y": "2021-10", "users": float64(1)},
				{"dob_y": "2022-10", "users": float64(2)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "dob_y", Descending: false}},
			},
		},
	}

	nilValues := []testCase{
		{
			name: "nil in group key single value",
			sourceAttributes: []simpleAttribute{
				{ident: "thing"},
				{ident: "name"},
			},
			group: []simpleAttribute{{
				ident: "thing",
			}},
			outAttributes: []simpleAttribute{{
				ident: "users",
				expr:  "count(name)",
			}},

			in: []simpleRow{
				{"thing": "A", "name": "Ana"},
				{"name": "John"},
				{"name": "Jane"},
			},

			out: []simpleRow{
				{"thing": nil, "users": float64(2)},
				{"thing": "A", "users": float64(1)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "thing", Descending: false}},
			},
		},
		{
			name: "nil in group key multiple value",
			sourceAttributes: []simpleAttribute{
				{ident: "thing"},
				{ident: "another"},
				{ident: "name"},
			},
			group: []simpleAttribute{{
				ident: "thing",
			}, {
				ident: "another",
			}},
			outAttributes: []simpleAttribute{{
				ident: "users",
				expr:  "count(name)",
			}},

			in: []simpleRow{
				{"thing": "A", "another": "A", "name": "Ana"},
				{"thing": "A", "name": "Ana"},
				{"another": "A", "name": "Ana"},
				{"name": "John"},
				{"name": "Jane"},
			},

			out: []simpleRow{
				{"thing": nil, "another": nil, "users": float64(2)},
				{"thing": nil, "another": "A", "users": float64(1)},
				{"thing": "A", "another": nil, "users": float64(1)},
				{"thing": "A", "another": "A", "users": float64(1)},
			},

			f: internalFilter{
				orderBy: filter.SortExprSet{{Column: "thing", Descending: false}, {Column: "another", Descending: false}},
			},
		},
	}

	batches := [][]testCase{
		baseBehavior,
		filtering,
		sorting,
		exprGroups,
		nilValues,
	}

	for _, batch := range batches {
		for _, tc := range batch {
			t.Run(tc.name, func(t *testing.T) {
				bootstrapAggregate(t, func(ctx context.Context, t *testing.T, sa *Aggregate, b Buffer) {
					for _, r := range tc.in {
						require.NoError(t, b.Add(ctx, r))
					}
					sa.Ident = tc.name
					sa.SourceAttributes = saToMapping(tc.sourceAttributes...)
					sa.Group = saToMapping(tc.group...)
					sa.OutAttributes = saToMapping(tc.outAttributes...)
					sa.filter = tc.f

					aa, err := sa.iterator(ctx, b)
					require.NoError(t, err)

					i := 0
					for aa.Next(ctx) {
						out := simpleRow{}
						require.NoError(t, aa.Scan(out))
						require.Equal(t, tc.out[i], out)
						i++
					}
					require.NoError(t, aa.Err())
					require.Equal(t, len(tc.out), i)
				})
			})
		}
	}
}

func TestStepAggregateValidation(t *testing.T) {
	ctx := context.Background()

	basicAttrs := []simpleAttribute{
		{ident: "k1"},
		{ident: "k2"},
		{ident: "v1"},
		{ident: "txt"},
	}

	run := func(t *testing.T, groups []simpleAttribute, attr []simpleAttribute) (err error) {
		sa := &Aggregate{
			Ident:            "agg",
			SourceAttributes: saToMapping(basicAttrs...),
			Group:            saToMapping(groups...),
			OutAttributes:    saToMapping(attr...),
		}

		return sa.dryrun(ctx)
	}

	runF := func(t *testing.T, f internalFilter, groups []simpleAttribute, attr []simpleAttribute) (err error) {
		sa := &Aggregate{
			Ident:            "agg",
			SourceAttributes: saToMapping(basicAttrs...),
			Group:            saToMapping(groups...),
			OutAttributes:    saToMapping(attr...),
			Filter:           f,
		}

		return sa.dryrun(ctx)
	}

	groups := []simpleAttribute{{
		ident: "k1",
	}}

	aggregates := []simpleAttribute{{
		ident: "v1",
		expr:  "sum(v1)",
	}}

	t.Run("no group attrs", func(t *testing.T) {
		err := run(t, nil, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no group attributes specified")
	})
	t.Run("no aggregates", func(t *testing.T) {
		err := run(t, groups, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no output attributes specified")
	})

	t.Run("group ident doesn't exist", func(t *testing.T) {
		groups := []simpleAttribute{{
			ident: "i_not_real",
		}}

		err := run(t, groups, aggregates)
		require.Error(t, err)
		require.Contains(t, err.Error(), "i_not_real")
	})

	t.Run("group func ident doesn't exist", func(t *testing.T) {
		groups := []simpleAttribute{{
			ident: "month(i_not_real)",
		}}

		err := run(t, groups, aggregates)
		require.Error(t, err)
		require.Contains(t, err.Error(), "i_not_real")
	})

	t.Run("aggregate func ident does not exist", func(t *testing.T) {
		aggregates := []simpleAttribute{{
			ident: "i_not_here",
			expr:  "sum(i_not_here)",
		}}

		err := run(t, groups, aggregates)
		require.Error(t, err)
		require.Contains(t, err.Error(), "i_not_here")
	})

	t.Run("sort ident does not exist", func(t *testing.T) {
		err := runF(t, internalFilter{orderBy: filter.SortExprSet{{Column: "i_not_yes"}}}, groups, aggregates)
		require.Error(t, err)
		require.Contains(t, err.Error(), "i_not_yes")
	})
}

func TestStepAggregate_cursorCollect_forward(t *testing.T) {
	tcc := []struct {
		name          string
		ss            filter.SortExprSet
		in            simpleRow
		group         []simpleAttribute
		outAttributes []simpleAttribute

		out func() *filter.PagingCursor
		err bool
	}{
		{
			name: "simple",
			in:   simpleRow{"pk1": 1, "f1": "v1"},
			group: []simpleAttribute{{
				ident: "pk1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "f1",
			}},
			out: func() *filter.PagingCursor {
				pc := &filter.PagingCursor{}
				pc.Set("pk1", 1, false)
				return pc
			},
		},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {

			def := Aggregate{
				filter: internalFilter{
					orderBy: c.ss,
				},
				Group:         saToMapping(c.group...),
				OutAttributes: saToMapping(c.outAttributes...),
			}

			out, err := (&aggregate{def: def}).ForwardCursor(c.in)
			require.NoError(t, err)

			require.Equal(t, c.out(), out)
		})
	}
}

func TestStepAggregate_cursorCollect_back(t *testing.T) {
	tcc := []struct {
		name          string
		ss            filter.SortExprSet
		in            simpleRow
		group         []simpleAttribute
		outAttributes []simpleAttribute

		out func() *filter.PagingCursor
		err bool
	}{
		{
			name: "simple",
			in:   simpleRow{"pk1": 1, "f1": "v1"},
			group: []simpleAttribute{{
				ident: "pk1",
			}},
			outAttributes: []simpleAttribute{{
				ident: "f1",
			}},
			out: func() *filter.PagingCursor {
				pc := &filter.PagingCursor{}
				pc.Set("pk1", 1, false)
				pc.ROrder = true
				return pc
			},
		},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {

			def := Aggregate{
				filter: internalFilter{
					orderBy: c.ss,
				},
				Group:         saToMapping(c.group...),
				OutAttributes: saToMapping(c.outAttributes...),
			}

			out, err := (&aggregate{def: def}).BackCursor(c.in)
			require.NoError(t, err)

			require.Equal(t, c.out(), out)
		})
	}
}

func TestStepAggregate_more(t *testing.T) {
	basicAttrs := []simpleAttribute{
		{ident: "k1"},
		{ident: "k2"},
		{ident: "v1"},
		{ident: "txt"},
	}

	tcc := []struct {
		name string
		in   []simpleRow

		group            []simpleAttribute
		outAttributes    []simpleAttribute
		sourceAttributes []simpleAttribute

		def *Aggregate

		out1 []simpleRow
		out2 []simpleRow
	}{
		{
			name:             "multiple keys",
			sourceAttributes: basicAttrs,
			in: []simpleRow{
				{"k1": "a", "k2": "a", "v1": 10, "txt": "foo"},
				{"k1": "a", "k2": "a", "v1": 2, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},
				{"k1": "a", "k2": "b", "v1": 3, "txt": "fas"},

				// ---
				{"k1": "b", "k2": "a", "v1": 20, "txt": "fas"},
				{"k1": "b", "k2": "a", "v1": 31, "txt": "fas"},
			},

			out1: []simpleRow{
				{"k1": "a", "k2": "a", "v1": float64(12)},
			},
			out2: []simpleRow{
				{"k1": "a", "k2": "b", "v1": float64(6)},
				{"k1": "b", "k2": "a", "v1": float64(51)},
			},

			def: &Aggregate{},

			group: []simpleAttribute{{
				ident: "k1",
			}, {
				ident: "k2",
			}},
			outAttributes: []simpleAttribute{{
				ident: "v1",
				expr:  "sum(v1)",
			},
			},
		},
	}

	ctx := context.Background()
	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			buff := InMemoryBuffer()
			for _, r := range tc.in {
				require.NoError(t, buff.Add(ctx, r))
			}

			d := tc.def
			d.Group = saToMapping(tc.group...)
			d.OutAttributes = saToMapping(tc.outAttributes...)
			d.SourceAttributes = saToMapping(tc.sourceAttributes...)
			for _, k := range tc.group {
				d.filter.orderBy = append(d.filter.orderBy, &filter.SortExpr{Column: k.ident})
			}

			aa, err := d.iterator(ctx, buff)
			require.NoError(t, err)

			require.True(t, aa.Next(ctx))
			out := simpleRow{}
			require.NoError(t, aa.Err())
			require.NoError(t, aa.Scan(out))
			require.Equal(t, tc.out1[0], out)

			buff.Seek(ctx, 0)
			require.NoError(t, aa.More(0, out))

			i := 0
			for aa.Next(ctx) {
				out := simpleRow{}
				require.NoError(t, aa.Err())
				require.NoError(t, aa.Scan(out))

				require.Equal(t, tc.out2[i], out)

				i++
			}
			require.NoError(t, aa.Err())
			require.Equal(t, len(tc.out2), i)
		})
	}
}
