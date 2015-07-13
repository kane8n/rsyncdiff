package main

import (
	"io/ioutil"
	"os"
	"path"

	difflib "github.com/pmezard/go-difflib/difflib"
)

type OutputFunction func(diffTargets []string) (int, error)
type DiffOutputter struct {
	function OutputFunction
}

func (d *DiffOutputter) CreateFunction() {
	if isVD {
		d.function = VimDiffOutput
		return
	}

	if isC {
		d.function = ContextOutput
		return
	}

	d.function = UnifiedOutput
}

func (d *DiffOutputter) Output(diffTargets []string) (int, error) {
	return d.function(diffTargets)
}

type DiffInfo struct {
	A        []string
	B        []string
	FromFile string
	ToFile   string
}

func NewDiffInfo(target string) (*DiffInfo, error) {
	from := path.Join(rsyncFrom, target)
	to := path.Join(rsyncTo, target)
	fromExist := isFileExist(from)
	fromIsDiffTarget := isDiffTarget(from)
	toExist := isFileExist(to)
	toIsDiffTarget := isDiffTarget(to)
	fromLines := ""
	if fromExist && fromIsDiffTarget {
		bytes, err := ioutil.ReadFile(from)
		if err != nil {
			return nil, err
		}
		fromLines = string(bytes)
	}
	toLines := ""
	if toExist && toIsDiffTarget {
		bytes, err := ioutil.ReadFile(to)
		if err != nil {
			return nil, err
		}
		toLines = string(bytes)
	}

	info := &DiffInfo{}
	info.A = difflib.SplitLines(toLines)
	info.B = difflib.SplitLines(fromLines)
	info.FromFile = from
	info.ToFile = to
	return info, nil
}

func UnifiedOutput(diffTargets []string) (int, error) {
	tmp, err := ioutil.TempFile(os.TempDir(), "rsyncdiff")
	if err != nil {
		return 1, err
	}
	defer os.Remove(tmp.Name())

	for _, target := range diffTargets {
		diffInfo, err := NewDiffInfo(target)
		if err != nil {
			return 1, err
		}
		diff := difflib.UnifiedDiff{
			A:        diffInfo.A,
			B:        diffInfo.B,
			FromFile: diffInfo.FromFile,
			ToFile:   diffInfo.ToFile,
		}
		result, err := difflib.GetUnifiedDiffString(diff)
		if err != nil {
			return 1, err
		}
		_, err = tmp.WriteString(result)
		if err != nil {
			return 1, err
		}
		tmp.Sync()
	}

	cmd := &Command{}
	cmd.Command = "cat"
	cmd.Options = []string{tmp.Name()}
	return cmd.Run()
}

func ContextOutput(diffTargets []string) (int, error) {
	tmp, err := ioutil.TempFile(os.TempDir(), "rsyncdiff")
	if err != nil {
		return 1, err
	}
	defer os.Remove(tmp.Name())

	for _, target := range diffTargets {
		diffInfo, err := NewDiffInfo(target)
		if err != nil {
			return 1, err
		}
		diff := difflib.ContextDiff{
			A:        diffInfo.A,
			B:        diffInfo.B,
			FromFile: diffInfo.FromFile,
			ToFile:   diffInfo.ToFile,
		}
		result, err := difflib.GetContextDiffString(diff)
		if err != nil {
			return 1, err
		}
		_, err = tmp.WriteString(result)
		if err != nil {
			return 1, err
		}
		tmp.Sync()
	}

	cmd := &Command{}
	cmd.Command = "cat"
	cmd.Options = []string{tmp.Name()}
	return cmd.Run()
}

func VimDiffOutput(diffTargets []string) (int, error) {
	pairs := getTargetFilePair(diffTargets)
	for _, pair := range pairs {
		cmd := &Command{}
		cmd.Command = "vimdiff"
		cmd.Options = []string{pair[0], pair[1]}
		exit, err := cmd.Run()
		if err != nil {
			return exit, err
		}
	}
	return 0, nil
}

func getTargetFilePair(diffTargets []string) [][]string {
	pairs := make([][]string, 0, len(diffTargets))
	for _, target := range diffTargets {
		from := path.Join(rsyncFrom, target)
		to := path.Join(rsyncTo, target)
		fromExist := isFileExist(from)
		fromIsDiffTarget := isDiffTarget(from)
		toExist := isFileExist(to)
		toIsDiffTarget := isDiffTarget(to)
		if fromExist && fromIsDiffTarget && toExist && toIsDiffTarget {
			pairs = append(pairs, []string{from, to})
			continue
		}
		if fromExist && fromIsDiffTarget {
			pairs = append(pairs, []string{from, ""})
			continue
		}
		if toExist && toIsDiffTarget {
			pairs = append(pairs, []string{"", to})
			continue
		}
	}
	return pairs
}

func isDiffTarget(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !(fi.IsDir() || isBinary(filename))
}

func isBinary(filename string) bool {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return false
	}
	zeros := 0
	for _, b := range bytes {
		if b == 0 {
			zeros++
		}
		if zeros >= 4 {
			return true
		}
	}
	return false
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isSpecifyTarget(filename string) bool {
	for _, target := range targetFiles {
		if target == filename {
			return true
		}
	}
	return false
}
