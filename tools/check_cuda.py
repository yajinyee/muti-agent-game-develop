import torch
print("PyTorch version:", torch.__version__)
print("CUDA available:", torch.cuda.is_available())
print("CUDA version (torch):", torch.version.cuda)
if torch.cuda.is_available():
    print("GPU:", torch.cuda.get_device_name(0))
    print("VRAM:", torch.cuda.get_device_properties(0).total_memory // 1024**2, "MB")
else:
    print("CUDA NOT available - checking why...")
    try:
        torch.cuda.init()
    except Exception as e:
        print("cuda.init() error:", e)
