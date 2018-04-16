package services

import (
	"context"

	"github.com/welaw/welaw/pkg/errs"
	"github.com/welaw/welaw/pkg/permissions"
	"github.com/welaw/welaw/proto"
)

func (svc service) GetServerStats(ctx context.Context) (*proto.ServerStats, error) {
	stats, err := svc.db.GetServerStats()
	if err != nil {
		return nil, err
	}
	return stats, nil

}

func (svc service) LoadRepos(ctx context.Context, opts *proto.LoadReposOptions) (*proto.LoadReposReply, error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, errs.ErrUnauthorized
	}
	perm, err := svc.hasPermission(uid, permissions.OpReposLoad, nil)
	if err != nil {
		return nil, err
	}
	if !perm {
		return nil, errs.ErrUnauthorized
	}
	// rm local repos dir
	// copy files from storage to local
	return nil, nil
}

// copy local repos to storage
func (svc service) SaveRepos(ctx context.Context, opts *proto.SaveReposOptions) (*proto.SaveReposReply, error) {
	uid, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, errs.ErrUnauthorized
	}
	if perm, err := svc.hasPermission(uid, permissions.OpReposSave, nil); err != nil {
		return nil, err
	} else if !perm {
		return nil, errs.ErrUnauthorized
	}

	switch {
	case opts == nil:
		return nil, errs.BadRequest("opts not found")
	case opts.ReqType == proto.SaveReposOptions_LOAD:
		return svc.loadRepos(ctx, opts)
	case opts.ReqType == proto.SaveReposOptions_SAVE:
		return svc.saveRepos(ctx, opts)
	}
	return nil, errs.BadRequest("unknown opts")
}

func (svc service) loadRepos(ctx context.Context, opts *proto.SaveReposOptions) (*proto.SaveReposReply, error) {
	//zipfile := "repos.zip"
	//r, err := svc.storageClient.Bucket(svc.Opts.DefaultBucketName).Object(zipfile).NewReader(ctx)
	//if err != nil {
	//return nil, err
	//}
	//defer r.Close()

	//outdir := "."
	//outFilePath := filepath.Join(outdir, zipfile)

	//data, err := ioutil.ReadAll(r)
	//if err != nil {
	//return nil, err
	//}

	//err = ioutil.WriteFile(outFilePath, data, 0644)
	//if err != nil {
	//return nil, err
	//}

	//progress := func(archivePath string) {
	//fmt.Println(archivePath)
	//}
	//reposDir := "."
	//err = zip.UnarchiveFile(outFilePath, reposDir, progress)
	//if err != nil {
	//return nil, err
	//}

	return nil, nil
}

func (svc service) saveRepos(ctx context.Context, opts *proto.SaveReposOptions) (*proto.SaveReposReply, error) {
	//zipfile := "repos.zip"
	//outdir := "."
	//// tmp dir
	//outFilePath := filepath.Join(outdir, zipfile)
	//progress := func(archivePath string) {
	//fmt.Println(archivePath)
	//}
	//err := zip.ArchiveFile(svc.vc.GetPath(), outFilePath, progress)
	//if err != nil {
	//return nil, err
	//}

	//f, err := ioutil.ReadFile(outFilePath)
	//if err != nil {
	//return nil, err
	//}

	//w := svc.storageClient.Bucket(svc.Opts.DefaultBucketName).Object(zipfile).NewWriter(ctx)
	//w.ObjectAttrs = storage.ObjectAttrs{
	//Bucket:      svc.Opts.DefaultBucketName,
	//ContentType: "application/zip",
	//Name:        zipfile,
	//}

	//_, err = w.Write(f)
	//if err != nil {
	//return nil, err
	//}
	//err = w.Close()
	//if err != nil {
	//return nil, err
	//}

	//err = os.Remove(outFilePath)
	//if err != nil {
	//return nil, err
	//}

	return &proto.SaveReposReply{}, nil
}
