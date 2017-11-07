package supervisor

import (
	"encoding/json"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/wallaby/util"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	/*
			Process State:
		state:		NULL --> ProcStarted --> ProcStopped --> ProcCleaned
		action:		    start            stop           clean
	*/
	ProcNull    = "ProcNull"
	ProcStarted = "ProcStarted"
	ProcStopped = "ProcStopped"
	ProcCleaned = "ProcCleaned"
)

type ProcInfo struct {
	Pid             int
	Status          string
	CWD             string
	StartCmd        string
	StopCmd         string
	StopCmdTimeOut  uint
	CheckCmd        string
	CheckCmdTimeOut int
	CleanCmd        string
	NeedSupervisor  bool
}

type Proc struct {
	ProcInfo
	sentKillSignal bool
	startCommand   *exec.Cmd
	mutex          sync.Mutex
}

/*
	run : cd $cwd; $cmd
	needSupervisor : automatically restart if true
*/
func NewProc(procInfo ProcInfo) *Proc {
	// clean Pid
	procInfo.Pid = -1
	return &Proc{
		ProcInfo:       procInfo,
		sentKillSignal: false,
		startCommand:   nil,
	}
}

func executeCommand(cwd string, command string) *exec.Cmd {
	//  /bin/sh -c "cd $cwd; $cmd"
	cmdStr := "cd " + cwd + ";" + command
	countlog.Info("event!Fork", "cmd", cmdStr)
	return exec.Command("/bin/sh", "-c", cmdStr)
}

func dirExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}

func (p *Proc) Fork() error {
	if !dirExists(p.CWD) {
		return util.NOTFOUND
	}
	cmd := executeCommand(p.CWD, p.StartCmd)
	err := cmd.Start()
	if err != nil {
		countlog.Error("event!Fork fail", "proc", err)
		return err
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Status = ProcStarted
	p.Pid = p.startCommand.Process.Pid
	p.startCommand = cmd
	p.sentKillSignal = false

	countlog.Info("event!Proc", "fork", p.Pid)

	if p.NeedSupervisor {
		go p.supervise()
	} else {
		go p.waitKilled()
	}

	return err
}

func (p *Proc) Kill() error {
	if p.Status != ProcStarted {
		return nil
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.sentKillSignal = true
	var err error
	countlog.Info("event!Proc", "kill", p.Pid)
	if p.StopCmd != "" {
		countlog.Info("event!Proc", "stop", p.StopCmd)
		err = p.stop()
	} else {
		err = p.kill()
	}
	if err != nil {
		countlog.Error("event!Fork fail", "proc", err)
	}
	return err
}

func (p *Proc) stop() error {
	countlog.Info("event!Proc", "stop", p.StopCmd)
	cmd := executeCommand(p.CWD, p.StopCmd)
	err := cmd.Start()
	if err != nil {
		countlog.Warn("event!Proc", "stop fail", err)
	}
	time.AfterFunc(time.Duration(p.StopCmdTimeOut)*time.Second, func() {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		if p.Status == ProcStopped {
			countlog.Info("event!Proc", "killed", p.CWD)
		} else {
			countlog.Info("event!Proc", "killing", p.CWD)
			p.startCommand.Process.Kill()
		}
	})
	return err
}

func (p *Proc) kill() error {
	return p.startCommand.Process.Kill()
}

func (p *Proc) Json() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	bin, err := json.Marshal(p.ProcInfo)
	if err == nil {
		return string(bin)
	} else {
		return ""
	}
}

func (p *Proc) supervise() error {
	err := p.waitKilled()
	// TODO: check wait status
	if err != nil { // TODO: or other ExitCode
		countlog.Error("event!Proc", "proc", err)
	}
	if p.sentKillSignal {
		return nil
	} else {
		// restart
		return p.Fork()
	}
}

func (p *Proc) waitKilled() error {
	countlog.Info("event!Proc", "wait", p.Pid)
	err := p.startCommand.Wait()
	//The returned error is nil if the startCommand runs,
	// has no problems copying stdin, stdout, and stderr,
	// and exits with a zero exit status.
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Status = ProcStopped
	return err
}

type ProcMgr struct {
	procList map[string]*Proc // CWD --> *Proc
	mutex    sync.Mutex
}

func (mgr *ProcMgr) StartProc(procinfo ProcInfo) error {
	// Check: proc is ProcStopped or ProcNull
	var proc *Proc
	var ok bool
	{
		mgr.mutex.Lock()
		defer mgr.mutex.Unlock()
		proc, ok = mgr.procList[procinfo.CWD]
		if ok {
			if proc.Status != ProcStopped && proc.Status != ProcNull {
				return util.ERRSTATUS
			}
		}
		proc = NewProc(procinfo)
		mgr.procList[procinfo.CWD] = proc
	}
	return proc.Fork()
}

func (mgr *ProcMgr) StopProc(cwd string, pid int) error {
	// Check: proc is ProcStopped or ProcNull
	var proc *Proc
	var ok bool
	{
		mgr.mutex.Lock()
		defer mgr.mutex.Unlock()
		proc, ok = mgr.procList[cwd]
		if !ok || proc.Status != ProcStarted {
			return util.ERRSTATUS
		}
	}
	return proc.Kill()
}

func (mgr *ProcMgr) List() []ProcInfo {
	// TODO:
	return nil
}

func (mgr *ProcMgr) CleanProc(cwd string) error {
	// Check: proc is ProcStarted
	// TODO:
	return nil
}
