package generator

import "testing"

var aliasTests = []struct {
	alias  string
	result bool
}{
	{"simplealias", true},
	{"simple-alias", true},
	{"simple.alias", true},
	{"simple/alias", true},
	{"simple@alias", true},
	{"simple@alias.com", true},
	{".simplealias", false},
	{"/simplealias", false},
}

func TestIsAlias(t *testing.T) {
	for _, at := range aliasTests {
		b := IsAlias(at.alias)
		if b != at.result {
			t.Errorf("is alias (%s) should be %t", at.alias, at.result)
		}
	}

}

var fileTests = []struct {
	alias  string
	result bool
}{
	{"simplefile", false},
	{"git@simplefile", false},
	{"http://simplefile", false},
	{".simplefile", true},
	{"./simplefile", true},
	{"/simplefile", true},
	{"../simplefile", true},
	{"../testdata/generator/simple", true},
	{"file://simplefile", true},
}

func TestIsFile(t *testing.T) {
	for _, ft := range fileTests {
		b := IsFilePath(ft.alias)
		if b != ft.result {
			t.Errorf("is file (%s) should be %t", ft.alias, ft.result)
		}
	}

}

var repoTests = []struct {
	alias  string
	result bool
}{
	{"simplerepo", false},
	{".simplerepo", false},
	{"./simplerepo", false},
	{"/simplerepo", false},
	{"../simplerepo", false},
	{"file://simplerepo", false},
	{"git@simplerepo", false},
	{"http://simplerepo", true},
	{"http://simplerepo.com", true},
	{"git@simplerepo.org/slug", true},
	{"git@github.com:brainicorn/skelp.git", true},
}

func TestIsRepo(t *testing.T) {
	for _, rt := range repoTests {
		b := IsRepoURL(rt.alias)
		if b != rt.result {
			t.Errorf("is repo (%s) should be %t", rt.alias, rt.result)
		}
	}

}
