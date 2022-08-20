package classpath

import (
	"os"
	"path/filepath"
)

type Classpath struct {
	bootClasspath Entry
	extClasspath  Entry
	userClasspath Entry
}

// Parse 使用-Xjre选项解析启动类路径和扩展类路径，使用 -classpath或-cp解析用户类路径
func Parse(jreOption, cpOption string) *Classpath {
	cp := &Classpath{}
	cp.parseBootAndExtClasspath(jreOption)
	cp.parseUserClasspath(cpOption)
	return cp
}

// ReadClass 读取类路径，依次从系统启动类、扩展类、用户定义类读取
func (self *Classpath) ReadClass(className string) ([]byte, Entry, error) {
	className = className + ".class"
	// 先从系统启动类路径查找
	if data, entry, err := self.bootClasspath.readClass(className); err == nil {
		return data, entry, nil
	}
	// 再从扩展类查找
	if data, entry, err := self.extClasspath.readClass(className); err == nil {
		return data, entry, nil
	}
	// 最后从用户指定目录加载
	return self.userClasspath.readClass(className)
}

func (self *Classpath) String() string {
	return self.userClasspath.String()
}

func (self *Classpath) parseBootAndExtClasspath(jreOption string) {
	jreDir := getJreDir(jreOption)
	// jre/lib/*
	jreLibPath := filepath.Join(jreDir, "lib", "*")
	self.bootClasspath = newWildcardEntry(jreLibPath)
	// jre/lib/ext/*
	jreExtPath := filepath.Join(jreDir, "lib", "ext", "*")
	self.extClasspath = newWildcardEntry(jreExtPath)
}

func (self *Classpath) parseUserClasspath(cpOption string) {
	if cpOption == "" {
		cpOption = "."
	}
	self.userClasspath = newEntry(cpOption)
}

func getJreDir(jreOption string) string {
	// 优先使用用户输入的-Xjre选项作为jre目录
	if jreOption != "" && exists(jreOption) {
		return jreOption
	}
	// 在当前mukluk下寻找jre目录
	if exists("./jre") {
		return "./jre"
	}
	// 尝试使用系统JAVA_HOME环境变量
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		return filepath.Join(jh, "jre")
	}
	// 都找不到的话提示
	panic("Can not find jre folder!")
}

// 判断目录是否存在
func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
