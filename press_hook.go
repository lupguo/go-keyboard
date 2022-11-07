package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/robotn/gohook"
	log "github.com/sirupsen/logrus"
	"github/lupingguo/go-keyborad/config"
)

// RStatus 运行状态
type RStatus bool

const (
	NoRunning RStatus = false
	Running   RStatus = true
)

// PressTask 按键任务
type PressTask struct {
	Name          string
	StartCh       chan bool
	StopCh        chan bool
	Lock          *sync.RWMutex
	RunningStatus RStatus

	// 按键配置
	Ticker   *time.Ticker
	Keyboard *config.Keyboard
}

// SetStatus 设置任务状态
func (t *PressTask) SetStatus(status RStatus) {
	t.Lock.Lock()
	defer t.Lock.Unlock()

	switch status {
	case Running:
		if t.RunningStatus == status {
			log.Warnf("go-keyboard[%s] is running, do not need start again!!", t.Name)
			return
		}
		t.RunningStatus = status
		t.StartCh <- true
	case NoRunning:
		if t.RunningStatus == status {
			log.Warnf("go-keyboard[%s] is stoped, do not need stop again!!", t.Name)
			return
		}
		t.RunningStatus = status
		t.StopCh <- true
	}
}

// Start 按键任务启动
func (t *PressTask) Start() {
	for {
		select {
		case <-t.StartCh:
			log.Infof("go-keyboard[%s] start....", t.Name)
			t.BeginLoopPress()
		}
	}
}

// BeginLoopPress 启动循环按键
func (t *PressTask) BeginLoopPress() {
	kbd := t.Keyboard
	t.Ticker = time.NewTicker(time.Duration(kbd.KeySleep) * time.Millisecond)

	// // pid by name
	// if processName := kbd.ProcessName; processName != "" {
	// 	ActivePidByName(processName)
	// }

	// 启动循环按键任务
	for {
		select {
		case <-t.Ticker.C:
			// robotgo.KeyToggle(kbd.KeyPress, "down")
			// robotgo.KeyToggle(kbd.KeyPress, "up")
			err := robotgo.KeyTap(kbd.KeyPress)
			if err != nil {
				log.Errorf("key tap [%s] got err: %s", kbd.KeyPress, err)
				return
			}
		case <-t.StopCh:
			log.Infof("receive stop sign, go-keyboard[%s] stop!", t.Name)
			t.Ticker.Stop()
			return
		}
	}
}

// ActivePidByName 通过名称激活窗口
func ActivePidByName(processName string) []int32 {
	fpids, err := robotgo.FindIds(processName)
	if err != nil {
		log.Warnf("can not found fpids by name[%s]", processName)
		return nil
	}
	if len(fpids) == 0 {
		log.Warnf("can not found fpids by name[%s], empty pids", processName)
		return nil
	}
	log.Infof("found fpids by name[%s], got pids: %v", processName, fpids)

	// active pid
	masterPID := fpids[0]
	if err := robotgo.ActivePID(masterPID); err != nil {
		log.Warnf("begin loop press, active pid [%d] got err: %s", masterPID, err)
		return nil
	}

	return fpids
}

// RegisterPressHook 注册按键回调
func RegisterPressHook() chan *PressTask {
	keyboards := config.GetKeyboards()

	// 注册按键开关(开)
	taskCh := make(chan *PressTask, len(keyboards))
	for _, keyboard := range keyboards {
		// 按键任务初始化
		task := &PressTask{
			Name:          keyboard.Name,
			StartCh:       make(chan bool),
			StopCh:        make(chan bool),
			Lock:          &sync.RWMutex{},
			RunningStatus: NoRunning,
			Keyboard:      keyboard,
		}

		// 按键任务监听
		keyDiff := keyboard.StartStopKeyDiff
		startKey, stopKey := keyboard.StartKey, keyboard.StopKey
		if keyDiff { // 启停按键不一致配置
			hook.Register(hook.KeyDown, startKey, func(e hook.Event) {
				task.SetStatus(Running)
			})

			hook.Register(hook.KeyDown, stopKey, func(e hook.Event) {
				task.SetStatus(NoRunning)
			})
		} else { // 启停按键一致的
			hook.Register(hook.KeyDown, startKey, func(e hook.Event) {
				if task.RunningStatus == NoRunning {
					task.SetStatus(Running)
				} else {
					task.SetStatus(NoRunning)
				}
			})
		}

		// 添加任务
		taskCh <- task
	}

	hook.Process(hook.Start())
	return taskCh
}

// DebugPrintKeyPress 调试鼠标、键盘事件
func DebugPrintKeyPress() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		fmt.Println("hook: ", ev)
	}
}
