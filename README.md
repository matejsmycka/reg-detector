# reg-detector

Fast simple regex matcher that should be used as simple manual checker for secrets in your file.
Useful for secret detection in your codebase. 
This tool might have a high false positive rate, which is dependant on the regexes you use. 

## Usage

```bash
Usage of ./reg-detector:
  -d string
    	Directory to search for files (recursively)
  -f string
    	File containing regular expressions to search for
```

To build the binary, run `make` in the root directory of the project, output dir is `built`.

## Example

```bash
$ ./reg-detector.elf -d /your_app -f regex_dir/ics.txt

```

## Regex directory

Regex for each use-case are included in `regex_dir` directory. You can use them as a base for your own regexes.

## Dev

```bash
make run ARGS="-d SCCM/ -f regex_dir/ics.txt"
```