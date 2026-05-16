#!/bin/bash
# no-direct-code.sh
# PreToolUse hook for vue-trainer agent.
# Blocks any file-editing tools and enforces instruction-only responses.

input=$(cat)
tool=$(echo "$input" | python3 -c "import sys,json; print(json.load(sys.stdin).get('toolName',''))" 2>/dev/null || echo "")

EDIT_TOOLS=(
  "str_replace_based_edit_tool"
  "replace_string_in_file"
  "multi_replace_string_in_file"
  "create_file"
  "write_file"
  "insert_edit_into_file"
  "edit_file"
)

for t in "${EDIT_TOOLS[@]}"; do
  if [[ "$tool" == "$t" ]]; then
    echo '{"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":"vue-trainer tidak boleh menulis kode langsung. Berikan instruksi TODO saja — biarkan user yang mengimplementasikan."}}'
    exit 2
  fi
done

exit 0
