### go 包安装问题

1. 相对路径导入的问题:
    1. 如果项目在GOPATH路径下(包括多级目录), 那么包导入的时候使用绝对路径, 从 GOPATH/src 开始算起. **此时不支持使用 ./chat 导入文件同级目录的包**
    2. 如果项目在非GOPATH路径及其子路径下, 那么使用 ./chat 导入文件同级目录的包, 使用绝对路径导入GOPATH路径下的包

2. 包管理神器: glide
    1. 安装:

        ``` sh
        go get github.com/Masterminds/glide
        go install github.com/Masterminds/glide
        或者:
        cd $GOPATH/src/github.com/Masterminds/glide
        go build
        cp glide /usr/local/bin. 
        # 直接在path中添加 $GOPATH/bin也行
        或者:
        sudo yum install glide

        #对于windows, go install 之后 glide.exe 在 $GOPATH/bin 下, 将 $GOPATH/bin 加入 path 变量
        ```
        
    2. 命令:

        ``` sh
        glide create|init 初始化项目并创建glide.yaml文件.
        glide get         获取单个包
            --all-dependencies   会下载所有关联的依赖包
            -s            删除所有版本控制，如.git
            -v            删除嵌套的vendor
        glide install     安装包
        glide update|up   更新包    
        glide config-wizard      启动引导向导
        glide novendor    列出除了vendor以外的所有目录 go test $(glide novendor)
        glide list        显示项目导入的所有包的按字母顺序排列的列表。
        glide help        打印 glide 帮助
        glide name        返回glide.yaml文件中列出的软件包的名称.  编写Glide脚本时，获取正在使用的包的名称。
        glide –version    显示版本信息
        #神器
        glide mirror      镜像提供了将 repo 位置替换为作为原始镜像的另一位置的能力
        #希望拥有连续集成（CI）系统的缓存时，或者需要在本地位置的依赖项上工作, 或者下载翻墙包时，非常有用
        
        例子:
        #下载指定版本的包
        glide get github.com/go-sql-driver/mysql#v1.2
        
        #测试所有包, 除了 vendor 下面的依赖项和依赖关系
        go test $(glide novendor)

        glide mirror set [original] [replacement]
        glide mirror set [original] [replacement] --vcs [type]
        glide mirror remove [original]
        glide mirror set https://github.com/example/foo https://git.example.com/example/foo.git
        glide mirror remove https://github.com/example/foo
        ```

    3. Glide 在windows10 64位系统上的bug修复方案
        问题(有可能会显示乱码): [ERROR] Unable to export dependencies to vendor directory: Error moving files: exit status 1. output: Access is denied. 0 dir(s) moved 
        
        解决方案: 
        ``` go
        // 找到这个文件 github.com/Masterminds/glide/path/winbug.go
        func CustomRemoveAll(p string) error {
            ...
            //主要修改这一行代码，将其替换为 
            // cmd := exec.Command("cmd.exe", "/c", "xcopy /s/y", o, n+"\\")
            cmd := exec.Command("cmd.exe", "/c", "rd", "/s", "/q", p) // 原始代码
            ...
        }
        
        // 在glide目录下重新编译: go build glide.go ,  go install
        // 将生成的glide.exe复制到 $GOPATH/bin 下, 此目录加入path环境变量
        // 更详细的讨论解决方案 https://github.com/Masterminds/glide/issues/873 (这是老毛病了)
        ```
        
    4. glide彻底解决go get golang.org/x/net 安装失败

        ``` c
        // 就是 go mirror 神器
        $ rm -rf ~/.glide
        $ mkdir -p ~/.glide
        $ glide mirror set https://golang.org/x/mobile https://github.com/golang/mobile --vcs git
        $ glide mirror set https://golang.org/x/crypto https://github.com/golang/crypto --vcs git
        $ glide mirror set https://golang.org/x/net https://github.com/golang/net --vcs git
        $ glide mirror set https://golang.org/x/tools https://github.com/golang/tools --vcs git
        $ glide mirror set https://golang.org/x/text https://github.com/golang/text --vcs git
        $ glide mirror set https://golang.org/x/image https://github.com/golang/image --vcs git
        $ glide mirror set https://golang.org/x/sys https://github.com/golang/sys --vcs git
        
        然后在项目中执行 glide init, 创建 yaml 文件, 下次使用 glide install 安装的时候就不会失败了
        ```
        
        
    5. glide 解决 golang.org/x/net/context 无法连接下载的问题


        有部分go软件包只依赖 net 包子目录中的 context 包, 这个官网没有单独的子目录连接, 使用:
        > $ glide mirror set https://golang.org/x/net/context https://github.com/golang/net --vcs git
        
        镜像至 net 包, 然后在 vendor 目录下, 将 golang.org/x/net/context.... 等等目录修改为 golang.org/x/net/...
        ps: 虽然可以自己建立一个 net/context 分支然后镜像, 但以后就不好更新同步了