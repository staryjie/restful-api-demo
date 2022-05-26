package impl

import (
	"context"
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
	rstmt, err := tx.Prepare(InsertResourceSQL)
	if err != nil {
		return err
	}
	defer rstmt.Close()

	_, err = rstmt.Exec(ins.Id, ins.Vendor, ins.Region, ins.CreateAt, ins.ExpireAt,
		ins.Type, ins.Name, ins.Description, ins.Status, ins.UpdateAt,
		ins.SyncAt, ins.Account, ins.PublicIP, ins.PrivateIP)
	if err != nil {
		return err
	}

	// 插入Describe数据
	dstmt, err := tx.Prepare(InsertDescribeSQL)
	if err != nil {
		return err
	}
	defer dstmt.Close()

	_, err = dstmt.Exec(ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec,
		ins.OSType, ins.OSName, ins.SerialNumber)
	if err != nil {
		return err
	}

	return nil
}
