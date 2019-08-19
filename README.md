# fish-history-gc

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
