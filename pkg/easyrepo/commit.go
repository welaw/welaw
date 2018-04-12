package easyrepo

import (
	git "gopkg.in/libgit2/git2go.v25"
	//git "gopkg.in/libgit2/git2go.v26"
	//git "github.com/libgit2/git2go"
)

const (
	createRepoMsg   = "Create Repo"
	createBranchMsg = "Create Branch"
)

type Commit struct {
	Parent   string
	Filename string
	Body     []byte
	Msg      string
	Sig      *git.Signature
}

func CreateCommit(repo *git.Repository, branch string, c *Commit) (*git.Oid, error) {
	// get parent commit
	parent, err := GetHead(repo, branch)
	if err != nil {
		return nil, err
	}
	ie, err := entryForBody(repo, c.Body, c.Filename)
	if err != nil {
		return nil, err
	}
	idx, err := repo.Index()
	if err != nil {
		return nil, err
	}
	idx.Add(ie)
	err = idx.Write()
	if err != nil {
		return nil, err
	}
	treeId, err := idx.WriteTree()
	if err != nil {
		return nil, err
	}
	tree, err := repo.LookupTree(treeId)
	if err != nil {
		return nil, err
	}
	defer tree.Free()
	commitId, err := repo.CreateCommit("refs/heads/"+branch, c.Sig, c.Sig, c.Msg, tree, parent)
	if err != nil {
		return nil, err
	}
	return commitId, nil
}

func CreateFirstCommit(repo *git.Repository, c *Commit) (*git.Oid, error) {
	var err error
	var ie *git.IndexEntry
	if ie, err = entryForBody(repo, c.Body, c.Filename); err != nil {
		return nil, err
	}
	var idx *git.Index
	if idx, err = repo.Index(); err != nil {
		return nil, err
	}
	idx.Add(ie)
	if err = idx.Write(); err != nil {
		return nil, err
	}
	var treeId *git.Oid
	if treeId, err = idx.WriteTree(); err != nil {
		return nil, err
	}
	var tree *git.Tree
	if tree, err = repo.LookupTree(treeId); err != nil {
		return nil, err
	}
	defer tree.Free()
	oid, err := repo.CreateCommit("refs/heads/master", c.Sig, c.Sig, c.Msg, tree)
	return oid, err
}

func GetCommit(repo *git.Repository, hash string) (*git.Commit, error) {
	oid, err := git.NewOid(hash)
	if err != nil {
		return nil, err
	}
	return repo.LookupCommit(oid)
}

func GetEasyCommitByHash(repo *git.Repository, hash, filename string) (*Commit, error) {
	oid, err := git.NewOid(hash)
	if err != nil {
		return nil, err
	}
	c, err := repo.LookupCommit(oid)
	if err != nil {
		return nil, err
	}
	msg := c.Message()
	tree, err := c.Tree()
	if err != nil {
		return nil, err
	}
	treeEntry := tree.EntryByName(filename)
	blob, err := repo.LookupBlob(treeEntry.Id)
	if err != nil {
		return nil, err
	}
	return &Commit{
		Msg:  msg,
		Body: blob.Contents(),
	}, nil
}

func GetEasyCommit(repo *git.Repository, c *git.Commit, filename string) (*Commit, error) {
	msg := c.Message()
	//odb, err := repo.Odb()
	tree, err := c.Tree()
	if err != nil {
		return nil, err
	}
	treeEntry := tree.EntryByName(filename)
	blob, err := repo.LookupBlob(treeEntry.Id)
	if err != nil {
		return nil, err
	}
	return &Commit{
		Msg:  msg,
		Body: blob.Contents(),
	}, nil
}

func entryForBody(repo *git.Repository, body []byte, filename string) (ie *git.IndexEntry, err error) {
	oid, err := repo.CreateBlobFromBuffer(body)
	if err != nil {
		return
	}
	return &git.IndexEntry{
		Mode: git.FilemodeBlob,
		Id:   oid,
		Path: filename,
	}, nil
}

// DeleteCommit will delete the given hash
//func (lr *lawRepo) DeleteCommit(branch, hash string) error {
//var err error

//oid, err := git.NewOid(hash)
//if err != nil {
//return err
//}
//_, err = lr.repo.LookupCommit(oid)
//if err != nil {
//return err
//}
// get parent commit
//opts, err := git.DefaultRebaseOptions()
//repo.InitRebase()
//return nil
