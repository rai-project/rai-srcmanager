package srcmanager

// a lot of this is extracted from https://github.com/Masterminds/vcs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type GitRepo struct {
	remote, local  string
	RemoteLocation string
	Logger         *log.Logger
}

// Logger is where you can provide a logger, implementing the log.Logger interface,
// where verbose output from each VCS will be written. The default logger does
// not log data. To log data supply your own logger or change the output location
// of the provided logger.
var Logger *log.Logger

// Remote retrieves the remote location for a repo.
func (g *GitRepo) Remote() string {
	return g.remote
}

// LocalPath retrieves the local file system location for a repo.
func (g *GitRepo) LocalPath() string {
	return g.local
}

func (g *GitRepo) setRemote(remote string) {
	g.remote = remote
}

func (g *GitRepo) setLocalPath(local string) {
	g.local = local
}

func (g GitRepo) run(cmd string, args ...string) ([]byte, error) {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	g.log(out)
	if err != nil {
		err = fmt.Errorf("%s: %s", out, err)
	}
	return out, err
}

func mergeEnvLists(in, out []string) []string {
NextVar:
	for _, inkv := range in {
		k := strings.SplitAfterN(inkv, "=", 2)[0]
		for i, outkv := range out {
			if strings.HasPrefix(outkv, k) {
				out[i] = inkv
				continue NextVar
			}
		}
		out = append(out, inkv)
	}
	return out
}

func envForDir(dir string) []string {
	env := os.Environ()
	return mergeEnvLists([]string{"PWD=" + dir}, env)
}

func (g *GitRepo) RunFromDir(cmd string, args ...string) ([]byte, error) {
	c := exec.Command(cmd, args...)
	c.Dir = g.local
	c.Env = envForDir(c.Dir)
	out, err := c.CombinedOutput()
	return out, err
}

func (g *GitRepo) log(v interface{}) {
	g.Logger.Printf("%s", v)
}

// NewGitRepo creates a new instance of GitRepo. The remote and local directories
// need to be passed in.
func NewGitRepo(remote, local string) (*GitRepo, error) {
	ins := depInstalled("git")
	if !ins {
		return nil, NewLocalError("git is not installed", nil, "")
	}
	r := &GitRepo{}
	r.setRemote(remote)
	r.setLocalPath(local)
	r.RemoteLocation = "origin"
	r.Logger = Logger

	// Make sure the local Git repo is configured the same as the remote when
	// A remote value was passed in.
	if r.CheckLocal() == true {
		c := exec.Command("git", "config", "--get", "remote.origin.url")
		c.Dir = local
		c.Env = envForDir(c.Dir)
		out, err := c.CombinedOutput()
		if err != nil {
			return nil, NewLocalError("Unable to retrieve local repo information", err, string(out))
		}

		localRemote := strings.TrimSpace(string(out))
		if remote != "" && localRemote != remote {
			return nil, ErrWrongRemote
		}

		// If no remote was passed in but one is configured for the locally
		// checked out Git repo use that one.
		if remote == "" && localRemote != "" {
			r.setRemote(localRemote)
		}
	}

	return r, nil
}

// Get is used to perform an initial clone of a repository.
func (s *GitRepo) Clone() error {
	out, err := s.run("git", "clone", s.Remote(), s.LocalPath())

	// There are some windows cases where Git cannot create the parent directory,
	// if it does not already exist, to the location it's trying to create the
	// repo. Catch that error and try to handle it.
	if err != nil && s.isUnableToCreateDir(err) {

		basePath := filepath.Dir(filepath.FromSlash(s.LocalPath()))
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			err = os.MkdirAll(basePath, 0755)
			if err != nil {
				return NewLocalError("Unable to create directory", err, "")
			}

			out, err = s.run("git", "clone", s.Remote(), s.LocalPath())
			if err != nil {
				return NewRemoteError("Unable to get repository", err, string(out))
			}
			return err
		}

	} else if err != nil {
		return NewRemoteError("Unable to get repository", err, string(out))
	}

	return nil
}

// Init initializes a git repository at local location.
func (s *GitRepo) Init() error {
	out, err := s.run("git", "init", s.LocalPath())

	// There are some windows cases where Git cannot create the parent directory,
	// if it does not already exist, to the location it's trying to create the
	// repo. Catch that error and try to handle it.
	if err != nil && s.isUnableToCreateDir(err) {

		basePath := filepath.Dir(filepath.FromSlash(s.LocalPath()))
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			err = os.MkdirAll(basePath, 0755)
			if err != nil {
				return NewLocalError("Unable to initialize repository", err, "")
			}

			out, err = s.run("git", "init", s.LocalPath())
			if err != nil {
				return NewLocalError("Unable to initialize repository", err, string(out))
			}
			return nil
		}

	} else if err != nil {
		return NewLocalError("Unable to initialize repository", err, string(out))
	}

	return nil
}

// Commit performs an Git commit.
func (s *GitRepo) Commit(message string) error {
	if !s.IsDirty() {
		return nil
	}
	out, err := s.RunFromDir("git", "add", ".")
	if err != nil {
		return NewRemoteError("Unable to add files to repository", err, string(out))
	}
	out, err = s.RunFromDir("git", "commit", "-a", "-m", message)
	if err != nil {
		return NewRemoteError("Unable to commit to repository", err, string(out))
	}
	out, err = s.RunFromDir("git", "push")
	if err != nil {
		return NewRemoteError("Unable to push commits", err, string(out))
	}
	return nil
}

// Update performs an Git fetch and pull to an existing checkout.
func (s *GitRepo) Update() error {
	// Perform a fetch to make sure everything is up to date.
	out, err := s.RunFromDir("git", "fetch", s.RemoteLocation, "--prune")
	if err != nil {
		return NewRemoteError("Unable to update repository", err, string(out))
	}

	// When in a detached head state, such as when an individual commit is checked
	// out do not attempt a pull. It will cause an error.
	detached, err := isDetachedHead(s.LocalPath())
	if err != nil {
		return NewLocalError("Unable to update repository", err, "")
	}

	if detached == true {
		return nil
	}

	out, err = s.RunFromDir("git", "pull", "--rebase", "origin", "master")
	if err != nil {
		return NewRemoteError("Unable to update repository", err, string(out))
	}
	/*
		files, err := glob.Glob(filepath.Join(s.LocalPath(), "*.go"))
		if err == nil && len(files) > 0 {
			out, err = s.RunFromDir("go", "generate", ".")
			if err != nil {
				return NewRemoteError("Unable to perform a go generate on repository repository", err, string(out))
			}
		}
	*/
	s.Logger.Info("Updated " + s.LocalPath())
	return nil
}

// CheckLocal verifies the local location is a Git repo.
func (s *GitRepo) CheckLocal() bool {
	if _, err := os.Stat(s.LocalPath() + "/.git"); err == nil {
		return true
	}

	return false
}

// isDetachedHead will detect if git repo is in "detached head" state.
func isDetachedHead(dir string) (bool, error) {
	p := filepath.Join(dir, ".git", "HEAD")
	contents, err := ioutil.ReadFile(p)
	if err != nil {
		return false, err
	}

	contents = bytes.TrimSpace(contents)
	if bytes.HasPrefix(contents, []byte("ref: ")) {
		return false, nil
	}

	return true, nil
}

const longForm = "2006-01-02 15:04:05 -0700"

// Date retrieves the date on the latest commit.
func (s *GitRepo) Date() (time.Time, error) {
	out, err := s.RunFromDir("git", "log", "-1", "--date=iso", "--pretty=format:%cd")
	if err != nil {
		return time.Time{}, NewLocalError("Unable to retrieve revision date", err, string(out))
	}
	t, err := time.Parse(longForm, string(out))
	if err != nil {
		return time.Time{}, NewLocalError("Unable to retrieve revision date", err, string(out))
	}
	return t, nil
}

// IsDirty returns if the checkout has been modified from the checked
// out reference.
func (s *GitRepo) IsDirty() bool {
	out, err := s.RunFromDir("git", "diff")
	if err != nil {
		log.Fatal(err)
	}
	return (err != nil) || (len(out) != 0)
}

// isUnableToCreateDir checks for an error in Init() to see if an error
// where the parent directory of the VCS local path doesn't exist. This is
// done in a multi-lingual manner.
func (s *GitRepo) isUnableToCreateDir(err error) bool {
	msg := err.Error()
	if strings.HasPrefix(msg, "could not create work tree dir") ||
		strings.HasPrefix(msg, "不能创建工作区目录") ||
		strings.HasPrefix(msg, "no s'ha pogut crear el directori d'arbre de treball") ||
		strings.HasPrefix(msg, "impossible de créer le répertoire de la copie de travail") ||
		strings.HasPrefix(msg, "kunde inte skapa arbetskatalogen") ||
		(strings.HasPrefix(msg, "Konnte Arbeitsverzeichnis") && strings.Contains(msg, "nicht erstellen")) ||
		(strings.HasPrefix(msg, "작업 디렉터리를") && strings.Contains(msg, "만들 수 없습니다")) {
		return true
	}

	return false
}
func depInstalled(name string) bool {
	if _, err := exec.LookPath(name); err != nil {
		return false
	}

	return true
}

func init() {
	// Initialize the logger to one that does not actually log anywhere. This is
	// to be overridden by the package user by setting vcs.Logger to a different
	// logger.
	Logger = log.New()
	Logger.Level = log.DebugLevel
}
