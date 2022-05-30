package impl

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/staryjie/restful-api-demo/apps/host"
)

// 完成对象和数据库之间的转换

func (i *HostServiceImpl) save(ctx context.Context, ins *host.Host) error {
	// 默认值填充
	ins.InjectDefault()

	var (
		err error
	)

	// 把数据入库到 resource表和host表
	// 一次需要往2个表录入数据, 我们需要2个操作 要么都成功，要么都失败, 事务的逻辑
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Start TX error, %s", err)
	}

	// Defer处理事务的提交方式
	//   1.无错误，则Commit事务
	//   2.有错误，则Rollback事务
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				i.l.Errorf("Rollback error, %s", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				i.l.Errorf("Commit error, %s", err)
			}
		}
	}()

	// 插入Resource数据
	rstmt, err := tx.PrepareContext(ctx, InsertResourceSQL)
	if err != nil {
		return err
	}
	defer rstmt.Close()

	_, err = rstmt.ExecContext(ctx, ins.Id, ins.Vendor, ins.Region, ins.CreateAt, ins.ExpireAt,
		ins.Type, ins.Name, ins.Description, ins.Status, ins.UpdateAt,
		ins.SyncAt, ins.Account, ins.PublicIP, ins.PrivateIP)
	if err != nil {
		return err
	}

	// 插入Describe数据
	dstmt, err := tx.PrepareContext(ctx, InsertDescribeSQL)
	if err != nil {
		return err
	}
	defer dstmt.Close()

	_, err = dstmt.ExecContext(ctx, ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec,
		ins.OSType, ins.OSName, ins.SerialNumber)
	if err != nil {
		return err
	}

	return nil
}

func (i *HostServiceImpl) update(ctx context.Context, ins *host.Host) error {
	var (
		err       error
		resStmt   *sql.Stmt
		hostStemt *sql.Stmt
	)
	// 开始事务
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer处理事务的提交方式
	//   1.无错误，则Commit事务
	//   2.有错误，则Rollback事务
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				i.l.Errorf("Rollback error, %s", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				i.l.Errorf("Commit error, %s", err)
			}
		}
	}()

	// 更新 Resource 表
	resStmt, err = tx.PrepareContext(ctx, updateResourceSQL)
	if err != nil {
		return err
	}

	// 执行SQL
	_, err = resStmt.ExecContext(ctx, ins.Vendor, ins.Region, ins.ExpireAt, ins.Name, ins.Description, ins.Id)
	if err != nil {
		return err
	}

	// 更新Host表
	hostStemt, err = tx.PrepareContext(ctx, updateHostSQL)
	if err != nil {
		return err
	}

	// 执行SQL
	_, err = hostStemt.ExecContext(ctx, ins.CPU, ins.Memory, ins.Id)
	if err != nil {
		return err
	}

	return nil
}

func (i *HostServiceImpl) delete(ctx context.Context, ins *host.Host) error {
	var (
		err       error
		resStmt   *sql.Stmt
		hostStemt *sql.Stmt
	)

	// 初始化一个事务，所有的操作都在这个事务中执行
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				i.l.Errorf("Rollback error, %s", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				i.l.Errorf("Commit error, %s", err)
			}
		}
	}()

	// 删除Resource表中的数据
	resStmt, err = tx.PrepareContext(ctx, deleteResourceSQl)
	if err != nil {
		return err
	}

	defer resStmt.Close()

	_, err = resStmt.ExecContext(ctx, ins.Id)
	if err != nil {
		return err
	}

	// 删除Host表中的数据
	hostStemt, err = tx.PrepareContext(ctx, deleteHostSQL)
	if err != nil {
		return err
	}

	defer hostStemt.Close()

	_, err = hostStemt.ExecContext(ctx, ins.Id)
	if err != nil {
		return err
	}

	return nil
}
