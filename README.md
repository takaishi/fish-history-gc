# fish-history-gc

Remove duplicate entries from fish history file.

## Example fish config

```
function history-merge --on-event fish_postexec
  history --save
  history --merge
  fish-history-gc -overwrite
end
```

## Usage

```
Usage: fish-history-gc [OPTIONS] [/fish/history/path]
  -overwrite
        Overwrite entries
```

## Install

```
$ brew tap takaishi/homebrew-fomulas
$ brew install fish-history-gc
```