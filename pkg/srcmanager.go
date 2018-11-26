package srcmanager

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/Unknwon/com"
	glob "github.com/mattn/go-zglob"
	log "github.com/sirupsen/logrus"
)

var (
	gopath       = com.GetGOPATHs()[0]
	repositories = _escFSMustString(false, "/repositories")
	Verbose      = false
)

func RepositoryURLs(isPublic bool) ([]string, error) {
	var repos []string
	for _, repo := range strings.Split(strings.TrimSpace(string(repositories)), "\n") {
		if strings.HasPrefix(repo, "#") || strings.TrimSpace(repo) == "" {
			continue
		}
		if !isPublic && strings.Contains(repo, "[private]") {
			continue
		}
		repo = strings.Replace(repo, "[private]", "", -1)
		repos = append(repos, repo)
	}
	// probably not needed because of the TrimSpace
	if len(repos) > 0 && repos[len(repos)-1] == "" {
		repos = repos[:len(repos)-1]
	}
	return repos, nil
}

var re = regexp.MustCompile("^github.com/(.+)$")

func githubURL(isPublic bool, url string) string {
	if isPublic {
		return re.ReplaceAllString(url, "https://github.com/${1}.git")
	}
	return re.ReplaceAllString(url, "git@github.com:${1}.git")
}

func Commit(isPublic bool, message string) error {
	rawURLs, err := RepositoryURLs(isPublic)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(rawURLs))
	for _, rawURL := range rawURLs {
		go func(rawURL string) {
			defer wg.Done()
			// defer log.Debug("Processed " + rawURL)
			cloneURL := githubURL(isPublic, rawURL)
			targetDir, err := getSrcPath(rawURL)
			if err != nil {
				log.WithError(err).Error("Cannot get source path for " + rawURL)
				return
			}

			repo, err := NewGitRepo(cloneURL, targetDir)
			if err != nil {
				log.WithError(err).Error("Cannot get git repo which targets ", targetDir, " with "+cloneURL)
				return
			}

			if !repo.CheckLocal() {
				log.Error("The directory " + targetDir + " does not exist. Run update to pull the latest repos.")
				return
			}
			if !repo.IsDirty() {
				return
			}
			log.Println("Committing " + targetDir + "...")
			err = repo.Commit(message)
			if err != nil {
				log.WithError(err).Error("Failed to commit.")
			}
			return
		}(rawURL)
	}
	wg.Wait()
	return nil
}

func Update(isPublic bool) error {
	rawURLs, err := RepositoryURLs(isPublic)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(rawURLs))
	for _, rawURL := range rawURLs {
		go func(rawURL string) {
			defer wg.Done()
			// defer log.Debug("Processed " + rawURL)
			cloneURL := githubURL(isPublic, rawURL)
			targetDir, err := getSrcPath(rawURL)
			if err != nil {
				log.WithError(err).Error("Cannot get source path for " + rawURL)
				return
			}

			if com.IsDir(targetDir) && !com.IsDir(filepath.Join(targetDir, ".git")) {
				return
			}

			repo, err := NewGitRepo(cloneURL, targetDir)
			if err != nil {
				log.WithError(err).Error("Cannot get git repo which targets ", targetDir, " with "+cloneURL)
				return
			}
			if repo.CheckLocal() {
				if err := repo.Update(); err != nil {
					log.WithError(err).Error("Cannot update ", targetDir, " with "+cloneURL)
					return
				}
			} else {
				if err := repo.Clone(); err != nil {
					log.WithError(err).Error("Cannot clone ", cloneURL, " into "+targetDir)
					return
				}
			}
			//cmd := exec.Command("go", "generate", targetDir)
			//cmd.CombinedOutput()
		}(rawURL)
	}
	wg.Wait()
	return nil
}

func Dirty(isPublic bool) error {
	rawURLs, err := RepositoryURLs(isPublic)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(rawURLs))
	for _, rawURL := range rawURLs {
		go func(rawURL string) {
			defer wg.Done()

			cloneURL := githubURL(isPublic, rawURL)
			targetDir, err := getSrcPath(rawURL)
			if err != nil {
				log.WithError(err).Error("Cannot get source path for " + rawURL)
				return
			}

			if com.IsDir(targetDir) && !com.IsDir(filepath.Join(targetDir, ".git")) {
				return
			}

			repo, err := NewGitRepo(cloneURL, targetDir)
			if err != nil {
				log.WithError(err).Error("Cannot get git repo which targets ", targetDir, " with "+cloneURL)
				return
			}

			if repo.IsDirty() {
				println(repo.LocalPath())
			}

		}(rawURL)
	}

	wg.Wait() // Dont return until all repos have been examined, and all messages have been printed

	return nil
}

func GoGet(isPublic bool) error {
	rawURLs, err := RepositoryURLs(isPublic)
	if err != nil {
		return err
	}

	importChan := make(chan string)
	done := make(chan struct{})

	go func() {
		importCache := map[string]bool{}
		var mutex sync.Mutex
		var wg sync.WaitGroup
		for impt := range importChan {
			wg.Add(1)
			func(impt string) {
				defer wg.Done()
				mutex.Lock()
				if _, ok := importCache[impt]; ok {
					mutex.Unlock()
					return
				}
				importCache[impt] = true
				mutex.Unlock()

				args := []string{
					"get",
					"-u",
				}
				if Verbose {
					args = append(args, "-v")
				}
				args = append(args, impt)
				log.Debug("Performing a go " + strings.Join(args, " "))
				cmd := exec.Command("go", args...)
				buf, err := cmd.CombinedOutput()
				if err != nil {
					log.WithError(err).Error("Failed to run go " + strings.Join(args, " "))
					return
				}
				log.Debugf(string(buf))
			}(impt)
		}
		wg.Wait()
		close(done)
	}()

	var wg sync.WaitGroup
	wg.Add(len(rawURLs))
	for _, rawURL := range rawURLs {
		go func(rawURL string) {
			defer wg.Done()

			srcDir, err := getSrcPath(rawURL)
			if err != nil {
				log.WithError(err).Error("Cannot get source path for " + rawURL)
				return
			}
			extraDependenciesFile := filepath.Join(srcDir, "dependencies")
			if com.IsFile(extraDependenciesFile) {
				buf, err := ioutil.ReadFile(extraDependenciesFile)
				if err == nil {
					deps := strings.Split(string(buf), "\n")
					for _, dep := range deps {
						dep := strings.TrimSpace(dep)
						if dep != "" {
							importChan <- dep
						}
					}
				}
			}

			fset := token.NewFileSet()
			files1, err := glob.Glob(filepath.Join(srcDir, "*.go"))
			if err != nil {
				files1 = []string{}
			}
			files2, err := glob.Glob(filepath.Join(srcDir, "**", "*.go"))
			if err != nil {
				files2 = []string{}
			}

			files := append(files1, files2...)
			for _, file := range files {
				file, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
				if err != nil {
					log.WithError(err).Error("Failed to parse directory " + srcDir)
					return
				}

				for _, impt := range file.Imports {
					path := strings.Trim(impt.Path.Value, `"`)
					if isStandardImport(path) {
						continue
					}
					if strings.Contains(impt.Path.Value, "rai-") {
						continue
					}
					importChan <- path
				}
			}
		}(rawURL)
	}

	wg.Wait() // Dont return until all repos have been examined, and all messages have been printed
	close(importChan)
	<-done
	return nil
}

func BumpVersion(isPublic bool) error {
	rawURLs, err := RepositoryURLs(isPublic)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(rawURLs))
	for _, rawURL := range rawURLs {
		go func(rawURL string) {
			defer wg.Done()
			targetDir, err := getSrcPath(rawURL)
			if err != nil {
				log.WithError(err).Error("Cannot get source path for " + rawURL)
				return
			}
			cfg := filepath.Join(targetDir, ".bumpversion.cfg")
			if !com.IsExist(cfg) {
				// log.Info(".bumpversion.cfg does not exist for " + rawURL)
				return
			}

			// bumpversion patch --commit && git push && git push --tags
			args := []string{
				"patch",
				"--commit",
			}
			cmd := exec.Command("bumpversion", args...)
			cmd.Dir = targetDir
			buf, err := cmd.CombinedOutput()
			if err != nil {
				log.WithError(err).Error("Failed to run bumpversion " + strings.Join(args, " ") + " in " + targetDir)
				return
			}
			log.Debugf(string(buf))

			args = []string{
				"push",
			}
			cmd = exec.Command("git", args...)
			cmd.Dir = targetDir
			buf, err = cmd.CombinedOutput()
			if err != nil {
				log.WithError(err).Error("Failed to run git " + strings.Join(args, " ") + " in " + targetDir)
				return
			}
			log.Debugf(string(buf))

			args = []string{
				"push",
				"--tags",
			}
			cmd = exec.Command("git", args...)
			cmd.Dir = targetDir
			buf, err = cmd.CombinedOutput()
			if err != nil {
				log.WithError(err).Error("Failed to run git " + strings.Join(args, " ") + " in " + targetDir)
				return
			}
			log.Debugf(string(buf))

		}(rawURL)
	}
	wg.Wait()
	return nil
}

// we assume that code will start with a domain name (dot in the first element).
func isStandardImport(path string) bool {
	i := strings.Index(path, "/")
	if i < 0 {
		i = len(path)
	}
	elem := path[:i]
	return !strings.Contains(elem, ".")
}

func getSrcPath(importPath string) (appPath string, err error) {
	paths := com.GetGOPATHs()
	for _, p := range paths {
		d := filepath.Join(p, "src", importPath)
		if com.IsExist(d) {
			appPath = d
			break
		}
	}

	if len(appPath) == 0 {
		appPath = filepath.Join(gopath, "src", importPath)
	}

	return appPath, nil
}
