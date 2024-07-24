# bravewaldo

Purpose: Bravewaldo is a fun and quirky command-line tool for processing and transforming Markdown files. It's like having a little digital helper that can spruce up your Markdown content with a few simple commands. Whether you want to wrap URLs, rewrite links, or do some other Markdown magic, Bravewaldo has got your back!

## Example Usage

```bash
# Basic usage
bravewaldo <command> [flags]

# Enable verbose mode for detailed output
bravewaldo <command> -v

# Specify log format (text or json)
bravewaldo <command> --log-format=json

# Use a custom config file
bravewaldo <command> --config=/path/to/config.yaml
```

## Install bravewaldo

On macOS/Linux:
```bash
brew install gkwa/homebrew-tools/bravewaldo
```

On Windows:
```powershell
TBD
```

## Cheatsheet

Here's a quick overview of the different commands and their functionality:

- `bravewaldo core1`: Wraps URLs in the input Markdown file with pipe characters (|).
- `bravewaldo core2`: Converts Markdown to formatted Markdown using the Goldmark library.
- `bravewaldo core3`: Converts Markdown headings to ATX style using the Goldmark library.
- `bravewaldo core4`: Wraps URLs in the input Markdown file using a custom renderer.
- `bravewaldo core5`: Processes URLs in the input Markdown file and writes the output to a file.
- `bravewaldo core6`: Extracts AutoLink URLs from the input Markdown file and prints them.
- `bravewaldo core7`: Extracts AutoLink URLs from the input Markdown file and prints them (same as core6).
- `bravewaldo core8`: Converts Markdown to formatted Markdown using the Goldmark library (similar to core2).
- `bravewaldo core9`: Converts Markdown headings to ATX style using the Goldmark library (similar to core3).
- `bravewaldo core10`: Rewrites URLs in the input Markdown file based on a predefined URL map.
- `bravewaldo core11`: Processes URLs in the input Markdown file, replacing them with friendly names if found in the URL map.

Each command has its own set of flags and options, so feel free to explore and experiment with different combinations to unlock the full potential of Bravewaldo!
