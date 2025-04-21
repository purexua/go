# 组织一个Go模块

> 原文：[Organizing a Go module](https://go.dev/doc/modules/layout)
> 译者：[purexua](https://github.com/purexua)

[TOC]

开发者初学 Go 时常见的疑问是“我该如何组织我的 Go 项目？”，从文件和文件夹的布局角度来说。本文件的目的是提供一些指导方针，以帮助回答这个问题。为了充分利用本文件，请确保你已经通过阅读 [教程](https://go.dev/doc/tutorial/create-module) 和 [管理模块源](https://go.dev/doc/modules/managing-source) 熟悉了Go模块的基础知识。

Go项目可以包括包、命令行程序或两者的组合。本指南按项目类型组织。

### [基本包](https://go.dev/doc/modules/layout#basic-package)

一个基本的Go包所有代码都在项目的根目录中。项目由一个模块组成，该模块包含一个包。包名与模块名的最后一个路径组件匹配。对于只需要一个Go文件的非常简单的包，项目结构如下：

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
```

*[在此文档中，文件/包名称完全是任意的]*

假设此目录上传到GitHub仓库 `github.com/someuser/modname`，`go.mod `文件中的 `module `行应写为 `module github.com/someuser/modname`。

`modname.go `中的代码通过以下方式声明包：

```
package modname

// ... package code here
```

用户可以通过在Go代码中使用以下命令来导入此包并依赖它：

```
import "github.com/someuser/modname"
```

一个Go包可以被拆分为多个文件，所有文件都位于同一目录下，例如：

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth.go
  auth_test.go
  hash.go
  hash_test.go
```

目录中的所有文件声明  `package modname`。

### [基本命令](https://go.dev/doc/modules/layout#basic-command)

一个基本的可执行程序（或命令行工具）根据其复杂性和代码大小进行结构化。最简单的程序可以由一个定义了  `func main` 的单个 Go 文件组成。较大的程序可以将代码分散在多个文件中，所有文件都声明  `package main`：

```
project-root-directory/
  go.mod
  auth.go
  auth_test.go
  client.go
  main.go
```

这里 `main.go` 文件包含 `func main`，但这只是一种约定。"main" 文件也可以称为 `modname.go`（对于适当的 `modname `值）或其他任何名称。

假设这个目录被上传到 `github.com/someuser/modname` 的GitHub仓库，`go.mod `文件中的 `module` 行应该这样写：

```
module github.com/someuser/modname
```

并且用户应该能够使用以下方式在自己的机器上安装它：

```
$ go install github.com/someuser/modname@latest
```

### [带有支持包的包或命令](https://go.dev/doc/modules/layout#package-or-command-with-supporting-packages)

较大的包或命令可能从将一些功能拆分到支持包中受益。最初，建议将这些包放置在名为 `internal` 的目录中；[这可以防止](https://pkg.go.dev/cmd/go#hdr-Internal_Directories)其他模块依赖于我们不一定要公开和支持的包。由于其他项目无法从我们的 `internal` 目录导入代码，我们可以自由地重构其API，通常可以随意移动事物而不会破坏外部用户。因此，包的项目结构如下：

```
project-root-directory/
  internal/
    auth/
      auth.go
      auth_test.go
    hash/
      hash.go
      hash_test.go
  go.mod
  modname.go
  modname_test.go
```

`modname.go` 文件声明 `package modname`，`auth.go` 声明 `package auth` 等等。`modname.go` 可以如下导入 `auth` 包：

```
import "github.com/someuser/modname/internal/auth"
```

命令在`内部`目录中带有支持包的布局非常相似，除了根目录中的文件（们）声明`package main`。

### [多个包](https://go.dev/doc/modules/layout#multiple-packages)

一个模块可以由多个可导入的包组成；每个包都有自己的目录，并且可以按层次结构组织。以下是一个示例项目结构：

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth/
    auth.go
    auth_test.go
    token/
      token.go
      token_test.go
  hash/
    hash.go
  internal/
    trace/
      trace.go
```

作为提醒，我们假设`go.mod`中的`module`行表示：

```
module github.com/someuser/modname
```

`modname` 包位于根目录，声明 `package modname`，可以被用户通过以下方式导入：

```
import "github.com/someuser/modname"
```

子包可以按以下方式由用户导入：

```
import "github.com/someuser/modname/auth"
import "github.com/someuser/modname/auth/token"
import "github.com/someuser/modname/hash"
```

包 `trace` 位于 `internal/trace` 中，无法在此模块外部导入。建议尽可能将包保留在 `internal` 中。

### [多个命令](https://go.dev/doc/modules/layout#multiple-commands)

同一存储库中的多个程序通常会有独立的目录：

```
project-root-directory/
  go.mod
  internal/
    ... shared internal packages
  prog1/
    main.go
  prog2/
    main.go
```

在每个目录中，程序的Go文件声明`package main`。顶级`internal`目录可以包含由存储库中所有命令使用的共享包。

用户可以按照以下方式安装这些程序：

```
$ go install github.com/someuser/modname/prog1@latest
$ go install github.com/someuser/modname/prog2@latest
```

一个常见的约定是将所有命令放入一个`cmd`目录中；虽然在一个仅由命令组成的仓库中这并非严格必要，但在一个既有命令又有可导入包的混合仓库中，这非常有用，正如我们接下来将要讨论的。

### [同一存储库中的包和命令](https://go.dev/doc/modules/layout#packages-and-commands-in-the-same-repository)

有时，一个仓库将提供可导入的包和具有相关功能的可安装命令。以下是一个此类仓库的示例项目结构：

```
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth/
    auth.go
    auth_test.go
  internal/
    ... internal packages
  cmd/
    prog1/
      main.go
    prog2/
      main.go
```

假设此模块名为 `github.com/someuser/modname`，用户现在可以从中导入包：

```
import "github.com/someuser/modname"
import "github.com/someuser/modname/auth"
```

从其中安装程序：

```
$ go install github.com/someuser/modname/cmd/prog1@latest
$ go install github.com/someuser/modname/cmd/prog2@latest
```

### [服务器项目](https://go.dev/doc/modules/layout#server-project)

Go 是实现 *服务器* 的常见语言选择。由于服务器开发涉及许多方面，如协议（REST？gRPC？）、部署、前端文件、容器化、脚本等，这类项目的结构存在很大差异。我们将在此处重点关注用 Go 编写的项目部分。

服务器项目通常不会有用于导出的包，因为服务器通常是一个自包含的二进制文件（或一组二进制文件）。因此，建议将实现服务器逻辑的Go包保存在`internal`目录中。此外，由于项目可能还有许多包含非Go文件的其它目录，将所有Go命令保存在一个`cmd`目录中是个好主意：

```
project-root-directory/
  go.mod
  internal/
    auth/
      ...
    metrics/
      ...
    model/
      ...
  cmd/
    api-server/
      main.go
    metrics-analyzer/
      main.go
    ...
  ... the project's other directories with non-Go code
```

如果服务器仓库增长到包含对其他项目有共享价值的包，最好将这些包拆分到独立的模块中。