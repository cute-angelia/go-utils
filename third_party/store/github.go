package store

import (
	"github.com/google/go-github/github"
	"context"
	"golang.org/x/oauth2"
	"io/ioutil"
	"time"
	"fmt"
)

type GithubInterface interface {
	// 创建文件
	CreateFile(file []byte, destpath string, sha *string) error
	// 获取文件信息
	GetFile(destpath string) (*github.RepositoryContent, error)
	// 更新文件
	UpdateFile(file []byte, destpath string) error
	// 删除文件
	DeleteFile(filepath string) error

	GetFileByte(filepath string) ([]byte, error)
}

func NewGithub(config GithubConfig) *GithubApi {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GithubApi{
		Token:  config.Token,
		Owner:  config.Owner,
		Repo:   config.Repo,
		Branch: config.Branch,
		Client: client,
		Ctx:    ctx,
	}
}

type GithubApi struct {
	Token  string
	Owner  string
	Repo   string
	Branch string
	Client *github.Client
	Ctx    context.Context
}

func (self *GithubApi) CreateFile(file []byte, destpath string, sha *string) error {
	opts := &github.RepositoryContentFileOptions{
		Message:   github.String(time.Now().Format("2006-01-02 15:04:05")),
		Content:   file,
		Branch:    github.String(self.Branch),
		Committer: &github.CommitAuthor{Name: github.String("a ghost"), Email: github.String("ghost@ghost.com")},
	}
	//_, _, err := client.Repositories.DeleteFile(ctx, OWNER, REPO, filepath, opts)

	err := fmt.Errorf("")

	if sha != nil {
		opts.SHA = sha
		_, _, err = self.Client.Repositories.UpdateFile(self.Ctx, self.Owner, self.Repo, destpath, opts)
	} else {
		_, _, err = self.Client.Repositories.CreateFile(self.Ctx, self.Owner, self.Repo, destpath, opts)
	}

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (self *GithubApi) GetFile(destpath string) (*github.RepositoryContent, error) {
	getOpts := &github.RepositoryContentGetOptions{
		Ref: self.Branch,
	}
	if respContent, _, _, err := self.Client.Repositories.GetContents(self.Ctx, self.Owner, self.Repo, destpath, getOpts); err != nil {
		return nil, err
	} else {
		return respContent, nil
	}
}

func (self *GithubApi) UpdateFile(file []byte, destpath string) error {
	if c, err := self.GetFile(destpath); err != nil {
		return err
	} else {
		sha := c.GetSHA()

		if err := self.CreateFile(file, destpath, &sha); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (self *GithubApi) DeleteFile(filepath string) error {
	if c, err := self.GetFile(filepath); err != nil {
		return err
	} else {
		sha := c.GetSHA()

		opts := &github.RepositoryContentFileOptions{
			Message:   github.String(time.Now().Format("2006-01-02 15:04:05")),
			Branch:    github.String(self.Branch),
			SHA:       &sha,
			Committer: &github.CommitAuthor{Name: github.String("a ghost"), Email: github.String("ghost@ghost.com")},
		}

		if _, _, err := self.Client.Repositories.DeleteFile(self.Ctx, self.Owner, self.Repo, filepath, opts); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (self *GithubApi) GetFileByte(filepath string) ([]byte, error) {
	if fileContent, err := ioutil.ReadFile(filepath); err == nil {
		return nil, err
	} else {
		return fileContent, nil
	}
}
