# goonwin

**goonwin** is a collection of core utility tools reimplemented in Go, specifically tailored for Windows environments.  
The goal is to provide native Go versions of classic UNIX utilities.

> ⚠️ **Note:** The included `shred` utility is **not recommended for secure file deletion** since I have no idea what I'm doing

---

## Features

Currently implemented(?) utilities:

- **`cat`-like reader** – Reads and prints contents of files.
- **`wget`-like downloader** – Downloads files from URLs to the local filesystem.
- **`shred`** – Attempts to overwrite and remove files or directories (⚠️).

TODO:
- Additional reimplementations of common GNU/coreutils functionality in Go for Windows.
