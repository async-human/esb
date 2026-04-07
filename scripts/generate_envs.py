#!/usr/bin/env python3
import sys
import argparse
from pathlib import Path

def main():
    # 1. Настройка аргументов командной строки
    parser = argparse.ArgumentParser(description="Генерация локальных .env файлов из глобального.")
    parser.add_argument("deploy_dir", type=str, help="Путь к директории deployment")
    parser.add_argument("global_env", type=str, help="Путь к глобальному .env файлу")
    args = parser.parse_args()

    # Используем pathlib для кроссплатформенной работы с путями (слэши / и \)
    deploy_dir = Path(args.deploy_dir)
    global_env_path = Path(args.global_env)
    core_dir = deploy_dir / "core"

    # 2. Проверки
    if not global_env_path.is_file():
        print(f"ОШИБКА: Глобальный файл {global_env_path} не найден!", file=sys.stderr)
        sys.exit(1)
        
    if not core_dir.is_dir():
        print(f"ОШИБКА: Директория компонентов {core_dir} не найдена!", file=sys.stderr)
        sys.exit(1)

    # 3. Читаем глобальный .env файл один раз
    # Игнорируем пустые строки и комментарии
    with open(global_env_path, 'r', encoding='utf-8') as f:
        global_lines = [
            line.strip() for line in f 
            if line.strip() and not line.strip().startswith('#')
        ]

    # 4. Извлекаем общие переменные (COMPOSE_, GLOBAL_)
    common_vars = [
        line for line in global_lines 
        if line.startswith(('COMPOSE_', 'GLOBAL_'))
    ]

    print("==> Генерация локальных .env файлов...")

    # 5. Проходим по всем папкам внутри core/
    for service_dir in core_dir.iterdir():
        if not service_dir.is_dir():
            continue

        service_name = service_dir.name
        local_env_path = service_dir / ".env"
        
        # Формируем префикс для поиска (например, "postgres" -> "POSTGRES_")
        prefix = f"{service_name.upper()}_"
        
        # Ищем переменные конкретного сервиса
        service_vars = [line for line in global_lines if line.startswith(prefix)]

        # 6. Записываем локальный .env файл
        # encoding='utf-8' важен для Windows, чтобы избежать проблем с кодировками
        with open(local_env_path, 'w', encoding='utf-8') as f:
            f.write("# АВТОГЕНЕРИРОВАНО. НЕ РЕДАКТИРОВАТЬ!\n")
            
            if common_vars:
                f.write("\n".join(common_vars) + "\n")
                
            if service_vars:
                f.write("\n".join(service_vars) + "\n")

        print(f"  Создан: {local_env_path}")

if __name__ == "__main__":
    main()