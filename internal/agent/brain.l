# Decide which agent should handle a task based on type
(de decide-agent (Type)
   (case Type
      ("code" "dev-agent")
      ("research" "research-agent")
      ("audit" "security-agent")
      (T "general-agent") ) )

# Main entry point to be called from Go
(prinl (decide-agent (opt)))
(bye)
