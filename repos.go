package main

import (
	"doc"
	"gnd.la/log"
	"gnd.la/util"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	gitPath  string
	reposTxt = util.RelativePath("repos.txt")
)

func repoPath(repo string) string {
	if p := strings.Index(repo, "://"); p >= 0 {
		return repo[p+3:]
	}
	return repo
}

func Repos() []string {
	data, _ := ioutil.ReadFile(reposTxt)
	if data != nil {
		var repos []string
		for _, v := range strings.Split(string(data), "\n") {
			if tr := strings.TrimSpace(v); tr != "" {
				repos = append(repos, tr)
			}
		}
		return repos
	}
	return nil
}

func StartUpdatingRepos() {
	go UpdateRepos()
	go func() {
		for _ = range time.Tick(30 * time.Minute) {
			UpdateRepos()
		}
	}()
}

func UpdateRepos() {
	for _, v := range Repos() {
		p := filepath.Join(doc.SourceDir, filepath.FromSlash(repoPath(v)))
		if _, err := os.Stat(p); err == nil {
			cmd := exec.Command(gitPath, "pull")
			cmd.Dir = p
			log.Debugf("updating %s in %s", v, p)
			if err := cmd.Run(); err != nil {
				log.Errorf("can't update repo %s: %s", v, err)
			}
		} else {
			parent := filepath.Dir(p)
			if err := os.MkdirAll(parent, 0755); err != nil {
				panic(err)
			}
			cmd := exec.Command(gitPath, "clone", v)
			cmd.Dir = parent
			log.Debugf("cloning %s to %s", v, p)
			if err := cmd.Run(); err != nil {
				log.Errorf("can't clone repo %s: %s", v, err)
			}
		}
	}
}

func init() {
	var err error
	gitPath, err = exec.LookPath("git")
	if err != nil {
		panic(err)
	}
}
