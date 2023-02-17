package slurm_test

import (
	"reflect"
	"testing"

	"github.com/CLIP-HPC/SlurmCommander/internal/slurm"
)

type gresTest []struct {
	testName  string
	input     string
	expect    int
	expectMap slurm.GresMap
}

var (
	gresTestTable = gresTest{
		{
			testName:  "GRES-empty",
			input:     "",
			expect:    0,
			expectMap: slurm.GresMap{},
		},
		{
			testName:  "GRES-junk: asdf123:123:123:123",
			input:     "asdf123:123:123:123",
			expect:    0,
			expectMap: slurm.GresMap{},
		},
		{
			testName:  "GRES-simple: gpu:8(S:0-1)",
			input:     "gpu:8(S:0-1)",
			expect:    8,
			expectMap: slurm.GresMap{"": 8},
		},
		{
			testName:  "GRES: gpu:P100:8(S:0-1)",
			input:     "gpu:P100:8(S:0-1)",
			expect:    8,
			expectMap: slurm.GresMap{"P100": 8},
		},
		{
			testName:  "GRES_USED: gpu:P100:2(IDX:3,7)",
			input:     "gpu:P100:2(IDX:3,7)",
			expect:    2,
			expectMap: slurm.GresMap{"P100": 2},
		},
		{
			testName:  "GRES: gpu:p100:6(S:0),gpu:rtx:2(S:0)",
			input:     "gpu:p100:6(S:0),gpu:rtx:2(S:0)",
			expect:    8,
			expectMap: slurm.GresMap{"p100": 6, "rtx": 2},
		},
		{
			testName:  "GRES_USED: gpu:p100:0(IDX:N/A),gpu:rtx:0(IDX:N/A)",
			input:     "gpu:p100:0(IDX:N/A),gpu:rtx:0(IDX:N/A)",
			expect:    0,
			expectMap: slurm.GresMap{"p100": 0, "rtx": 0},
		},
		{
			testName:  "GRES_USED: gpu:p100:2(IDX:0-1),gpu:rtx:1(IDX:7)",
			input:     "gpu:p100:2(IDX:0-1),gpu:rtx:1(IDX:7)",
			expect:    3,
			expectMap: slurm.GresMap{"p100": 2, "rtx": 1},
		},
	}
)

func TestParseGRES(t *testing.T) {
	for i, v := range gresTestTable {
		t.Logf("Running test %d : %q\n", i, v.testName)
		rez := *slurm.ParseGRES(v.input)
		t.Logf("Expect: %d Got: %d\n", v.expect, rez)
		if rez != v.expect {
			t.Fatal("FAILED !!!")
		}
	}
}

func TestParseGRESAll(t *testing.T) {
	for i, v := range gresTestTable {
		t.Logf("Running test %d : %q\n", i, v.testName)
		rez := *slurm.ParseGRESAll(v.input)
		t.Logf("Expect: %#v Got: %#v\n", v.expectMap, rez)
		if !reflect.DeepEqual(rez, v.expectMap) {
			t.Fatal("FAILED !!!")
		}
	}
}
