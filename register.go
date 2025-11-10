package zkregister

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func RegisterZK(ctx context.Context, conn *zk.Conn, l *slog.Logger, service ServiceInfo) {
	instance := GetInstance(service)
	go register(ctx, conn, l, instance)
}

func register(ctx context.Context, conn *zk.Conn, l *slog.Logger, instance Instance) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	servicePath := "/services"
	parentPath := "/services/" + instance.Name
	path := parentPath + "/" + instance.ID
	for {
		select {
		case <-ctx.Done():
			l.Info("register context done")
			err := conn.Delete(path, -1)
			if err != nil {
				l.Error("delete instance path failed", "err", err)
			}

			children, _, _ := conn.Children(parentPath)
			if len(children) == 0 {
				err := conn.Delete(parentPath, -1)
				if err != nil {
					l.Error("delete parent path failed", "err", err)
				}
			}

			return
		case <-ticker.C:
			l.Info("register instance to zookeeper")

			if exists, _, _ := conn.Exists(servicePath); !exists {
				_, err := conn.Create(servicePath, nil, 0, zk.WorldACL(zk.PermAll))
				if err != nil {
					l.Error("create service path failed", "err", err)
					continue
				}
			}

			exists, _, _ := conn.Exists(parentPath)
			if !exists {
				_, err := conn.Create(parentPath, nil, 0, zk.WorldACL(zk.PermAll))
				if err != nil {
					l.Error("create parent path failed", "err", err)
					continue
				}
			}

			instance.RegistrationTimeUTC = time.Now().UnixMilli()
			data, err := json.Marshal(instance)
			if err != nil {
				l.Error("marshal instance failed", "err", err)
				continue
			}

			exists, _, _ = conn.Exists(path)
			if exists {
				_, err = conn.Set(path, data, -1)
				if err != nil {
					l.Error("set instance path failed", "err", err)
					continue
				}
				continue
			}
			_, err = conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
			if err != nil {
				l.Error("create instance path failed", "err", err)
				continue
			}
		}
	}
}
