#!/bin/bash
set -euo pipefail

BOXED="${BOXED:-./boxed}"

check_git_repo() {
    if ! git rev-parse --git-dir &> /dev/null; then
        echo "Error: Not a git repository"
        exit 1
    fi
}

get_current_branch() {
    git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown"
}

get_commit_hash() {
    git rev-parse --short HEAD 2>/dev/null || echo "unknown"
}

get_uncommitted_changes() {
    local modified=$(git status --porcelain 2>/dev/null | grep "^ M" | wc -l)
    local untracked=$(git status --porcelain 2>/dev/null | grep "^??" | wc -l)
    local staged=$(git status --porcelain 2>/dev/null | grep "^M" | wc -l)

    echo "${staged} staged, ${modified} modified, ${untracked} untracked"
}

get_commits_ahead_behind() {
    local upstream=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null)

    if [ -z "$upstream" ]; then
        echo "no upstream"
        return
    fi

    local ahead=$(git rev-list --count @{u}..HEAD 2>/dev/null || echo "0")
    local behind=$(git rev-list --count HEAD..@{u} 2>/dev/null || echo "0")

    echo "$ahead ahead, $behind behind"
}

get_last_commit_message() {
    git log -1 --pretty=format:"%s" 2>/dev/null | head -c 50
}

get_last_commit_author() {
    git log -1 --pretty=format:"%an" 2>/dev/null
}

determine_status() {
    local changes="$1"
    local ahead_behind="$2"

    if [[ "$changes" =~ "0 staged, 0 modified, 0 untracked" ]] && \
       [[ "$ahead_behind" =~ "0 ahead, 0 behind" || "$ahead_behind" == "no upstream" ]]; then
        echo "success"
    elif [[ "$ahead_behind" =~ "no upstream" ]]; then
        echo "warning"
    else
        echo "info"
    fi
}

main() {
    check_git_repo

    branch=$(get_current_branch)
    commit=$(get_commit_hash)
    changes=$(get_uncommitted_changes)
    ahead_behind=$(get_commits_ahead_behind)
    last_message=$(get_last_commit_message)
    last_author=$(get_last_commit_author)

    box_type=$(determine_status "$changes" "$ahead_behind")

    case "$box_type" in
        success)
            subtitle="Clean Working Tree"
            ;;
        warning)
            subtitle="No Upstream Branch"
            ;;
        *)
            subtitle="Changes Detected"
            ;;
    esac

    $BOXED "$box_type" \
        --title "Git Summary" \
        --subtitle "$subtitle" \
        --kv "Branch=$branch" \
        --kv "Commit=$commit" \
        --kv "Status=$changes" \
        --kv "Sync=$ahead_behind" \
        --kv "Last=$last_message" \
        --kv "Author=$last_author" \
        --footer "$(git config --get remote.origin.url 2>/dev/null || echo 'No remote configured')"
}

main "$@"
