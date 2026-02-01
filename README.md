# gotodo 🚀

A lightweight, color-coded CLI task manager built in Go, designed specifically for Arch Linux power users. Features a terminal interface and a live-updating desktop widget.

## ✨ Features
- **CRUD Operations**: Add, list, update, and delete tasks.
- **Search**: Case-insensitive filtering to find tasks fast.
- **Data Integrity**: Status validation (`todo`, `in-progress`, `done`, `blocked`).
- **ANSI Colors**: Beautiful terminal output for quick status recognition.
- **Desktop Widget**: Integrated support for Conky (Wayland/X11 compatible).



---

## 📥 Installation

### 1. Prerequisites
- **Go** (for building): `sudo pacman -S go`
- **Conky** (optional, for widget): `sudo pacman -S conky`

### 2. Build and Install
Clone the repo and build the binary:
```bash
git clone [https://github.com/YOUR_USERNAME/gotodo.git](https://github.com/YOUR_USERNAME/gotodo.git)
cd gotodo
go build -o ~/.local/bin/gotodo
```

3. Setup Path
Ensure ~/.local/bin is in your $PATH. Add this to your ~/.bashrc:


```bash
export PATH=$PATH:/home/$(whoami)/.local/bin
```

## 🛠 Usage

You can interact with **gotodo** using the following commands:

| Command | Arguments | Description |
| :--- | :--- | :--- |
| `add` | `"Task description"` | Creates a new task with 'Todo' status |
| `list` | *None* | Displays all tasks with color-coded status and logs |
| `update` | `[ID] [Status]` | Updates a task (Options: todo, in-progress, done, blocked) |
| `delete` | `[ID]` | Permanently removes a task and re-indexes the list |
| `search` | `"Keyword"` | Finds tasks containing a specific word (case-insensitive) |
| `widget` | *None* | Output formatted specifically for Conky desktop display |
| `help` | *None* | Displays the built-in help menu |

**Example:**
`gotodo update 1 in-progress`

> **Note:** The `update` command is strictly validated. If you try to use a status 
> other than the four allowed types, the program will provide a helpful error 
) and list the valid options.

## 🖼 Desktop Widget (Conky)
To display your tasks on your wallpaper:

Copy the provided gotodo.conf to ~/.config/conky/.

Run Conky:

```bash
conky -c ~/.config/conky/gotodo.conf &
```

## 📄 License
MIT
