package pidfile

import (
	"fmt"
	"os"
	"strconv"
)

type PIDFile struct {
	pid     int
	path    string
	process *os.Process
}

func NewPIDFile(path string) (p *PIDFile, compare bool) {
	p = new(PIDFile)
	pid := os.Getpid()
	p.path = path
	if p.isExists() {
		if err := p.read(); err != nil {
			p.delete()
		} else {
			if p.pid != pid {
				return p, false
			} else {
				return p, true
			}
		}
	}
	p.pid = pid
	p.create()
	return p, true
}

func (p *PIDFile) Delete() {
	p.delete()
}

func (p *PIDFile) Send(sign os.Signal) error {
	if err := p.process.Signal(sign); err != nil {
		return err
	}
	return nil
}

func (p *PIDFile) Kill() error {
	if err := p.process.Kill(); err != nil {
		return err
	}
	return nil
}

func (p *PIDFile) isExists() bool {
	_, err := os.Lstat(p.path)
	if err != nil {
		return os.IsNotExist(err)
	}
	return true
}

func (p *PIDFile) read() error {
	file, err := os.Open(p.path)
	if err != nil {
		return err
	}
	defer file.Close()
	strPid := ""
	fmt.Fscan(file, &strPid)
	p.pid, err = strconv.Atoi(strPid)
	if err != nil {
		return err
	}
	p.process, err = os.FindProcess(p.pid)
	if err != nil {
		return err
	}
	return nil
}

func (p *PIDFile) delete() {
	os.Remove(p.path)
}

func (p *PIDFile) create() error {
	file, err := os.Create(p.path)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(strconv.Itoa(p.pid))
	return nil
}
