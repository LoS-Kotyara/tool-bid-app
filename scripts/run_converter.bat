@echo off
chcp 65001 >nul
title Конвертер заявок Excel в Markdown (tool-bid-app)

REM Проверка: передан ли хотя бы один файл
if "%~1"=="" (
    echo [!] Перетащите один или несколько файлов Excel *.xlsx на этот файл для конвертации.
    echo.
    pause
    exit /b 1
)

REM Проверка наличия Python в системе
where python >nul 2>&1
if %errorlevel% neq 0 (
    echo [X] Ошибка: Python не найден в системе или не добавлен в PATH!
    echo.
    pause
    exit /b 1
)

set "SCRIPT_PATH=%~dp0excel_to_md.py"

REM Проверка наличия Python-скрипта
if not exist "%SCRIPT_PATH%" (
    echo [X] Ошибка: Скрипт не найден по пути:
    echo     "%SCRIPT_PATH%"
    echo.
    pause
    exit /b 1
)

set /a SUCCESS_COUNT=0
set /a TOTAL_COUNT=0

:LOOP
if "%~1"=="" goto END

set /a TOTAL_COUNT+=1
echo ========================================================
echo [%TOTAL_COUNT%] Обработка файла: "%~nx1"
echo Путь: "%~f1"
echo ========================================================

python "%SCRIPT_PATH%" "%~f1"
if %errorlevel% equ 0 (
    set /a SUCCESS_COUNT+=1
) else (
    echo [!] Ошибка при конвертации: "%~nx1"
)

echo.
shift
goto LOOP

:END
echo ========================================================
echo Обработка завершена. Успешно: %SUCCESS_COUNT% из %TOTAL_COUNT%.
echo ========================================================
echo.
pause