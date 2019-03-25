# goInstallBinariesTool

此程序意在帮助国内童鞋在安装go-vim插件后，解决":GoInstallBinaries"失败（没有翻墙）的问题。

此程序利用[gopm](https://github.com/gpmgo/gopm)工具下载相关源码包， goInstallBinariesTool 会自动第一个安装 gopm 命令。

安装 goInstallBinaries Tool ，并运行程序即可：（Windows如下）
```
go get -u github.com/davidhenrygao/goInstallBinariesTool

goInstallBinariesTool.exe
```

程序默认读取pkgfile文件，解析相关要下载的源码包，并进行下载安装。你可以通过运行"goInstallBinariesTool xxx"来指定读取文件（例如你可以将go-vim插件中的go.vim拷贝到本目录下，然后运行"goInstallBinariesTool go.vim"即可）。

pkgfile文件必须满足如下格式：
packages = ["xxx","yyy",...]
程序会从中解析出xxx，yyy等包名

注：
    指定go.vim文件运行程序有可能出错，出错原因会进行打印，请根据打印自行解决某个包的按错出错问题。
    如安装“github.com/kisielk/errcheck”包时，会提示缺少“github.com/kisielk/gotool”包，则应该先安装该包。（你可以手动进行安装，也可以在读取文件中“github.com/kisielk/errcheck”包行前加入一行“github.com/kisielk/gotool”，然后在重新运行程序）
