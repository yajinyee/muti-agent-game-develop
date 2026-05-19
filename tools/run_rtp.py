import subprocess, sys
result = subprocess.run(
    [sys.executable, 'tools/simulate_rtp.py'],
    capture_output=True, encoding='utf-8', errors='replace'
)
# 只顯示最後 20 行
lines = (result.stdout + result.stderr).splitlines()
for line in lines[-20:]:
    print(line)
