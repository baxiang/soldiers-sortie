package workflow

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"userService/pkg/camunda"
	camundapb "userService/pkg/camunda/pb"
	"userService/pkg/common"
	camundamodel "userService/pkg/model/camunda"
	insmodel "userService/pkg/model/institution"
	mchtmodel "userService/pkg/model/merchant"
	termmodel "userService/pkg/model/term"
	"userService/pkg/model/user"
	"userService/pkg/pb"
	"userService/pkg/util"

	"google.golang.org/grpc/metadata"
)

const (
	InvalidParam  = "InvalidParamError"
	AlreadyExists = "AlreadyExistsError"
	NotFound      = "NotFoundError"
	INTERNAL      = "InternalServerError"
)

type service struct{}

func (s *service) ListTask(ctx context.Context, in *pb.ListTaskRequest) (*pb.ListTaskReply, error) {
	if in.Size == 0 {
		in.Size = 10
	}
	if in.Page == 0 {
		in.Page = 1
	}
	reply := new(pb.ListTaskReply)
	id := getUserIdFromContext(ctx)
	if id < 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}
	all := common.Enforcer.GetImplicitRolesForUser(fmt.Sprintf("user:%d", id))
	roleIds := make([]int64, 0)
	for _, a := range all {
		if strings.HasPrefix(a, "role") {
			ss := strings.Split(a, ":")
			i, _ := strconv.Atoi(ss[1])
			roleIds = append(roleIds, int64(i))
		}
	}

	db := common.DB

	tasks, count, err := camundamodel.FindTaskByRoles(db, roleIds, in.Page, in.Size)
	if err != nil {
		return nil, err
	}
	pbTasks := make([]*pb.TaskField, len(tasks))
	for i := range tasks {
		pbTasks[i] = &pb.TaskField{
			TaskId:        tasks[i].TaskId,
			Title:         tasks[i].Title,
			CurrentNode:   tasks[i].CurrentNode,
			CamundaTaskId: tasks[i].CamundaTaskId,
			InstanceId:    tasks[i].InstanceId,
			EndFlag:       *tasks[i].EndFlag,
			WorkflowName:  tasks[i].WorkflowName,
			Username:      tasks[i].UserName,
			DataId:        tasks[i].DataId,
			CreatedAt:     tasks[i].CreatedAt.Format("2006-01-02 15:03:04"),
			UpdatedAt:     tasks[i].UpdatedAt.Format("2006-01-02 15:03:04"),
		}
	}

	return &pb.ListTaskReply{
		Page:  in.Page,
		Size:  in.Size,
		Count: count,
		Items: pbTasks,
	}, nil
}

func (s *service) HandleTask(ctx context.Context, in *pb.HandleTaskRequest) (*pb.HandleTaskReply, error) {
	reply := new(pb.HandleTaskReply)
	if in.TaskId == "" || in.Result == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamError",
			Description: "result和taskId不能为空",
		}
		return reply, nil
	}
	id := getUserIdFromContext(ctx)
	if id < 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}
	taskId, _ := strconv.ParseInt(in.TaskId, 10, 64)
	db := common.DB.Begin()
	defer db.Rollback()
	task, err := camundamodel.FindTaskById(db, taskId)
	if err != nil {
		return nil, err
	}
	if task == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到任务",
		}
		return reply, nil
	}

	client := camunda.Get()
	values := make(map[string]*camundapb.Variable)
	values["result"] = &camundapb.Variable{
		Value: in.Result,
		Type:  "string",
	}

	// 保存备注
	err = camundamodel.SaveRemark(db, &camundamodel.Action{
		Comment:    in.Remark,
		Action:     in.Result,
		TaskId:     taskId,
		InstanceId: task.InstanceId,
		UserId:     id,
		RoleId:     task.RoleId,
	})
	if err != nil {
		return nil, err
	}
	// 完成任务
	completeTaskRes, err := client.Task.Complete(ctx, &camundapb.CompleteTaskReq{
		Id: task.CamundaTaskId,
		Body: &camundapb.CompleteTaskReqBody{
			Variables: values,
		},
	})
	if err != nil {
		return nil, err
	}

	if camunda.CheckError(completeTaskRes) {
		reply.Err = camunda.TransError(completeTaskRes)
		return reply, nil
	}
	// 修改当前状态为结束
	t := new(camundamodel.Task)
	endFlag := true
	t.EndFlag = &endFlag
	err = camundamodel.UpdateTask(db, &camundamodel.Task{
		TaskId: task.TaskId,
	}, t)
	if err != nil {
		return nil, err
	}

	// 查询 camunda instance id
	instance, err := camundamodel.FindProcessInstanceById(db, task.InstanceId)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusInternalServerError,
			Message:     INTERNAL,
			Description: "找不到任务流程",
		}
		return reply, nil
	}

	// 查询下一个任务
	listTaskRes, err := client.Task.GetList(ctx, &camundapb.GetListTaskReq{
		ProcessInstanceId: instance.CamundaInstanceId,
	})
	if err != nil {
		return nil, err
	}
	if camunda.CheckError(listTaskRes) {
		reply.Err = camunda.TransError(listTaskRes)
		return reply, nil
	}

	if len(listTaskRes.Tasks) != 0 {
		// 有下个任务，保存新的任务节点
		for _, t := range listTaskRes.Tasks {
			formValue, err := client.Task.GetFormValue(ctx, &camundapb.GetFormValueRequest{
				TaskId:    t.Id,
				ValueName: "status",
			})
			if err != nil {
				return nil, err
			}

			endFlag := false
			role, err := user.FindRole(db, t.Assignee)
			if err != nil {
				return nil, err
			}
			if role == nil {
				reply.Err = &pb.Error{
					Code:        http.StatusInternalServerError,
					Message:     INTERNAL,
					Description: "角色不存在",
				}
				return reply, nil
			}

			err = camundamodel.SaveTask(db, &camundamodel.Task{
				Title:         task.Title,
				RoleId:        role.ID,
				CurrentNode:   t.Name,
				CamundaTaskId: t.Id,
				InstanceId:    instance.InstanceId,
				EndFlag:       &endFlag,
				WorkflowName:  instance.WorkflowName,
				UserName:      instance.UserName,
				DataId:        instance.DataId,
			})
			if err != nil {
				return nil, err
			}

			// 修改状态
			if t.FormKey == "ins" {
				err = insmodel.UpdateInstitution(db, &insmodel.InstitutionInfo{InsIdCd: instance.DataId}, &insmodel.InstitutionInfo{InsSta: formValue.Value})
				if err != nil {
					return nil, err
				}
			} else if t.FormKey == "mcht" || t.FormKey == "term" {
				err = mchtmodel.UpdateMerchant(db, &mchtmodel.MerchantInfo{MchtCd: instance.DataId}, &mchtmodel.MerchantInfo{SystemFlag: formValue.Value})
				if err != nil {
					return nil, err
				}
				infos, _, err := termmodel.QueryTermInfo(db, &termmodel.Info{MchtCd: instance.DataId}, nil, nil, 1, 9999)
				if err != nil {
					return nil, err
				}
				for _, info := range infos {
					if info.SystemFlag != "01" && info.SystemFlag != "00" && info.SystemFlag != "13" {
						err = termmodel.UpdateTerm(db, &termmodel.Info{MchtCd: info.MchtCd, TermId: info.TermId}, &termmodel.Info{SystemFlag: formValue.Value})
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}
	db.Commit()

	return reply, nil
}

func (s *service) Start(ctx context.Context, in *pb.StartWorkflowRequest) (*pb.StartWorkflowReply, error) {
	reply := new(pb.StartWorkflowReply)

	if in.Name == "" || in.Type == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "参数不能为空",
		}
		return reply, nil
	}

	id := getUserIdFromContext(ctx)
	if id < 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}

	db := common.DB.Begin()
	defer db.Rollback()

	// 查询用户信息
	userInfo, err := user.FindUserByID(db, id)
	if err != nil {
		return nil, err
	}

	if userInfo == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamError",
			Description: "用户不存在",
		}
		return reply, nil
	}

	processes, err := camundamodel.QueryProcessDefinition(db, &camundamodel.ProcessDefinition{
		Name: in.Type,
	})
	if err != nil {
		return nil, err
	}
	if len(processes) == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "InvalidParamError",
			Description: "工作流类型不存在",
		}
		return reply, nil
	}
	// 查询是否有未完成的任务
	instances, err := camundamodel.QueryProcessInstance(db, &camundamodel.ProcessInstance{
		DataId: in.DataId,
		EndFlag: sql.NullInt64{
			Int64: 0,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}
	var instanceId string

	client := camunda.Get()
	if len(instances) == 0 {
		startProcessInstanceRes, err := client.ProcessDefinition.Start(ctx, &camundapb.StartProcessDefinitionReq{
			Id: processes[0].Id,
			Body: &camundapb.StartProcessDefinitionReqBody{
				BusinessKey: in.Name,
			},
		})
		if err != nil {
			return nil, err
		}

		if camunda.CheckError(startProcessInstanceRes) {
			reply.Err = camunda.TransError(startProcessInstanceRes)
			return reply, nil
		}

		err = camundamodel.SaveProcessInstance(db, &camundamodel.ProcessInstance{
			CamundaInstanceId: startProcessInstanceRes.Item.Id,
			Title:             in.Name,
			DataId:            in.DataId,
			UserId:            id,
			UserName:          userInfo.UserName,
			WorkflowName:      in.Type,
			EndFlag: sql.NullInt64{
				Int64: 0,
				Valid: true,
			},
		})
		if err != nil {
			return nil, err
		}
		instanceId = startProcessInstanceRes.Item.Id
	} else {
		instanceId = instances[0].CamundaInstanceId
	}

	// 获取任务
	taskListRes, err := client.Task.GetList(ctx, &camundapb.GetListTaskReq{
		ProcessInstanceId: instanceId,
	})
	if err != nil {
		return nil, err
	}

	if camunda.CheckError(taskListRes) {
		reply.Err = camunda.TransError(taskListRes)
		return reply, nil
	}

	// 写入任务列表
	if len(taskListRes.Tasks) == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     "WorkflowError",
			Description: "没有任务需要执行",
		}
		return reply, nil
	}

	if len(instances) == 0 {
		instance, err := camundamodel.FindProcessInstanceByCamundaInstanceId(db, instanceId)
		if err != nil {
			return nil, err
		}
		if instance == nil {
			reply.Err = &pb.Error{
				Code:        http.StatusInternalServerError,
				Message:     INTERNAL,
				Description: "找不到任务流程流程",
			}
			return reply, nil
		}

		for _, task := range taskListRes.Tasks {
			formValue, err := client.Task.GetFormValue(ctx, &camundapb.GetFormValueRequest{
				TaskId:    task.Id,
				ValueName: "status",
			})
			if err != nil {
				return nil, err
			}

			var endFlag = false
			role, err := user.FindRole(db, task.Assignee)
			if err != nil {
				return nil, err
			}
			if role == nil {
				reply.Err = &pb.Error{
					Code:        http.StatusInternalServerError,
					Message:     INTERNAL,
					Description: "角色不存在",
				}
				return reply, nil
			}
			err = camundamodel.SaveTask(db, &camundamodel.Task{
				Title:         in.Name,
				RoleId:        role.ID,
				CurrentNode:   task.Name,
				CamundaTaskId: task.Id,
				InstanceId:    instance.InstanceId,
				EndFlag:       &endFlag,
				WorkflowName:  in.Type,
				UserName:      userInfo.UserName,
				DataId:        instance.DataId,
			})
			if err != nil {
				return nil, err
			}

			// 修改商户或者机构状态
			if task.FormKey == "ins" {
				err = insmodel.UpdateInstitution(db, &insmodel.InstitutionInfo{InsIdCd: instance.DataId}, &insmodel.InstitutionInfo{InsSta: formValue.Value})
				if err != nil {
					return nil, err
				}
			} else if task.FormKey == "mcht" || task.FormKey == "term" {
				err = mchtmodel.UpdateMerchant(db, &mchtmodel.MerchantInfo{MchtCd: instance.DataId}, &mchtmodel.MerchantInfo{SystemFlag: formValue.Value})
				if err != nil {
					return nil, err
				}
				infos, _, err := termmodel.QueryTermInfo(db, &termmodel.Info{MchtCd: instance.DataId}, nil, nil, 1, 9999)
				if err != nil {
					return nil, err
				}
				for _, info := range infos {
					if info.SystemFlag != "01" && info.SystemFlag != "00" {
						err = termmodel.UpdateTerm(db, &termmodel.Info{MchtCd: info.MchtCd, TermId: info.TermId}, &termmodel.Info{SystemFlag: formValue.Value})
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}
	db.Commit()
	task, err := camundamodel.FindTaskByCamundaId(common.DB, taskListRes.Tasks[0].Id)
	if err != nil {
		return nil, err
	}
	reply.TaskId = task.TaskId
	return reply, nil
}

func (s *service) ListRemark(ctx context.Context, in *pb.ListRemarkRequest) (*pb.ListRemarkReply, error) {
	if in.Size == 0 {
		in.Size = 10
	}
	if in.Page == 0 {
		in.Page = 1
	}

	query := new(camundamodel.Action)
	if in.Item != nil {
		query.ActionId = in.Item.ActionId
		query.Comment = in.Item.Comment
		query.TaskId = in.Item.TaskId
		query.Action = in.Item.Action
		query.InstanceId = in.Item.InstanceId
	}
	db := common.DB
	items, count, err := camundamodel.QueryRemark(db, query, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	docs := make([]*pb.RemarkField, len(items))
	for i := range items {
		docs[i] = &pb.RemarkField{
			ActionId:   items[i].ActionId,
			Action:     items[i].Action.Action,
			Comment:    items[i].Comment,
			TaskId:     items[i].TaskId,
			InstanceId: items[i].InstanceId,
			Username:   items[i].UserName,
			CreatedAt:  items[i].CreatedAt.Format(util.TimePattern),
			UpdatedAt:  items[i].UpdatedAt.Format(util.TimePattern),
		}
	}

	return &pb.ListRemarkReply{
		Items: docs,
		Count: count,
		Page:  in.Page,
		Size:  in.Size,
	}, nil
}

func getUserIdFromContext(ctx context.Context) int64 {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return -1
	}
	ids := md.Get("userid")
	if !ok || len(ids) == 0 {
		return -1
	}

	str := ids[0]
	id, _ := strconv.ParseInt(str, 10, 64)
	return id
}
