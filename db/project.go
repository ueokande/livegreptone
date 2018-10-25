package db

import (
	"context"
	"encoding/json"

	"github.com/coreos/etcd/clientv3"
	"github.com/ueokande/livegreptone"
)

func (d *model) GetProject(ctx context.Context, id string) (livegreptone.Project, error) {
	key := ProjectKey(id)
	resp, err := d.etcd.Get(ctx, key)
	if err != nil {
		return livegreptone.Project{}, err
	}
	if resp.Count == 0 {
		return livegreptone.Project{}, ErrNotFound
	}

	var p livegreptone.Project
	err = json.Unmarshal([]byte(resp.Kvs[0].Value), &p)
	if err != nil {
		return livegreptone.Project{}, err
	}
	return p, nil
}

func (d *model) GetProjectIds(ctx context.Context) ([]string, error) {
	resp, err := d.etcd.Get(ctx, ProjectKeyPrefix,
		clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(resp.Kvs))
	for i, kv := range resp.Kvs {
		ids[i] = string(kv.Key)[len(ProjectKeyPrefix):]
	}
	return ids, nil
}

func (d *model) UpdateProject(ctx context.Context, project livegreptone.Project) error {
	key := ProjectKey(project.Id)
	value, err := json.Marshal(project)
	if err != nil {
		return err
	}
	_, err = d.etcd.Put(ctx, key, string(value))
	if err != nil {
		return err
	}
	return nil
}

func (d *model) RemoveProject(ctx context.Context, id string) error {
	key := ProjectKey(id)
	_, err := d.etcd.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
