package task

import (
	"database/sql"
	"strconv"
	"userService/pkg/camunda"
	"userService/pkg/camunda/pb"
	"userService/pkg/common"
	camundamodel "userService/pkg/model/camunda"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	addIns               = "add_ins"
	addMcht              = "add_mcht"
	addTerm              = "add_term"
	deleteIns            = "delete_ins"
	deleteMcht           = "delete_mcht"
	deleteTerm           = "delete_term"
	updateIns            = "update_ins"
	cancelUpdateIns      = "cancel_update_ins"
	updateMcht           = "update_mcht"
	cancelUpdateMcht     = "cancel_update_mcht"
	insUnregister        = "ins_unregister"
	cancelInsUnregister  = "cancel_ins_unregister"
	freezeMcht           = "freeze_mcht"
	cancelFreezeMcht     = "cancel_freeze_mcht"
	unfreezeMcht         = "unfreeze_mcht"
	cancelUnfreezeMcht   = "cancel_unfreeze_mcht"
	mchtUnregister       = "mcht_unregister"
	cancelMchtUnregister = "cancel_mcht_unregister"
	termUnregister       = "term_unregister"
	cancelTermUnregister = "cancel_term_unregister"
)

var topics = []string{
	addIns,
	addMcht,
	addTerm,
	deleteIns,
	deleteMcht,
	deleteTerm,
	updateIns,
	cancelUpdateIns,
	updateMcht,
	cancelUpdateMcht,
	insUnregister,
	cancelInsUnregister,
	freezeMcht,
	cancelFreezeMcht,
	unfreezeMcht,
	cancelUnfreezeMcht,
	mchtUnregister,
	cancelMchtUnregister,
	termUnregister,
	cancelTermUnregister,
}

func RunServiceTask(format string, workerNum int) {
	ch := make(chan int, workerNum)
	ctx := context.TODO()
	for i := 0; i < workerNum; i++ {
		go finishRegister(ctx, i, ch)
	}
	c := cron.New()
	_ = c.AddFunc(format, func() {
		ch <- 1
	})
	c.Start()
}

func finishRegister(ctx context.Context, workerId int, ch <-chan int) {
	id := strconv.Itoa(workerId)
	for {
		<-ch
		client := camunda.Get()
		for _, t := range topics {
			func(topic string) {
				db := common.DB.Begin()
				defer db.Rollback()
				logrus.Debugln("拉取外部任务")
				resp, err := client.ExternalTask.FetchAndLock(ctx, &pb.FetchAndLockExternalTaskReq{
					WorkerId: id,
					MaxTasks: 1,
					Topics: []*pb.Topic{
						{
							TopicName:    topic,
							LockDuration: 100000,
						},
					},
				})
				if err != nil {
					logrus.Errorln(err)
					return
				}
				if camunda.CheckError(resp) {
					logrus.Errorln(camunda.TransError(resp))
					return
				}
				if len(resp.Item) == 0 {
					return
				}

				item := resp.Item[0]
				// 查询机构id
				instance, err := camundamodel.FindProcessInstanceByCamundaInstanceId(db, item.ProcessInstanceId)
				if err != nil {
					logrus.Errorln(err)
					return
				}
				if instance == nil {
					logrus.Errorln("process %s not found", item.ProcessInstanceId)
					return
				}
				switch topic {
				case addIns:
					// 机构注册
					err = institutionRegister(db, instance)
				case addMcht:
					// 商户注册
					err = merchantRegister(db, instance)
				case addTerm:
					// 终端新增
					err = termAdd(db, instance)
				case deleteIns:
					// 机构删除
					err = deleteInstitution(db, instance)
				case deleteMcht:
					// 删除商户
					err = deleteMerchant(db, instance)
				case deleteTerm:
					err = termDelete(db, instance)
				case updateIns:
					// 更新机构
					err = institutionUpdate(db, instance)
				case cancelUpdateIns:
					// 取消更新机构
					err = institutionUpdateCancel(db, instance)
				case updateMcht:
					// 更新商户
					err = merchantUpdate(db, instance)
				case cancelUpdateMcht:
					// 取消更新商户
					err = merchantUpdateCancel(db, instance)
				case insUnregister:
					// 机构注销
					err = institutionUnRegister(db, instance)
				case cancelInsUnregister:
					// 取消机构注销
					err = institutionCancelUnRegister(db, instance)
				case freezeMcht:
					// 商户冻结
					err = merchantFreeze(db, instance)
				case cancelFreezeMcht:
					// 取消商户冻结
					err = cancelMerchantFreeze(db, instance)
				case unfreezeMcht:
					// 商户解冻
					err = merchantUnFreeze(db, instance)
				case cancelUnfreezeMcht:
					// 取消商户解冻
					err = cancelMerchantUnFreeze(db, instance)
				case mchtUnregister:
					// 商户注销
					err = merchantUnregister(db, instance)
				case cancelMchtUnregister:
					// 取消商户注销
					err = cancelMerchantUnregister(db, instance)
				case termUnregister:
					// 终端注销
					err = termInfoUnregister(db, instance)
				case cancelTermUnregister:
					// 取消终端注销
					err = cancelTermInfoUnregister(db, instance)
				}
				if err != nil {
					logrus.Errorln(err)
					return
				}
				// 将工作流最终状态修改为已完成
				err = camundamodel.UpdateProcessInstance(db, instance.InstanceId, &camundamodel.ProcessInstance{
					EndFlag: sql.NullInt64{
						Int64: 1,
						Valid: true,
					},
				})
				if err != nil {
					logrus.Errorln(err)
					return
				}

				completeResp, err := client.ExternalTask.Complete(ctx, &pb.CompleteExternalTaskReq{
					Id: resp.Item[0].Id,
					Body: &pb.CompleteExternalTaskReqBody{
						WorkerId: id,
					},
				})
				if err != nil {
					logrus.Errorln(err)
					return
				}
				if camunda.CheckError(completeResp) {
					logrus.Errorln(camunda.TransError(completeResp))
					return
				}
				db.Commit()
				logrus.Infoln("任务完成: ", resp.Item[0].Id, ", topic:", topic)

			}(t)
		}
	}
}
