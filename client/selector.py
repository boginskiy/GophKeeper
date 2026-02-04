#!/usr/bin/env python3
import tkinter as tk
from tkinter import filedialog
import sys

root = tk.Tk()
root.withdraw()

file_path = filedialog.askopenfilename(
    title="Select a file",
    filetypes=[("All files", "*.*")]
)

if file_path:
    print(file_path)