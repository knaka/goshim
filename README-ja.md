---
title: goadhoc(1) ドキュメント（日本語）
---

goadhoc(1)

# NAME

goadhoc - ソースに変更があった Go プログラムだけを再コンパイルして透過的に実行するプログラム

# INSTALLATION

```
$ go install github.com/knaka/goadhoc/cmd/goadhoc@latest
```

# SYNOPSIS

設定ファイル（ `~/.config/goadhoc.toml` ）が存在しなければ、初回実行時に自動的に作成する。

デフォルトでは、 `~/src/go/` が対象プロジェクトになっている。

以下を実行すると、 `~/src/go/src/cmd/` 以下のディレクトリ名と同名のシンボリックリンクを `$GOBIN` 以下に作成する。

  $ goadhoc install

上記で作成されたシンボリックリンクを実行すると、対象のプログラムをビルドして `$GOBIN/.goadhoc/` 以下にインストールし、実行の際に指定された引数を渡してそのバイナリを透過的に実行する。

前回実行時以降にソースが修正されていた場合は、自動的に再ビルドして実行する。

# DESCRIPTION

# OPTIONS
