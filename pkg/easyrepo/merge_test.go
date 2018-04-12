package easyrepo

//git "github.com/libgit2/git2go"

//git "github.com/libgit2/git2go"

//func TestMergeBranches(t *testing.T) {
//dir, err := ioutil.TempDir("", "testrepo")
//require.NoError(t, err)
////defer os.RemoveAll(dir)

//repo, err := CreateRepo(dir)
//require.NoError(t, err)
//require.NotNil(t, repo)

//c1, err := CreateFirstCommit(repo, &Commit{
//Filename: "test.txt",
//Body:     []byte("test body\n"),
//Msg:      "Initial commit",
//Sig: &git.Signature{
//Name:  "testuser",
//Email: "testuser@welaw.org",
//When:  time.Now(),
//},
//})
//require.NoError(t, err)
//require.NotNil(t, c1)

//// create our branch
//br, err := CreateBranch(repo, c1.String(), "our-branch")
//require.NoError(t, err)
//require.NotNil(t, br)

//c2, err := CreateCommit(repo, "our-branch", &Commit{
//Parent:   c1.String(),
//Filename: "test.txt",
//Msg:      "our branch commit message",
//Body: []byte(`this is a first commit
//on our branch and
//this is last line
//`),
//Sig: &git.Signature{
//Name:  "testuser",
//Email: "testuser@welaw.org",
//When:  time.Now(),
//},
//})
//require.NoError(t, err)
//require.NotNil(t, c2)

//// create their branch
//br, err = CreateBranch(repo, c2.String(), "their-branch")
//require.NoError(t, err)
//require.NotNil(t, br)

//c3, err := CreateCommit(repo, "their-branch", &Commit{
//Parent:   c2.String(),
//Filename: "test.txt",
//Msg:      "their branch commit message",
//Body: []byte(`now this is first on
//their branch and
//this is last line
//`),
//Sig: &git.Signature{
//Name:  "testuser",
//Email: "testuser@welaw.org",
//When:  time.Now(),
//},
//})
//require.NoError(t, err)
//require.NotNil(t, c3)

//// ------
////c4, err := CreateCommit(repo, "their-branch", &Commit{
////Parent:   c3.String(),
////Filename: "test.txt",
////Msg:      "their 2nd branch commit message",
////Body:     []byte(`another`),
////Sig: &git.Signature{
////Name:  "testuser",
////Email: "testuser@welaw.org",
////When:  time.Now(),
////},
////})
////require.NoError(t, err)
////require.NotNil(t, c4)
//// ------

//// merge them
//fmt.Printf("dir: %v\n", dir)
//hash, err := MergeBranches(repo, "their-branch", "our-branch")
//require.NoError(t, err)
//require.NotNil(t, hash)
//}

//func TestMergeCommits(t *testing.T) {
//dir, err := ioutil.TempDir("", "testrepo")
//require.NoError(t, err)
////defer os.RemoveAll(dir)

//repo, err := CreateRepo(dir)
//require.NoError(t, err)
//require.NotNil(t, repo)

//c1, err := CreateFirstCommit(repo, &Commit{
//Filename: "test.txt",
//Body:     []byte("test body"),
//Msg:      "commit message",
//Sig: &git.Signature{
//Name:  "testuser",
//Email: "testuser@welaw.org",
//When:  time.Now(),
//},
//})
//require.NoError(t, err)
//require.NotNil(t, c1)

//m2, err := CreateCommit(repo, "master", &Commit{
//Parent:   c1.String(),
//Filename: "test.txt",
//Body:     []byte("second on master\n"),
//Msg:      "second commit on master branch",
//Sig: &git.Signature{
//Name:  "testuser",
//Email: "testuser@welaw.org",
//When:  time.Now(),
//},
//})
//require.NoError(t, err)
//require.NotNil(t, m2)

//br, err := CreateBranch(repo, c1.String(), "test-branch")
//require.NoError(t, err)
//require.NotNil(t, br)

//c2, err := CreateCommit(repo, "test-branch", &Commit{
//Parent:   c1.String(),
//Filename: "test.txt",
//Body:     []byte("changed\n"),
//Msg:      "first commit on branch",
//Sig: &git.Signature{
//Name:  "testuser",
//Email: "testuser@welaw.org",
//When:  time.Now(),
//},
//})
//require.NoError(t, err)
//require.NotNil(t, c2)

//com1, err := repo.LookupCommit(m2)
//require.NoError(t, err)

//com2, err := repo.LookupCommit(c2)
//require.NoError(t, err)

//mergeOpts, err := git.DefaultMergeOptions()
//require.NoError(t, err)
//mc, err := repo.MergeCommits(com1, com2, &mergeOpts)
////mc, err := repo.MergeTrees(com1, com2, &mergeOpts)
//require.NoError(t, err)
//require.False(t, mc.HasConflicts())

//treeId, err := mc.WriteTreeTo(repo)
////err = mc.Write()
//require.NoError(t, err)

////err = mc.Write()
////require.NoError(t, err)

////treeId, err := mc.WriteTree()
////require.NoError(t, err)

//tree, err := repo.LookupTree(treeId)
//require.NoError(t, err)
//defer tree.Free()

//sig := com2.Author()
//_, err = repo.CreateCommit("refs/heads/test-branch", sig, sig, "merge commit message", tree, com2)
//require.NoError(t, err)
//}
