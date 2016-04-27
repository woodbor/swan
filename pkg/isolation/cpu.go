package isolation

import "os/exec"
import "io/ioutil"
import "strconv"


type CpuShares struct {
 cgroupName string
 cpuShares string
 cgCpus string

}

func NewCpuShares(nameOfTheCgroup string, myCpuShares string, myCpus string ) *CpuShares{
	return &CpuShares{cgroupName: nameOfTheCgroup, cpuShares: myCpuShares, cgCpus: myCpus}
}

func (cpu *CpuShares) Create() error {
	return nil
}

func (cpu *CpuShares) Delete() error {
	controllerName := "cpu"
        cmd := exec.Command("sh", "-c", "cgdelete -g "+controllerName+":"+cpu.cgroupName)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}


func (cpu *CpuShares) Isolate(PID int) error {
     	// 1.a Create cpu cgroup
	controllerName := "cpu"
	cgroupName := "A"
        cmd := exec.Command("sh", "-c", "cgcreate -g "+controllerName+":"+cgroupName)

	err := cmd.Run()
	if err != nil {
		return err
	}

     	// 1.b Set cpu cgroup shares

	cpuShares := "1024"

        cmd = exec.Command("sh", "-c", "cgset -r cpu.shares="+cpuShares+" "+ cgroupName)

	err = cmd.Run()
	if err != nil {
		return err
	}


     	// 4. Set PID to cgroups

	//Associate task with the cgroup
	//cgclassify seems to exit with error so temporarily using file io


        strPID :=strconv.Itoa(PID)
	d := []byte(strPID)
	err = ioutil.WriteFile("/sys/fs/cgroup/"+controllerName+"/"+cgroupName+"/tasks", d, 0644)

	if err != nil {
		panic(err.Error())
	}




//        cmd = exec.Command("sh", "-c", "cgclassify -g cpu:A " + string(PID))
//
//	err = cmd.Run()
//	if err != nil {
//		panic("Cgclassify failed: " + err.Error())
//	}


	return nil
}


//	//Write CPU shares
//        strShares :=strconv.Itoa(1024)
//	c := []byte(strShares)
//	err2 := ioutil.WriteFile("/tmp/x.mri", c, 0644)
//
//	if err2 != nil {
//		panic(err2.Error())
//	}
//
//	//Associate task with the cgroup
//        strPID :=strconv.Itoa(PID)
//	d := []byte(strPID)
//	err3 := ioutil.WriteFile("/tmp/y.mri", d, 0644)
//
//	if err3 != nil {
//		panic(err3.Error())
//	}
//
//
//
//
//
//func (cpu *CpuShares) Isolate(PID int) error {
//        cmd := exec.Command("sh", "-c", "touch /tmp/x.tmp")
//
//	err := cmd.Run()
//	if err != nil {
//		panic(err.Error())
//	}
//
//	return nil
//}

