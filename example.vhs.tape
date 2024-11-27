Output example.gif

Require spin
Require ping

Set Shell "bash"
Set FontSize 16
Set Width 800
Set Height 300
Set Framerate 24
Set Theme "Cobalt2"

Type "spin"
Sleep 500ms
Enter
Sleep 3.2s
Ctrl+C

Sleep 3s

Type "ping 1.1.1.1 -c 10 | spin"
Sleep 500ms
Enter
Sleep 10s

Sleep 3s

