---
title: goshim(1) Document
---

> Document in Japanese: <!-- mdpplink href=./README-ja.md -->[goshim(1) ドキュメント（日本語）](./README-ja.md)<!-- /mdpplink -->

goshim(1)

# NAME

<!--
goshim - ソースに変更があった Go プログラムだけを再コンパイルして実行するプログラム
-->

goshim - Re-compiles only updated Go-source codes and execute transparently

# INSTALLATION

  $ go install github.com/knaka/goshim/cmd/goshim@latest

# SYNOPSIS

<!-- 
設定ファイル（ `~/.config/goshim.toml` ）が存在しなければ、初回実行時に自動的に作成する。
-->

If the configuration file (`~/.config/goshim.toml`) does not exist, it is automatically created at the first invocation. 

<!-- 
デフォルトでは、 `~/src/go/` が対象プロジェクトになっている。
-->

By default, `~/src/go/` is the target project.

<!-- 
以下を実行すると、 `~/src/go/src/cmd/` 以下のディレクトリ名と同名のシンボリックリンクを `$GOBIN` 以下に作成する。
-->

The following will create a symbolic links under `$GOBIN` with the same names as the directory names under `~/src/go/src/cmd/`.

```
$ goshim install
```

<!--
上記で作成されたシンボリックリンクを実行すると、対象のプログラムをビルドして `$GOBIN/.goshim/` 以下にインストールし、実行の際に指定された引数を渡してそのバイナリを透過的に実行する。
-->

Executing the symbolic link created above will build and install the target program into `$GOBIN/.goshim/` and transparently execute its binary passing the arguments specified at the time of execution.

<!-- 
前回実行時以降にソースが修正されていた場合は、自動的に再ビルドして実行する。
-->

If the source code has been modified since the last run, it is automatically rebuilt and executed.

# DESCRIPTION

# OPTIONS
