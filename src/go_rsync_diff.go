package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	flag "github.com/docker/docker/pkg/mflag"
)

var builddate string

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	rsyncFrom string
	rsyncTo   string
)

var (
	targetFiles      arrayFlags
	rsyncExclude     string
	rsyncExcludeFrom string
	isVD             bool
	isLess           bool
	isColor          bool
	isC              bool
	isU              bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, ``)
		fmt.Fprintln(os.Stderr, "Difference verification tool of rsync command")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "BuildDate:%s\n", builddate)
		fmt.Fprintf(os.Stderr, "%s [OPTIONS] RSYNC-FROM RSYNC-TO\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Var(&targetFiles, []string{"t", "-target-file"}, "Difference acquisition object file")
	flag.StringVar(&rsyncExclude, []string{"e", "-exclude"}, "", "rsync exclude option")
	flag.StringVar(&rsyncExcludeFrom, []string{"-exclude-from"}, "", "rsync exlude from option")
	flag.BoolVar(&isVD, []string{"v", "-vimdiff"}, false, "Produce a vimdiff. Specify also t option")
	flag.BoolVar(&isLess, []string{"l", "-less"}, false, "using less for output")
	flag.BoolVar(&isColor, []string{"r", "-colordiff"}, false, "using colordiff for output")
	flag.BoolVar(&isC, []string{"c", "-context-diff"}, false, "Produce a context format diff")
	flag.BoolVar(&isU, []string{"-unified-diff"}, true, "Produce a unified format diff (default)")
	flag.Parse()
}

func parseRsyncOutput(rsyncOutput []string) []string {
	diffTargets := make([]string, 0, len(rsyncOutput))
	reg := regexp.MustCompile(".*(\\++|\\.+|deleting)\\ ")
	for _, line := range rsyncOutput {
		filename := reg.ReplaceAllString(strings.TrimRight(line, "\n"), "")
		if len(targetFiles) != 0 && !isSpecifyTarget(filename) {
			continue
		}
		diffTargets = append(diffTargets, filename)
	}
	return diffTargets
}

func exitOnUsage(msg string, err error) {
	if err != nil {
		flag.Usage()
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func optCheck(args []string) {
	var err error

	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	_, err = exec.LookPath("rsync")
	exitOnUsage("install of rsync is necessary", err)

	_, err = exec.LookPath("cat")
	exitOnUsage("install of cat is necessary", err)

	if isColor {
		_, err = exec.LookPath("colordiff")
		exitOnUsage("install of colordiff is necessary", err)
	}

	if isVD && len(targetFiles) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if isVD {
		_, err = exec.LookPath("vimdiff")
		exitOnUsage("install of vimdiff is necessary", err)
	}
}

func exitOnErr(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func main() {
	args := flag.Args()
	optCheck(args)
	rsyncFrom = args[0]
	rsyncTo = args[1]

	cmd := &Command{}
	cmd.Command = "rsync"
	cmd.Options = []string{rsyncFrom, rsyncTo}
	out, err := cmd.Output()
	exitOnErr("rsync proccess", err)
	lines := strings.Split(string(out), "\n")
	diffTargets := parseRsyncOutput(lines)

	outputter := new(DiffOutputter)
	outputter.CreateFunction()
	exit, err := outputter.Output(diffTargets)
	exitOnErr("", err)
	os.Exit(exit)
}
