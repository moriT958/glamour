# NoteBox 向けの fork

[NoteBox](https://github.com/moriT958/NoteBox) が使うための [charmbracelet/glamour](https://github.com/charmbracelet/glamour) の fork です。

- 元にしたバージョン: `v2.0.1`
- ブランチ: `notebox`
- module パスは上流と同じ `charm.land/glamour/v2` のままにしてあり、NoteBox からは go.mod の `replace` でタグを指定して参照します

  ```
  replace charm.land/glamour/v2 => github.com/moriT958/glamour/v2 v2.0.1-notebox.1
  ```

## 上流からの変更

### `WithCodeBlockRenderer`

言語名を持つ fenced code block (```` ``` ```` と `~~~` のどちらも) の中身の描画を、呼び出し側で差し替えられるようにするオプションです。

```go
glamour.WithCodeBlockRenderer(func(language, code string) (string, bool) {
	// true を返すと、返した文字列がコードブロックの中身として表示される
	// false を返すと、通常どおりシンタックスハイライトして表示される
})
```

- `language` は、シンタックスハイライトに渡されるのと同じ言語名 (info string の最初の語) がそのまま渡されます
- 言語名のないコードブロックと、字下げによるコードブロックでは呼ばれません
- 差し替えた中身は、通常のコードブロックと同じ字下げ (引用やリストの中を含む) で、プレーンテキストと同じ文字色で表示されます

変更したファイルは `glamour.go`、`ansi/renderer.go`、`ansi/codeblock.go` と、テストの `codeblock_renderer_test.go` です。

## ブランチとタグ

| ブランチ | 中身 |
|---|---|
| `main` | 上流の `main` と同じ。上流と同期するだけで、変更は加えない |
| `notebox` | 上流のタグを土台に、上の変更を加えたもの |

NoteBox が参照するコミットには、`v<上流のバージョン>-notebox.<連番>` の形式でタグを付けます (例: `v2.0.1-notebox.1`)。`notebox` ブランチを rebase しても、タグが残っている限り NoteBox が参照する古いコミットは消えません。

- `notebox` ブランチに変更を加えたら、連番を 1 つ増やしたタグを付けます
- NoteBox で glamour を新しいバージョンに上げるときは、`notebox` ブランチを上流の新しいタグの上に rebase し、連番を 1 から振り直したタグを付けます

## 上流への提案

`WithCodeBlockRenderer` を上流に提案するときは、上流の最新の `main` からブランチを作り、`notebox` ブランチのうち `NOTEBOX.md` 以外のコミットを cherry-pick します。
