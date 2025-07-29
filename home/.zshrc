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

pdf_hx() {
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


pdf_search() {
  local search_dir="${1:-.}"
  local query="$2"

  if ! command -v fzf >/dev/null 2>&1; then
    echo "Error: fzf is required but not installed. Install with: brew install fzf"
    return 1
  fi

  if ! command -v pdfgrep >/dev/null 2>&1; then
    echo "Error: pdfgrep is required. Install with: brew install pdfgrep"
    return 1
  fi

  if ! command -v pdftotext >/dev/null 2>&1; then
    echo "Error: pdftotext (from poppler) is required. Install with: brew install poppler"
    return 1
  fi

  if [[ -z "$query" ]]; then
    echo "Usage: pdf_search [directory] <search_query>"
    return 1
  fi

  # Escape query for safe highlighting
  local esc_query
  esc_query=$(printf '%s\n' "$query" | sed 's/[]\/$*.^[]/\\&/g')

  pdfgrep -r -n -H --include="*.pdf" "$query" "$search_dir" 2>/dev/null | \
  awk -F':' '{printf "%s\t%s\t%s\n", $1, $2, $3}' | \
  fzf --height=80% \
      --layout=reverse \
      --border \
      --delimiter='\t' \
      --with-nth=1,2,3 \
      --preview '
        file=$(echo {1})
        page=$(echo {2})
        line=$(echo {3})

        tmpfile=$(mktemp)
        pdftotext -f "$page" -l "$page" "$file" "$tmpfile" 2>/dev/null

        if [[ -s "$tmpfile" ]]; then
          echo "== PDF: $(basename "$file") =="
          echo "== Page: $page =="
          echo ""
          grep -i -C 3 --color=never -- "$line" "$tmpfile" | \
            sed "s/'"$esc_query"'/$(printf '"'"'\033[1;31m'"'"')&$(printf '"'"'\033[0m'"'"')/Ig"
        else
          echo "No text extracted from page $page"
        fi

        rm -f "$tmpfile"
      ' \
      --preview-window=right:60%:wrap \
      --bind="enter:execute(open {1})" \
      --bind="ctrl-o:execute(open -a Preview {1})" \
      --header="Enter: open file | Ctrl-O: open with Preview.app" \
      --prompt="PDF Search ($query): "
}

# Quick PDF info function
pdf_info() {
    local pdf_file="$1"
    
    if [[ -z "$pdf_file" ]]; then
        echo "Usage: pdf_info <pdf_file>"
        return 1
    fi
    
    if [[ ! -f "$pdf_file" ]]; then
        echo "Error: File '$pdf_file' not found"
        return 1
    fi
    
    echo "PDF Information for: $pdf_file"
    echo "===================="
    
    # Basic file info
    ls -lh "$pdf_file"
    echo ""
    
    # PDF metadata if available
    if command -v pdfinfo >/dev/null 2>&1; then
        pdfinfo "$pdf_file" 2>/dev/null
    elif command -v mdls >/dev/null 2>&1; then
        mdls "$pdf_file" | grep -E "(kMDItemTitle|kMDItemAuthors|kMDItemNumberOfPages|kMDItemCreationDate)"
    fi
}

export NVM_DIR="$HOME/.config/nvm"
export PATH="$HOME/.local/bin:$PATH"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion


