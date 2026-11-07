@echo off
chcp 65001 >nul
title Конвертер заявок Excel в Markdown (tool-bid-app)

if "%~1"=="" (
    echo [!] Перетащите один или несколько файлов Excel (.xlsx) на этот файл для конвертации.
    echo.
    pause
    exit /b 1
)

python "%~dp0excel_to_md.py" %*

echo.
pause
