# This is a VHS (https://github.com/charmbracelet/vhs) tape file.

# +-----------------------------------------------------------+
# |                         Settings                          |
# +-----------------------------------------------------------+

Output assets/demo.gif

Set Shell "bash"
Set Theme "Catppuccin Mocha"
Set FontSize 22
Set Width 1300
Set Height 650

# +-----------------------------------------------------------+
# |                 Required programs on PATH                 |
# +-----------------------------------------------------------+

Require adless

# +-----------------------------------------------------------+
# |                          START                            |
# +-----------------------------------------------------------+

Type "We currently don't block any domains.."
Sleep 2s
Backspace 38

Type "cat /etc/hosts"
Sleep 1s
Enter
Sleep 2s

Hide
Backspace 10
Type "clear"
Enter
Show

Type "Alright, it's time to enable domain blocking.."
Sleep 2s
Backspace 46

Type "sudo adless enable"
Sleep 1s
Enter
Sleep 4s

Type "cat /etc/hosts"
Sleep 1s
Enter
Sleep 3s
