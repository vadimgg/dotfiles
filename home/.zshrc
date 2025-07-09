# Aliases for common dirs
alias home="cd ~"

# System Aliases
alias ..="cd .."
alias x="exit"
alias cat='bat'
alias ls='lsd'
alias ll='ls -al'

# Git Aliases
alias add="git add"
alias commit="git commit"
alias pull="git pull"
alias stat="git status"
alias gdiff="git diff HEAD"
alias vdiff="git difftool HEAD"
alias log="git log --color --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit"
alias cfg="git --git-dir=$HOME/dotfiles/ --work-tree=$HOME"
alias push="git push"
alias g="lazygit"

export EDITOR="$(which hx)"
export VISUAL="$(which hx)"
export XDG_CONFIG_HOME="$HOME/.config"


function c() {
	local tmp="$(mktemp -t "yazi-cwd.XXXXXX")" cwd
	yazi "$@" --cwd-file="$tmp"
	if cwd="$(command cat -- "$tmp")" && [ -n "$cwd" ] && [ "$cwd" != "$PWD" ]; then
		builtin cd -- "$cwd"
	fi
	rm -f -- "$tmp"
}

pdfhx() {
  if [[ -z "$1" ]]; then
    echo "Usage: pdfhx <file.pdf>"
    return 1
  fi

  local file="$1"
  if [[ ! -f "$file" ]]; then
    echo "Error: File '$file' not found"
    return 1
  fi

  local tmpfile
  tmpfile=$(mktemp /tmp/pdfhx_XXXXXX) || {
    echo "Error: Failed to create temp file"
    return 1
  }

  local txtfile="${tmpfile}.txt"

  pdftotext -layout "$file" "$txtfile" || {
    echo "Error: pdftotext failed"
    return 1
  }

  hx "$txtfile"
  rm -f "$txtfile"
}
