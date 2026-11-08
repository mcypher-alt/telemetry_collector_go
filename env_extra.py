import os
Import("env")

env_file = os.path.join(env.get("PROJECT_DIR"), ".env")

if os.path.isfile(env_file):
    with open(env_file) as f:
        for line in f:
            line = line.strip()
            # Пропускаем комментарии и пустые строки
            if line and not line.startswith("#") and "=" in line:
                key, val = line.split("=", 1)
                val = val.strip().strip("'\"")
                # Передаем макрос компилятору GCC со строковыми кавычками
                env.Append(CPPDEFINES=[(key, f'\\"{val}\\"')])
